package schemas

import "database/sql"

type GetOrderByIDRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type GetOrderByONRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type GetOrderListRequest struct {
	Limit  int `form:"limit,default=10"`
	Offset int `form:"offset,default=1"`
}

type GetValidPriceRequest struct {
	UserID           int32  `form:"user_id" binding:"required"`
	OfferID          string `form:"offer_id" binding:"required"`
	PaymentGatewayID int32  `form:"payment_gateway_id" binding:"required"`
	PlatformID       int32  `form:"platform_id,default=4"`
	CurrencyCode     string `form:"currency_code,default=IDR"`
	DiscountCode     string `form:"discount_code"`
	CountryCode      string `form:"country_code,default=ID"`
}

type GeoInfo struct {
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	CountryCode string  `json:"country_code" binding:"required"`
	CountryName string  `json:"country_name" binding:"required"`
	ClientIP    string  `json:"client_ip"`
}

type CheckPricingRequest struct {
	UserID           int32   `json:"user_id"`
	OfferID          []int32 `json:"offer_id"`
	PaymentGatewayID int32   `json:"payment_gateway_id"`
	PlatformID       int32   `json:"platform_id"`
	DiscountCode     string  `json:"discount_code"`
	CurrencyCode     string  `json:"currency_code"`
	IsRenewal        bool    `json:"is_renewal,default=false"`
	GeoInfo          GeoInfo
}

type CheckoutRequest struct {
	OfferID          []int32 `json:"offer_id" binding:"required"`
	DiscountCode     string  `json:"discount_code"`
	CurrencyCode     string  `json:"currency_code" binding:"required"`
	PlatformID       int32   `json:"platform_id" binding:"required"`
	PaymentGatewayID int32   `json:"payment_gateway_id" binding:"required"`
	GeoInfo          GeoInfo `json:"geo_info" binding:"required"`
	Signature        string  `json:"signature" binding:"required"`
}

type CheckoutTxParams struct {
	UserID           int32   `json:"user_id"`
	OfferID          []int32 `json:"offer_id"`
	DiscountCode     string  `json:"discount_code"`
	CurrencyCode     string  `json:"currency_code"`
	PlatformID       int32   `json:"platform_id"`
	PaymentGatewayID int32   `json:"payment_gateway_id"`
	ClientID         int32   `json:"client_id"`
	GeoInfo          GeoInfo
}

type PaymentTxParams struct {
	OrderID       int32 `json:"order_id"`
	OrderStatus   int32 `json:"order_status"`
	PaymentStatus int32 `json:"payment_status"`
}

type CompleteRequest struct {
	OrderID int64 `json:"order_id"`
}

type CompleteOrderRequest struct {
	ID                     int32        `json:"id"`
	CreatedAt              sql.NullTime `json:"created_at"`
	ModifiedAt             sql.NullTime `json:"modified_at"`
	StatusCode             string       `json:"status_code"`
	StatusMessage          string       `json:"status_message"`
	SignatureKey           string       `json:"signature_key"`
	Bank                   string       `json:"bank"`
	FraudStatus            string       `json:"fraud_status"`
	PaymentType            string       `json:"payment_type"`
	OrderID                string       `json:"order_id"`
	TransactionID          string       `json:"transaction_id"`
	TransactionStatus      string       `json:"transaction_status"`
	GrossAmount            string       `json:"gross_amount"`
	MaskedCard             string       `json:"masked_card"`
	Currency               string       `json:"currency"`
	CardType               string       `json:"card_type"`
	ChannelResponseCode    string       `json:"channel_response_code"`
	ChannelResponseMessage string       `json:"channel_response_message"`
	ApprovalCode           string       `json:"approval_code"`
}

type ItemRemoteCheckout struct {
	ID         int32   `json:"id"`
	SubTotal   float32 `json:"sub_total"`
	GrandTotal float32 `json:"grand_total"`
	DiscountID int32   `json:"discount_id"`
}

type RemoteCheckoutRequest struct {
	SubTotal          float32 `json:"sub_total" binding:"required"`
	GrandTotal        float32 `json:"grand_total" binding:"required"`
	SenderEmail       string  `json:"sender_email" binding:"required"`
	ReceiverEmail     string  `json:"receiver_email"binding:"required"`
	UserID            int32   `json:"user_id"`
	RemoteServiceName string  `json:"remote_service_name"`
	CurrencyCode      string  `json:"currency_code"`
	Item              []ItemRemoteCheckout
	UserMessage       string `json:"user_message"`
	RemoteOrderNumber string `json:"remote_order_number"`
}
