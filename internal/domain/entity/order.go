package entity

import "time"

type Order struct {
	RestaurantID int `json:"restaurant_id"`
	Items []int `json:"items"`
	Priority int `json:"priority"`
	MaxWait int `json:"max_wait"`
	CreateTime time.Time `json:"created_time"`
}
