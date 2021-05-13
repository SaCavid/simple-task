package handlers

import (
	"fmt"
	"github.com/SaCavid/simple-task/models"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func (srv *Server) BulkUpdateBalances() {

	for {

		if !srv.Balance {
			time.Sleep(1 * time.Second)
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

			if err := srv.Repo.Db.Exec(fmt.Sprintf("UPDATE users AS u SET balance = data.a FROM (VALUES %s) AS data(user_id, a) WHERE u.user_id = data.user_id", strings.Join(value, ","))).Error; err != nil {
				log.Println(err)
				continue
			}

			balancesList = balancesList[count:]
		}

		// log.Println("Rows inserted:", len(balancesList))
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

	transactions := make([]models.Data, 0)

	err = srv.Repo.Db.Find(&transactions).Error
	if err != nil {
		return err
	}

	srv.Mu.Lock()
	for _, v := range transactions {
		srv.TransactionIds[v.TransactionId] = ""
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
