package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type (
	User struct {
		gorm.Model
		UserId  string `gorm:"index"`
		Balance float64
	}

	Data struct {
		gorm.Model
		UserId        string
		State         bool
		Source        int
		Amount        float64
		TransactionId string `gorm:"index"`
	}

	JsonData struct {
		State         string `json:"state"`
		Source        string `json:"source"`
		Amount        string `json:"amount"`
		TransactionId string `json:"transactionId"`
	}

	Response struct {
		Error   bool        `json:"error"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func (d JsonData) ValidateData() error {

	if d.Amount == "" {
		return fmt.Errorf("amount cant be null")
	}

	if d.TransactionId == "" {
		return fmt.Errorf("transaction id cant be null")
	}

	return nil
}
