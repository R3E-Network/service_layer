package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Port           int    `json:"port"`
	Host           string `json:"host"`
	TLSCertPath    string `json:"tlsCertPath"`
	TLSKeyPath     string `json:"tlsKeyPath"`
	EnableTLS      bool   `json:"enableTls"`
	ReadTimeoutSec int    `json:"readTimeoutSec"`
	ReadTimeout    time.Duration   `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration   `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration   `mapstructure:"idle_timeout"`
	RateLimit      RateLimitConfig `mapstructure:"rate_limit"`
	CORS           CORSConfig      `mapstructure:"cors"`
	Mode           string          `mapstructure:"mode"`
	Timeout        int             `mapstructure:"timeout"`
}

// BlockchainConfig contains Neo N3 blockchain settings
type BlockchainConfig struct {
	Network          string   `json:"network" mapstructure:"network"`
	RPCEndpoint      string   `json:"rpcEndpoint" mapstructure:"rpc_endpoint"`
	WSEndpoint       string   `json:"wsEndpoint" mapstructure:"ws_endpoint"`
	RPCEndpoints     []string `json:"rpcEndpoints" mapstructure:"rpc_endpoints"`
	NetworkMagic     uint32   `json:"networkMagic" mapstructure:"network_magic"`
	WalletPath       string   `json:"walletPath" mapstructure:"wallet_path"`
	WalletPassword   string   `json:"walletPassword" mapstructure:"wallet_password"`
	AccountAddress   string   `json:"accountAddress" mapstructure:"account_address"`
	GasBankContract  string   `json:"gasBankContract" mapstructure:"gas_bank_contract"`
	OracleContract   string   `json:"oracleContract" mapstructure:"oracle_contract"`
	PriceFeedTimeout int      `json:"priceFeedTimeout" mapstructure:"price_feed_timeout"`
}

// TEEConfig defines the configuration for Trusted Execution Environment
type TEEConfig struct {
	Provider          string      `mapstructure:"provider"`
	EnableAttestation bool        `mapstructure:"enable_attestation"`
	Azure             AzureConfig `mapstructure:"azure"`
}

// AzureConfig defines the configuration for Azure Confidential Computing
type AzureConfig struct {
	ClientID                string       `mapstructure:"client_id"`
	ClientSecret            string       `mapstructure:"client_secret"`
	TenantID                string       `mapstructure:"tenant_id"`
	SubscriptionID          string       `mapstructure:"subscription_id"`
	ResourceGroup           string       `mapstructure:"resource_group"`
	AttestationProviderName string       `mapstructure:"attestation_provider_name"`
	AttestationEndpoint     string       `mapstructure:"attestation_endpoint"`
	AllowedSGXEnclaves      []string     `mapstructure:"allowed_sgx_enclaves"`
	Region                  string       `mapstructure:"region"`
	Runtime                 RuntimeConfig `mapstructure:"runtime"`
}

// RuntimeConfig contains runtime-specific configuration for TEE
type RuntimeConfig struct {
	JSMemoryLimit     int   `mapstructure:"js_memory_limit"`
	ExecutionTimeout  int   `mapstructure:"execution_timeout"`
}

// GasBankConfig contains gas management settings
type GasBankConfig struct {
	MinimumGasBalance float64 `json:"minimumGasBalance"`
	AutoRefill        bool    `json:"autoRefill"`
	RefillAmount      float64 `json:"refillAmount"`
}

