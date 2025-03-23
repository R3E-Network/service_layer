package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Function management handlers
func (s *Server) handleCreateFunction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var request struct {
		Name       string   `json:"name"`
		Code       string   `json:"code"`
		SecretRefs []string `json:"secretRefs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Register function with TEE manager
	functionID := generateID()
	err := s.teeManager.RegisterFunction(functionID, request.Code, request.SecretRefs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id":   functionID,
		"name": request.Name,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (s *Server) handleListFunctions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	// In a real implementation, this would fetch functions from storage
	functions := []map[string]interface{}{
		{
			"id":   "func1",
			"name": "Example Function",
		},
	}

	respondWithJSON(w, http.StatusOK, functions)
}

func (s *Server) handleGetFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	functionID := vars["id"]

	// In a real implementation, this would fetch the function from storage
	function := map[string]interface{}{
		"id":         functionID,
		"name":       "Example Function",
		"created_at": "2023-05-10T15:04:05Z",
	}

	respondWithJSON(w, http.StatusOK, function)
}

func (s *Server) handleUpdateFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	functionID := vars["id"]

	var request struct {
		Name       string   `json:"name"`
		Code       string   `json:"code"`
		SecretRefs []string `json:"secretRefs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Update function in TEE
	err := s.teeManager.RegisterFunction(functionID, request.Code, request.SecretRefs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":   functionID,
		"name": request.Name,
	})
}

func (s *Server) handleDeleteFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	functionID := vars["id"]

	// In a real implementation, this would delete the function

	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
		"id":     functionID,
	})
}

func (s *Server) handleExecuteFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	functionID := vars["id"]
	userID := r.Context().Value("userID").(string)

	var request struct {
		Params map[string]interface{} `json:"params"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Allocate gas for execution
	estimatedGas := 5.0
	allocationID, err := s.gasBankSvc.AllocateGas(userID, functionID, estimatedGas)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Execute function in TEE
	result, err := s.teeManager.ExecuteFunction(r.Context(), functionID, request.Params)
	if err != nil {
		// Refund gas if execution fails
		s.gasBankSvc.FinalizeGasUsage(allocationID, 0)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Finalize gas usage
	actualGas := 3.0
	s.gasBankSvc.FinalizeGasUsage(allocationID, actualGas)

	respondWithJSON(w, http.StatusOK, result)
}

// Secret management handlers
func (s *Server) handleCreateSecret(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var request struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// In a real implementation, this would securely store the secret in the TEE
	secretID := generateID()

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"id":   secretID,
		"name": request.Name,
	})
}

func (s *Server) handleListSecrets(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	// In a real implementation, this would fetch secrets metadata (not values)
	secrets := []map[string]string{
		{
			"id":   "secret1",
			"name": "API Key",
		},
	}

	respondWithJSON(w, http.StatusOK, secrets)
}

func (s *Server) handleDeleteSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secretID := vars["id"]

	// In a real implementation, this would delete the secret

	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
		"id":     secretID,
	})
}

// Trigger management handlers
func (s *Server) handleCreateTrigger(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var request struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		FunctionID string `json:"functionId"`
		Schedule   string `json:"schedule,omitempty"`
		Condition  string `json:"condition,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// In a real implementation, this would create a trigger in the trigger service
	triggerID := generateID()

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"id":   triggerID,
		"name": request.Name,
	})
}

func (s *Server) handleListTriggers(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	// In a real implementation, this would fetch triggers from storage
	triggers := []map[string]string{
		{
			"id":   "trigger1",
			"name": "Daily Update",
			"type": "schedule",
		},
	}

	respondWithJSON(w, http.StatusOK, triggers)
}

func (s *Server) handleGetTrigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	triggerID := vars["id"]

	// In a real implementation, this would fetch the trigger from storage
	trigger := map[string]string{
		"id":          triggerID,
		"name":        "Daily Update",
		"type":        "schedule",
		"function_id": "func1",
		"schedule":    "0 0 * * *",
	}

	respondWithJSON(w, http.StatusOK, trigger)
}

func (s *Server) handleUpdateTrigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	triggerID := vars["id"]

	var request struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		FunctionID string `json:"functionId"`
		Schedule   string `json:"schedule,omitempty"`
		Condition  string `json:"condition,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// In a real implementation, this would update the trigger

	respondWithJSON(w, http.StatusOK, map[string]string{
		"id":   triggerID,
		"name": request.Name,
	})
}

func (s *Server) handleDeleteTrigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	triggerID := vars["id"]

	// In a real implementation, this would delete the trigger

	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
		"id":     triggerID,
	})
}

// Price feed handlers
func (s *Server) handleListPriceFeeds(w http.ResponseWriter, r *http.Request) {
	prices := s.priceFeedSvc.GetAllPrices()

	// Convert to simple format for response
	result := make([]map[string]interface{}, 0)
	for token, price := range prices {
		result = append(result, map[string]interface{}{
			"token":     token,
			"price":     price.Price,
			"timestamp": price.Timestamp,
		})
	}

	respondWithJSON(w, http.StatusOK, result)
}

func (s *Server) handleGetPriceFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	price, err := s.priceFeedSvc.GetPrice(token)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, price)
}

func (s *Server) handleCreatePriceFeed(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var request struct {
		Token  string `json:"token"`
		Source string `json:"source"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// In a real implementation, this would add a custom price feed

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"token":  request.Token,
		"status": "added",
	})
}

// GasBank handlers
func (s *Server) handleGetGasBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	balance, err := s.gasBankSvc.GetBalance(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"balance": balance,
	})
}

func (s *Server) handleDepositGas(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var request struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if request.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "Amount must be positive")
		return
	}

	if err := s.gasBankSvc.DepositGas(userID, request.Amount); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	balance, _ := s.gasBankSvc.GetBalance(userID)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":   userID,
		"deposited": request.Amount,
		"balance":   balance,
	})
}

func (s *Server) handleWithdrawGas(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var request struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if request.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "Amount must be positive")
		return
	}

	if err := s.gasBankSvc.WithdrawGas(userID, request.Amount); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	balance, _ := s.gasBankSvc.GetBalance(userID)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":   userID,
		"withdrawn": request.Amount,
		"balance":   balance,
	})
}

// Utility functions
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func generateID() string {
	// In a real implementation, this would generate a unique ID
	// using UUID or another algorithm
	return "id-" + time.Now().Format("20060102150405")
}
