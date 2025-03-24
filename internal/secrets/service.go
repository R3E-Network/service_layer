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
	"strconv"
	"sync"
	"time"

	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/internal/models"
)

// Service manages secure storage of user secrets
type Service struct {
	cfg      *config.TEEConfig
	secrets  map[string]*encryptedSecret
	metadata map[int]*models.Secret
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
		metadata: make(map[int]*models.Secret),
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
func (s *Service) StoreSecret(secretID, ownerID, name, value string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Encrypt the secret value
	encryptedValue, err := s.encrypt(value)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt secret: %v", err)
	}

	// Create a new encrypted secret
	secret := &encryptedSecret{
		Value:     encryptedValue,
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	s.secrets[secretID] = secret

	// Parse secretID as integer or generate a new one
	secretIDInt := 0
	if secretID != "" {
		var err error
		secretIDInt, err = strconv.Atoi(secretID)
		if err != nil {
			return "", fmt.Errorf("invalid secret ID format: %v", err)
		}
	} else {
		// Generate a new ID - using timestamp for simplicity
		secretIDInt = int(time.Now().Unix())
		secretID = strconv.Itoa(secretIDInt)
	}

	// Store metadata separately
	metadata := &models.Secret{
		ID:        secretIDInt,
		UserID:    0, // Using 0 as a placeholder for the UserID since we're using ownerID as a string
		Name:      name,
		Value:     "", // We don't store the actual value in metadata
		Version:   1,  // Initial version
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.CreatedAt,
	}
	s.metadata[secretIDInt] = metadata

	log.Printf("Secret stored: %s (owner: %s)", secretID, ownerID)

	return secretID, nil
}

// RetrieveSecret decrypts and returns a secret value
func (s *Service) RetrieveSecret(secretID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// We don't need to convert secretID to integer for retrieving from s.secrets
	secret, exists := s.secrets[secretID]
	if !exists {
		return "", fmt.Errorf("secret with ID %s does not exist", secretID)
	}

	// Decrypt the secret value
	decryptedValue, err := s.decrypt(secret.Value)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %v", err)
	}

	return decryptedValue, nil
}

// DeleteSecret removes a secret
func (s *Service) DeleteSecret(secretID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	secretIDInt, err := strconv.Atoi(secretID)
	if err != nil {
		return fmt.Errorf("invalid secret ID format: %v", err)
	}

	if _, exists := s.secrets[secretID]; !exists {
		return fmt.Errorf("secret with ID %s does not exist", secretID)
	}

	delete(s.secrets, secretID)
	delete(s.metadata, secretIDInt)

	log.Printf("Secret deleted: %s", secretID)

	return nil
}

// ListSecrets returns metadata for all secrets owned by a user
func (s *Service) ListSecrets(ownerID string) []*models.Secret {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secrets := make([]*models.Secret, 0)

	for secretID, metadata := range s.metadata {
		secret, exists := s.secrets[strconv.Itoa(secretID)]
		if exists && secret.OwnerID == ownerID {
			secrets = append(secrets, metadata)
		}
	}

	return secrets
}

// GetSecretMetadata returns metadata for a single secret
func (s *Service) GetSecretMetadata(secretID string) (*models.Secret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secretIDInt, err := strconv.Atoi(secretID)
	if err != nil {
		return nil, fmt.Errorf("invalid secret ID format: %v", err)
	}

	metadata, exists := s.metadata[secretIDInt]
	if !exists {
		return nil, fmt.Errorf("secret with ID %s does not exist", secretID)
	}

	return metadata, nil
}

// encrypt encrypts a secret value
func (s *Service) encrypt(value string) ([]byte, error) {
	// Get the master key
	masterKey, exists := s.keyring["master"]
	if !exists {
		return nil, errors.New("encryption keyring not initialized")
	}

	// Create AES cipher
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the value
	encryptedValue := aesGCM.Seal(nil, nonce, []byte(value), nil)

	return encryptedValue, nil
}

// decrypt decrypts a secret value
func (s *Service) decrypt(encryptedValue []byte) (string, error) {
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
	plaintext, err := aesGCM.Open(nil, encryptedValue[:aesGCM.NonceSize()], encryptedValue[aesGCM.NonceSize():], nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
