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
	msg                    = `{"state": "win", "amount": "10.15", "transactionId": "Same identification"}`
	errorUsedTransactionId = `{"state": "win", "amount": "10.15", "transactionId": "Same identification"}`
	errorNoTransactionId   = `{"state": "win", "amount": "10.15", "transactionId": ""}`
	errorNoState           = `{"state": "", "amount": "10.15", "transactionId": "Some identification"}`
	errorState             = `{"state": "error-state", "amount": "10.15", "transactionId": "Some identification 2"}`
	errorNullAmount        = `{"state": "win", "amount": "", "transactionId": "Some identification 3"}`
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
	h.notRegistered(e)
}

func (h *Server) notAcceptableSourceType(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg))
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
		log.Println("Testing not logged. Expected Code: 403. Got:", err.Error())
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorUsedTransactionId))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec = httptest.NewRecorder()
	c2 := e.NewContext(req, rec)

	err = h.Handler(c2)
	if err != nil {
		log.Println("Testing repeat transaction id. Expected Code: 406. Got:", err.Error())
	}
}

func (h *Server) notRegistered(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg))
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
