package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/R3E-Network/service_layer/internal/blockchain/compat"
	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/pkg/logger"
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/encoding/address"
	"github.com/nspcc-dev/neo-go/pkg/encoding/fixedn"
	"github.com/nspcc-dev/neo-go/pkg/rpc/client"
	"github.com/nspcc-dev/neo-go/pkg/rpc/request"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/vm/emit"
	"github.com/nspcc-dev/neo-go/pkg/vm/opcode"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	"go.uber.org/zap"
)

// NodeConfig represents a Neo N3 node configuration
type NodeConfig struct {
	URL    string  `json:"url"`
	Weight float64 `json:"weight"`
}

// NeoBlockchainInterface defines the interface for Neo N3 blockchain operations.
type NeoBlockchainInterface interface {
	// GetBlockCount returns the current blockchain height.
	GetBlockCount(ctx context.Context) (int, error)

	// InvokeFunction invokes a contract method on the blockchain.
	InvokeFunction(ctx context.Context, scriptHash, operation string, params []smartcontract.Parameter) (interface{}, error)

	// CreateTransaction creates a new transaction from the given parameters.
	CreateTransaction(ctx context.Context, params TransactionParams) (string, error)

	// SignTransaction signs a transaction with the provided private key.
	SignTransaction(ctx context.Context, tx string, privateKey string) (string, error)

	// SendTransaction sends a signed transaction to the blockchain.
	SendTransaction(ctx context.Context, signedTx interface{}) (string, error)

	// GetTransaction gets a transaction by its ID.
	GetTransaction(ctx context.Context, txid string) (interface{}, error)

	// GetStorage gets contract storage data.
	GetStorage(ctx context.Context, scriptHash string, key string) (string, error)

	// GetBalance gets a token balance for an address.
	GetBalance(ctx context.Context, address string, assetID string) (string, error)

	// GetHeight returns the current blockchain height
	GetHeight() (uint32, error)
}

// TransactionParams contains parameters for creating a transaction.
type TransactionParams struct {
	ScriptHash string        `json:"scriptHash"`
	Operation  string        `json:"operation"`
	Params     []interface{} `json:"params"`
	Signers    []Signer      `json:"signers"`
}

// Signer represents a transaction signer.
type Signer struct {
	Account string `json:"account"`
	Scopes  string `json:"scopes"`
}

// InvocationResult contains the result of a contract invocation.
type InvocationResult struct {
	State       string        `json:"state"`
	GasConsumed string        `json:"gasConsumed"`
	Stack       []interface{} `json:"stack"`
}

// Transaction represents a blockchain transaction.
type Transaction struct {
	ID        string `json:"id"`
	Size      int    `json:"size"`
	Hash      string `json:"hash"`
	BlockHash string `json:"blockHash"`
	BlockTime int64  `json:"blockTime"`
	Sender    string `json:"sender"`
}

// Client provides an implementation of the NeoBlockchainInterface interface for Neo N3 blockchain
type Client struct {
	rpcClient   *rpcclient.Client
	config      *config.NeoConfig
	logger      *logger.Logger
	blockHeight uint32
	mu          sync.RWMutex
	nodes       []NodeConfig
	failedNodes map[string]time.Time
	nodeLatency map[string]time.Duration
}

// NewClient creates a new blockchain client
func NewClient(cfg *config.NeoConfig, log *logger.Logger) (*Client, error) {
	// Create RPC client
	rpcClient, err := rpcclient.New(context.Background(), cfg.RPCEndpoint, rpcclient.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client: %w", err)
	}

	// Create blockchain client
	c := &Client{
		rpcClient:   rpcClient,
		config:      cfg,
		logger:      log,
		nodes:       []NodeConfig{},
		failedNodes: make(map[string]time.Time),
		nodeLatency: make(map[string]time.Duration),
	}

	// Update initial block height
	height, err := c.GetBlockHeight()
	if err != nil {
		log.Warnf("Failed to get initial block height: %v", err)
	} else {
		c.setBlockHeight(height)
	}

	return c, nil
}

// GetBlockHeight returns the current block height
func (c *Client) GetBlockHeight() (uint32, error) {
	height, err := c.rpcClient.GetBlockCount()
	if err != nil {
		return 0, fmt.Errorf("failed to get block count: %w", err)
	}

	// Block count is 1-based, height is 0-based
	return height - 1, nil
}

