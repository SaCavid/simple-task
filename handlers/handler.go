package handlers

import (
	"fmt"
	"github.com/SaCavid/simple-task/models"
	"github.com/SaCavid/simple-task/service"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"sync"
	"time"
)

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

func (h *Server) Handler(c echo.Context) error {

	jd := new(models.JsonData)

	if err := c.Bind(&jd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "bad request"})
	}

	if h.CheckTransactionId(jd.TransactionId) {
		return echo.NewHTTPError(http.StatusNotAcceptable, &models.Response{Error: true, Message: fmt.Sprintf("this transaction id already used")})
	}

	// Save transaction id not to use again ever if its failed
	h.SaveTransactionId(jd.TransactionId)

	var s SourceType
	jd.Source = c.Request().Header.Get("Source-Type")
	i, err := s.IndexOf(jd.Source)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	id := c.Request().Header.Get("Authorization")

	if err := jd.ValidateData(); err != nil {
		// create clean data for transaction information to save
		data := models.Data{
			UserId:        id,
			State:         false,
			Source:        i,
			Status:        2, // error . saved for unique transaction id.
			Amount:        0,
			TransactionId: jd.TransactionId,
		}
		data.CreatedAt = time.Now()
		data.UpdatedAt = time.Now()

		h.SaveTransaction(data)
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error()})
	}

	a, err := strconv.ParseFloat(jd.Amount, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
	}

	data := models.Data{
		UserId:        id,
		State:         false,
		Source:        i,
		Status:        2, // error . saved for unique transaction id. not to allow repeat
		Amount:        a,
		TransactionId: jd.TransactionId,
	}
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	if id == "" {
		h.SaveTransaction(data)
		return echo.NewHTTPError(http.StatusForbidden, &models.Response{Error: true, Message: "not logged"})
	}

	if id != "only_for_testing_benchmark" {
		if !h.CheckUser(id) {
			h.SaveTransaction(data)
			return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "user didnt registered"})
		}
	}
	var balance float64
	switch jd.State {
	case "win":

		balance, err = h.UserWin(id, &data)
		if err != nil {
			data.Status = 2
			h.SaveTransaction(data)
			return echo.NewHTTPError(http.StatusInternalServerError, &models.Response{Error: true, Message: err.Error()})
		}

		break
	case "lose":

		balance, err = h.UserLost(id, &data)
		if err != nil {
			data.Status = 2
			h.SaveTransaction(data)
			return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: err.Error() + " " + jd.State + "-->" + jd.Amount + " Balance:" + fmt.Sprintf("%.2f", balance)})
		}

		break
	default:

		h.SaveTransaction(data)
		return echo.NewHTTPError(http.StatusBadRequest, &models.Response{Error: true, Message: "error with state"})
	}

	return c.JSON(201, &models.Response{Message: "transaction processed", Data: "Balance:" + fmt.Sprintf("%.2f", balance)})
}
