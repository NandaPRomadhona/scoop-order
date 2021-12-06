package schemas

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
	CountryCode      string `form:"discount_code,default=ID"`
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
	GeoInfo          GeoInfo
}

type CheckoutRequest struct {
	OfferID          []int32 `json:"offer_id" binding:"required"`
	DiscountCode     string  `json:"discount_code"`
	CurrencyCode     string  `json:"currency_code" binding:"required"`
	PlatformID       int32   `json:"platform_id" binding:"required"`
	PaymentGatewayID int32   `json:"payment_gateway_id" binding:"required"`
	GeoInfo          GeoInfo `json:"geo_info" binding:"required"`
}

type CheckoutTxParams struct {
	UserID           int32   `json:"user_id"`
	OfferID          []int32 `json:"offer_id"`
	DiscountCode     string  `json:"discount_code"`
	CurrencyCode     string  `json:"currency_code"`
	PlatformID       int32   `json:"platform_id"`
	PaymentGatewayID int32   `json:"payment_gateway_id"`
	GeoInfo          GeoInfo
}

type PaymentTxParams struct{
	OrderID int32 `json:"order_id"`
	OrderStatus int32 `json:"order_status"`
	PaymentStatus int32 `json:"payment_status"`
}

type CompleteRequest struct {
	OrderID int64 `json:"order_id"`
}