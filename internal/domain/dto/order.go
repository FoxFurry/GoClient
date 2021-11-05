package dto

import "github.com/foxfurry/go_client/internal/domain/entity"

type Order struct {
	ClientID int `json:"client_id"`
	Orders []entity.Order `json:"orders"`
}
