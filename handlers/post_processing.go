package handlers

import (
	"github.com/SaCavid/simple-task/models"
	"log"
	"os"
	"strconv"
	"time"
)

// Post processing task:
// Every N minutes 10 latest odd records must be canceled and balance should be corrected by the application.
// Cancelled records shouldn't be processed twice.
func (h *Server) PostProcessing() {

	t := os.Getenv("N_MINUTES")

	m, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		log.Println(err)
		m = 10
	}

	for {
		time.Sleep(time.Duration(m) * time.Minute)

		var data []models.Data

		// get latest 10 odd records
		err := h.Repo.Db.Table("data").Where("MOD (id, 2) = 1").Order("id  DESC").Limit("10").Find(&data).Error
		if err != nil {
			log.Println("Post Processing:", err)
			continue
		}

		for _, v := range data {

			// check if its not canceled before or not transaction record with error
			if v.Status == 1 {
				h.Mu.Lock()
				b := h.UserBalances[v.UserId]

				if v.State { // win transaction
					if b.Amount-v.Amount < 0 {
						log.Println("Cancel not accepted. balance cant be negative.")
						h.Mu.Unlock()
						continue
					}

					b.Amount = b.Amount - v.Amount
					b.Saved = true
					h.UserBalances[v.UserId] = b
					h.Balance = true
				} else { // lose transaction
					b.Amount = b.Amount + v.Amount
					b.Saved = true
					h.UserBalances[v.UserId] = b
					h.Balance = true
				}

				h.Mu.Unlock()
				v.Status = 3 // transaction status canceled - 3

				err = h.Repo.Db.Save(&v).Error
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}
