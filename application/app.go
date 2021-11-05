package application

import (
	"context"
	"encoding/json"
	"github.com/foxfurry/go_client/internal/domain/dto"
	"github.com/foxfurry/go_client/internal/http/controller"
	"github.com/foxfurry/go_client/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type IApp interface {
	Start()
	Shutdown()
}

type clientApp struct {
	server *http.Server
}

func CreateApp(ctx context.Context) IApp {
	appHandler := gin.Default()

	menu := initialize(ctx)

	ctrl := controller.NewClientController(viper.GetInt("table_num"), menu.RestaurantsData)
	ctrl.RegisterClientRouter(appHandler)

	app := clientApp{
		server: &http.Server{
			Addr:    viper.GetString("client_host"),
			Handler: appHandler,
		},
	}

	ctrl.StartService(ctx)

	return &app
}

func initialize(ctx context.Context) dto.Menu {
	deliveryHost := viper.GetString("delivery_host")

	logger.LogMessageF("Trying to reach delivery server on: %v", deliveryHost)

	waitConnection(deliveryHost, ctx)

	req, err := http.Get("http://" + deliveryHost + "/menu")
	if err != nil {
		logger.LogPanic(err.Error())
	}

	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.LogPanic(err.Error())
	}

	bodyMenu := dto.Menu{}

	if err = json.Unmarshal(bodyData, &bodyMenu); err != nil {
		logger.LogPanic(err.Error())
	}

	logger.LogMessageF("Successfully received data! Restaurants available %d", bodyMenu.Restaurants)

	return bodyMenu
}

func waitConnection(host string, ctx context.Context) {
	tryCount := 1
	dialTick := time.Tick(time.Second * time.Duration(viper.GetInt("dial_timeout")))

	for {
		select {
		case <-dialTick:
			conn, err := net.Dial("tcp", host)

			if err == nil {
				conn.Close()
				return
			}

			logger.LogMessageF("Could not reach delivery server. Retrying %d", tryCount)
			tryCount++

		case <-ctx.Done():
			logger.LogMessage("Stopping dialing")
			return
		}
	}
}

func (d *clientApp) Start() {
	logger.LogMessage("Starting client server!")

	if err := d.server.ListenAndServe(); err != http.ErrServerClosed {
		logger.LogPanicF("Unexpected error while running server: %v", err)
	}

}

func (d *clientApp) Shutdown() {
	if err := d.server.Shutdown(context.Background()); err != nil {
		logger.LogPanicF("Unexpected error while closing server: %v", err)
	}
	logger.LogMessage("Server terminated successfully")
}
