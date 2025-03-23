package triggers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/your-org/neo-oracle/internal/functions"
	"github.com/your-org/neo-oracle/internal/models"
	"github.com/your-org/neo-oracle/internal/pricefeed"
)

// Service manages automatic function execution triggers
type Service struct {
	executor    *functions.Executor
	priceFeed   *pricefeed.Service
	triggers    map[string]*models.Trigger
	cronRunner  *cron.Cron
	priceAlerts map[string]map[string][]priceAlert
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

type priceAlert struct {
	triggerID  string
	condition  string
	threshold  float64
	comparison string
}

// NewService creates a new trigger service
func NewService(executor *functions.Executor, priceFeed *pricefeed.Service) *Service {
	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		executor:    executor,
		priceFeed:   priceFeed,
		triggers:    make(map[string]*models.Trigger),
		cronRunner:  cron.New(),
		priceAlerts: make(map[string]map[string][]priceAlert),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start initializes and runs the trigger service
func (s *Service) Start(ctx context.Context) error {
	log.Println("Starting Trigger service...")

	// Start cron scheduler
	s.cronRunner.Start()

	// Start price alert monitor
	go s.monitorPriceAlerts()

	return nil
}

// Stop shuts down the trigger service
func (s *Service) Stop() error {
	log.Println("Stopping Trigger service...")

	// Stop cron scheduler
	cronCtx := s.cronRunner.Stop()
	<-cronCtx.Done()

	// Stop price monitor
	s.cancel()

	return nil
}

// Name returns the service name
func (s *Service) Name() string {
	return "Triggers"
}

// monitorPriceAlerts checks price conditions for alerts
func (s *Service) monitorPriceAlerts() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkPriceAlerts()
		case <-s.ctx.Done():
			return
		}
	}
}

// checkPriceAlerts evaluates price conditions
func (s *Service) checkPriceAlerts() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get current prices
	prices := s.priceFeed.GetAllPrices()

	// Check each token's alerts
	for token, price := range prices {
		// Skip if no alerts for this token
		tokenAlerts, exists := s.priceAlerts[token]
		if !exists {
			continue
		}

		// Check user alerts for this token
		for userID, alerts := range tokenAlerts {
			for _, alert := range alerts {
				triggered := false

				// Evaluate condition
				switch alert.comparison {
				case "above":
					triggered = price.Price > alert.threshold
				case "below":
					triggered = price.Price < alert.threshold
				}

				if triggered {
					trigger, exists := s.triggers[alert.triggerID]
					if exists {
						log.Printf("Price alert triggered: %s %s %.2f (current: %.2f)",
							token, alert.comparison, alert.threshold, price.Price)

						// Execute the function
						go s.executeTriggeredFunction(trigger)
					}
				}
			}
		}
	}
}

// executeTriggeredFunction runs a function based on a trigger
func (s *Service) executeTriggeredFunction(trigger *models.Trigger) {
	result, err := s.executor.ExecuteFunction(context.Background(), trigger.FunctionID, trigger.Parameters)
	if err != nil {
		log.Printf("Error executing triggered function %s: %v", trigger.FunctionID, err)
		return
	}

	log.Printf("Trigger %s executed function %s successfully", trigger.ID, trigger.FunctionID)
}

// CreateTrigger adds a new trigger
func (s *Service) CreateTrigger(trigger *models.Trigger) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.triggers[trigger.ID]; exists {
		return fmt.Errorf("trigger with ID %s already exists", trigger.ID)
	}

	// Set up the trigger based on type
	switch trigger.Type {
	case models.TriggerTypeSchedule:
		return s.setupScheduleTrigger(trigger)
	case models.TriggerTypePriceAlert:
		return s.setupPriceAlertTrigger(trigger)
	default:
		return fmt.Errorf("unsupported trigger type: %s", trigger.Type)
	}
}

