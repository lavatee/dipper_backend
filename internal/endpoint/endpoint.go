package endpoint

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lavatee/dipper_backend/internal/service"
)

type Endpoint struct {
	services *service.Service
	botToken string
}

func NewEndpoint(services *service.Service, botToken string) *Endpoint {
	return &Endpoint{
		services: services,
		botToken: botToken,
	}
}

func (e *Endpoint) InitRoutes() *gin.Engine {
	router := gin.New()
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
	{
		api := router.Group("/api", e.Middleware)
		api.POST("/login", e.LogIn)
		api.POST("/taps-batch", e.TapsBatch)
		api.GET("/ref-users", e.GetRefUsers)
		api.GET("/upgrades", e.GetAllUpgrades)
		api.GET("/user-upgrades", e.GetUserUpgrades)
		api.POST("/buy-upgrade/by-coins", e.BuyUpgradeByCoins)
	}
	return router
}
