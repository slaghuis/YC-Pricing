package router

import (
  "net/http"
  "github.com/gin-gonic/gin"
	"github.com/slaghuis/YC-Pricing/pkg/handler"
  "github.com/slaghuis/YC-Pricing/pkg/config"
  "github.com/slaghuis/YC-Pricing/pkg/middleware"
)

func NewRouter(h *handler.Handler, cfg config.JwtConfigurations) *gin.Engine {
    r := gin.Default()

    // Health
    r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

    v1 := r.Group("/api/v1")

    // Example protected route for your own testing
    auth := v1.Group("/")
    auth.Use(middleware.RequireAccessToken(cfg))

    auth.POST("/skus", h.CreateSKUHandler())
    auth.PUT("/skus", h.UpdateSKUByCode())
    auth.POST("/pricing-rules", h.AddPricingRuleHandler())
    auth.POST("/calculate", h.LookupPriceHandler())

    auth.POST("/seed", h.SeedData).Use(middleware.RBACMiddleware(middleware.RoleSysAdmin))

//    To Add
//        list rules, deactivate rule, audit price calculations

    return r
}
