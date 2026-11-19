package rest

import (
	"ca/internal/transport/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *api) register(c *gin.Context) {
	var user dto.UserCreateRequest

	if err := c.ShouldBind(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err := api.userService.Register(c, &user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
