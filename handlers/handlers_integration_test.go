package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/knblkxx/golangcalcV3/structs"
)

func TestCalculatorHandler_Post(t *testing.T) {
	reqBody := structs.Request{Expression: "1+1"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBuffer(b))
	rw := httptest.NewRecorder()

	CalculatorHandler(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("POST /calculate returned status %d; want %d", rw.Code, http.StatusOK)
	}

	var resp structs.ResponseOK
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Invalid JSON response: %v", err)
	}
	if resp.Id == 0 {
		t.Errorf("Expected non-zero Id in response, got %d", resp.Id)
	}
}

func TestExpressionAnswer_NotFound(t *testing.T) {
	// Setup router with no entries
	router := mux.NewRouter()
	router.HandleFunc("/expression/{id}", ExpressionAnswer)

	req := httptest.NewRequest(http.MethodGet, "/expression/99999", nil)
	rw := httptest.NewRecorder()

	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("GET non-existent id returned status %d; want %d", rw.Code, http.StatusOK)
	}
	var resp map[string]structs.ResponseResult
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
	res := resp["expression"]
	if res.Status != "not found" {
		t.Errorf("Expected status 'not found', got %q", res.Status)
	}
}
