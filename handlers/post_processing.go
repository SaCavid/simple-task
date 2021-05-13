package handlers

import (
	"github.com/SaCavid/simple-task/models"
	"log"
	"os"
	"strconv"
	"time"
)

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
