package pricing

import (
    "context"
    "fmt"
    "time"

    "github.com/slaghuis/YC-Pricing/pkg/models"
)

// PricingContext carries all inputs for calculation
type PricingContext struct {
    SKU            models.SKU
    TargetCurrency string
    RegionCode     string
    Quantity       int
    CustomerSegment string
}

// PricingResult carries outputs
type PricingResult struct {
    OriginalPrice float64   `json:"original_price"`
    FinalPrice    float64   `json:"final_price"`
    CurrencyCode  string    `json:"currency_code"`
    AppliedRules  []uint    `json:"applied_rules"`
    CalculatedAt  time.Time `json:"calculated_at"`
}


// --- Pipeline Interfaces ---

type RuleEvaluator interface {
    Evaluate(ctx context.Context, pc *PricingContext, pr *PricingResult) error
}

// --- Concrete Evaluators ---

// Currency conversion
type CurrencyConverter struct {
    Rates map[string]float64 // could be loaded from DB
}

func (c *CurrencyConverter) Evaluate(ctx context.Context, pc *PricingContext, pr *PricingResult) error {
    fmt.Printf("Currency Convertor calling\n")
    rate := c.Rates[pc.TargetCurrency]
    pr.FinalPrice = pc.SKU.BasePrice * rate
    pr.CurrencyCode = pc.TargetCurrency
    return nil
}

// Discounts
type DiscountEvaluator struct {
    Discounts []models.Discount
}

func (d *DiscountEvaluator) Evaluate(ctx context.Context, pc *PricingContext, pr *PricingResult) error {
    fmt.Printf("Discount Evaluator calling\n")
    for _, disc := range d.Discounts {
        if disc.AppliesTo == "sku" && disc.Value > 0 {
            if disc.DiscountType == "percentage" {
                pr.FinalPrice *= (1 - disc.Value/100)
            } else if disc.DiscountType == "fixed_amount" {
                pr.FinalPrice -= disc.Value
            }
            pr.AppliedRules = append(pr.AppliedRules, disc.RuleID)
        }
    }
    return nil
}

// Taxes
type TaxEvaluator struct {
    Taxes []models.TaxRule
}

func (t *TaxEvaluator) Evaluate(ctx context.Context, pc *PricingContext, pr *PricingResult) error {
    fmt.Printf("Tax Evaluator calling\n")
    for _, tax := range t.Taxes {
        if tax.RegionCode == pc.RegionCode {
            pr.FinalPrice *= (1 + tax.Rate/100)
            pr.AppliedRules = append(pr.AppliedRules, tax.RuleID)
        }
    }
    return nil
}

// --- Orchestrator ---

type PricingEngine struct {
    Evaluators []RuleEvaluator
}

func (pe *PricingEngine) Calculate(ctx context.Context, pc *PricingContext) (*PricingResult, error) {
    result := &PricingResult{
        OriginalPrice: pc.SKU.BasePrice,
        FinalPrice: pc.SKU.BasePrice,
        CurrencyCode: pc.SKU.CurrencyCode,
        AppliedRules: []uint{},
        CalculatedAt: time.Now(),
    }

    for _, evaluator := range pe.Evaluators {
        if err := evaluator.Evaluate(ctx, pc, result); err != nil {
            return nil, err
        }
    }
    return result, nil
}
