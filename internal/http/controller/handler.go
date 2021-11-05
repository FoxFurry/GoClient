package controller

import (
	"context"
	"github.com/foxfurry/go_client/internal/domain/dto"
	"github.com/foxfurry/go_client/internal/domain/entity"
	"github.com/foxfurry/go_client/internal/service/client"
	"github.com/gin-gonic/gin"
	"log"
)

type IController interface {
	RegisterClientRouter(c *gin.Engine)
	StartService(ctx context.Context)
}

type ClientController struct {
	clientService client.IClient
}

func NewClientController(tableNum int, restaurantsData []entity.Restaurant) IController {
	return &ClientController{clientService: client.NewClientService(tableNum, restaurantsData)}
}

func (ctrl *ClientController) distribution(c *gin.Context) {
	var data dto.Distribution

	if err := c.ShouldBindJSON(&data); err != nil {
		log.Panic(err)
	}

	ctrl.clientService.Distribute(data.ClientID)

	c.Status(200)
}

func (ctrl *ClientController) StartService(ctx context.Context) {
	go ctrl.clientService.Start(ctx)
}