// PriceFeedConfig contains price feed service settings
type PriceFeedConfig struct {
	UpdateIntervalSec int      `json:"updateIntervalSec"`
	DataSources       []string `json:"dataSources"`
	SupportedTokens   []string `json:"supportedTokens"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	EnableFileLogging     bool   `json:"enableFileLogging"`
	LogFilePath           string `json:"logFilePath"`
	EnableDebugLogs       bool   `json:"enableDebugLogs"`
	RotationIntervalHours int    `json:"rotationIntervalHours"`
	MaxLogFiles           int    `json:"maxLogFiles"`
	Level                 string `mapstructure:"level"`
	Format                string `mapstructure:"format"`
	Output                string `mapstructure:"output"`
	FilePrefix            string `mapstructure:"file_prefix"`
}

// MetricsConfig contains monitoring settings
type MetricsConfig struct {
	Enabled       bool   `json:"enabled"`
	ListenAddress string `mapstructure:"listenAddress"`
}

// Config represents the application configuration
type Config struct {
	Environment string           `mapstructure:"environment"`
	Server      ServerConfig     `mapstructure:"server"`
	Database    DatabaseConfig   `mapstructure:"database"`
	Blockchain  BlockchainConfig `mapstructure:"blockchain"`
	Functions   FunctionsConfig  `mapstructure:"functions"`
	Secrets     SecretsConfig    `mapstructure:"secrets"`
	Oracle      OracleConfig     `mapstructure:"oracle"`
	PriceFeed   PriceFeedConfig  `mapstructure:"price_feed"`
	Automation  AutomationConfig `mapstructure:"automation"`
	GasBank     GasBankConfig    `mapstructure:"gas_bank"`
	TEE         TEEConfig        `mapstructure:"tee"`
	Monitoring  MonitoringConfig `mapstructure:"monitoring"`
	Auth        AuthConfig       `mapstructure:"auth"`
	Services    ServicesConfig   `mapstructure:"services"`
	Neo         NeoConfig        `mapstructure:"neo"`
	Features    FeaturesConfig   `mapstructure:"features"`
	Logging     LoggingConfig    `mapstructure:"logging"`
	Metrics     MetricsConfig    `mapstructure:"metrics"`
	Health      HealthConfig     `mapstructure:"health"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	DSN             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	SSLMode         string `mapstructure:"sslmode"`
}

// FunctionsConfig represents the functions configuration
type FunctionsConfig struct {
	MaxMemory        int `mapstructure:"max_memory"`
	ExecutionTimeout int `mapstructure:"execution_timeout"`
	MaxConcurrency   int `mapstructure:"max_concurrency"`
}

// SecretsConfig represents the secrets configuration
type SecretsConfig struct {
	KMSProvider string `mapstructure:"kms_provider"`
	KeyID       string `mapstructure:"key_id"`
	Region      string `mapstructure:"region"`
}

// OracleConfig represents the oracle service configuration
type OracleConfig struct {
	UpdateInterval int `mapstructure:"update_interval"`
	MaxDataSources int `mapstructure:"max_data_sources"`
}

// AutomationConfig represents the automation service configuration
type AutomationConfig struct {
	MaxTriggers int `mapstructure:"max_triggers"`
	MinInterval int `mapstructure:"min_interval"`
}

// MonitoringConfig represents the monitoring configuration
type MonitoringConfig struct {
	Enabled         bool             `mapstructure:"enabled"`
	PrometheusPort  int              `mapstructure:"prometheus_port"`
	MetricsEndpoint string           `mapstructure:"metrics_endpoint"`
	Prometheus      PrometheusConfig `mapstructure:"prometheus"`
	Logging         LoggingConfig    `mapstructure:"logging"`
}

// PrometheusConfig contains Prometheus-specific configuration
type PrometheusConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

// RateLimitConfig represents the rate limiting configuration
type RateLimitConfig struct {
	Enabled        bool  `mapstructure:"enabled"`
	RequestsPerIP  int   `mapstructure:"requests_per_ip"`
	RequestsPerKey int   `mapstructure:"requests_per_key"`
	BurstIP        int   `mapstructure:"burst_ip"`
	BurstKey       int   `mapstructure:"burst_key"`
	TimeWindowSec  int64 `mapstructure:"time_window_sec"`
}

// CORSConfig represents the CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	MaxAge         int      `mapstructure:"max_age"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	Secret             string `mapstructure:"secret"`
	JWTSecret          string `mapstructure:"jwt_secret"`
	AccessTokenTTL     int    `mapstructure:"access_token_ttl"`
	RefreshTokenTTL    int    `mapstructure:"refresh_token_ttl"`
	TokenExpiry        int    `mapstructure:"token_expiry"`
	RefreshTokenExpiry int    `mapstructure:"refresh_token_expiry"`
	EnableAPIKeys      bool   `mapstructure:"enable_api_keys"`
	APIKeyPrefix       string `mapstructure:"api_key_prefix"`
	APIKeyLength       int    `mapstructure:"api_key_length"`
	APIKeyTTL          int    `mapstructure:"api_key_ttl"`
}

