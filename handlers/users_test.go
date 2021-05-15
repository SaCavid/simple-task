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
	h2 := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	e2 := echo.New()

	h2.registerUser(e2)
}

func (h *Server) registerUser(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing register. Expected Code: 201. Got:", err.Error())
	}
}
