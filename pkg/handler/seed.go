package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/slaghuis/YC-Pricing/pkg/models"
)

func (h *Handler) SeedData(c *gin.Context) {
    currencies := []models.Currency{
        {Code: "USD", Name: "US Dollar", ExchangeRateToBase: 0.054},
        {Code: "EUR", Name: "Euro", ExchangeRateToBase: 0.052},
        {Code: "ZAR", Name: "South African Rand", ExchangeRateToBase: 1.0},
    }
    for _, c := range currencies {
        h.db.FirstOrCreate(&c, models.Currency{Code: c.Code})
    }

    regions := []models.Region{
        {Code: "ZA", Name: "South Africa"},
        {Code: "US-CA", Name: "California, USA"},
    }
    for _, r := range regions {
        h.db.FirstOrCreate(&r, models.Region{Code: r.Code})
    }

    sku := models.SKU{
        Code: "SKU-001", Name: "Test Product",
        BasePrice: 100.0, CurrencyCode: "ZAR",
    }
    h.db.FirstOrCreate(&sku, models.SKU{Code: sku.Code})

    c.JSON(http.StatusOK, gin.H{
      "status": "succsess",
      "message": "database seeded successfully",
    })
}