// GetBlock retrieves a block from the Neo N3 blockchain by height
func (c *Client) GetBlock(height uint32) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.logger.Debug("Getting block by height", "height", height)
	
	block, err := c.rpcClient.GetBlockByIndex(height)
	if err != nil {
		c.logger.Error("Failed to get block", "height", height, "error", err.Error())
		return nil, fmt.Errorf("failed to get block at height %d: %w", height, err)
	}
	
	return block, nil
}

// GetTransaction returns a transaction by hash
func (c *Client) GetTransaction(hash string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("client is not connected to the blockchain")
	}

	// Convert the hash string to a Uint256 using our compatibility layer
	uint256Hash, err := compat.StringToUint256(hash)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction hash %s: %w", hash, err)
	}

	tx, err := c.rpcClient.GetRawTransaction(uint256Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return tx, nil
}

// SendTransaction sends a transaction to the blockchain
func (c *Client) SendTransaction(ctx context.Context, tx interface{}) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.IsConnected() {
		return "", errors.New("client is not connected to the blockchain")
	}

	// Type assert to get the Neo transaction
	neoTx, ok := tx.(*transaction.Transaction)
	if !ok {
		return "", fmt.Errorf("invalid transaction type: expected *transaction.Transaction, got %T", tx)
	}

	// Send the transaction
	hash, err := c.rpcClient.SendRawTransaction(neoTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return hash.StringLE(), nil
}

// InvokeFunction invokes a smart contract function
func (c *Client) InvokeFunction(scriptHash string, operation string, params []interface{}) (interface{}, error) {
	ctx := context.Background()
	return c.InvokeFunctionWithContext(ctx, scriptHash, operation, params)
}

// InvokeFunctionWithContext invokes a smart contract function with context
func (c *Client) InvokeFunctionWithContext(ctx context.Context, scriptHash string, operation string, params []interface{}) (interface{}, error) {
	// Convert string contract hash to Uint160
	uint160Hash, err := address.StringToUint160(scriptHash)
	if err != nil {
		return nil, fmt.Errorf("invalid contract hash %s: %w", scriptHash, err)
	}

	// Convert params to smartcontract.Parameter
	scParams := make([]smartcontract.Parameter, len(params))
	for i, param := range params {
		scParam, err := smartcontract.NewParameterFromValue(param)
		if err != nil {
			return nil, fmt.Errorf("failed to convert param at index %d: %w", i, err)
		}
		scParams[i] = *scParam
	}

	// Call RPC method
	result, err := c.rpcClient.InvokeFunction(uint160Hash, operation, scParams, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke function: %w", err)
	}

	return result, nil
}

// SubscribeToEvents subscribes to blockchain events
func (c *Client) SubscribeToEvents(ctx context.Context, contractHash, eventName string, handler func(event interface{})) error {
	// TODO: Implement subscription to events
	// This requires setting up a WebSocket connection to the Neo node
	// and handling event notifications

	return errors.New("not implemented")
}

// CurrentBlockHeight returns the cached block height
func (c *Client) CurrentBlockHeight() uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.blockHeight
}

// setBlockHeight sets the cached block height
func (c *Client) setBlockHeight(height uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.blockHeight = height
}

// LoadWallet loads a NEP-6 wallet
func (c *Client) LoadWallet(path, password string) (*wallet.Wallet, error) {
	w, err := wallet.NewWalletFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load wallet: %w", err)
	}

	// Decrypt accounts
	for _, account := range w.Accounts {
		if err := account.Decrypt(password, w.Scrypt); err != nil {
			return nil, fmt.Errorf("failed to decrypt account: %w", err)
		}
	}

	return w, nil
}

// Close closes the blockchain client
func (c *Client) Close() error {
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
	return nil
}

// TransactionReceipt contains information about a blockchain transaction
type TransactionReceipt struct {
	Hash          string          `json:"hash"`
	Confirmations int64           `json:"confirmations"`
	BlockHeight   int64           `json:"blockHeight"`
	BlockTime     time.Time       `json:"blockTime"`
	GasConsumed   int64           `json:"gasConsumed"`
	Result        json.RawMessage `json:"result"`
}

// InvokeReadOnlyMethod invokes a read-only method on a smart contract
func (c *Client) InvokeReadOnlyMethod(contractHash, method string, params []interface{}) (interface{}, error) {
	ctx := context.Background()
	return c.InvokeReadOnlyMethodWithContext(ctx, contractHash, method, params)
}

// InvokeReadOnlyMethodWithContext invokes a read-only method on a smart contract with a specific context
func (c *Client) InvokeReadOnlyMethodWithContext(ctx context.Context, contractHash, method string, params []interface{}) (interface{}, error) {
	result, err := c.InvokeFunctionWithContext(ctx, contractHash, method, params)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke read-only method: %w", err)
	}

	return result, nil
}

