package handlers

import (
	"github.com/SaCavid/simple-task/models"
	"log"
)

func (h *Server) CreateData(data *models.Data) error {

	if err := h.Repo.Db.Create(data).Error; err != nil {
		log.Println(err)
		return err
	}

	return nil
}
