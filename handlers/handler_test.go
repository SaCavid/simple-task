package handlers

import (
	"github.com/SaCavid/simple-task/models"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	msg1                   = `{"state": "win", "amount": "10.15", "transactionId": "Same identification 1"}`
	msg2                   = `{"state": "win", "amount": "10.15", "transactionId": "Same identification 2"}`
	msg3                   = `{"state": "win", "amount": "10.15", "transactionId": "Same identification 3"}`
	errorUsedTransactionId = `{"state": "win", "amount": "10.15", "transactionId": "Same identification 1"}`
	errorNoTransactionId   = `{"state": "win", "amount": "10.15", "transactionId": ""}`
	errorNoState           = `{"state": "", "amount": "10.15", "transactionId": "Some identification"}`
	errorState             = `{"state": "error-state", "amount": "10.15", "transactionId": "Some identification 4"}`
	errorNullAmount        = `{"state": "win", "amount": "", "transactionId": "Some identification 5"}`
)

func TestServer_Handler(t *testing.T) {

	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	e := echo.New()

	h.notAcceptableSourceType(e)
	h.noState(e)
	h.wrongState(e)
	h.noTransactionId(e)
	h.noAmount(e)
	h.sameTransactionId(e)
	h.notLogged(e)
	h.notRegistered(e)
}

func (h *Server) notAcceptableSourceType(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "not-source")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing not acceptable source type. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) noState(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNoState))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null state. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) wrongState(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorState))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing wrong state. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) noTransactionId(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNoTransactionId))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null transaction Id. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) noAmount(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNullAmount))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null amount. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) sameTransactionId(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorUsedTransactionId))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing repeat transaction id. Expected Code: 406. Got:", err.Error())
	}
}

func (h *Server) notLogged(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg2))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing not logged. Expected Code: 403. Got:", err.Error())
	}
}

func (h *Server) notRegistered(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg3))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "not-registered-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing not registered. Expected Code: 400. Got:", err.Error())
	}
}

func (h *Server) benchmarkRegistered(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg3))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "not-registered-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h.Handler(c)
}

// Benchmark
func BenchmarkServer_Handler(b *testing.B) {

	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	e := echo.New()

	// run the function N times
	N := 100000
	for i := 0; i < N; i++ {
		h.benchmarkRegistered(e)
	}
}
