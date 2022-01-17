package schemas

import (
	"database/sql"
	repository2 "scoop-order/repository"
)

type Offer struct {
	OfferID    int32   `json:"offer_id"`
	Name       string  `json:"name"`
	BasePrice  float64 `json:"base_price"`
	FinalPrice float64 `json:"final_price"`
	Discount   float64 `json:"discount"`
}

type CheckPricingResponse struct {
	CurrencyCode    string `json:"currency_code"`
	Offers          []Offer
	TotalBasePrice  float64 `json:"total_base_price"`
	TotalDiscount   float64 `json:"total_discount"`
	TotalFinalPrice float64 `json:"total_final_price"`
}

type CheckoutTxResult struct {
	//Order       database.OrderPublished
	UserID           int32                            `json:"user_id"`
	PaymentGatewayID int32                            `json:"payment_gateway_id"`
	Currency         string                           `json:"currency"`
	TotalAmount      float64                          `json:"total_amount"`
	FinalAmount      float64                          `json:"final_amount"`
	OrderID          int32                            `json:"order_id"`
	OrderNumber      int                              `json:"order_number"`
	OrderStatus      int32                            `json:"order_status"`
	IsActive         bool                             `json:"is_active"`
	Orderline        []repository2.OrderLinePublished `json:"orderline"`
	IsDiscount       bool                             `json:"is_discount"`
}

type OrderResponse struct {
	Order   CheckoutTxResult       `json:"order"`
	Payment map[string]interface{} `json:"payment"`
}

type PaymentResponse struct {
	Created              sql.NullTime   `json:"created"`
	Modified             sql.NullTime   `json:"modified"`
	ID                   int32          `json:"id"`
	OrderID              int32          `json:"order_id"`
	UserID               int32          `json:"user_id"`
	PaymentgatewayID     int32          `json:"paymentgateway_id"`
	CurrencyCode         string         `json:"currency_code"`
	Amount               float64        `json:"amount"`
	PaymentStatus        int32          `json:"payment_status"`
	IsActive             bool           `json:"is_active"`
	IsTestPayment        bool           `json:"is_test_payment"`
	PaymentDatetime      sql.NullTime   `json:"payment_datetime"`
	FinancialArchiveDate sql.NullTime   `json:"financial_archive_date"`
	IsTrial              bool           `json:"is_trial"`
	MerchantParams       sql.NullString `json:"merchant_params"`
}

type PaymentGatewayResponse struct {
	ID             int32           `json:"id"`
	Name           string          `json:"name"`
	IsActive       bool            `json:"is_active"`
	BaseCurrencyID int32           `json:"base_currency_id"`
	MinimalAmount  sql.NullFloat64 `json:"minimal_amount"`
	IsRenewal      bool            `json:"is_renewal"`
	PaymentGroup   string          `json:"payment_group"`
}
