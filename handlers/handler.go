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
	Balance      bool // true --> there is not saved balance
	UserBalances map[string]models.Balance

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

	jd.Source = c.Request().Header.Get("Source-Type")

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
		var s SourceType

		i, err := s.IndexOf(jd.Source)
		if err != nil {
			log.Println(jd.Source, err)
			return err
		}

		log.Println("not logged")
		data := models.Data{
			UserId:        "",
			State:         false,
			Source:        i,
			Status:        2, // error . saved for unique transaction id.
			Amount:        0,
			TransactionId: "",
		}

		srv.SaveTransaction(data)
		return echo.NewHTTPError(http.StatusForbidden, &models.Response{Error: true, Message: "not logged"})
	}

	if !srv.CheckUser(id) {
		log.Println("user id didnt registered")
		var s SourceType

		i, err := s.IndexOf(jd.Source)
		if err != nil {
			log.Println(jd.Source, err)
			return err
		}

		log.Println("not logged")
		data := models.Data{
			UserId:        "",
			State:         false,
			Source:        i,
			Status:        2, // error . saved for unique transaction id.
			Amount:        0,
			TransactionId: "",
		}

		srv.SaveTransaction(data)
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "user didnt registered"})
	}

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
			time.Sleep(10 * time.Second)
			continue
		}
		count := len(srv.Transactions)

		if count > 500 {
			count = 500
		}

		transactionsList := srv.Transactions[:count]
		srv.Transactions = srv.Transactions[count:]
		srv.Mu.Unlock()

		tx := srv.Repo.Db.Begin()
		err := tx.Error
		if err != nil {
			log.Println(err)
			continue
		}

		var value []string
		var values []interface{}
		for _, data := range transactionsList {
			value = append(value, "(?,?,?,?,?,?,?,?,?)")
			values = append(values, data.CreatedAt)
			values = append(values, data.UpdatedAt)
			values = append(values, data.DeletedAt)
			values = append(values, data.UserId)
			values = append(values, data.State)
			values = append(values, data.Status)
			values = append(values, data.Source)
			values = append(values, data.Amount)
			values = append(values, data.TransactionId)
		}

		stmt := fmt.Sprintf("INSERT INTO data (created_at, updated_at, deleted_at, user_id, state, status, source, amount, transaction_id) VALUES %s", strings.Join(value, ","))
		err = tx.Exec(stmt, values...).Error
		if err != nil {
			tx.Rollback()
			log.Println(err)
		}

		err = tx.Commit().Error
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Rows inserted:", len(values)/8)
	}
}

func (srv *Server) BulkUpdateBalances() {

	for {

		if !srv.Balance {
			time.Sleep(10 * time.Second)
			continue
		}

		type balance struct {
			UserId string
			Amount float64
		}

		balancesList := make([]balance, 0)

		srv.Mu.Lock()
		for k, v := range srv.UserBalances {
			if v.Saved {
				s := balance{
					UserId: k,
					Amount: v.Amount,
				}
				balancesList = append(balancesList, s)
			}
			v.Saved = false
		}
		srv.Balance = false
		srv.Mu.Unlock()

		for {
			count := len(balancesList)

			if count == 0 {
				break
			}

			if count > 500 {
				count = 500
			}

			chunkList := balancesList[:count]
			var value []string
			for _, data := range chunkList {
				value = append(value, fmt.Sprintf("('%s',%.2f)", data.UserId, data.Amount))
			}

			srv.Repo.Db.LogMode(true)
			if err := srv.Repo.Db.Exec(fmt.Sprintf("UPDATE users AS u SET balance = data.a FROM (VALUES %s) AS data(user_id, a) WHERE t.user_id = data.user_id", strings.Join(value, ","))).Error; err != nil {
				log.Println(err)
				continue
			}

			balancesList = balancesList[count:]
		}

		log.Println("Rows inserted:", len(balancesList))
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
	b := models.Balance{}
	b.Amount = 0
	b.Saved = false
	srv.UserBalances[id] = b
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
		b := models.Balance{
			Amount: v.Balance,
			Saved:  false,
		}
		srv.UserBalances[v.UserId] = b
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
	b := srv.UserBalances[id]
	b.Amount = b.Amount + a
	b.Saved = true     // not saved
	srv.Balance = true // not saved balance in map
	srv.UserBalances[id] = b
	srv.Mu.Unlock()
	mData := models.Data{}

	mData.CreatedAt = time.Now()
	mData.UpdatedAt = time.Now()

	mData.UserId = id
	mData.TransactionId = d.TransactionId
	mData.State = true
	mData.Amount = a
	mData.Status = 1

	mData.Source = i

	srv.SaveTransaction(mData)
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
	b := srv.UserBalances[id]
	if (b.Amount - a) < 0 {
		srv.Mu.Unlock()
		return b.Amount, fmt.Errorf("not enough user balance")
	}

	b.Amount = b.Amount - a
	b.Saved = true     // user balance not saved
	srv.Balance = true // not saved balance in map
	srv.UserBalances[id] = b
	srv.Mu.Unlock()
	mData := models.Data{}

	mData.CreatedAt = time.Now()
	mData.UpdatedAt = time.Now()

	mData.UserId = id
	mData.TransactionId = d.TransactionId
	mData.State = false
	mData.Amount = a
	mData.Source = i
	mData.Status = 1

	srv.SaveTransaction(mData)
	//err = srv.CreateData(mData)
	//if err != nil {
	//	log.Println(err)
	//	return 0, err
	//}

	return b.Amount, nil
}

func (srv *Server) PostProcessing() {

	t := os.Getenv("N_MINUTES")

	m, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		log.Println(err)
		m = 10
	}

	for {
		time.Sleep(time.Duration(m) * time.Minute)

		var data []models.Data

		err := srv.Repo.Db.Table("data").Where("MOD (id, 2) = 1").Order("id  DESC").Limit("10").Find(&data).Error
		if err != nil {
			log.Println("Post Processing:", err)
			continue
		}

		for _, v := range data {

			if v.Status == 1 { // if its not canceled before or not transaction record with error
				srv.Mu.Lock()
				b := srv.UserBalances[v.UserId]

				if v.State { // win transaction
					if b.Amount-v.Amount < 0 {
						log.Println("Cancel not accepted. balance cant be negative.")
						srv.Mu.Unlock()
						continue
					}

					b.Amount = b.Amount - v.Amount
					b.Saved = true
					srv.UserBalances[v.UserId] = b
					srv.Balance = true
				} else { // lose transaction
					b.Amount = b.Amount + v.Amount
					b.Saved = true
					srv.UserBalances[v.UserId] = b
					srv.Balance = true
				}

				srv.Mu.Unlock()
				v.Status = 3

				err = srv.Repo.Db.Save(&v).Error
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}
