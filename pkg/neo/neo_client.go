package neo

import (
	"fmt"
	"time"

	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"github.com/joeqian10/neo3-gogogo/sc"
	"github.com/joeqian10/neo3-gogogo/tx"
	"github.com/joeqian10/neo3-gogogo/wallet"
)

// Client provides extended functionality for Neo N3 blockchain
type Client struct {
	RpcClient    *rpc.RpcClient
	NetworkMagic uint32
	Account      *wallet.Account
}

// NewClient creates a new Neo N3 client
func NewClient(rpcEndpoint string, networkMagic uint32) *Client {
	return &Client{
		RpcClient:    rpc.NewClient(rpcEndpoint),
		NetworkMagic: networkMagic,
	}
}

// SetAccount configures the client with a wallet account
func (c *Client) SetAccount(account *wallet.Account) {
	c.Account = account
}

// GetNeoBalance retrieves NEO balance for an address
func (c *Client) GetNeoBalance(address string) (float64, error) {
	scriptHash, err := helper.AddressToScriptHash(address, c.NetworkMagic)
	if err != nil {
		return 0, err
	}

	neoAssetHash, err := helper.UInt160FromString("ef4073a0f2b305a38ec4050e4d3d28bc40ea63f5")
	if err != nil {
		return 0, err
	}

	response, err := c.RpcClient.InvokeFunction(
		neoAssetHash.String(),
		"balanceOf",
		[]sc.ContractParameter{
			{
				Type:  sc.Hash160,
				Value: scriptHash,
			},
		},
		nil,
	)

	if err != nil {
		return 0, err
	}

	if response.State == "FAULT" {
		return 0, fmt.Errorf("contract execution failed: %s", response.Exception)
	}

	if len(response.Stack) == 0 {
		return 0, fmt.Errorf("empty response stack")
	}

	balance, err := sc.ParseNeoVMStack(response.Stack[0])
	if err != nil {
		return 0, err
	}

	balanceInt, ok := balance.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid balance format")
	}

	return float64(balanceInt), nil
}

// GetGasBalance retrieves GAS balance for an address
func (c *Client) GetGasBalance(address string) (float64, error) {
	scriptHash, err := helper.AddressToScriptHash(address, c.NetworkMagic)
	if err != nil {
		return 0, err
	}

	gasAssetHash, err := helper.UInt160FromString("d2a4cff31913016155e38e474a2c06d08be276cf")
	if err != nil {
		return 0, err
	}

	response, err := c.RpcClient.InvokeFunction(
		gasAssetHash.String(),
		"balanceOf",
		[]sc.ContractParameter{
			{
				Type:  sc.Hash160,
				Value: scriptHash,
			},
		},
		nil,
	)

	if err != nil {
		return 0, err
	}

	if response.State == "FAULT" {
		return 0, fmt.Errorf("contract execution failed: %s", response.Exception)
	}

	if len(response.Stack) == 0 {
		return 0, fmt.Errorf("empty response stack")
	}

	balance, err := sc.ParseNeoVMStack(response.Stack[0])
	if err != nil {
		return 0, err
	}

	balanceInt, ok := balance.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid balance format")
	}

	// GAS has 8 decimals
	return float64(balanceInt) / 100000000, nil
}

// TransferGas sends GAS tokens
func (c *Client) TransferGas(toAddress string, amount float64) (string, error) {
	if c.Account == nil {
		return "", fmt.Errorf("account not set")
	}

	// Convert Neo N3 address to scripthash
	toScriptHash, err := helper.AddressToScriptHash(toAddress, c.NetworkMagic)
	if err != nil {
		return "", err
	}

	// GAS asset hash on Neo N3
	gasAssetHash, err := helper.UInt160FromString("d2a4cff31913016155e38e474a2c06d08be276cf")
	if err != nil {
		return "", err
	}

	// Convert amount to integer (GAS has 8 decimals)
	amountInt := int64(amount * 100000000)

	// Create a transaction
	tb := tx.NewTransactionBuilder(c.NetworkMagic)

	// Add transfer script
	script, err := sc.MakeScript(
		gasAssetHash,
		"transfer",
		[]sc.ContractParameter{
			{Type: sc.Hash160, Value: c.Account.ScriptHash},
			{Type: sc.Hash160, Value: toScriptHash},
			{Type: sc.Integer, Value: amountInt},
			{Type: sc.Any, Value: nil},
		},
	)
	if err != nil {
		return "", err
	}

	tb.Script = script

	// Get current block count
	blockCount, err := c.RpcClient.GetBlockCount()
	if err != nil {
		return "", err
	}

	// Set validUntilBlock
	tb.ValidUntilBlock = blockCount.Result + 100

	// Create transaction
	transaction := tb.GetTransaction()

	// Calculate network fee
	networkFee, err := c.CalculateNetworkFee(transaction)
	if err != nil {
		return "", err
	}
	transaction.NetworkFee = networkFee

	// Calculate system fee
	invokeResult, err := c.RpcClient.InvokeScript(helper.BytesToHex(script), nil)
	if err != nil {
		return "", err
	}
	systemFee := invokeResult.GasConsumed
	transaction.SystemFee = int64(systemFee)

	// Sign transaction
	privateKey := c.Account.KeyPair.PrivateKey
	signature, err := crypto.Sign(transaction.GetHashData(), privateKey)
	if err != nil {
		return "", err
	}

	// Add signature to witnesses
	witness := tx.Witness{
		InvocationScript:   crypto.CreateSignatureInvocationScript(signature),
		VerificationScript: c.Account.Contract.Script,
	}

	transaction.Witnesses = []*tx.Witness{&witness}

	// Send transaction
	response, err := c.RpcClient.SendRawTransaction(transaction.ToByteArray())
	if err != nil {
		return "", err
	}

	if !response.Result {
		return "", fmt.Errorf("transaction rejected: %v", response.Error)
	}

	// Return the transaction hash
	return transaction.GetHash().String(), nil
}

// CalculateNetworkFee estimates the network fee for a transaction
func (c *Client) CalculateNetworkFee(transaction *tx.Transaction) (int64, error) {
	// In a real implementation, this would calculate the network fee accurately
	// For simplicity, we'll use a fixed fee
	return 1000000, nil
}

// WaitForTransaction waits for a transaction to be confirmed
func (c *Client) WaitForTransaction(txHash string, timeout time.Duration) (bool, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		response, err := c.RpcClient.GetRawTransaction(txHash)
		if err != nil {
			// Transaction not found yet, wait and retry
			time.Sleep(2 * time.Second)
			continue
		}

		if response.Error != nil {
			return false, fmt.Errorf("error getting transaction: %v", response.Error)
		}

		// Transaction found
		if response.Result.Confirmations > 0 {
			return true, nil
		}

		time.Sleep(2 * time.Second)
	}

	return false, fmt.Errorf("transaction confirmation timeout")
}
