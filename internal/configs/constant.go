package configs

const (
	OrderTemp               = 10000
	OrderNew                = 10001
	OrderWaitingForPayment  = 20001
	OrderPaymentInProcess   = 20002
	OrderPaymentBilled      = 20003
	OrderWaitingForDelivery = 30001
	OrderDeliveryInProcess  = 30002
	OrderDelivered          = 30003
	OrderWaitingForRefund   = 40001
	OrderRefundInProcess    = 40002
	OrderRefunded           = 40003
	OrderCancelled          = 50000
	OrderPaymentError       = 50001
	OrderDeliveryError      = 50002
	OrderRefundError        = 50003
	OrderExpired            = 50004
	OrderComplete           = 90000
)

const (
	WaitingForPayment = 20001
	PaymentInProcess  = 20002
	PaymentBilled     = 20003
	PaymentRestored   = 20004
	PaymentCancelled  = 50000
	PaymentError      = 50001
	PaymentExpired    = 50004
	PaymentDenied     = 50003
)

const (
	DiscountAllowedJoin    = 1
	DiscountNotAllowedJoin = 2
)

const (
	DiscountOffer   = 1
	DiscountOrder   = 2
	DiscountPGOffer = 3
	DiscountPGOrder = 4
	DiscountCode    = 5
)

const (
	DiscountRuleAmount    = 1
	DiscountRulPercentage = 2
	DiscountRulePrice     = 3
	// Discount Rule:
	// 1. BY AMOUNT
	// 2. BY PERCENTAGE
	// 3. BY TO AMOUNT #flush all price, with discount price
)

const (
	OfferTypeSingle       = 1
	OfferTypeSubscription = 2
	OfferTypeBundle       = 3
	OfferTypeBuffet       = 4
)

const (
	OfferStatusNew              = 1
	OfferStatusWaitingForReview = 2
	OfferStatusInReview         = 3
	OfferStatusReject           = 4
	OfferStatusApprove          = 5
	OfferStatusPrepareForSale   = 6
	OfferStatusReadyToSale      = 7
	OfferNotForSale             = 8
)

const (
	ItemTypeMagazine  = "magazine"
	ItemTypeBook      = "book"
	ItemTypeNewspaper = "newspaper"
	ItemTypeBonusItem = "bonus item"
	ItemTypeAudioBook = "audio book"
)