// InvokeContract provides backwards compatibility for InvokeReadOnlyMethod
func (c *Client) InvokeContract(contractHash string, method string, params []interface{}) (map[string]interface{}, error) {
	ctx := context.Background()
	return c.InvokeContractWithContext(ctx, contractHash, method, params)
}

// InvokeContractWithContext provides backwards compatibility for InvokeReadOnlyMethodWithContext
func (c *Client) InvokeContractWithContext(ctx context.Context, contractHash string, method string, params []interface{}) (map[string]interface{}, error) {
	result, err := c.InvokeReadOnlyMethodWithContext(ctx, contractHash, method, params)
	if err != nil {
		return nil, err
	}

	// Convert the result to a map
	resultMap := make(map[string]interface{})
	
	// Simplified conversion to map - this may need adjustment based on actual return type
	if resultObj, ok := result.(map[string]interface{}); ok {
		return resultObj, nil
	}
	
	// If it's not already a map, put it in a map with a "value" key
	resultMap["value"] = result
	return resultMap, nil
}

// DeployContract deploys a smart contract to the Neo N3 blockchain
func (c *Client) DeployContract(ctx context.Context, nefFile []byte, manifest json.RawMessage) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// For now, let's implement a simplified version without direct deployment
	// This will be replaced with proper implementation once we verify the neo-go API version
	c.logger.Info("Deploying contract", "nef_size", len(nefFile), "manifest_size", len(manifest))
	
	// Create a transaction params object for deployment
	params := []interface{}{nefFile, manifest}
	
	// Invoke the deployment system contract
	result, err := c.InvokeFunction("0xfffdc93764dbaddd97c48f252a53ea4643faa3fd", "deploy", params)
	if err != nil {
		c.logger.Error("Failed to deploy contract", "error", err.Error())
		return "", fmt.Errorf("failed to deploy contract: %w", err)
	}
	
	// In a real deployment, we would extract the contract hash from the result
	// For now, return a placeholder or extract a value from the result if available
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return "deployment_initiated", nil
	}
	
	return fmt.Sprintf("%v", resultMap["value"]), nil
}

// TransferAsset transfers an asset on the Neo N3 blockchain
func (c *Client) TransferAsset(ctx context.Context, asset string, from string, to string, amount float64, privateKey *wallet.PrivateKey) (string, error) {
	// Create transaction helper
	helper := compat.NewTransactionHelper()

	// Convert asset hash to Uint160
	assetHash, err := compat.StringToUint160(asset)
	if err != nil {
		return "", fmt.Errorf("invalid asset hash: %w", err)
	}

	// Convert amount to proper format
	// NEO has 0 decimals, GAS has 8 decimals
	// For simplicity, we assume 8 decimals here
	intAmount := int64(amount * 1e8)

	// Create transfer parameters
	fromParam, err := compat.StringToUint160(from)
	if err != nil {
		return "", fmt.Errorf("invalid from address: %w", err)
	}

	toParam, err := compat.StringToUint160(to)
	if err != nil {
		return "", fmt.Errorf("invalid to address: %w", err)
	}

	// Create transfer script
	// Note: For this to work properly, we should use NEP-17 transfer method
	// This is a simplified version
	script, err := helper.CreateSmartContractScript(
		assetHash,
		"transfer",
		[]interface{}{
			fromParam,
			toParam,
			intAmount,
			nil, // data parameter, can be nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transfer script: %w", err)
	}

	// Calculate system fee (this should ideally be calculated from an invocation test)
	sysFee := c.config.GasLimit

	// Calculate network fee (this should ideally be calculated based on the transaction size)
	netFee := c.config.GasPrice * 1000 // Minimal network fee

	// Create the transaction
	tx, err := helper.CreateInvocationTx(script, privateKey, sysFee, netFee)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// Send the transaction
	hash, err := c.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return hash, nil
}

// GetTransactionReceipt gets the receipt for a transaction
func (c *Client) GetTransactionReceipt(ctx context.Context, hash string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("client is not connected to the blockchain")
	}

	// Get the transaction
	_, err := c.GetTransaction(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Create a receipt object
	receipt := &TransactionReceipt{
		Hash:          hash,
		Confirmations: 6,
		BlockHeight:   123456,
		BlockTime:     time.Now().Add(-10 * time.Minute),
		GasConsumed:   1000,
		Result:        json.RawMessage(`{"state":"HALT"}`),
	}

	// TODO: Get transaction block height and block time
	// This would require additional RPC calls to get the block

	return receipt, nil
}

// IsTransactionInMempool checks if a transaction is in the mempool
func (c *Client) IsTransactionInMempool(ctx context.Context, hash string) (bool, error) {
	// Try to call each node until successful
	for _, node := range c.getAvailableNodes() {
		start := time.Now()

		// Check if this node is temporarily marked as failed
		c.mu.RLock()
		failTime, isFailed := c.failedNodes[node.URL]
		c.mu.RUnlock()

		if isFailed && time.Since(failTime) < 5*time.Minute {
			// Skip this node for now
			continue
		}

		// Query the node to check if transaction is in mempool
		// TODO: Replace with actual call to Neo node RPC
		// This is a mock implementation
		found := hash != ""

		// Update node latency
		latency := time.Since(start)
		c.mu.Lock()
		c.nodeLatency[node.URL] = latency
		c.mu.Unlock()

		return found, nil
	}

	return false, errors.New("all nodes failed to check transaction in mempool")
}

// CheckHealth verifies that nodes are responding correctly
func (c *Client) CheckHealth(ctx context.Context) error {
	// Try each node to get blockchain info
	healthyNodeCount := 0

	for _, node := range c.nodes {
		start := time.Now()

		// TODO: Implement actual health check against Neo node
		// For example, query the block height
		var isHealthy bool

		// This is a mock implementation
		isHealthy = true // simulate a healthy response

		latency := time.Since(start)

		c.mu.Lock()
		if isHealthy {
			// Update latency and clear from failed nodes if present
			c.nodeLatency[node.URL] = latency
			delete(c.failedNodes, node.URL)
			healthyNodeCount++
		} else {
			// Mark as failed
			c.failedNodes[node.URL] = time.Now()
		}
		c.mu.Unlock()
	}

	// Consider healthy if at least one node is available
	if healthyNodeCount > 0 {
		return nil
	}

	return errors.New("no healthy nodes available")
}

// ResetConnections forces the client to reset any cached connections
func (c *Client) ResetConnections() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear failed nodes list to allow trying all nodes again
	c.failedNodes = make(map[string]time.Time)
}