// ServicesConfig contains configuration for external services
type ServicesConfig struct {
	TokenPriceAPI string       `mapstructure:"token_price_api"`
	GasAPI        string       `mapstructure:"gas_api"`
	Functions     FunctionsApi `mapstructure:"functions"`
	GasBank       GasBankApi   `mapstructure:"gas_bank"`
	Oracle        OracleApi    `mapstructure:"oracle"`
	PriceFeed     PriceFeedApi `mapstructure:"price_feed"`
	Secrets       SecretsApi   `mapstructure:"secrets"`
}

// FunctionsApi contains API configuration for functions service
type FunctionsApi struct {
	Endpoint          string `mapstructure:"endpoint"`
	Timeout           int    `mapstructure:"timeout"`
	MaxSourceCodeSize int    `mapstructure:"max_source_code_size"`
}

// GasBankApi contains API configuration for gas bank service
type GasBankApi struct {
	Endpoint      string  `mapstructure:"endpoint"`
	Timeout       int     `mapstructure:"timeout"`
	MinDeposit    float64 `mapstructure:"min_deposit"`
	MaxWithdrawal float64 `mapstructure:"max_withdrawal"`
	GasReserve    string  `mapstructure:"gas_reserve"`
}

// OracleApi contains API configuration for oracle service
type OracleApi struct {
	Endpoint       string `mapstructure:"endpoint"`
	Timeout        int    `mapstructure:"timeout"`
	RequestTimeout int    `mapstructure:"request_timeout"`
	NumWorkers     int    `mapstructure:"num_workers"`
	SigningKey     string `mapstructure:"signing_key"`
}

// PriceFeedApi contains API configuration for price feed service
type PriceFeedApi struct {
	Endpoint                  string  `mapstructure:"endpoint"`
	Timeout                   int     `mapstructure:"timeout"`
	NumWorkers                int     `mapstructure:"num_workers"`
	DefaultUpdateInterval     string  `mapstructure:"default_update_interval"`
	DefaultDeviationThreshold float64 `mapstructure:"default_deviation_threshold"`
	DefaultHeartbeatInterval  string  `mapstructure:"default_heartbeat_interval"`
	CoinMarketCapAPIKey       string  `mapstructure:"coin_market_cap_api_key"`
}

// SecretsApi contains API configuration for secrets service
type SecretsApi struct {
	Endpoint          string `mapstructure:"endpoint"`
	Timeout           int    `mapstructure:"timeout"`
	MaxSecretsPerUser int    `mapstructure:"max_secrets_per_user"`
	MaxSecretSize     int    `mapstructure:"max_secret_size"`
}

// NeoConfig contains Neo N3 blockchain configuration
type NeoConfig struct {
	NetworkID        int     `mapstructure:"network_id"`
	ChainID          int     `mapstructure:"chain_id"`
	Network          string  `mapstructure:"network"`
	RPCEndpoint      string  `mapstructure:"rpc_endpoint"`
	WSEndpoint       string  `mapstructure:"ws_endpoint"`
	GasBankContract  string  `mapstructure:"gas_bank_contract"`
	OracleContract   string  `mapstructure:"oracle_contract"`
	PriceFeedTimeout int     `mapstructure:"price_feed_timeout"`
	Confirmations    int64   `mapstructure:"confirmations"`
	GasLimit         int64   `mapstructure:"gas_limit"`
	GasPrice         int64   `mapstructure:"gas_price"`
}

// NodeConfig contains configuration for a blockchain node
type NodeConfig struct {
	URL      string `mapstructure:"url"`
	Priority int    `mapstructure:"priority"`
}

// StringURL returns the URL as a string for compatibility
func (n NodeConfig) StringURL() string {
	return n.URL
}

// NamedNodeConfig allows initialization from string
type NamedNodeConfig string

// StringToNodeConfig converts a string URL to a NodeConfig
func StringToNodeConfig(url string) NodeConfig {
	return NodeConfig{
		URL:      url,
		Priority: 0,
	}
}

