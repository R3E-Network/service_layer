package blockchain

import (
	"context"
	"encoding/json"
)

// ClientAdapter adapts the Client struct to implement the BlockchainClient interface
type ClientAdapter struct {
	client *Client
}

// NewClientAdapter creates a new adapter that makes a Client implement BlockchainClient
func NewClientAdapter(client *Client) BlockchainClient {
	return &ClientAdapter{client: client}
}

// GetBlockHeight returns the current block height
func (a *ClientAdapter) GetBlockHeight() (uint32, error) {
	return a.client.GetHeight()
}

// GetBlock returns a block by height
func (a *ClientAdapter) GetBlock(height uint32) (interface{}, error) {
	return a.client.GetBlock(height)
}

// GetTransaction returns a transaction by hash
func (a *ClientAdapter) GetTransaction(hash string) (interface{}, error) {
	return a.client.GetTransaction(hash)
}

// SendTransaction sends a transaction to the blockchain
func (a *ClientAdapter) SendTransaction(tx interface{}) (string, error) {
	// Type assertion may be needed depending on the actual implementation
	// This is a simplified version
	return "", nil
}

// InvokeContract invokes a contract on the blockchain
func (a *ClientAdapter) InvokeContract(contractHash string, method string, params []interface{}) (map[string]interface{}, error) {
	return a.client.InvokeContract(contractHash, method, params)
}

// DeployContract deploys a contract to the blockchain
func (a *ClientAdapter) DeployContract(ctx context.Context, nefFile []byte, manifest json.RawMessage) (string, error) {
	// Placeholder for implementation
	// Will need specific handling based on the actual Client implementation
	return "", nil
}

// SubscribeToEvents subscribes to blockchain events
func (a *ClientAdapter) SubscribeToEvents(ctx context.Context, contractHash, eventName string, handler func(event interface{})) error {
	return a.client.SubscribeToEvents(ctx, contractHash, eventName, handler)
}

// GetTransactionReceipt gets a transaction receipt
func (a *ClientAdapter) GetTransactionReceipt(ctx context.Context, hash string) (interface{}, error) {
	return a.client.GetTransactionReceipt(ctx, hash)
}

// IsTransactionInMempool checks if a transaction is in the mempool
func (a *ClientAdapter) IsTransactionInMempool(ctx context.Context, hash string) (bool, error) {
	return a.client.IsTransactionInMempool(ctx, hash)
}

// CheckHealth checks if the blockchain client is healthy
func (a *ClientAdapter) CheckHealth(ctx context.Context) error {
	return a.client.CheckHealth(ctx)
}

// ResetConnections resets connections to the blockchain
func (a *ClientAdapter) ResetConnections() {
	a.client.ResetConnections()
}

// Close closes the blockchain client
func (a *ClientAdapter) Close() error {
	return a.client.Close()
}
