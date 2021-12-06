package databases

import (
	"database/sql"
	"encoding/json"
	"time"
)

type CoreCampaign struct {
	Created     sql.NullTime   `json:"created"`
	Modified    sql.NullTime   `json:"modified"`
	ID          int32          `json:"id"`
	Name        sql.NullString `json:"name"`
	Description sql.NullString `json:"description"`
	IsActive    sql.NullBool   `json:"is_active"`
	StartDate   sql.NullTime   `json:"start_date"`
	EndDate     sql.NullTime   `json:"end_date"`
	TotalCost   sql.NullString `json:"total_cost"`
}

type CoreDiscount struct {
	Created              sql.NullTime    `json:"created"`
	Modified             sql.NullTime    `json:"modified"`
	ID                   int32           `json:"id"`
	Name                 string          `json:"name"`
	TagName              sql.NullString  `json:"tag_name"`
	Description          sql.NullString  `json:"description"`
	CampaignID           int32           `json:"campaign_id"`
	ValidTo              time.Time       `json:"valid_to"`
	ValidFrom            time.Time       `json:"valid_from"`
	DiscountRule         sql.NullInt32   `json:"discount_rule"`
	DiscountType         sql.NullInt32   `json:"discount_type"`
	DiscountStatus       sql.NullInt32   `json:"discount_status"`
	DiscountScheduleType sql.NullInt32   `json:"discount_schedule_type"`
	IsActive             sql.NullBool    `json:"is_active"`
	DiscountUsd          float64         `json:"discount_usd"`
	DiscountIdr          float64         `json:"discount_idr"`
	DiscountPoint        sql.NullInt32   `json:"discount_point"`
	MinUsdOrderPrice     sql.NullFloat64 `json:"min_usd_order_price"`
	MaxUsdOrderPrice     sql.NullFloat64 `json:"max_usd_order_price"`
	MinIdrOrderPrice     sql.NullFloat64 `json:"min_idr_order_price"`
	MaxIdrOrderPrice     sql.NullFloat64 `json:"max_idr_order_price"`
	PredefinedGroup      sql.NullInt32   `json:"predefined_group"`
	VendorParticipation  sql.NullInt32   `json:"vendor_participation"`
	PartnerParticipation sql.NullInt32   `json:"partner_participation"`
	SalesRecognition     sql.NullInt32   `json:"sales_recognition"`
	BinCodes             []string        `json:"bin_codes"`
	//TrialTime            sql.NullInt64   `json:"trial_time"`
}

type CoreDiscountcode struct {
	Created      sql.NullTime  `json:"created"`
	Modified     sql.NullTime  `json:"modified"`
	ID           int32         `json:"id"`
	Code         string        `json:"code"`
	MaxUses      sql.NullInt32 `json:"max_uses"`
	CurrentUses  sql.NullInt32 `json:"current_uses"`
	IsActive     sql.NullBool  `json:"is_active"`
	DiscountType sql.NullInt32 `json:"discount_type"`
	DiscountID   int32         `json:"discount_id"`
	IsForNewUser sql.NullBool  `json:"is_for_new_user"`
}

type CoreOrder struct {
	Created           time.Time      `json:"created"`
	Modified          time.Time      `json:"modified"`
	ID                int32          `json:"id"`
	OrderNumber       int            `json:"order_number"`
	TotalAmount       float64        `json:"total_amount"`
	FinalAmount       float64        `json:"final_amount"`
	UserID            int32          `json:"user_id"`
	ClientID          sql.NullInt32  `json:"client_id"`
	PartnerID         sql.NullInt32  `json:"partner_id"`
	IsActive          sql.NullBool   `json:"is_active"`
	PointReward       sql.NullInt32  `json:"point_reward"`
	CurrencyCode      string         `json:"currency_code"`
	PaymentgatewayID  int32          `json:"paymentgateway_id"`
	TierCode          sql.NullString `json:"tier_code"`
	PlatformID        sql.NullInt32  `json:"platform_id"`
	TemporderID       sql.NullInt32  `json:"temporder_id"`
	OrderStatus       sql.NullInt32  `json:"order_status"`
	RemoteOrderNumber sql.NullString `json:"remote_order_number"`
	IsRenewal         sql.NullBool   `json:"is_renewal"`
}

