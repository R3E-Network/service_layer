package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/your-org/neo-oracle/internal/config"
)

// Service handles application metrics and monitoring
type Service struct {
	cfg      *config.MetricsConfig
	server   *http.Server
	registry *prometheus.Registry

	// Metrics
	functionCalls      *prometheus.CounterVec
	functionDuration   *prometheus.HistogramVec
	gasBankOperations  *prometheus.CounterVec
	blockchainRequests *prometheus.CounterVec
	apiRequests        *prometheus.CounterVec
	errorCount         *prometheus.CounterVec
}

// NewService creates a new metrics service
func NewService(cfg *config.MetricsConfig) *Service {
	registry := prometheus.NewRegistry()

	s := &Service{
		cfg:      cfg,
		registry: registry,
	}

	// Initialize metrics
	s.initializeMetrics()

	return s
}

// initializeMetrics sets up Prometheus metrics
func (s *Service) initializeMetrics() {
	// Function execution metrics
	s.functionCalls = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "neo_oracle_function_calls_total",
			Help: "Total number of JavaScript function calls",
		},
		[]string{"function_id", "status"},
	)
	s.registry.MustRegister(s.functionCalls)

	s.functionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "neo_oracle_function_duration_seconds",
			Help:    "Histogram of function execution times",
			Buckets: prometheus.ExponentialBuckets(0.01, 2, 10),
		},
		[]string{"function_id"},
	)
	s.registry.MustRegister(s.functionDuration)

	// GasBank metrics
	s.gasBankOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "neo_oracle_gasbank_operations_total",
			Help: "Total number of GasBank operations",
		},
		[]string{"operation", "status"},
	)
	s.registry.MustRegister(s.gasBankOperations)

	// Blockchain metrics
	s.blockchainRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "neo_oracle_blockchain_requests_total",
			Help: "Total number of blockchain requests",
		},
		[]string{"type", "status"},
	)
	s.registry.MustRegister(s.blockchainRequests)

	// API metrics
	s.apiRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "neo_oracle_api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"method", "path", "status_code"},
	)
	s.registry.MustRegister(s.apiRequests)

	// Error metrics
	s.errorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "neo_oracle_errors_total",
			Help: "Total number of errors",
		},
		[]string{"service", "type"},
	)
	s.registry.MustRegister(s.errorCount)
}

// Start initializes the metrics server
func (s *Service) Start(ctx context.Context) error {
	if !s.cfg.Enabled {
		log.Println("Metrics collection disabled")
		return nil
	}

	// Create HTTP server for metrics endpoint
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))

	addr := s.cfg.ListenAddress
	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting metrics server on %s", addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	return nil
}

// Stop shuts down the metrics server
func (s *Service) Stop() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Println("Stopping metrics server...")
		return s.server.Shutdown(ctx)
	}

	return nil
}

// Name returns the service name
func (s *Service) Name() string {
	return "Metrics"
}

// RecordFunctionCall records a function execution
func (s *Service) RecordFunctionCall(functionID string, status string, duration time.Duration) {
	s.functionCalls.WithLabelValues(functionID, status).Inc()
	s.functionDuration.WithLabelValues(functionID).Observe(duration.Seconds())
}

// RecordGasBankOperation records a GasBank operation
func (s *Service) RecordGasBankOperation(operation string, status string) {
	s.gasBankOperations.WithLabelValues(operation, status).Inc()
}

// RecordBlockchainRequest records a blockchain interaction
func (s *Service) RecordBlockchainRequest(requestType string, status string) {
	s.blockchainRequests.WithLabelValues(requestType, status).Inc()
}

// RecordAPIRequest records an API request
func (s *Service) RecordAPIRequest(method string, path string, statusCode int) {
	s.apiRequests.WithLabelValues(method, path, fmt.Sprintf("%d", statusCode)).Inc()
}

// RecordError increments the error counter
func (s *Service) RecordError(service string, errorType string) {
	s.errorCount.WithLabelValues(service, errorType).Inc()
}
