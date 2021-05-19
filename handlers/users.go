package handlers

import (
	"fmt"
	"github.com/SaCavid/simple-task/models"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
	"time"
)

// simple registration handler
// Example registration json object:
// { "UserId":"NewId"}
func (h *Server) Register(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	// empty user id not allowed
	if user.UserId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "user id can't be null"})
	}

	// check if user already registered or not
	if h.CheckUser(user.UserId) {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "user already registered"})
	}

	// add user to database
	err := h.Repo.Db.Create(user).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
	}

	// add user to map for further use
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

// update user balances if get true in Server.Balance
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

// check if user already registered and exists or not
// can be improved adding database check and expire time
func (h *Server) CheckUser(id string) bool {
	h.Mu.Lock()
	_, ok := h.UserBalances[id]
	h.Mu.Unlock()
	return ok
}

// add new user to map Server.UserBalances
func (h *Server) AddUser(id string) {
	h.Mu.Lock()
	b := models.Balance{}
	b.Amount = 0
	b.Saved = false
	h.UserBalances[id] = b
	h.Mu.Unlock()
}

// get all data in server startup
func (h *Server) FetchData() error {

	// get all users information for further use
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

	// get all transactions information. not to allow repeating transaction id
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

// win state transaction
func (h *Server) UserWin(id string, d *models.Data) (float64, error) {

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
	return b.Amount, nil
}

// lose state transaction
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