type CoreOrderline struct {
	Created               time.Time       `json:"created"`
	Modified              time.Time       `json:"modified"`
	ID                    int32           `json:"id"`
	Name                  sql.NullString  `json:"name"`
	OfferID               int32           `json:"offer_id"`
	IsActive              sql.NullBool    `json:"is_active"`
	IsFree                sql.NullBool    `json:"is_free"`
	IsDiscount            sql.NullBool    `json:"is_discount"`
	UserID                sql.NullInt32   `json:"user_id"`
	CampaignID            sql.NullInt32   `json:"campaign_id"`
	OrderID               sql.NullInt32   `json:"order_id"`
	Quantity              sql.NullInt32   `json:"quantity"`
	OrderlineStatus       sql.NullInt32   `json:"orderline_status"`
	CurrencyCode          sql.NullString  `json:"currency_code"`
	Price                 sql.NullFloat64 `json:"price"`
	FinalPrice            sql.NullFloat64 `json:"final_price"`
	LocalizedCurrencyCode sql.NullString  `json:"localized_currency_code"`
	LocalizedFinalPrice   sql.NullFloat64 `json:"localized_final_price"`
	IsTrial               sql.NullBool    `json:"is_trial"`
}

type CoreOffer struct {
	Created            time.Time       `json:"created"`
	Modified           time.Time       `json:"modified"`
	ID                 int32           `json:"id"`
	Name               sql.NullString  `json:"name"`
	OfferStatus        int16           `json:"offer_status"`
	SortPriority       int16           `json:"sort_priority"`
	IsActive           sql.NullBool    `json:"is_active"`
	OfferTypeID        int32           `json:"offer_type_id"`
	ExclusiveClients   []string        `json:"exclusive_clients"`
	IsFree             sql.NullBool    `json:"is_free"`
	OfferCode          sql.NullString  `json:"offer_code"`
	ItemCode           sql.NullString  `json:"item_code"`
	PriceUsd           float64         `json:"price_usd"`
	PriceIdr           float64         `json:"price_idr"`
	PricePoint         sql.NullInt32   `json:"price_point"`
	DiscountID         []string        `json:"discount_id"`
	DiscountTag        sql.NullString  `json:"discount_tag"`
	DiscountName       sql.NullString  `json:"discount_name"`
	DiscountPriceUsd   sql.NullFloat64 `json:"discount_price_usd"`
	DiscountPriceIdr   sql.NullFloat64 `json:"discount_price_idr"`
	DiscountPricePoint sql.NullInt32   `json:"discount_price_point"`
	IsDiscount         sql.NullBool    `json:"is_discount"`
	ImageHighres       sql.NullString  `json:"image_highres"`
	ImageNormal        sql.NullString  `json:"image_normal"`
	VendorPriceUsd     sql.NullFloat64 `json:"vendor_price_usd"`
	VendorPriceIdr     sql.NullFloat64 `json:"vendor_price_idr"`
	VendorPricePoint   sql.NullInt32   `json:"vendor_price_point"`
	LongName           sql.NullString  `json:"long_name"`
	TierID             sql.NullInt32   `json:"tier_id"`
	TierCode           sql.NullString  `json:"tier_code"`
	Currency           sql.NullString  `json:"currency"`
	DiscountTierID     sql.NullInt32   `json:"discount_tier_id"`
	DiscountTierCode   sql.NullString  `json:"discount_tier_code"`
	DiscountTierPrice  sql.NullFloat64 `json:"discount_tier_price"`
}

type CorePlatform struct {
	Created     sql.NullTime   `json:"created"`
	Modified    sql.NullTime   `json:"modified"`
	ID          int32          `json:"id"`
	Name        sql.NullString `json:"name"`
	Description sql.NullString `json:"description"`
	IsActive    bool           `json:"is_active"`
}

