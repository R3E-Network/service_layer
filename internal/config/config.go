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

// Config contains all configuration for the service layer
type Config struct {
	Server     ServerConfig     `json:"server"`
	Blockchain BlockchainConfig `json:"blockchain"`
	TEE        TEEConfig        `json:"tee"`
	GasBank    GasBankConfig    `json:"gasbank"`
	PriceFeed  PriceFeedConfig  `json:"pricefeed"`
	Logging    LoggingConfig    `json:"logging"`
	Metrics    MetricsConfig    `json:"metrics"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Port           int    `json:"port"`
	Host           string `json:"host"`
	TLSCertPath    string `json:"tlsCertPath"`
	TLSKeyPath     string `json:"tlsKeyPath"`
	EnableTLS      bool   `json:"enableTls"`
	ReadTimeoutSec int    `json:"readTimeoutSec"`
}

// BlockchainConfig contains Neo N3 blockchain settings
type BlockchainConfig struct {
	RPCEndpoints     []string `json:"rpcEndpoints"`
	NetworkMagic     uint32   `json:"networkMagic"`
	WalletPath       string   `json:"walletPath"`
	WalletPassword   string   `json:"walletPassword"`
	AccountAddress   string   `json:"accountAddress"`
	GasBankContract  string   `json:"gasBankContract"`
	OracleContract   string   `json:"oracleContract"`
	PriceFeedTimeout int      `json:"priceFeedTimeout"`
}

// TEEConfig contains Trusted Execution Environment settings
type TEEConfig struct {
	Provider            string `json:"provider"` // azure, aws, etc.
	AzureAttestationURL string `json:"azureAttestationUrl"`
	EnclaveImageID      string `json:"enclaveImageId"`
	JSRuntimePath       string `json:"jsRuntimePath"`
	SecretsStoragePath  string `json:"secretsStoragePath"`
	MaxMemoryMB         int    `json:"maxMemoryMb"`
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
}

// MetricsConfig contains monitoring settings
type MetricsConfig struct {
	Enabled       bool   `json:"enabled"`
	ListenAddress string `json:"listenAddress"`
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
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Host         string          `mapstructure:"host"`
	Port         int             `mapstructure:"port"`
	ReadTimeout  time.Duration   `mapstructure:"read_timeout"`
	WriteTimeout time.Duration   `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration   `mapstructure:"idle_timeout"`
	RateLimit    RateLimitConfig `mapstructure:"rate_limit"`
	CORS         CORSConfig      `mapstructure:"cors"`
	Mode         string          `mapstructure:"mode"`
	Timeout      int             `mapstructure:"timeout"`
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

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePrefix string `mapstructure:"file_prefix"`
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
	RpcUrl             string       `mapstructure:"rpc_url"`
	RPCURL             string       `mapstructure:"rpc_url"` // Alias for RpcUrl for backward compatibility
	NetworkFee         string       `mapstructure:"network_fee"`
	SystemFee          string       `mapstructure:"system_fee"`
	GasToken           string       `mapstructure:"gas_token"`
	ContractAddress    string       `mapstructure:"contract_address"`
	PrivateKey         string       `mapstructure:"private_key"`
	Network            string       `mapstructure:"network"`
	GasLimit           int          `mapstructure:"gas_limit"`
	GasPrice           int          `mapstructure:"gas_price"`
	Nodes              []NodeConfig `mapstructure:"nodes"`
	WalletPath         string       `mapstructure:"wallet_path"`
	URLs               []string     `mapstructure:"urls"`
	Confirmations      int          `mapstructure:"confirmations"`
	ConfirmationsInt64 int64        `mapstructure:"confirmations"` // Int64 version for compatibility
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
	// Convert Confirmations to ConfirmationsInt64
	c.ConfirmationsInt64 = int64(c.Confirmations)

	// Convert URLs to Nodes if needed
	if len(c.Nodes) == 0 && len(c.URLs) > 0 {
		c.Nodes = StringsToNodeConfigs(c.URLs)
	}
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
			Network:     "testnet",
			RPCEndpoint: "http://localhost:10332",
			WSEndpoint:  "ws://localhost:10334",
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
			UpdateInterval: 60, // seconds
			Sources: map[string]string{
				"default": "https://api.coingecko.com/api/v3",
			},
		},
		Automation: AutomationConfig{
			MaxTriggers: 100,
			MinInterval: 5, // seconds
		},
		GasBank: GasBankConfig{
			MaxWithdrawal: "10000",
			FeePercentage: 1,
		},
		TEE: TEEConfig{
			Enabled:  true,
			Provider: "azure",
			Azure: AzureConfig{
				ClientID:           "",
				ClientSecret:       "",
				TenantID:           "",
				SubscriptionID:     "",
				ResourceGroupName:  "",
				AttestationURL:     "",
				ConfidentialLedger: "",
				EnclaveSeal:        "",
				Runtime: RuntimeConfig{
					JSMemoryLimit:    128, // MB
					ExecutionTimeout: 30,  // seconds
					MaxConcurrency:   10,
					MaxCodeSize:      1024, // KB
					EnableNetworking: false,
					AllowedHosts:     []string{},
				},
			},
			Runtime: RuntimeConfig{
				JSMemoryLimit:    128, // MB
				ExecutionTimeout: 30,  // seconds
				MaxConcurrency:   10,
				MaxCodeSize:      1024, // KB
				EnableNetworking: false,
				AllowedHosts:     []string{},
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
			RpcUrl:          "",
			NetworkFee:      "",
			SystemFee:       "",
			GasToken:        "",
			ContractAddress: "",
			PrivateKey:      "",
			Network:         "",
			GasLimit:        0,
			GasPrice:        0,
			Nodes:           []NodeConfig{},
			WalletPath:      "",
			URLs:            []string{},
			Confirmations:   0,
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
		Server: ServerConfig{
			Port:           8080,
			Host:           "0.0.0.0",
			ReadTimeoutSec: 30,
		},
		Blockchain: BlockchainConfig{
			RPCEndpoints:     []string{"http://localhost:10332"},
			NetworkMagic:     860833102, // Neo N3 TestNet
			PriceFeedTimeout: 60,
		},
		TEE: TEEConfig{
			Provider:      "simulation",
			MaxMemoryMB:   512,
			JSRuntimePath: "./jsruntime",
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
		},
		Metrics: MetricsConfig{
			Enabled:       true,
			ListenAddress: ":9090",
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
