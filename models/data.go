package models

import "fmt"

type (
	Data struct {
		State         string `json:"state"`
		Amount        string `json:"amount"`
		TransactionId string `json:"transactionId"`
	}

	Response struct {
		Error   bool        `json:"error"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func (d Data) ValidateData() error {

	if d.Amount == "" {
		return fmt.Errorf("amount cant be null")
	}

	if d.TransactionId == "" {
		return fmt.Errorf("transaction id cant be null")
	}

	return nil
}
