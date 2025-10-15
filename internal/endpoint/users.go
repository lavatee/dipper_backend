package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogInRequest struct {
	Ref string `json:"ref"`
}

func (e *Endpoint) LogIn(c *gin.Context) {
	var req LogInRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := e.GetUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user.ByRef = req.Ref
	user, err = e.services.Users.Login(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (e *Endpoint) TapsBatch(c *gin.Context) {
	user, err := e.GetUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := e.services.Users.TapsBatch(user.TelegramID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (e *Endpoint) GetRefUsers(c *gin.Context) {
	user, err := e.GetUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	users, err := e.services.Users.GetRefUsers(user.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
