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

	Balance struct {
		Amount float64
		Saved  bool // if true this balance didnt saved to the database
	}

	Data struct {
		gorm.Model
		UserId        string
		State         bool    // transaction win - lose state
		Status        uint8   // operation status processed -1 / error denied -2 / canceled -3 / cancel denied -4 and etc
		Source        int     // source of operation
		Amount        float64 // amount of operation
		TransactionId string  `gorm:"index"` // unique transaction id
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

	if d.State == "" {
		return fmt.Errorf("null state")
	}

	if d.State != "win" {
		if d.State != "lose" {
			return fmt.Errorf("wrong state")
		}
	}

	if d.TransactionId == "" {
		return fmt.Errorf("transaction id cant be null")
	}

	if d.Amount == "" {
		return fmt.Errorf("amount cant be null")
	}

	return nil
}
