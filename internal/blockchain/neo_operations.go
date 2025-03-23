package blockchain

import (
	"fmt"
	"log"
	"math/big"

	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/sc"
	"github.com/joeqian10/neo3-gogogo/tx"
	"github.com/your-org/neo-oracle/internal/models"
)

// NeoOperations provides Neo N3-specific blockchain operations
type NeoOperations struct {
	client *Client
}

// NewNeoOperations creates a new Neo operations handler
func NewNeoOperations(client *Client) *NeoOperations {
	return &NeoOperations{
		client: client,
	}
}

// UpdatePriceOracle sends price updates to the blockchain oracle contract
func (n *NeoOperations) UpdatePriceOracle(prices map[string]*models.PriceData) error {
	if !n.client.IsConnected() {
		return fmt.Errorf("not connected to blockchain")
	}

	if n.client.account == nil {
		return fmt.Errorf("no account loaded")
	}

	// Get oracle contract hash
	oracleContractHash, err := helper.UInt160FromString(n.client.cfg.OracleContract)
	if err != nil {
		return fmt.Errorf("invalid oracle contract hash: %w", err)
	}

	// Build transaction for each price
	for token, priceData := range prices {
		// Convert price to fixed point integer (8 decimals)
		priceBig := big.NewFloat(priceData.Price)
		priceBig = priceBig.Mul(priceBig, big.NewFloat(100000000)) // 8 decimals
		priceInt, _ := priceBig.Int64()

		// Create script for UpdatePrice method
		script, err := sc.MakeScript(
			oracleContractHash,
			"UpdatePrice",
			[]sc.ContractParameter{
				{Type: sc.String, Value: token},
				{Type: sc.Integer, Value: priceInt},
				{Type: sc.Integer, Value: priceData.Timestamp.Unix()},
			},
		)
		if err != nil {
			log.Printf("Error creating price update script for %s: %v", token, err)
			continue
		}

		// Create transaction
		tx, err := n.createAndSignTransaction(script)
		if err != nil {
			log.Printf("Error creating transaction for %s price update: %v", token, err)
			continue
		}

		// Send transaction
		response, err := n.client.rpcClient.SendRawTransaction(tx.ToByteArray())
		if err != nil {
			log.Printf("Error sending transaction for %s price update: %v", token, err)
			continue
		}

		if !response.Result {
			log.Printf("Transaction rejected for %s price update: %v", token, response.Error)
			continue
		}

		log.Printf("Price updated for %s: %.8f (tx: %s)",
			token, priceData.Price, tx.GetHash().String())
	}

	return nil
}

// RecordFunctionExecution logs a function execution on the blockchain
func (n *NeoOperations) RecordFunctionExecution(functionID string, result string) error {
	if !n.client.IsConnected() {
		return fmt.Errorf("not connected to blockchain")
	}

	if n.client.account == nil {
		return fmt.Errorf("no account loaded")
	}

	// Get oracle contract hash
	oracleContractHash, err := helper.UInt160FromString(n.client.cfg.OracleContract)
	if err != nil {
		return fmt.Errorf("invalid oracle contract hash: %w", err)
	}

	// Create script for RecordFunctionExecution method
	script, err := sc.MakeScript(
		oracleContractHash,
		"RecordFunctionExecution",
		[]sc.ContractParameter{
			{Type: sc.String, Value: functionID},
			{Type: sc.String, Value: result},
		},
	)
	if err != nil {
		return fmt.Errorf("error creating script: %w", err)
	}

	// Create transaction
	tx, err := n.createAndSignTransaction(script)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	// Send transaction
	response, err := n.client.rpcClient.SendRawTransaction(tx.ToByteArray())
	if err != nil {
		return fmt.Errorf("error sending transaction: %w", err)
	}

	if !response.Result {
		return fmt.Errorf("transaction rejected: %v", response.Error)
	}

	log.Printf("Function execution recorded on blockchain: %s (tx: %s)",
		functionID, tx.GetHash().String())

	return nil
}

