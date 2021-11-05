package controller

import "github.com/gin-gonic/gin"

func (d *ClientController) RegisterClientRouter(c *gin.Engine) {
	c.POST("/distribution", d.distribution)
}
