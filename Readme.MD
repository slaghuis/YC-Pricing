

How It Works
PricingContext: carries SKU, target currency, region, quantity, etc.
Pipeline of evaluators: each evaluator applies its logic (currency conversion, discounts, promotions, taxes).
PricingEngine: orchestrates evaluators in order of priority.
PricingResult: final price + audit trail of applied rules.


Design Notes
JSONB for conditions/actions: gives flexibility to encode complex rules without schema churn. You can later add a rules engine layer to interpret them.
priority field: ensures deterministic application order (e.g., tax after discount).
audit table: lets you debug pricing decisions and support customer service queries.
currency normalization: keep one base currency (say USD) and store exchange rates relative to it. Convert at calculation time.

Example Flow
Fetch product base price.
Convert to target currency.
Apply active discounts/promotions (ordered by priority).
Apply tax rules based on region.
Persist calculation in price_calcula

A rule condition DSL (Domain-Specific Language) is essentially a structured way of describing when a pricing rule should apply. Instead of hardcoding conditions in Go, you store them in the database (often as JSON) and interpret them at runtime. This makes your pricing service flexible: business users or admins can define new rules without redeploying code.

Eample DSL Imagine you want a rule: “10% off if quantity ≥ 3 and customer segment = VIP.”

{
  "conditions": {
    "quantity": { "operator": ">=", "value": 3 },
    "customer_segment": { "operator": "=", "value": "VIP" }
  },
  "actions": {
    "discount": { "type": "percentage", "value": 10 }
  }
}


or

Another rule: “Apply 15% VAT if region = ZA.”

{
  "conditions": {
    "region_code": { "operator": "=", "value": "ZA" }
  },
  "actions": {
    "tax": { "type": "percentage", "value": 15 }
  }
}

or add 5 percent discount for a product SKU over a given period.
{
  "name": "Spring Promo SKU-001",
  "rule_type": "promotion",
  "priority": 20,
  "conditions": {
    "sku_code": { "operator": "=", "value": "SKU-001" }
  },
  "actions": {
    "discount": { "type": "percentage", "value": 5 }
  },
  "start_date": "2026-04-15T00:00:00Z",
  "end_date": "2026-05-15T23:59:59Z"
}



Can be extended in future to combine multiple conditions (e.g. sku_code AND customer_segment) so you can run promotions like “5% off SKU-001 for VIP customers only”?  
Can be extended to set the discount on the product code, covering all sku's relating to the product
