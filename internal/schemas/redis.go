package schemas

type DiscountTemp struct {
	DiscountID int32 `json:"discount_id"`
	DiscountName string `json:"discount_name"`
	CurrencyCode string `json:"currency_code"`
	DiscountCode string `json:"discount_code"`
	DiscountType int32 `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	RawPrice float64 `json:"raw_price"`
	FinalPrice float64 `json:"final_price"`
}