// AllocateGasBankFunds requests gas allocation from the GasBank contract
func (n *NeoOperations) AllocateGasBankFunds(userAddress string, functionID string, amount float64) error {
	if !n.client.IsConnected() {
		return fmt.Errorf("not connected to blockchain")
	}

	if n.client.account == nil {
		return fmt.Errorf("no account loaded")
	}

	// Get gasbank contract hash
	gasBankHash, err := helper.UInt160FromString(n.client.cfg.GasBankContract)
	if err != nil {
		return fmt.Errorf("invalid gasbank contract hash: %w", err)
	}

	// Convert user address to script hash
	userScriptHash, err := helper.AddressToScriptHash(userAddress, n.client.cfg.NetworkMagic)
	if err != nil {
		return fmt.Errorf("invalid user address: %w", err)
	}

	// Convert amount to fixed point integer (8 decimals)
	amountBig := big.NewFloat(amount)
	amountBig = amountBig.Mul(amountBig, big.NewFloat(100000000)) // 8 decimals
	amountInt, _ := amountBig.Int64()

	// Create script for AllocateGas method
	script, err := sc.MakeScript(
		gasBankHash,
		"AllocateGas",
		[]sc.ContractParameter{
			{Type: sc.Hash160, Value: userScriptHash},
			{Type: sc.String, Value: functionID},
			{Type: sc.Integer, Value: amountInt},
		},
	)
	if err != nil {
		return fmt.Errorf("error creating script: %w", err)
	}

	// Create transaction
	tx, err := n.createAndSignTransaction(script)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	// Send transaction
	response, err := n.client.rpcClient.SendRawTransaction(tx.ToByteArray())
	if err != nil {
		return fmt.Errorf("error sending transaction: %w", err)
	}

	if !response.Result {
		return fmt.Errorf("transaction rejected: %v", response.Error)
	}

	log.Printf("Gas allocated for user %s, function %s: %.8f GAS (tx: %s)",
		userAddress, functionID, amount, tx.GetHash().String())

	return nil
}

// createAndSignTransaction creates and signs a Neo N3 transaction
func (n *NeoOperations) createAndSignTransaction(script []byte) (*tx.Transaction, error) {
	// Create transaction builder
	tb := tx.NewTransactionBuilder(n.client.cfg.NetworkMagic)
	tb.Script = script

	// Get current block count
	blockCount, err := n.client.rpcClient.GetBlockCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get block count: %w", err)
	}

	// Set validUntilBlock (current + 100)
	tb.ValidUntilBlock = blockCount.Result + 100

	// Create transaction
	transaction := tb.GetTransaction()

	// Calculate network fee
	networkFee, err := n.calculateNetworkFee(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate network fee: %w", err)
	}
	transaction.NetworkFee = networkFee

	// Calculate system fee
	invokeResult, err := n.client.rpcClient.InvokeScript(helper.BytesToHex(script), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke script: %w", err)
	}
	systemFee := invokeResult.GasConsumed
	transaction.SystemFee = int64(systemFee)

	// Sign transaction
	privateKey := n.client.account.KeyPair.PrivateKey
	signature, err := crypto.Sign(transaction.GetHashData(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Add signature to witness
	witness := tx.Witness{
		InvocationScript:   crypto.CreateSignatureInvocationScript(signature),
		VerificationScript: n.client.account.Contract.Script,
	}
	transaction.Witnesses = []*tx.Witness{&witness}

	return transaction, nil
}

// calculateNetworkFee estimates the network fee for a transaction
func (n *NeoOperations) calculateNetworkFee(transaction *tx.Transaction) (int64, error) {
	// In a real implementation, this would calculate the network fee
	// based on the transaction size and the number of signatures

	// For now, use a fixed fee
	return 1000000, nil // 0.01 GAS (8 decimals)
}