// StringsToNodeConfigs converts a slice of string URLs to NodeConfigs
func StringsToNodeConfigs(urls []string) []NodeConfig {
	nodes := make([]NodeConfig, len(urls))
	for i, url := range urls {
		nodes[i] = StringToNodeConfig(url)
	}
	return nodes
}

// ToBlockchainNodeConfig converts config.NodeConfig to blockchain.NodeConfig
func (n NodeConfig) ToBlockchainNodeConfig() interface{} {
	return struct {
		URL    string  `json:"url"`
		Weight float64 `json:"weight"`
	}{
		URL:    n.URL,
		Weight: float64(n.Priority),
	}
}

// SetupDefaultValues initializes configuration with default values when loading
func (c *NeoConfig) SetupDefaultValues() {
}

// FeaturesConfig contains feature flag configuration
type FeaturesConfig struct {
	EnableGasBank         bool `mapstructure:"enable_gas_bank"`
	GasBank               bool `mapstructure:"enable_gas_bank"` // Alias for EnableGasBank
	EnableOracle          bool `mapstructure:"enable_oracle"`
	Oracle                bool `mapstructure:"enable_oracle"` // Alias for EnableOracle
	EnablePriceFeed       bool `mapstructure:"enable_price_feed"`
	PriceFeed             bool `mapstructure:"enable_price_feed"` // Alias for EnablePriceFeed
	EnableSecrets         bool `mapstructure:"enable_secrets"`
	Secrets               bool `mapstructure:"enable_secrets"` // Alias for EnableSecrets
	EnableFunctions       bool `mapstructure:"enable_functions"`
	Functions             bool `mapstructure:"enable_functions"` // Alias for EnableFunctions
	EnableEvents          bool `mapstructure:"enable_events"`
	Events                bool `mapstructure:"enable_events"` // Alias for EnableEvents
	EnableTEE             bool `mapstructure:"enable_tee"`
	TEE                   bool `mapstructure:"enable_tee"` // Alias for EnableTEE
	EnableAutomation      bool `mapstructure:"enable_automation"`
	Automation            bool `mapstructure:"enable_automation"` // Alias for EnableAutomation
	EnableRandomGenerator bool `mapstructure:"enable_random_generator"`
	RandomGenerator       bool `mapstructure:"enable_random_generator"` // Alias for EnableRandomGenerator
}

// HealthConfig represents health monitoring configuration
type HealthConfig struct {
	CheckIntervalSec int `mapstructure:"check_interval_sec"`
	MaxRetries       int `mapstructure:"max_retries"`
	RetryDelaySec    int `mapstructure:"retry_delay_sec"`
}

