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
	errorNullAmount        = `{"state": "win", "amount": "", "transactionId": "Some identification"}`
)

func TestServer_Handler(t *testing.T) {
	notAcceptableSourceType()
	noState()
	noTransactionId()
	noAmount()
	sameTransactionId()
	notRegistered()
}

func notAcceptableSourceType() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "not-source")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing not registered state. Expected Code: 400. Got:", err.Error())
	}

}

func noState() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNoState))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null state. Expected Code: 400. Got:", err.Error())
	}

}

func noTransactionId() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNoTransactionId))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null transaction Id. Expected Code: 400. Got:", err.Error())
	}

}

func noAmount() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNullAmount))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null amount. Expected Code: 400. Got:", err.Error())
	}

}

func sameTransactionId() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorUsedTransactionId))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

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

func notRegistered() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "not-registered-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing not registered. Expected Code: 400. Got:", err.Error())
	}
}
