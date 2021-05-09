package handlers

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"simple-task/models"
	"strconv"
	"sync"
)

//const (
//	game = iota
//	server
//	payment
//	client
//)

//var (
//	SourceTypes = [...]string{"game", "server", "payment", "client"}
//)

type Server struct {
	Mu sync.Mutex

	// For faster Transaction id check - must be unique id -- Better to use Redis
	TransactionIds map[string]string

	// For faster user balance check -- Better to use Redis
	UserBalances map[string]float64
}

func (srv *Server) Handler(c echo.Context) error {
	d := new(models.Data)
	// log.Println(c.Request().Header.Get("Content-Length"))
	// log.Println(c.Request().Header.Get("Source-Type"))
	if err := c.Bind(&d); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	if err := d.ValidateData(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	if srv.CheckTransactionId(d.TransactionId) {
		return echo.NewHTTPError(http.StatusNotAcceptable, &models.Response{Error: true, Message: fmt.Sprintf("this transaction id already used")})
	}

	switch d.State {
	case "win":

		err := srv.UserWin("user id", d.Amount)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
		}

		break
	case "lost":

		err := srv.UserLost("user id", d.Amount)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
		}

		break
	default:
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "error with state"})
	}

	return c.JSON(http.StatusOK, &models.Response{Message: "transaction processed"})
}

func (srv *Server) CheckTransactionId(id string) bool {
	srv.Mu.Lock()
	_, ok := srv.TransactionIds[id]
	srv.Mu.Unlock()
	return ok
}

func (srv *Server) SaveTransactionId(id string) {
	srv.Mu.Lock()
	srv.TransactionIds[id] = id
	srv.Mu.Unlock()
}

func (srv *Server) CheckUser(id string) bool {
	srv.Mu.Lock()
	_, ok := srv.UserBalances[id]
	srv.Mu.Unlock()
	return ok
}

func (srv *Server) AddUser(id string) {
	srv.Mu.Lock()
	srv.UserBalances[id] = 0
	srv.Mu.Unlock()
}

func (srv *Server) SaveUser(id string, balance float64) {
	srv.Mu.Lock()
	srv.UserBalances[id] = balance
	srv.Mu.Unlock()
}

func (srv *Server) UserWin(id, amount string) error {

	a, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}

	srv.Mu.Lock()
	balance := srv.UserBalances[id]
	srv.UserBalances[id] = balance + a
	srv.Mu.Unlock()

	return nil
}

func (srv *Server) UserLost(id, amount string) error {

	a, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}

	srv.Mu.Lock()
	balance := srv.UserBalances[id]
	if (balance - a) < 0 {
		srv.Mu.Unlock()
		return fmt.Errorf("not enough user balance")
	}

	srv.UserBalances[id] = balance - a
	srv.Mu.Unlock()

	return nil
}