// New creates a new config instance with default values
func New() *Config {
	return &Config{
		Environment: "development",
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  120 * time.Second,
			RateLimit: RateLimitConfig{
				Enabled:        true,
				RequestsPerIP:  100,  // 100 requests per minute for IP-based limiting
				RequestsPerKey: 1000, // 1000 requests per minute for API key-based limiting
				BurstIP:        20,   // Allow bursts of up to 20 requests
				BurstKey:       100,  // Allow bursts of up to 100 requests for API keys
				TimeWindowSec:  60,   // 1 minute window
			},
			CORS: CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders: []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
				MaxAge:         86400,
			},
			Mode:    "development",
			Timeout: 60,
		},
		Database: DatabaseConfig{
			Driver:          "postgres",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300, // 5 minutes
		},
		Blockchain: BlockchainConfig{
			Network:          "",
			RPCEndpoint:      "",
			WSEndpoint:       "",
			RPCEndpoints:     []string{"http://localhost:10332"},
			NetworkMagic:     860833102, // Neo N3 TestNet
			WalletPath:       "./wallet.json",
			WalletPassword:   "",
			AccountAddress:   "",
			GasBankContract:  "0x1234567890123456789012345678901234567890",
			OracleContract:   "0x0987654321098765432109876543210987654321",
			PriceFeedTimeout: 60,
		},
		Functions: FunctionsConfig{
			MaxMemory:        128, // MB
			ExecutionTimeout: 30,  // seconds
			MaxConcurrency:   10,
		},
		Oracle: OracleConfig{
			UpdateInterval: 60, // seconds
			MaxDataSources: 100,
		},
		PriceFeed: PriceFeedConfig{
			UpdateIntervalSec: 60, // seconds
			DataSources:       []string{"coinmarketcap", "coingecko"},
			SupportedTokens:   []string{"NEO", "GAS", "ETH", "BTC"},
		},
		Automation: AutomationConfig{
			MaxTriggers: 100,
			MinInterval: 5, // seconds
		},
		GasBank: GasBankConfig{
			MinimumGasBalance: 10.0,
			AutoRefill:        true,
			RefillAmount:      50.0,
		},
		TEE: TEEConfig{
			Provider:            "simulation",
			EnableAttestation:   false,
			Azure: AzureConfig{
				ClientID:                "",
				ClientSecret:            "",
				TenantID:                "",
				SubscriptionID:          "",
				ResourceGroup:           "",
				AttestationProviderName: "",
				AttestationEndpoint:     "",
				AllowedSGXEnclaves:      []string{},
				Region:                  "",
				Runtime: RuntimeConfig{
					JSMemoryLimit:     0,
					ExecutionTimeout:  0,
				},
			},
		},
		Monitoring: MonitoringConfig{
			Enabled:        true,
			PrometheusPort: 9090,
		},
		Auth: AuthConfig{
			Secret:             "",
			JWTSecret:          "default-jwt-secret-key",
			AccessTokenTTL:     3600,
			RefreshTokenTTL:    86400,
			TokenExpiry:        3600,
			RefreshTokenExpiry: 86400,
			EnableAPIKeys:      false,
			APIKeyPrefix:       "",
			APIKeyLength:       16,
			APIKeyTTL:          3600,
		},
		Services: ServicesConfig{
			TokenPriceAPI: "",
			GasAPI:        "",
			Functions: FunctionsApi{
				Endpoint:          "http://localhost:8081",
				Timeout:           10,
				MaxSourceCodeSize: 1024,
			},
			GasBank: GasBankApi{
				Endpoint:      "http://localhost:8082",
				Timeout:       10,
				MinDeposit:    100,
				MaxWithdrawal: 10000,
				GasReserve:    "100",
			},
			Oracle: OracleApi{
				Endpoint:       "http://localhost:8083",
				Timeout:        10,
				RequestTimeout: 5,
				NumWorkers:     5,
				SigningKey:     "default-signing-key",
			},
			PriceFeed: PriceFeedApi{
				Endpoint:                  "http://localhost:8084",
				Timeout:                   10,
				NumWorkers:                5,
				DefaultUpdateInterval:     "60",
				DefaultDeviationThreshold: 10,
				DefaultHeartbeatInterval:  "30",
				CoinMarketCapAPIKey:       "default-coin-market-cap-api-key",
			},
			Secrets: SecretsApi{
				Endpoint:          "http://localhost:8085",
				Timeout:           10,
				MaxSecretsPerUser: 100,
				MaxSecretSize:     1024,
			},
		},
		Neo: NeoConfig{
			NetworkID:        1, // MainNet
			ChainID:          1, // MainNet
			Network:          "",
			RPCEndpoint:      "http://localhost:10332",
			WSEndpoint:       "ws://localhost:10332/ws", 
			GasBankContract:  "0x1234567890123456789012345678901234567890",
			OracleContract:   "0x0987654321098765432109876543210987654321",
			PriceFeedTimeout: 60,
			Confirmations:    1,
			GasLimit:         0,
			GasPrice:         0,
		},
		Features: FeaturesConfig{
			EnableGasBank:         false,
			EnableOracle:          false,
			EnablePriceFeed:       false,
			EnableSecrets:         false,
			EnableFunctions:       false,
			EnableEvents:          false,
			EnableTEE:             false,
			EnableAutomation:      false,
			EnableRandomGenerator: false,
		},
		Logging: LoggingConfig{
			EnableFileLogging:     true,
			LogFilePath:           "./logs/neo-oracle.log",
			EnableDebugLogs:       false,
			RotationIntervalHours: 24,
			MaxLogFiles:           7,
			Level:                "info",
			Format:               "json",
			Output:               "file",
			FilePrefix:           "service-layer",
		},
		Metrics: MetricsConfig{
			Enabled:       true,
			ListenAddress: ":9090",
		},
		Health: HealthConfig{
			CheckIntervalSec: 10,
			MaxRetries:       3,
			RetryDelaySec:    5,
		},
	}
}

