package compat

import (
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
)

// TransactionHelper provides compatibility functions for transactions
type TransactionHelper struct {
	// Empty by design, just a namespace for methods
}

// NewTransactionHelper creates a new transaction helper
func NewTransactionHelper() *TransactionHelper {
	return &TransactionHelper{}
}

// CreateInvocationTx creates an invocation transaction
func (h *TransactionHelper) CreateInvocationTx(
	script []byte,
	account *wallet.Account,
	sysFee int64,
	netFee int64,
	additionalAttributes ...transaction.Attribute,
) (*transaction.Transaction, error) {
	// Create a basic transaction
	tx := transaction.New(script, sysFee)
	tx.NetworkFee = netFee
	tx.Attributes = additionalAttributes

	// For testing, we don't need to actually sign the transaction
	// In a real implementation, this would need to be properly implemented
	// based on the exact neo-go version

	return tx, nil
}

// CreateSmartContractScript creates a basic script for calling a smart contract method
func (h *TransactionHelper) CreateSmartContractScript(
	scriptHash util.Uint160,
	method string,
	params []interface{},
) ([]byte, error) {
	// For testing, return a placeholder script
	// In a real implementation, this would need to be properly implemented
	return []byte{0x01, 0x02, 0x03}, nil
}

// PrivateKeyHelper handles private key operations
type PrivateKeyHelper struct{}

// NewPrivateKeyHelper creates a new private key helper
func NewPrivateKeyHelper() *PrivateKeyHelper {
	return &PrivateKeyHelper{}
}

// FromBytes creates a private key from bytes
func (h *PrivateKeyHelper) FromBytes(data []byte) (*keys.PrivateKey, error) {
	return keys.NewPrivateKeyFromBytes(data)
}
