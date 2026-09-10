package handler

import (
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
  "gorm.io/datatypes"

  "github.com/slaghuis/YC-Pricing/pkg/models"
  "github.com/slaghuis/YC-Pricing/pkg/pricing"
)

// POST /skus
func (h *Handler) CreateSKUHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Code         string  `json:"code" binding:"required"`
            Name         string  `json:"name"`
            BasePrice    float64 `json:"base_price" binding:"required"`
            CurrencyCode string  `json:"currency_code" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
          c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "bad request: " + err.Error(),
          })
          return
        }

        sku := models.SKU{
            Code:         input.Code,
            Name:         input.Name,
            BasePrice:    input.BasePrice,
            CurrencyCode: input.CurrencyCode,
        }

        if err := h.db.Create(&sku).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
              "status": "error",
              "message": "sku not created: " + err.Error(),
            })
            return
        }

        c.JSON(http.StatusCreated, sku)
    }
}

func (h *Handler) UpdateSKUByCode() gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Code         string  `json:"code" binding:"required"`
            Name         string  `json:"name"`
            BasePrice    float64 `json:"base_price" binding:"required"`
            CurrencyCode string  `json:"currency_code" binding:"required"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
          c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "bad request: " + err.Error(),
          })
          return
        }

        var sku models.SKU
        h.db.Where("code = ?", input.Code).First(&sku)

        sku.Code          = input.Code
        sku.Name          = input.Name
        sku.BasePrice     = input.BasePrice
        sku.CurrencyCode  = input.CurrencyCode

        if err := h.db.Save(&sku).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
              "status": "error",
              "message": "sku not created: " + err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, sku)
    }
}


// POST /pricing-rules
func (h *Handler) AddPricingRuleHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Name       string                 `json:"name"`
            RuleType   string                 `json:"rule_type" binding:"required"`
            Priority   int                    `json:"priority"`
            Conditions map[string]interface{} `json:"conditions"`
            Actions    map[string]interface{} `json:"actions"`
            StartDate  *time.Time             `json:"start_date"`
            EndDate    *time.Time             `json:"end_date"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
          c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "bad request: " + err.Error(),
          })
          return
        }

        rule := models.PricingRule{
            Name:       input.Name,
            RuleType:   input.RuleType,
            Priority:   input.Priority,
            IsActive:   true,
            StartDate:  input.StartDate,
            EndDate:    input.EndDate,
            Conditions: datatypes.JSONMap(input.Conditions), // cast to JSON
            Actions:    datatypes.JSONMap(input.Actions),    // cast to JSON
        }

        if err := h.db.Create(&rule).Error; err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{
            "status": "error",
            "message": "pricing rule not added: " + err.Error(),
          })
          return
        }

        c.JSON(http.StatusCreated, rule)
    }
}

// POST /checkout/price
func (h *Handler) LookupPriceHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            SKUCode        string `json:"sku_code" binding:"required"`
            Quantity       int    `json:"quantity"`
            RegionCode     string `json:"region_code"`
            CustomerSegment string `json:"customer_segment"`
            TargetCurrency string `json:"target_currency"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
          c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "bad request: " + err.Error(),
          })
          return
        }

//        h.LogJson(input)

        var sku models.SKU
        if err := h.db.Where("code = ?", input.SKUCode).First(&sku).Error; err != nil {
          c.JSON(http.StatusNotFound, gin.H{
            "status": "error",
            "message": "sku not found: " + err.Error(),
          })
          return
        }

        pc := &pricing.PricingContext{
            SKU:            sku,
            Quantity:       input.Quantity,
            RegionCode:     input.RegionCode,
            CustomerSegment: input.CustomerSegment,
            TargetCurrency: input.TargetCurrency,
        }

        result, err := pricing.CalculatePrice(c, h.db, pc)
        if err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{
            "status": "error",
            "message": "error calculating price: " + err.Error(),
          })
          return
        }

        c.JSON(http.StatusOK, result)
    }
}






/*
PART OF CALCULATE PRICE

pc := &PricingContext{
    SKU: models.SKU{BasePrice: 100.0, CurrencyCode: "USD"},
    Quantity: 3,
    CustomerSegment: "VIP",
    RegionCode: "ZA",
    TargetCurrency: "USD",
}

result, err := CalculatePrice(context.Background(), db, pc)
if err != nil {
    fmt.Println("Error:", err)
}

fmt.Printf("Final Price: %.2f %s\nApplied Rules: %v\n",
    result.FinalPrice, result.CurrencyCode, result.AppliedRules)

*/
