package dto

import "github.com/foxfurry/go_client/internal/domain/entity"

type Menu struct{
	Restaurants int `json:"restaurants"`
	RestaurantsData []entity.Restaurant `json:"restaurants_data"`
}
