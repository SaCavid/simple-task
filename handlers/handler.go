package handlers

import (
	"../models"
	"../service"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type SourceType int

var (
	// must be unique names
	// index must be same in constants
	SourceTypes = [...]string{"game", "server", "payment"}
)

const (
	game SourceType = iota
	server
	payment
)

// Example usage
//var d SourceType = game
//d.String()
func (s SourceType) String() string {
	return SourceTypes[s]
}

func (s SourceType) IndexOf(name string) (int, error) {

	for k, v := range SourceTypes {
		if v == name {
			return k, nil
		}
	}

	return -1, fmt.Errorf("not acceptable source type")
}

type Server struct {
	Mu sync.Mutex

	// For faster Transaction id check - must be unique id -- Better to use Redis
	TransactionIds map[string]string

	// For faster user balance check -- Better to use Redis
	UserBalances map[string]float64

	Repo *service.TaskRepository
}

func (srv *Server) Handler(c echo.Context) error {
	jd := new(models.JsonData)
	log.Println(c.Request().Header.Get("Content-Length"))
	log.Println(c.Request().Header.Get("Source-Type"))
	if err := c.Bind(&jd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	if err := jd.ValidateData(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	if srv.CheckTransactionId(jd.TransactionId) {
		return echo.NewHTTPError(http.StatusNotAcceptable, &models.Response{Error: true, Message: fmt.Sprintf("this transaction id already used")})
	}
	jd.Source = c.Request().Header.Get("Source-Type")
	log.Println(jd.Source)
	switch jd.State {
	case "win":

		err := srv.UserWin("user id", jd)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
		}

		break
	case "lost":

		err := srv.UserLost("user id", jd)
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

func (srv *Server) UserWin(id string, d *models.JsonData) error {

	a, err := strconv.ParseFloat(d.Amount, 64)
	if err != nil {
		return err
	}
	var s SourceType

	i, err := s.IndexOf(d.Source)
	if err != nil {
		log.Println(d.Source, err)
		return err
	}

	srv.Mu.Lock()
	balance := srv.UserBalances[id]
	srv.UserBalances[id] = balance + a
	srv.Mu.Unlock()
	mData := new(models.Data)

	mData.TransactionId = d.TransactionId
	mData.State = true
	mData.Amount = a

	mData.Source = i

	err = srv.CreateData(mData)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (srv *Server) UserLost(id string, d *models.JsonData) error {

	a, err := strconv.ParseFloat(d.Amount, 64)
	if err != nil {
		return err
	}
	var s SourceType

	i, err := s.IndexOf(d.Source)
	if err != nil {
		log.Println(err)
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
	mData := new(models.Data)

	mData.TransactionId = d.TransactionId
	mData.State = false
	mData.Amount = a
	mData.Source = i

	err = srv.CreateData(mData)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (srv *Server) PostProcessing() {

}
