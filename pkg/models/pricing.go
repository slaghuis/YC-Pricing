package models

import (
    "time"
    "gorm.io/gorm"
    "gorm.io/datatypes"
)

// --- Core Entities ---

type SKU struct {
    ID           uint           `gorm:"primaryKey"`
    Code         string         `gorm:"uniqueIndex;not null"` // e.g. "SKU-12345"
    Name         string
    BasePrice    float64        `gorm:"not null"`
    CurrencyCode string         `gorm:"not null"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Currency struct {
    Code                string  `gorm:"primaryKey"` // e.g. "USD"
    Name                string
    ExchangeRateToBase  float64 `gorm:"not null"`
    UpdatedAt           time.Time
}

// --- Pricing Rules ---
type PricingRule struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string
    RuleType  string         `gorm:"not null"`
    Priority  int            `gorm:"default:100"`
    IsActive  bool           `gorm:"default:true"`
    StartDate *time.Time
    EndDate   *time.Time
    Conditions datatypes.JSONMap `gorm:"type:jsonb"`
    Actions    datatypes.JSONMap `gorm:"type:jsonb"`
    CreatedAt time.Time
    UpdatedAt time.Time
}


// --- Discounts & Promotions ---

type Discount struct {
    ID          uint   `gorm:"primaryKey"`
    RuleID      uint   `gorm:"not null"`
    DiscountType string `gorm:"not null"` // percentage, fixed_amount
    Value       float64
    AppliesTo   string // sku, category, cart_total
    Rule        PricingRule `gorm:"foreignKey:RuleID"`
}

type Promotion struct {
    ID            uint   `gorm:"primaryKey"`
    RuleID        uint   `gorm:"not null"`
    PromotionType string `gorm:"not null"` // buy_x_get_y, free_shipping, bundle_price
    Details       map[string]interface{} `gorm:"type:jsonb"`
    Rule          PricingRule `gorm:"foreignKey:RuleID"`
}

// --- Taxes ---

type Region struct {
    Code string `gorm:"primaryKey"` // e.g. "US-CA"
    Name string
}

type TaxRule struct {
    ID        uint   `gorm:"primaryKey"`
    RuleID    uint   `gorm:"not null"`
    RegionCode string `gorm:"not null"`
    TaxType   string  // VAT, sales_tax, GST
    Rate      float64
    Rule      PricingRule `gorm:"foreignKey:RuleID"`
    Region    Region      `gorm:"foreignKey:RegionCode"`
}

// --- Audit & History ---

type PriceCalculation struct {
    ID           uint   `gorm:"primaryKey"`
    SKUCode      string `gorm:"not null"`
    CurrencyCode string `gorm:"not null"`
    FinalPrice   float64
    AppliedRules map[string]interface{} `gorm:"type:jsonb"` // array of rule_ids
    CalculatedAt time.Time
}
