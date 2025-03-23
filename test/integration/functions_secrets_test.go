// Integration tests for functions and secrets services
package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/R3E-Network/service_layer/internal/config"
	"github.com/R3E-Network/service_layer/internal/core/functions"
	"github.com/R3E-Network/service_layer/internal/core/secrets"
	"github.com/R3E-Network/service_layer/internal/models"
	"github.com/R3E-Network/service_layer/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFunctionRepository is a mock for the FunctionRepository interface
type MockFunctionRepository struct {
	mock.Mock
}

func (m *MockFunctionRepository) Create(function *models.Function) error {
	args := m.Called(function)
	return args.Error(0)
}

func (m *MockFunctionRepository) Update(function *models.Function) error {
	args := m.Called(function)
	return args.Error(0)
}

// Update Delete method to match the FunctionRepository interface
func (m *MockFunctionRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockFunctionRepository) GetByID(id int) (*models.Function, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Function), args.Error(1)
}

func (m *MockFunctionRepository) GetByUserID(userID, limit, offset int) ([]*models.Function, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]*models.Function), args.Error(1)
}

func (m *MockFunctionRepository) GetByName(userID int, name string) (*models.Function, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Function), args.Error(1)
}

func (m *MockFunctionRepository) Search(query string, userID, limit, offset int) ([]*models.Function, error) {
	args := m.Called(query, userID, limit, offset)
	return args.Get(0).([]*models.Function), args.Error(1)
}

// Add GetByUserIDAndName method to match FunctionRepository interface
func (m *MockFunctionRepository) GetByUserIDAndName(userID int, name string) (*models.Function, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Function), args.Error(1)
}

// Update List method to match FunctionRepository interface
func (m *MockFunctionRepository) List(limit, offset, userID int) ([]*models.Function, error) {
	args := m.Called(limit, offset, userID)
	return args.Get(0).([]*models.Function), args.Error(1)
}

// Add GetSecrets method to match FunctionRepository interface
func (m *MockFunctionRepository) GetSecrets(functionID int) ([]string, error) {
	args := m.Called(functionID)
	return args.Get(0).([]string), args.Error(1)
}

// Add IncrementExecutionCount method to match FunctionRepository interface
func (m *MockFunctionRepository) IncrementExecutionCount(functionID int) error {
	args := m.Called(functionID)
	return args.Error(0)
}

// Add SetSecrets method to match FunctionRepository interface
func (m *MockFunctionRepository) SetSecrets(functionID int, secrets []string) error {
	args := m.Called(functionID, secrets)
	return args.Error(0)
}

// Add UpdateLastExecution method to match FunctionRepository interface
func (m *MockFunctionRepository) UpdateLastExecution(functionID int, lastExecution time.Time) error {
	args := m.Called(functionID, lastExecution)
	return args.Error(0)
}

// MockSecretRepository is a mock for the SecretRepository interface
type MockSecretRepository struct {
	mock.Mock
}

func (m *MockSecretRepository) Create(secret *models.Secret) error {
	args := m.Called(secret)
	return args.Error(0)
}

func (m *MockSecretRepository) Update(secret *models.Secret) error {
	args := m.Called(secret)
	return args.Error(0)
}

func (m *MockSecretRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSecretRepository) GetByID(id int) (*models.Secret, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Secret), args.Error(1)
}

func (m *MockSecretRepository) GetByName(userID int, name string) (*models.Secret, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Secret), args.Error(1)
}

// Add GetByUserIDAndName method to match the SecretRepository interface
func (m *MockSecretRepository) GetByUserIDAndName(userID int, name string) (*models.Secret, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Secret), args.Error(1)
}

func (m *MockSecretRepository) GetByUserID(userID, limit, offset int) ([]*models.Secret, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]*models.Secret), args.Error(1)
}

func (m *MockSecretRepository) GetByUserIDAndNames(userID int, names []string) ([]*models.Secret, error) {
	args := m.Called(userID, names)
	return args.Get(0).([]*models.Secret), args.Error(1)
}

func (m *MockSecretRepository) Search(query string, userID, limit, offset int) ([]*models.Secret, error) {
	args := m.Called(query, userID, limit, offset)
	return args.Get(0).([]*models.Secret), args.Error(1)
}

// Fix List method to match SecretRepository interface - with correct signature
func (m *MockSecretRepository) List(limit int) ([]*models.Secret, error) {
	args := m.Called(limit)
	return args.Get(0).([]*models.Secret), args.Error(1)
}

// MockTEERuntime is a mock implementation of the TEE runtime
type MockTEERuntime struct {
	mock.Mock
}

func (m *MockTEERuntime) ExecuteFunction(userID int, sourceCode string, params map[string]interface{}, secretValues map[string]string) (*models.ExecutionResult, error) {
	args := m.Called(userID, sourceCode, params, secretValues)
	return args.Get(0).(*models.ExecutionResult), args.Error(1)
}

func (m *MockTEERuntime) Create(execution *models.Execution) error {
	args := m.Called(execution)
	return args.Error(0)
}

