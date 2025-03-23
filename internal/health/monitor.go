package health

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/your-org/neo-oracle/internal/common"
	"github.com/your-org/neo-oracle/internal/config"
)

// ServiceStatus represents the health status of a service
type ServiceStatus struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	LastChecked time.Time `json:"lastChecked"`
	Message     string    `json:"message,omitempty"`
	Uptime      int64     `json:"uptime"` // seconds
}

// Monitor tracks service health and performs auto-recovery
type Monitor struct {
	services       map[string]common.Service
	statuses       map[string]*ServiceStatus
	checkInterval  time.Duration
	startTimes     map[string]time.Time
	recoveryPolicy map[string]RecoveryPolicy
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
}

// RecoveryPolicy defines how to handle service failures
type RecoveryPolicy struct {
	MaxRetries     int
	RetryDelay     time.Duration
	CurrentRetries int
}

// NewMonitor creates a new health monitor
func NewMonitor(cfg *config.HealthConfig) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())

	return &Monitor{
		services:       make(map[string]common.Service),
		statuses:       make(map[string]*ServiceStatus),
		startTimes:     make(map[string]time.Time),
		recoveryPolicy: make(map[string]RecoveryPolicy),
		checkInterval:  time.Duration(cfg.CheckIntervalSec) * time.Second,
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start initiates the health monitor
func (m *Monitor) Start(ctx context.Context) error {
	log.Println("Starting Health Monitor service...")

	// Start monitoring loop
	go m.monitorLoop()

	return nil
}

// Stop shuts down the health monitor
func (m *Monitor) Stop() error {
	log.Println("Stopping Health Monitor service...")
	m.cancel()
	return nil
}

// Name returns the service name
func (m *Monitor) Name() string {
	return "HealthMonitor"
}

// RegisterService adds a service to be monitored
func (m *Monitor) RegisterService(name string, service common.Service) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.services[name] = service
	m.statuses[name] = &ServiceStatus{
		Name:        name,
		Status:      "starting",
		LastChecked: time.Now(),
	}
	m.startTimes[name] = time.Now()

	// Set default recovery policy
	m.recoveryPolicy[name] = RecoveryPolicy{
		MaxRetries: 3,
		RetryDelay: 5 * time.Second,
	}

	log.Printf("Service %s registered for health monitoring", name)
}

// SetRecoveryPolicy configures recovery behavior for a service
func (m *Monitor) SetRecoveryPolicy(serviceName string, maxRetries int, retryDelay time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	policy := RecoveryPolicy{
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
	}

	m.recoveryPolicy[serviceName] = policy
}

// monitorLoop periodically checks service health
func (m *Monitor) monitorLoop() {
	ticker := time.NewTicker(m.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkServicesHealth()
		case <-m.ctx.Done():
			return
		}
	}
}

// checkServicesHealth verifies the health of all registered services
func (m *Monitor) checkServicesHealth() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, service := range m.services {
		status := m.statuses[name]
		startTime := m.startTimes[name]

		// Update uptime
		status.Uptime = int64(time.Since(startTime).Seconds())
		status.LastChecked = time.Now()

		// Check service health
		// In a real implementation, we would call a health check method
		// For now, we'll assume it's healthy if it exists
		if service != nil {
			status.Status = "healthy"
			status.Message = ""

			// Reset retry counter on success
			policy := m.recoveryPolicy[name]
			policy.CurrentRetries = 0
			m.recoveryPolicy[name] = policy
		} else {
			status.Status = "unhealthy"
			status.Message = "Service unavailable"

			// Attempt recovery if needed
			m.attemptServiceRecovery(name)
		}
	}
}

// attemptServiceRecovery tries to restart an unhealthy service
func (m *Monitor) attemptServiceRecovery(serviceName string) {
	policy := m.recoveryPolicy[serviceName]
	service := m.services[serviceName]

	// Check if we've exceeded max retries
	if policy.CurrentRetries >= policy.MaxRetries {
		log.Printf("Service %s recovery failed after %d attempts",
			serviceName, policy.CurrentRetries)
		return
	}

	// Increment retry counter
	policy.CurrentRetries++
	m.recoveryPolicy[serviceName] = policy

	log.Printf("Attempting to recover service %s (attempt %d/%d)",
		serviceName, policy.CurrentRetries, policy.MaxRetries)

	// Try to restart the service
	go func() {
		// Add delay before retry
		time.Sleep(policy.RetryDelay)

		// Stop the service first
		if err := service.Stop(); err != nil {
			log.Printf("Error stopping service %s: %v", serviceName, err)
		}

		// Start the service again
		if err := service.Start(context.Background()); err != nil {
			log.Printf("Error restarting service %s: %v", serviceName, err)
			return
		}

		log.Printf("Service %s successfully restarted", serviceName)

		// Update status
		m.mu.Lock()
		defer m.mu.Unlock()

		if status, exists := m.statuses[serviceName]; exists {
			status.Status = "recovered"
			status.Message = fmt.Sprintf("Recovered after %d attempts",
				policy.CurrentRetries)
			m.startTimes[serviceName] = time.Now() // Reset uptime
		}
	}()
}

// GetServiceStatus returns the current health status for a service
func (m *Monitor) GetServiceStatus(serviceName string) (*ServiceStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status, exists := m.statuses[serviceName]
	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	return status, nil
}

// GetAllServiceStatuses returns health information for all services
func (m *Monitor) GetAllServiceStatuses() []*ServiceStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]*ServiceStatus, 0, len(m.statuses))
	for _, status := range m.statuses {
		statuses = append(statuses, status)
	}

	return statuses
}