// ConnectionString returns a PostgreSQL connection string
func (c DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// RedisAddress returns the Redis server address
func (c RedisConfig) RedisAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Default config file path
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	// Create config instance with default values
	cfg := &Config{}

	// Try to load from config file
	if err := loadFromFile(configPath, cfg); err != nil {
		fmt.Printf("Warning: Could not load config file: %v\n", err)
	}

	// Override with environment variables
	if err := envdecode.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode environment variables: %w", err)
	}

	return cfg, nil
}

// loadFromFile loads configuration from YAML file
func loadFromFile(filePath string, cfg *Config) error {
	// Expand file path
	expandedPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// Read config file
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return err
	}

	// Parse YAML
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}

// IsDevelopment returns true if the server is in development mode
func (c ServerConfig) IsDevelopment() bool {
	return strings.ToLower(c.Mode) == "development"
}

// IsProduction returns true if the server is in production mode
func (c ServerConfig) IsProduction() bool {
	return strings.ToLower(c.Mode) == "production"
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// DefaultConfig creates a default configuration
func DefaultConfig() *Config {
	return &Config{
		Environment: "development",
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8000,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
			RateLimit: RateLimitConfig{
				Enabled:        true,
				RequestsPerIP:  100,
				RequestsPerKey: 1000,
				BurstIP:        5,
				BurstKey:       50,
				TimeWindowSec:  60,
			},
			CORS: CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders: []string{"Content-Type", "Authorization"},
				MaxAge:         86400,
			},
			Mode:    "development",
			Timeout: 30,
		},
		Blockchain: BlockchainConfig{
			Network:          "",
			RPCEndpoint:      "",
			WSEndpoint:       "",
			RPCEndpoints:     []string{"http://localhost:10332"},
			NetworkMagic:     860833102, // Neo N3 TestNet
			WalletPath:       "./wallet.json",
			WalletPassword:   "",
			AccountAddress:   "",
			GasBankContract:  "0x1234567890123456789012345678901234567890",
			OracleContract:   "0x0987654321098765432109876543210987654321",
			PriceFeedTimeout: 60,
		},
		TEE: TEEConfig{
			Provider:            "simulation",
			EnableAttestation:   false,
			Azure: AzureConfig{
				ClientID:                "",
				ClientSecret:            "",
				TenantID:                "",
				SubscriptionID:          "",
				ResourceGroup:           "",
				AttestationProviderName: "",
				AttestationEndpoint:     "",
				AllowedSGXEnclaves:      []string{},
				Region:                  "",
				Runtime: RuntimeConfig{
					JSMemoryLimit:     0,
					ExecutionTimeout:  0,
				},
			},
		},
		GasBank: GasBankConfig{
			MinimumGasBalance: 10.0,
			AutoRefill:        true,
			RefillAmount:      50.0,
		},
		PriceFeed: PriceFeedConfig{
			UpdateIntervalSec: 300, // 5 minutes
			DataSources:       []string{"coinmarketcap", "coingecko"},
			SupportedTokens:   []string{"NEO", "GAS", "ETH", "BTC"},
		},
		Logging: LoggingConfig{
			EnableFileLogging:     true,
			LogFilePath:           "./logs/neo-oracle.log",
			EnableDebugLogs:       false,
			RotationIntervalHours: 24,
			MaxLogFiles:           7,
			Level:                "info",
			Format:               "json",
			Output:               "file",
			FilePrefix:           "service-layer",
		},
		Metrics: MetricsConfig{
			Enabled:       true,
			ListenAddress: ":9090",
		},
		Health: HealthConfig{
			CheckIntervalSec: 10,
			MaxRetries:       3,
			RetryDelaySec:    5,
		},
	}
}

// SaveConfig saves the configuration to a file
func SaveConfig(cfg *Config, configPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}
