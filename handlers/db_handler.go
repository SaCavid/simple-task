package handlers

import (
	"../models"
	"log"
)

func (srv *Server) CreateData(data *models.Data) error {

	if err := srv.Repo.Db.Create(data).Error; err != nil {
		log.Println(err)
		return err
	}

	return nil
}
