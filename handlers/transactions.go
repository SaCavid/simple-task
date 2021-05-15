package handlers

import (
	"fmt"
	"github.com/SaCavid/simple-task/models"
	"log"
	"strings"
	"time"
)

// check transaction id to map
func (h *Server) CheckTransactionId(id string) bool {
	h.Mu.Lock()
	_, ok := h.TransactionIds[id]
	h.Mu.Unlock()
	return ok
}

// save transaction id to map
func (h *Server) SaveTransactionId(id string) {
	h.Mu.Lock()
	h.TransactionIds[id] = ""
	h.Mu.Unlock()
}

// save transaction record to temp map
func (h *Server) SaveTransaction(data models.Data) {
	h.Mu.Lock()
	h.Transactions = append(h.Transactions, data)
	h.Mu.Unlock()
}

// bulk insert transactions
func (h *Server) BulkInsertTransactions() {

	for {

		h.Mu.Lock()

		if len(h.Transactions) <= 0 {
			h.Mu.Unlock()
			time.Sleep(10 * time.Second)
			continue
		}
		count := len(h.Transactions)

		// maximum 500 rows per operation for safe database usage
		if count > 500 {
			count = 500
		}

		transactionsList := h.Transactions[:count]
		h.Mu.Unlock()

		tx := h.Repo.Db.Begin()
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

		// empty inserted transactions if not error
		h.Transactions = h.Transactions[count:]

		//	log.Println("Rows inserted:", len(values)/8)
	}
}
