package calculator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/maxpawgdbs/yandex-go/structs"
)

var ts *httptest.Server

// TestMain sets up a mock orchestrator server and adjusts delays.
func TestMain(m *testing.M) {
	// Zero out delays
	TIME_ADDITION_MS = 0
	TIME_SUBTRACTION_MS = 0
	TIME_MULTIPLICATIONS_MS = 0
	TIME_DIVISIONS_MS = 0

	// Start a test HTTP server to mock TaskURL
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req structs.AgentResponse
		json.NewDecoder(r.Body).Decode(&req)
		// compute result locally
		var res float64
		switch req.Operation {
		case "+":
			res = req.Arg1 + req.Arg2
		case "-":
			res = req.Arg1 - req.Arg2
		case "*":
			res = req.Arg1 * req.Arg2
		case "/":
			if req.Arg2 != 0 {
				res = req.Arg1 / req.Arg2
			}
		}
		json.NewEncoder(w).Encode(structs.AgentResult{Result: res})
	}))
	// point TaskURL to our test server
	TaskURL = ts.URL

	code := m.Run()
	ts.Close()
	os.Exit(code)
}

func TestCalcExpressionSimple(t *testing.T) {
	expr := "1+2"
	want := "3.000000"
	got, err := CalcExpression(expr)
	if err != nil {
		t.Fatalf("CalcExpression(%q) returned unexpected error: %v", expr, err)
	}
	if got != want {
		t.Errorf("CalcExpression(%q) = %q; want %q", expr, got, want)
	}
}

func TestCalcExpressionPrecedence(t *testing.T) {
	expr := "2+3*4"
	want := "14.000000"
	got, err := CalcExpression(expr)
	if err != nil {
		t.Fatalf("CalcExpression(%q) returned unexpected error: %v", expr, err)
	}
	if got != want {
		t.Errorf("CalcExpression(%q) = %q; want %q", expr, got, want)
	}
}

func TestCalcExpressionDivideByZero(t *testing.T) {
	expr := "1/0"
	_, err := CalcExpression(expr)
	if err == nil {
		t.Errorf("CalcExpression(%q) expected error for division by zero, got nil", expr)
	}
}
