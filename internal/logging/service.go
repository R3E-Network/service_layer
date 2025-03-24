package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/R3E-Network/service_layer/internal/config"
)

// Service manages application logging
type Service struct {
	cfg        *config.LoggingConfig
	logFile    *os.File
	stdLogger  *log.Logger
	fileLogger *log.Logger
	mu         sync.Mutex
}

// NewService creates a new logging service
func NewService(cfg *config.LoggingConfig) (*Service, error) {
	s := &Service{
		cfg: cfg,
	}

	// Initialize standard logger
	s.stdLogger = log.New(os.Stdout, "", log.LstdFlags)

	// Initialize file logger if enabled
	if cfg.EnableFileLogging {
		if err := s.setupFileLogger(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// setupFileLogger configures logging to a file
func (s *Service) setupFileLogger() error {
	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(s.cfg.LogFilePath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(s.cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	s.logFile = file

	// Create multi-writer for console and file
	multiWriter := io.MultiWriter(os.Stdout, file)
	s.fileLogger = log.New(multiWriter, "", log.LstdFlags)

	return nil
}

// Start initializes the logging service
func (s *Service) Start(ctx context.Context) error {
	log.SetOutput(os.Stdout)

	if s.cfg.EnableFileLogging {
		// Start log rotation if enabled
		if s.cfg.RotationIntervalHours > 0 {
			go s.rotateLogsPeriodically(ctx)
		}
	}

	log.Println("Logging service started")
	return nil
}

// Stop closes the logging service
func (s *Service) Stop() error {
	if s.logFile != nil {
		err := s.logFile.Close()
		s.logFile = nil
		return err
	}
	return nil
}

// Name returns the service name
func (s *Service) Name() string {
	return "Logging"
}

// rotateLogsPeriodically handles periodic log rotation
func (s *Service) rotateLogsPeriodically(ctx context.Context) {
	interval := time.Duration(s.cfg.RotationIntervalHours) * time.Hour
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.rotateLog(); err != nil {
				log.Printf("Error rotating logs: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// rotateLog performs log file rotation
func (s *Service) rotateLog() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.logFile == nil {
		return nil
	}

	// Close current log file
	s.logFile.Close()

	// Generate new filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	dir := filepath.Dir(s.cfg.LogFilePath)
	base := filepath.Base(s.cfg.LogFilePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	rotatedPath := filepath.Join(dir, fmt.Sprintf("%s-%s%s", name, timestamp, ext))

	// Rename current log file
	if err := os.Rename(s.cfg.LogFilePath, rotatedPath); err != nil {
		// If rename fails, open the original file again
		file, err := os.OpenFile(s.cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			s.logFile = file
			multiWriter := io.MultiWriter(os.Stdout, file)
			s.fileLogger = log.New(multiWriter, "", log.LstdFlags)
		}
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	// Open new log file
	file, err := os.OpenFile(s.cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// If opening fails, log to console only
		s.logFile = nil
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	s.logFile = file
	multiWriter := io.MultiWriter(os.Stdout, file)
	s.fileLogger = log.New(multiWriter, "", log.LstdFlags)

	log.Printf("Log rotated: %s", rotatedPath)

	// Clean up old logs if max count specified
	if s.cfg.MaxLogFiles > 0 {
		go s.cleanOldLogs()
	}

	return nil
}

// cleanOldLogs removes old log files beyond the max count
func (s *Service) cleanOldLogs() {
	dir := filepath.Dir(s.cfg.LogFilePath)
	base := filepath.Base(s.cfg.LogFilePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	// List all log files
	pattern := filepath.Join(dir, fmt.Sprintf("%s-*%s", name, ext))
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("Error listing log files: %v", err)
		return
	}

	// If we have fewer files than the max, do nothing
	if len(files) <= s.cfg.MaxLogFiles {
		return
	}

	// Sort files by modification time (oldest first)
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	fileInfos := make([]fileInfo, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			log.Printf("Error getting file info: %v", err)
			continue
		}
		fileInfos = append(fileInfos, fileInfo{
			path:    file,
			modTime: info.ModTime(),
		})
	}

	// Sort by modification time (oldest first)
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].modTime.Before(fileInfos[j].modTime)
	})

	// Delete oldest files until we're at the max
	toRemove := len(fileInfos) - s.cfg.MaxLogFiles
	for i := 0; i < toRemove; i++ {
		file := fileInfos[i].path
		if err := os.Remove(file); err != nil {
			log.Printf("Error removing old log file %s: %v", file, err)
		} else {
			log.Printf("Removed old log file: %s", file)
		}
	}
}

// Info logs an informational message
func (s *Service) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	s.log("INFO", msg)
}

// Error logs an error message
func (s *Service) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	s.log("ERROR", msg)
}

// Debug logs a debug message
func (s *Service) Debug(format string, v ...interface{}) {
	if !s.cfg.EnableDebugLogs {
		return
	}
	msg := fmt.Sprintf(format, v...)
	s.log("DEBUG", msg)
}

// log writes a message to the appropriate loggers
func (s *Service) log(level, msg string) {
	logMsg := fmt.Sprintf("[%s] %s", level, msg)

	if s.fileLogger != nil {
		s.fileLogger.Println(logMsg)
	} else {
		s.stdLogger.Println(logMsg)
	}
}