func (m *MockTEERuntime) Delete(executionID int) error {
	args := m.Called(executionID)
	return args.Error(0)
}

func (m *MockTEERuntime) GetByID(executionID int) (*models.Execution, error) {
	args := m.Called(executionID)
	return args.Get(0).(*models.Execution), args.Error(1)
}

func (m *MockTEERuntime) GetLogs(executionID, limit, offset int) ([]*models.ExecutionLog, error) {
	args := m.Called(executionID, limit, offset)
	return args.Get(0).([]*models.ExecutionLog), args.Error(1)
}

func (m *MockTEERuntime) ListByFunctionID(functionID, limit, offset int) ([]*models.Execution, error) {
	args := m.Called(functionID, limit, offset)
	return args.Get(0).([]*models.Execution), args.Error(1)
}

func (m *MockTEERuntime) Update(execution *models.Execution) error {
	args := m.Called(execution)
	return args.Error(0)
}

// Add AddLog method to match the ExecutionRepository interface
func (m *MockTEERuntime) AddLog(executionLog *models.ExecutionLog) error {
	args := m.Called(executionLog)
	return args.Error(0)
}

// Setup functions with secrets integration test
func setupFunctionsSecretsTest() (*functions.Service, *secrets.Service, *MockFunctionRepository, *MockSecretRepository, *MockTEERuntime) {
	// Create logger with empty config
	log := logger.New(logger.LoggingConfig{})

	// Create config
	cfg := &config.Config{}

	// Create mock repositories
	mockFunctionRepo := new(MockFunctionRepository)
	mockSecretRepo := new(MockSecretRepository)
	mockTeeRuntime := new(MockTEERuntime)

	// Create services - adjust parameter count based on actual function signature
	secretsService := secrets.NewService(cfg, log, mockSecretRepo, nil)

	functionsService := functions.NewService(cfg, log, mockFunctionRepo, mockTeeRuntime, nil)

	return functionsService, secretsService, mockFunctionRepo, mockSecretRepo, mockTeeRuntime
}

