package pricing

import "encoding/json"

func toFloat64(v interface{}) float64 {
    switch n := v.(type) {
    case float64:
        return n
    case int:
        return float64(n)
    case int64:
        return float64(n)
    case json.Number:
        f, _ := n.Float64()
        return f
    default:
        return 0
    }
}

func toInt(v interface{}) int {
    switch n := v.(type) {
    case int:
        return n
    case int64:
        return int(n)
    case float64:
        return int(n)
    case json.Number:
        i, _ := n.Int64()
        return int(i)
    default:
        return 0
    }
}