type CorePlatformsOffer struct {
	Created            sql.NullTime    `json:"created"`
	Modified           sql.NullTime    `json:"modified"`
	ID                 int32           `json:"id"`
	OfferID            int32           `json:"offer_id"`
	PlatformID         int32           `json:"platform_id"`
	TierID             sql.NullInt32   `json:"tier_id"`
	TierCode           sql.NullString  `json:"tier_code"`
	Currency           sql.NullString  `json:"currency"`
	PriceUsd           sql.NullString  `json:"price_usd"`
	PriceIdr           sql.NullString  `json:"price_idr"`
	PricePoint         sql.NullInt32   `json:"price_point"`
	DiscountTierID     sql.NullInt32   `json:"discount_tier_id"`
	DiscountTierCode   sql.NullString  `json:"discount_tier_code"`
	DiscountTag        sql.NullString  `json:"discount_tag"`
	DiscountName       sql.NullString  `json:"discount_name"`
	DiscountTierPrice  sql.NullString  `json:"discount_tier_price"`
	DiscountPriceUsd   sql.NullString  `json:"discount_price_usd"`
	DiscountPriceIdr   sql.NullString  `json:"discount_price_idr"`
	DiscountPricePoint sql.NullInt32   `json:"discount_price_point"`
	DiscountID         []string        `json:"discount_id"`
	Colors             json.RawMessage `json:"colors"`
}

type CoreDistributioncountry struct {
	Created   sql.NullTime   `json:"created"`
	Modified  sql.NullTime   `json:"modified"`
	ID        int32          `json:"id"`
	Name      sql.NullString `json:"name"`
	GroupType sql.NullInt32  `json:"group_type"`
	Countries []string       `json:"countries"`
	IsActive  sql.NullBool   `json:"is_active"`
	VendorID  sql.NullInt32  `json:"vendor_id"`
}

type CorePaymentgateway struct {
	Created                sql.NullTime   `json:"created"`
	Modified               sql.NullTime   `json:"modified"`
	ID                     int32          `json:"id"`
	Name                   string         `json:"name"`
	Mnc                    sql.NullString `json:"mnc"`
	Mcc                    sql.NullString `json:"mcc"`
	SmsNumber              sql.NullString `json:"sms_number"`
	PaymentFlowType        int16          `json:"payment_flow_type"`
	Description            sql.NullString `json:"description"`
	Slug                   sql.NullString `json:"slug"`
	Meta                   sql.NullString `json:"meta"`
	IconImageNormal        sql.NullString `json:"icon_image_normal"`
	IconImageHighres       sql.NullString `json:"icon_image_highres"`
	SortPriority           int16          `json:"sort_priority"`
	IsActive               sql.NullBool   `json:"is_active"`
	OrganizationID         sql.NullInt32  `json:"organization_id"`
	Clients                []string       `json:"clients"`
	BaseCurrencyID         int32          `json:"base_currency_id"`
	MinimalAmount          sql.NullString `json:"minimal_amount"`
	LowestSupportedVersion sql.NullString `json:"lowest_supported_version"`
	IsRenewal              sql.NullBool   `json:"is_renewal"`
	PaymentGroup           interface{}    `json:"payment_group"`
	MerchantCode           sql.NullString `json:"merchant_code"`
	MerchantKey            sql.NullString `json:"merchant_key"`
}

type CorePayment struct {
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

type CoreOrderDetail struct {
	Created           sql.NullTime   `json:"created"`
	Modified          sql.NullTime   `json:"modified"`
	ID                int32          `json:"id"`
	OrderID           sql.NullInt32  `json:"order_id"`
	UserID            sql.NullInt32  `json:"user_id"`
	UserEmail         sql.NullString `json:"user_email"`
	UserName          sql.NullString `json:"user_name"`
	UserStreetAddress sql.NullString `json:"user_street_address"`
	UserCity          sql.NullString `json:"user_city"`
	UserZipCode       sql.NullString `json:"user_zipcode"`
	UserState         sql.NullString `json:"user_state"`
	UserCountry       sql.NullString `json:"user_country"`
	Latitude          sql.NullString `json:"latitude"`
	Longitude         sql.NullString `json:"longitude"`
	Note              sql.NullString `json:"note"`
	IpAddress         sql.NullString `json:"ip_address"`
	OsVersion         sql.NullString `json:"os_version"`
	ClientVersion     sql.NullString `json:"client_version"`
	DeviceModel       sql.NullString `json:"device_model"`
	TemporderID       sql.NullInt32  `json:"temporder_id"`
}
