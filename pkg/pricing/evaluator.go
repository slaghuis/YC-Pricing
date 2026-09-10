package pricing

import (
    "context"
    "time"

    "gorm.io/gorm"
    "github.com/slaghuis/YC-Pricing/pkg/models"
)

// DSL structs
type Condition struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}

type Action struct {
    Type  string      `json:"type"`  // discount, tax, promotion
    Mode  string      `json:"mode"`  // percentage, fixed_amount
    Value interface{} `json:"value"`
}

type RuleDSL struct {
    Conditions []Condition `json:"conditions"`
    Actions    []Action    `json:"actions"`
}

// --- Rule Loader ---

func LoadActiveRules(db *gorm.DB) ([]models.PricingRule, error) {
    var rules []models.PricingRule
    err := db.Where("is_active = ?", true).
        Where("start_date IS NULL OR start_date <= ?", time.Now()).
        Where("end_date IS NULL OR end_date >= ?", time.Now()).
        Order("priority ASC").
        Find(&rules).Error
    return rules, err
}

// --- Rule Evaluator ---

func EvaluateRule(rule models.PricingRule, pc *PricingContext, pr *PricingResult) error {
  var dsl RuleDSL

  // Convert Conditions map into []Condition
  for field, raw := range rule.Conditions {
      condMap, ok := raw.(map[string]interface{})
      if !ok {
          continue
      }
      dsl.Conditions = append(dsl.Conditions, Condition{
          Field:    field,
          Operator: condMap["operator"].(string),
          Value:    condMap["value"],
      })
  }

  // Convert Actions map into []Action
  for field, raw := range rule.Actions {
      actMap, ok := raw.(map[string]interface{})
      if !ok {
          continue
      }
      dsl.Actions = append(dsl.Actions, Action{
          Type:  field,
          Mode:  actMap["type"].(string),
          Value: actMap["value"],
      })
  }

  // Check conditions
  for _, cond := range dsl.Conditions {
      if !evaluateCondition(cond, pc) {
          return nil // skip rule
      }
  }

  // Apply actions
  for _, act := range dsl.Actions {
      applyAction(act, pr)
      pr.AppliedRules = append(pr.AppliedRules, rule.ID)
  }

  return nil
}

func evaluateCondition(cond Condition, pc *PricingContext) bool {
    switch cond.Field {
    case "quantity":
        val := toInt(cond.Value)
        switch cond.Operator {
        case ">=": return pc.Quantity >= val
        case "<=": return pc.Quantity <= val
        case "=":  return pc.Quantity == val
        }
    case "customer_segment":
        val, ok := cond.Value.(string)
        if !ok {
            return false
        }
        return cond.Operator == "=" && pc.CustomerSegment == val
    case "region_code":
        val, ok := cond.Value.(string)
        if !ok {
            return false
        }
        return cond.Operator == "=" && pc.RegionCode == val
    case "sku_code":
        val, ok := cond.Value.(string)
        if !ok {
            return false
        }
        return cond.Operator == "=" && pc.SKU.Code == val
    }
    return false
}

/*
func evaluateCondition(cond Condition, pc *PricingContext) bool {
    switch cond.Field {
    case "quantity":
        val := toInt(cond.Value)

        switch cond.Operator {
        case ">=": return pc.Quantity >= val
        case "<=": return pc.Quantity <= val
        case "=":  return pc.Quantity == val
        }
    case "customer_segment":
        val, ok := cond.Value.(string)
        if !ok {
            return false
        }
        return cond.Operator == "=" && pc.CustomerSegment == val
    case "region_code":
        val, ok := cond.Value.(string)
        if !ok {
            return false
        }
        return cond.Operator == "=" && pc.RegionCode == val
    }
    return false
}
*/

func applyAction(act Action, pr *PricingResult) {
    switch act.Type {
    case "discount":
        val := toFloat64(act.Value)
        if act.Mode == "percentage" {
            pr.FinalPrice *= (1 - val/100)
        } else if act.Mode == "fixed_amount" {
            pr.FinalPrice -= val
        }
    case "tax":
        val := toFloat64(act.Value)
        if act.Mode == "percentage" {
            pr.FinalPrice *= (1 + val/100)
        }
    }
}

// --- Orchestrator ---

func CalculatePrice(ctx context.Context, db *gorm.DB, pc *PricingContext) (*PricingResult, error) {
    rules, err := LoadActiveRules(db)
    if err != nil {
        return nil, err
    }

    result := &PricingResult{
        OriginalPrice: pc.SKU.BasePrice,
        FinalPrice:   pc.SKU.BasePrice,
        CurrencyCode: pc.SKU.CurrencyCode,
        AppliedRules: []uint{},
        CalculatedAt: time.Now(),
    }

    for _, rule := range rules {
        if err := EvaluateRule(rule, pc, result); err != nil {
            return nil, err
        }
    }

    return result, nil
}
