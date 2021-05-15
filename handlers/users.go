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

func (h *Server) Register(c echo.Context) error {
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

	h.Mu.Lock()
	if len(h.UserBalances) >= int(m) {
		h.Mu.Unlock()
		return echo.NewHTTPError(http.StatusNotAcceptable, &models.Response{Error: true, Message: "maximum user count reached"})
	}
	h.Mu.Unlock()

	if h.CheckUser(user.UserId) {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "user already registered"})
	}

	err = h.Repo.Db.Create(user).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
	}

	h.AddUser(user.UserId)

	return c.JSON(http.StatusOK, &models.Response{Message: "user registered"})
}

func (h *Server) FetchUsersForTesting(c echo.Context) error {

	type User struct {
		UserId string
	}

	users := make([]User, 0)

	err := h.Repo.Db.Table("users").Find(&users).Error
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &models.Response{Message: "users", Data: users})
}

func (h *Server) BulkUpdateBalances() {

	for {

		if !h.Balance {
			time.Sleep(1 * time.Second)
			continue
		}

		type balance struct {
			UserId string
			Amount float64
		}

		balancesList := make([]balance, 0)

		h.Mu.Lock()
		for k, v := range h.UserBalances {
			if v.Saved {
				s := balance{
					UserId: k,
					Amount: v.Amount,
				}
				balancesList = append(balancesList, s)
			}
			v.Saved = false
		}
		h.Balance = false
		h.Mu.Unlock()

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

			if err := h.Repo.Db.Exec(fmt.Sprintf("UPDATE users AS u SET balance = data.a FROM (VALUES %s) AS data(user_id, a) WHERE u.user_id = data.user_id", strings.Join(value, ","))).Error; err != nil {
				log.Println(err)
				continue
			}

			balancesList = balancesList[count:]
		}

		// log.Println("Rows inserted:", len(balancesList))
	}
}

func (h *Server) CheckUser(id string) bool {
	h.Mu.Lock()
	_, ok := h.UserBalances[id]
	h.Mu.Unlock()
	return ok
}

func (h *Server) AddUser(id string) {
	h.Mu.Lock()
	b := models.Balance{}
	b.Amount = 0
	b.Saved = false
	h.UserBalances[id] = b
	h.Mu.Unlock()
}

func (h *Server) FetchUsers() error {

	users := make([]models.User, 0)

	err := h.Repo.Db.Find(&users).Error
	if err != nil {
		return err
	}

	h.Mu.Lock()
	for _, v := range users {
		b := models.Balance{
			Amount: v.Balance,
			Saved:  false,
		}
		h.UserBalances[v.UserId] = b
	}
	h.Mu.Unlock()

	transactions := make([]models.Data, 0)

	err = h.Repo.Db.Find(&transactions).Error
	if err != nil {
		return err
	}

	h.Mu.Lock()
	for _, v := range transactions {
		h.TransactionIds[v.TransactionId] = ""
	}
	h.Mu.Unlock()

	return nil
}

func (h *Server) UserWin(id string, d *models.Data) error {

	h.Mu.Lock()
	b := h.UserBalances[id]
	b.Amount = b.Amount + d.Amount
	b.Saved = true   // not saved
	h.Balance = true // not saved balance in map
	h.UserBalances[id] = b
	h.Mu.Unlock()

	d.State = true
	d.Status = 1
	h.SaveTransaction(*d)
	//err = srv.CreateData(mData)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	return nil
}

func (h *Server) UserLost(id string, d *models.Data) (float64, error) {

	h.Mu.Lock()
	b := h.UserBalances[id]
	if (b.Amount - d.Amount) < 0 {
		h.Mu.Unlock()
		return b.Amount, fmt.Errorf("not enough user balance")
	}

	b.Amount = b.Amount - d.Amount
	b.Saved = true   // user balance not saved
	h.Balance = true // not saved balance in map
	h.UserBalances[id] = b
	h.Mu.Unlock()
	d.Status = 1

	h.SaveTransaction(*d)
	//err = srv.CreateData(mData)
	//if err != nil {
	//	log.Println(err)
	//	return 0, err
	//}

	return b.Amount, nil
}
