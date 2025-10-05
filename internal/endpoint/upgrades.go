package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) GetAllUpgrades(c *gin.Context) {
	upgrades, err := e.services.Upgrades.GetAllUpgrades()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"upgrades": upgrades})
}

func (e *Endpoint) GetUserUpgrades(c *gin.Context) {
	user, err := e.GetUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	upgrades, err := e.services.Upgrades.GetUserUpgrades(user.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"upgrades": upgrades})
}

type BuyUpgradeRequest struct {
	UpgradeID int `json:"upgrade_id"`
}

func (e *Endpoint) BuyUpgradeByCoins(c *gin.Context) {
	var req BuyUpgradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := e.GetUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := e.services.Upgrades.BuyUpgradeByCoins(user.TelegramID, req.UpgradeID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
