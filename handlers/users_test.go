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

func TestServer_Register(t *testing.T) {
	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	e := echo.New()

	h.registerUser(e)
}

func (h *Server) registerUser(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing register. Expected Code: 400. Got:", err.Error())
	}
}