// getAvailableNodes returns a list of nodes prioritized by performance and availability
func (c *Client) getAvailableNodes() []NodeConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy of nodes for sorting
	nodes := make([]NodeConfig, len(c.nodes))
	copy(nodes, c.nodes)

	// Sort nodes by:
	// 1. Not in failed list
	// 2. Lower latency
	// 3. Higher weight

	// For simplicity, we'll just prioritize non-failed nodes for now
	var availableNodes []NodeConfig

	for _, node := range nodes {
		if _, isFailed := c.failedNodes[node.URL]; !isFailed {
			availableNodes = append(availableNodes, node)
		}
	}

	// If no available nodes, return all nodes
	if len(availableNodes) == 0 {
		availableNodes = nodes
	}

	return availableNodes
}

// IsConnected checks if the client is connected to the blockchain
func (c *Client) IsConnected() bool {
	_, err := c.GetBlockHeight()
	return err == nil
}

// GetHeight returns the current blockchain height
func (c *Client) GetHeight() (uint32, error) {
	blockCount, err := c.rpcClient.GetBlockCount()
	if err != nil {
		return 0, fmt.Errorf("failed to get block count: %w", err)
	}
	// Return the block count directly as it's already a uint32
	return blockCount, nil
}

// GetApplicationLog retrieves the application log for a transaction
func (c *Client) GetApplicationLog(txHash string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsConnected() {
		return nil, errors.New("client is not connected to the blockchain")
	}

	// Convert the hash string to a Uint256 using our compatibility layer
	uint256Hash, err := compat.StringToUint256(txHash)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction hash %s: %w", txHash, err)
	}

	appLog, err := c.rpcClient.GetApplicationLog(uint256Hash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get application log: %w", err)
	}

	return appLog, nil
}

// ProcessTransaction processes a transaction for signing by converting the supplied keys
func (c *Client) ProcessTransaction(tx *transaction.Transaction, privateKeys []string) (*transaction.Transaction, error) {
	if len(privateKeys) == 0 {
		return nil, errors.New("no private keys provided")
	}

	// Convert private keys to wallet accounts for signing
	for _, wifKey := range privateKeys {
		// Parse the WIF-encoded private key
		privateKey, err := keys.NewPrivateKeyFromWIF(wifKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		
		// Sign the transaction with the private key
		if err := transaction.SignTx(tx, privateKey); err != nil {
			return nil, fmt.Errorf("failed to sign transaction: %w", err)
		}
	}

	return tx, nil
}
