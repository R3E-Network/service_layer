package secrets

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/your-org/neo-oracle/internal/config"
	"github.com/your-org/neo-oracle/internal/models"
)

// Service manages secure storage of user secrets
type Service struct {
	cfg      *config.TEEConfig
	secrets  map[string]*encryptedSecret
	metadata map[string]*models.Secret
	keyring  map[string][]byte
	mu       sync.RWMutex
}

// encryptedSecret represents an encrypted secret value
type encryptedSecret struct {
	OwnerID   string
	Name      string
	Value     []byte
	Nonce     []byte
	CreatedAt time.Time
}

// NewService creates a new secrets service
func NewService(cfg *config.TEEConfig) *Service {
	return &Service{
		cfg:      cfg,
		secrets:  make(map[string]*encryptedSecret),
		metadata: make(map[string]*models.Secret),
		keyring:  make(map[string][]byte),
	}
}

// Start initializes the secrets service
func (s *Service) Start(ctx context.Context) error {
	log.Println("Starting Secrets service...")

	// Initialize encryption keys
	if err := s.initializeKeyring(); err != nil {
		return err
	}

	// In a real implementation, we would load secrets from persistent storage here

	return nil
}

// Stop shuts down the secrets service
func (s *Service) Stop() error {
	log.Println("Stopping Secrets service...")

	// In a real implementation, we would persist any pending changes

	return nil
}

// Name returns the service name
func (s *Service) Name() string {
	return "Secrets"
}

// initializeKeyring sets up encryption keys
func (s *Service) initializeKeyring() error {
	// In a real implementation, this would initialize encryption keys
	// potentially from a secure hardware module or from a secure storage
	// protected by the TEE

	// Generate a master key for demonstration
	masterKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, masterKey); err != nil {
		return err
	}

	s.keyring["master"] = masterKey
	log.Println("Encryption keyring initialized")

	return nil
}

// StoreSecret securely stores a secret
func (s *Service) StoreSecret(ownerID string, name string, value string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we have encryption keys
	masterKey, exists := s.keyring["master"]
	if !exists {
		return "", errors.New("encryption keyring not initialized")
	}

	// Create AES cipher
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the value
	encryptedValue := aesGCM.Seal(nil, nonce, []byte(value), nil)

	// Generate ID for the secret
	secretID := fmt.Sprintf("secret-%d", time.Now().UnixNano())

	// Store the encrypted secret
	secret := &encryptedSecret{
		OwnerID:   ownerID,
		Name:      name,
		Value:     encryptedValue,
		Nonce:     nonce,
		CreatedAt: time.Now(),
	}
	s.secrets[secretID] = secret

	// Store metadata separately
	metadata := &models.Secret{
		ID:        secretID,
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: secret.CreatedAt,
	}
	s.metadata[secretID] = metadata

	log.Printf("Secret stored: %s (owner: %s)", secretID, ownerID)

	return secretID, nil
}

// RetrieveSecret decrypts and returns a secret value
func (s *Service) RetrieveSecret(secretID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secret, exists := s.secrets[secretID]
	if !exists {
		return "", fmt.Errorf("secret with ID %s does not exist", secretID)
	}

	// Get the master key
	masterKey, exists := s.keyring["master"]
	if !exists {
		return "", errors.New("encryption keyring not initialized")
	}

	// Create AES cipher
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decrypt
	plaintext, err := aesGCM.Open(nil, secret.Nonce, secret.Value, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// DeleteSecret removes a secret
func (s *Service) DeleteSecret(secretID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.secrets[secretID]; !exists {
		return fmt.Errorf("secret with ID %s does not exist", secretID)
	}

	delete(s.secrets, secretID)
	delete(s.metadata, secretID)

	log.Printf("Secret deleted: %s", secretID)

	return nil
}

// ListSecrets returns metadata for all secrets owned by a user
func (s *Service) ListSecrets(ownerID string) []*models.Secret {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secrets := make([]*models.Secret, 0)

	for _, metadata := range s.metadata {
		if metadata.OwnerID == ownerID {
			secrets = append(secrets, metadata)
		}
	}

	return secrets
}

// GetSecretMetadata returns metadata for a single secret
func (s *Service) GetSecretMetadata(secretID string) (*models.Secret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metadata, exists := s.metadata[secretID]
	if !exists {
		return nil, fmt.Errorf("secret with ID %s does not exist", secretID)
	}

	return metadata, nil
}