// setupScheduleTrigger configures a time-based trigger
func (s *Service) setupScheduleTrigger(trigger *models.Trigger) error {
	if trigger.Schedule == "" {
		return fmt.Errorf("schedule cannot be empty for schedule trigger")
	}

	// Add to cron
	_, err := s.cronRunner.AddFunc(trigger.Schedule, func() {
		s.executeTriggeredFunction(trigger)
	})

	if err != nil {
		return fmt.Errorf("invalid cron schedule: %v", err)
	}

	s.triggers[trigger.ID] = trigger
	log.Printf("Schedule trigger created: %s (%s)", trigger.ID, trigger.Schedule)

	return nil
}

// setupPriceAlertTrigger configures a price-based trigger
func (s *Service) setupPriceAlertTrigger(trigger *models.Trigger) error {
	if trigger.Condition == "" {
		return fmt.Errorf("condition cannot be empty for price alert trigger")
	}

	// Parse condition (simplified for example)
	parts := strings.Split(trigger.Condition, " ")
	if len(parts) != 3 {
		return fmt.Errorf("invalid condition format: %s", trigger.Condition)
	}

	token := parts[0]
	comparison := parts[1]
	threshold, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return fmt.Errorf("invalid threshold in condition: %v", err)
	}

	// Create alert
	alert := priceAlert{
		triggerID:  trigger.ID,
		condition:  trigger.Condition,
		threshold:  threshold,
		comparison: comparison,
	}

	// Add to price alerts
	if _, exists := s.priceAlerts[token]; !exists {
		s.priceAlerts[token] = make(map[string][]priceAlert)
	}
	if _, exists := s.priceAlerts[token][trigger.UserID]; !exists {
		s.priceAlerts[token][trigger.UserID] = make([]priceAlert, 0)
	}

	s.priceAlerts[token][trigger.UserID] = append(s.priceAlerts[token][trigger.UserID], alert)
	s.triggers[trigger.ID] = trigger

	log.Printf("Price alert trigger created: %s (%s)", trigger.ID, trigger.Condition)

	return nil
}

// DeleteTrigger removes a trigger
func (s *Service) DeleteTrigger(triggerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trigger, exists := s.triggers[triggerID]
	if !exists {
		return fmt.Errorf("trigger with ID %s does not exist", triggerID)
	}

	delete(s.triggers, triggerID)

	// Handle specific trigger type cleanup
	switch trigger.Type {
	case models.TriggerTypePriceAlert:
		s.removePriceAlert(trigger)
	}

	log.Printf("Trigger deleted: %s", triggerID)

	return nil
}

// removePriceAlert cleans up price alert data
func (s *Service) removePriceAlert(trigger *models.Trigger) {
	// Parse condition to get token
	parts := strings.Split(trigger.Condition, " ")
	if len(parts) != 3 {
		return
	}

	token := parts[0]

	// Remove alert
	tokenAlerts, exists := s.priceAlerts[token]
	if !exists {
		return
	}

	userAlerts, exists := tokenAlerts[trigger.UserID]
	if !exists {
		return
	}

	// Filter out the alert
	newAlerts := make([]priceAlert, 0)
	for _, alert := range userAlerts {
		if alert.triggerID != trigger.ID {
			newAlerts = append(newAlerts, alert)
		}
	}

	if len(newAlerts) > 0 {
		s.priceAlerts[token][trigger.UserID] = newAlerts
	} else {
		delete(s.priceAlerts[token], trigger.UserID)
		if len(s.priceAlerts[token]) == 0 {
			delete(s.priceAlerts, token)
		}
	}
}

// GetTrigger retrieves a trigger by ID
func (s *Service) GetTrigger(triggerID string) (*models.Trigger, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	trigger, exists := s.triggers[triggerID]
	if !exists {
		return nil, fmt.Errorf("trigger with ID %s does not exist", triggerID)
	}

	return trigger, nil
}

// ListTriggers retrieves all triggers
func (s *Service) ListTriggers() []*models.Trigger {
	s.mu.RLock()
	defer s.mu.RUnlock()

	triggers := make([]*models.Trigger, 0, len(s.triggers))
	for _, trigger := range s.triggers {
		triggers = append(triggers, trigger)
	}

	return triggers
}
