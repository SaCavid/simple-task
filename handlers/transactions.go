package handlers

import (
	"fmt"
	"github.com/SaCavid/simple-task/models"
	"log"
	"strings"
	"time"
)

func (srv *Server) CheckTransactionId(id string) bool {
	srv.Mu.Lock()
	_, ok := srv.TransactionIds[id]
	srv.Mu.Unlock()
	return ok
}

func (srv *Server) SaveTransactionId(id string) {
	srv.Mu.Lock()
	srv.TransactionIds[id] = ""
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

		//	log.Println("Rows inserted:", len(values)/8)
	}
}