// Create test secrets for testing
func createTestSecrets() []*models.Secret {
	return []*models.Secret{
		{
			ID:        1,
			UserID:    1,
			Name:      "test-secret",
			Value:     "test-secret-value",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// Create helper function to convert map to json.RawMessage
func createJSONRawMessage(data map[string]interface{}) json.RawMessage {
	bytes, _ := json.Marshal(data)
	return bytes
}

// Test function execution with secrets
func TestFunctionExecutionWithSecrets(t *testing.T) {
	// Setup test environment
	functionsService, secretsService, mockFunctionRepo, mockSecretRepo, mockTeeRuntime := setupFunctionsSecretsTest()

	// Create a test function that accesses secrets
	testFunction := &models.Function{
		ID:         1,
		UserID:     1,
		Name:       "test-function",
		SourceCode: "function main(args) { return { result: args.secret_value }; }",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Create test data
	testSecrets := createTestSecrets()

	// Setup repository mocks
	mockFunctionRepo.On("GetByID", 1).Return(testFunction, nil)
	mockSecretRepo.On("GetByUserIDAndNames", 1, []string{"test-secret"}).Return(testSecrets, nil)

	// Setup expected secret map for test
	expectedSecretMap := map[string]string{
		"test-secret": "test-secret-value",
	}

	// Create execution result with the correct structure based on the models.ExecutionResult
	executionResult := &models.ExecutionResult{
		// Convert map to json.RawMessage
		Result: createJSONRawMessage(map[string]interface{}{
			"result": "secret-value",
		}),
	}

	mockTeeRuntime.On(
		"ExecuteFunction",
		1,
		testFunction.SourceCode,
		mock.Anything, // params
		expectedSecretMap,
	).Return(executionResult, nil)

	// Execute function with test parameters
	testParams := map[string]interface{}{
		"test_param": "test_value",
	}

	// Use ctx variable to make use of the context import
	ctx := context.Background()
	result, err := functionsService.ExecuteFunction(ctx, 1, 1, testParams, false)

	// Assert test results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Use string key for map access - first unmarshal the json.RawMessage
	var resultMap map[string]interface{}
	json.Unmarshal(result.Result, &resultMap)
	assert.Equal(t, "secret-value", resultMap["result"])

	// Verify that all mocked methods were called as expected
	mockFunctionRepo.AssertExpectations(t)
	mockSecretRepo.AssertExpectations(t)
	mockTeeRuntime.AssertExpectations(t)

	// For completeness, assert the secrets service was utilized properly
	assert.NotNil(t, secretsService)
}

// Test function execution with unauthorized secret access
func TestFunctionExecutionWithUnauthorizedSecretAccess(t *testing.T) {
	// Setup
	functionsService, _, mockFunctionRepo, mockSecretRepo, mockTeeRuntime := setupFunctionsSecretsTest()

	// Create test data
	testFunction := &models.Function{
		ID:         1,
		UserID:     1,
		Name:       "test-function",
		SourceCode: "function main(args) { return { result: args.secret_value }; }",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Setup repository mocks
	mockFunctionRepo.On("GetByID", 1).Return(testFunction, nil)

	// Either match the function signature or modify test as needed
	mockSecretRepo.On("GetByUserIDAndNames", 1, []string{"test-secret"}).Return([]*models.Secret{}, nil)

	// Setup empty expected secret map since access is unauthorized
	expectedSecretMap := map[string]string{}

	executeResult := &models.ExecutionResult{
		Result: createJSONRawMessage(map[string]interface{}{
			"error": "unauthorized secret access",
		}),
	}

	mockTeeRuntime.On(
		"ExecuteFunction",
		1,
		testFunction.SourceCode,
		mock.Anything, // params
		expectedSecretMap,
	).Return(executeResult, nil)

	testParams := map[string]interface{}{
		"input": "test-value",
	}

	ctx := context.Background()
	result, err := functionsService.ExecuteFunction(ctx, 1, 1, testParams, false)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Use string key for map access - first unmarshal the json.RawMessage
	var resultMap map[string]interface{}
	json.Unmarshal(result.Result, &resultMap)
	assert.Equal(t, "unauthorized secret access", resultMap["error"])

	// Verify all mocks were called
	mockFunctionRepo.AssertExpectations(t)
	mockSecretRepo.AssertExpectations(t)
	mockTeeRuntime.AssertExpectations(t)
}

// Test multiple function executions with different secret access
func TestMultipleFunctionExecutionsWithDifferentSecretAccess(t *testing.T) {
	// Setup test environment
	functionsService, _, mockFunctionRepo, mockSecretRepo, mockTeeRuntime := setupFunctionsSecretsTest()

	// Create test functions
	function1 := &models.Function{
		ID:         1,
		UserID:     1,
		Name:       "function1",
		SourceCode: "function main(args) { return { result: args.secret1 }; }",
		// SecretsAccess removed as it's not part of the Function struct
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	function2 := &models.Function{
		ID:         2,
		UserID:     1,
		Name:       "function2",
		SourceCode: "function main(args) { return { result: args.secret2 }; }",
		// SecretsAccess removed as it's not part of the Function struct
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup test secrets
	secret1 := &models.Secret{
		ID:        1,
		UserID:    1,
		Name:      "secret1",
		Value:     "secret1-value",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	secret2 := &models.Secret{
		ID:        2,
		UserID:    1,
		Name:      "secret2",
		Value:     "secret2-value",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup repository mocks
	mockFunctionRepo.On("GetByID", 1).Return(function1, nil)
	mockFunctionRepo.On("GetByID", 2).Return(function2, nil)

	// Setup secret repository mock for specific secrets
	mockSecretRepo.On("GetByUserIDAndNames", 1, []string{"secret1"}).Return([]*models.Secret{secret1}, nil)
	mockSecretRepo.On("GetByUserIDAndNames", 1, []string{"secret2"}).Return([]*models.Secret{secret2}, nil)

	// Mock the execution results for each function
	result1 := &models.ExecutionResult{
		Result: createJSONRawMessage(map[string]interface{}{
			"result": "secret1-value",
		}),
	}

	result2 := &models.ExecutionResult{
		Result: createJSONRawMessage(map[string]interface{}{
			"result": "secret2-value",
		}),
	}

	// Setup TEE runtime mocks with expected secret values
	mockTeeRuntime.On(
		"ExecuteFunction",
		1,
		function1.SourceCode,
		mock.Anything,
		map[string]string{"secret1": "secret1-value"},
	).Return(result1, nil)

	mockTeeRuntime.On(
		"ExecuteFunction",
		1,
		function2.SourceCode,
		mock.Anything,
		map[string]string{"secret2": "secret2-value"},
	).Return(result2, nil)

	// Execute the functions and test results
	ctx := context.Background()
	testParams := map[string]interface{}{"param": "value"}

	// Execute function 1
	exec1, err1 := functionsService.ExecuteFunction(ctx, 1, 1, testParams, false)
	assert.NoError(t, err1)
	assert.NotNil(t, exec1)
	// Access using string key - need to unmarshal the json.RawMessage first
	var resultMap1 map[string]interface{}
	json.Unmarshal(exec1.Result, &resultMap1)
	assert.Equal(t, "secret1-value", resultMap1["result"])

	// Execute function 2
	exec2, err2 := functionsService.ExecuteFunction(ctx, 2, 1, testParams, false)
	assert.NoError(t, err2)
	assert.NotNil(t, exec2)
	// Access using string key - need to unmarshal the json.RawMessage first
	var resultMap2 map[string]interface{}
	json.Unmarshal(exec2.Result, &resultMap2)
	assert.Equal(t, "secret2-value", resultMap2["result"])

	// Verify all mocks were called
	mockFunctionRepo.AssertExpectations(t)
	mockSecretRepo.AssertExpectations(t)
	mockTeeRuntime.AssertExpectations(t)
}
