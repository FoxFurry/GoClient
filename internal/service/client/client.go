package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/foxfurry/go_client/internal/domain/dto"
	"github.com/foxfurry/go_client/internal/domain/entity"
	"github.com/foxfurry/go_client/internal/infrastructure/logger"
	"github.com/foxfurry/go_client/internal/infrastructure/random"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	tableGenerateProbability = 10
)

type table struct {
	maxWaitTime  int
	generateTime time.Time
}

type IClient interface {
	Start(ctx context.Context)
	Distribute(id int)
}

type clientService struct {
	tables      []table
	restaurants []entity.Restaurant
}

func NewClientService(tableNum int, restaurantsData []entity.Restaurant) IClient {
	return &clientService{
		tables: make([]table, tableNum, tableNum),
		restaurants: restaurantsData,
	}
}

func (s *clientService) Start(ctx context.Context) {
	ticker := time.Tick(time.Second)

	for {
		select {
		case <-ticker:
			for idx, val := range s.tables {
				if val.maxWaitTime == 0 { // If current table is empty
					if random.CoinFlip(tableGenerateProbability) {
						order := s.generateOrder(idx)
						s.sendOrder(order)
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *clientService) Distribute(id int) {
	logger.LogTable(id, "Received order distribution")

	distributionTime := time.Now()

	var (
		waitTime     = float64(s.tables[id].maxWaitTime)
		deliveryTime = float64(distributionTime.Second() - s.tables[id].generateTime.Second())
	)

	if deliveryTime < waitTime {
		logger.LogTable(id, "5 START DELIVERY")
	} else if deliveryTime < waitTime * 1.1 {
		logger.LogTable(id, "4 START DELIVERY")
	} else if deliveryTime < waitTime * 1.2 {
		logger.LogTable(id, "3 START DELIVERY")
	} else if deliveryTime < waitTime * 1.3 {
		logger.LogTable(id, "2 START DELIVERY")
	} else if deliveryTime < waitTime * 1.4 {
		logger.LogTable(id, "1 START DELIVERY")
	} else {
		logger.LogTable(id, "0 STAR DELIVERY")
	}

	s.tables[id].maxWaitTime = 0
}

func (s *clientService) generateOrder(id int) *dto.Order {
	subOrderCount := rand.Intn(viper.GetInt("max_order_size")-1)+1

	newOrder := new(dto.Order)

	newOrder.ClientID = id
	newOrder.Orders, s.tables[id].maxWaitTime = s.generateSubordersn(subOrderCount)

	s.tables[id].generateTime = time.Now()

	logger.LogTableF(id, "Generated order with %d suborders and %d max wait time", subOrderCount,s.tables[id].maxWaitTime)

	return newOrder
}

func (s *clientService) generateSubordersn(orderCount int) ([]entity.Order, int) {
	var (
		maxTime     int
		resultOrder []entity.Order
	)

	for idx := 0; idx < orderCount; idx++ {
		restaurantID := 1
		if len(s.restaurants) > 1 {
			restaurantID = rand.Intn(len(s.restaurants)-1)+1
		}
		items, currentMaxTime := s.generateItemsn(restaurantID, viper.GetInt("max_order_items"))

		if currentMaxTime > maxTime {
			maxTime = currentMaxTime
		}

		resultOrder = append(resultOrder, entity.Order{
			RestaurantID: restaurantID,
			Items:        items,
			Priority:     1,
			MaxWait:      currentMaxTime,
			CreateTime:   time.Now(),
		})

	}

	return resultOrder, maxTime
}

func (s *clientService) generateItemsn(restaurantID, itemCount int) ([]int, int) {
	result := make([]int, 0, itemCount)
	maxWait := 0

	for idx := 0; idx < itemCount; idx++ {
		itemID := rand.Intn(len(s.restaurants[restaurantID-1].Menus))

		result = append(result, itemID)

		if s.restaurants[restaurantID-1].Menus[itemID].PreparationTime > maxWait {
			maxWait = s.restaurants[restaurantID-1].Menus[itemID].PreparationTime
		}
	}

	return result, int(float64(maxWait) * viper.GetFloat64("max_wait_multiplier"))
}

func (s *clientService) sendOrder(order *dto.Order) {
	jsonBody, err := json.Marshal(order)
	if err != nil {
		log.Panic(err)
	}
	contentType := "application/json"

	http.Post("http://" + viper.GetString("delivery_host")+"/order", contentType, bytes.NewReader(jsonBody))
}
