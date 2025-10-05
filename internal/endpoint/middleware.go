package endpoint

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lavatee/dipper_backend/internal/model"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func (e *Endpoint) Middleware(c *gin.Context) {
	initData := c.GetHeader("Authorization")
	if err := initdata.Validate(initData, e.botToken, 24*time.Hour); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid init data"})
		return
	}
	data, err := initdata.Parse(initData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid init data"})
		return
	}
	c.Set("userID", data.User.ID)
	c.Set("userUsername", data.User.Username)
}

func (e *Endpoint) GetUser(c *gin.Context) (model.User, error) {
	userID, ok := c.Get("userID")
	if !ok {
		return model.User{}, fmt.Errorf("userID not found")
	}
	userUsername, ok := c.Get("userUsername")
	if !ok {
		return model.User{}, fmt.Errorf("userUsername not found")
	}
	return model.User{
		TelegramID:       fmt.Sprint(userID),
		TelegramUsername: fmt.Sprint(userUsername),
	}, nil
}
