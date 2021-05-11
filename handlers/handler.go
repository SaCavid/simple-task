package handlers

import (
	"fmt"
	"github.com/SaCavid/simple-task/models"
	"github.com/SaCavid/simple-task/service"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SourceType int

const (
	game SourceType = iota
	server
	payment
)

var (
	// must be unique names
	// index must be same in constants
	SourceTypes = [...]string{"game", "server", "payment"}
)

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

	Transactions []models.Data
	Repo         *service.TaskRepository
}

func (srv *Server) Register(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	if user.UserId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "user id can't be null"})
	}

	maximumFakeUsers := os.Getenv("N_FAKE_USERS") // for testing

	m, err := strconv.ParseInt(maximumFakeUsers, 10, 64)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: "general error"})
	}

	srv.Mu.Lock()
	if len(srv.UserBalances) >= int(m) {
		srv.Mu.Unlock()
		return echo.NewHTTPError(http.StatusNotAcceptable, &models.Response{Error: true, Message: "maximum user count reached"})
	}
	srv.Mu.Unlock()

	if srv.CheckUser(user.UserId) {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "user already registered"})
	}

	err = srv.Repo.Db.Create(user).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
	}

	srv.AddUser(user.UserId)

	return c.JSON(http.StatusOK, &models.Response{Message: "user registered"})
}

func (srv *Server) FetchUsersForTesting(c echo.Context) error {

	type User struct {
		UserId string
	}

	users := make([]User, 0)

	err := srv.Repo.Db.Table("users").Find(&users).Error
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &models.Response{Message: "users", Data: users})
}

func (srv *Server) Handler(c echo.Context) error {
	jd := new(models.JsonData)

	//log.Println(c.Request().Header.Get("Content-Length"))
	//log.Println(c.Request().Header.Get("Source-Type"))
	if err := c.Bind(&jd); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	if err := jd.ValidateData(); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	if srv.CheckTransactionId(jd.TransactionId) {
		log.Println("transaction id already used")
		return echo.NewHTTPError(http.StatusNotAcceptable, &models.Response{Error: true, Message: fmt.Sprintf("this transaction id already used")})
	}

	// Save transaction id not to use again ever if its failed
	srv.SaveTransactionId(jd.TransactionId)
	id := c.Request().Header.Get("Authorization")
	if id == "" {
		log.Println("not logged")
		return echo.NewHTTPError(http.StatusForbidden, &models.Response{Error: true, Message: "not logged"})
	}

	if !srv.CheckUser(id) {
		log.Println("user id didnt registered")
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "user didnt registered"})
	}

	jd.Source = c.Request().Header.Get("Source-Type")

	switch jd.State {
	case "win":

		err := srv.UserWin(id, jd)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
		}

		break
	case "lose":

		balance, err := srv.UserLost(id, jd)
		if err != nil {
			log.Println(err, jd.State, "-->", jd.Amount, "User balance:", balance)
			return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
		}

		break
	default:

		log.Println("error with state", jd.State)
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

func (srv *Server) SaveTransaction(data models.Data) {
	srv.Mu.Lock()
	srv.Transactions = append(srv.Transactions, data)
	srv.Mu.Unlock()
}

func (srv *Server) BulkInsertTransactions() {

	for {

		srv.Mu.Lock()

		if len(srv.Transactions) <= 0 {
			srv.Mu.Unlock()
			time.Sleep(30 * time.Second)
			continue
		}
		count := len(srv.Transactions)

		if count > 500 {
			count = 500
		}

		transactionsList := srv.Transactions[:count]
		srv.Transactions = srv.Transactions[count:]
		srv.Mu.Unlock()

		srv.Repo.Db.Begin()

		srv.Repo.Db.LogMode(true)
		var value []string
		var values []interface{}
		for _, data := range transactionsList {
			value = append(value, "(?,?,?,?,?,?,?,?)")
			values = append(values, data.CreatedAt)
			values = append(values, data.UpdatedAt)
			values = append(values, data.DeletedAt)
			values = append(values, data.UserId)
			values = append(values, data.State)
			values = append(values, data.Source)
			values = append(values, data.Amount)
			values = append(values, data.TransactionId)
		}

		stmt := fmt.Sprintf("INSERT INTO data (created_at, updated_at, deleted_at, user_id, state, source, amount, transaction_id) VALUES %s", strings.Join(value, ","))
		err := srv.Repo.Db.Begin().Exec(stmt, values...).Error
		if err != nil {
			srv.Repo.Db.Begin().Rollback()
			log.Println(err)
		}

		srv.Repo.Db.Begin().Commit()
		log.Println("Rows inserted:", len(values)/8)
	}
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

func (srv *Server) FetchUsers() error {

	users := make([]models.User, 0)

	err := srv.Repo.Db.Find(&users).Error
	if err != nil {
		return err
	}

	srv.Mu.Lock()
	for _, v := range users {
		srv.UserBalances[v.UserId] = v.Balance
	}
	srv.Mu.Unlock()

	return nil
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

	mData.CreatedAt = time.Now()
	mData.UpdatedAt = time.Now()

	mData.UserId = id
	mData.TransactionId = d.TransactionId
	mData.State = true
	mData.Amount = a

	mData.Source = i

	srv.SaveTransaction(*mData)
	//err = srv.CreateData(mData)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	return nil
}

func (srv *Server) UserLost(id string, d *models.JsonData) (float64, error) {

	a, err := strconv.ParseFloat(d.Amount, 64)
	if err != nil {
		return 0, err
	}
	var s SourceType

	i, err := s.IndexOf(d.Source)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	srv.Mu.Lock()
	balance := srv.UserBalances[id]
	if (balance - a) < 0 {
		srv.Mu.Unlock()
		return balance, fmt.Errorf("not enough user balance")
	}

	srv.UserBalances[id] = balance - a
	srv.Mu.Unlock()
	mData := new(models.Data)

	mData.CreatedAt = time.Now()
	mData.UpdatedAt = time.Now()

	mData.UserId = id
	mData.TransactionId = d.TransactionId
	mData.State = false
	mData.Amount = a
	mData.Source = i

	srv.SaveTransaction(*mData)
	//err = srv.CreateData(mData)
	//if err != nil {
	//	log.Println(err)
	//	return 0, err
	//}

	return balance, nil
}

func (srv *Server) PostProcessing() {

}
