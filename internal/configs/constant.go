package configs

const (
	OrderTemp = 10000
	OrderNew = 10001
	OrderWaitingForPayment = 20001
	OrderPaymentInProcess = 20002
	OrderPaymentBilled = 20003
	OrderWaitingForDelivery = 30001
	OrderDeliveryInProcess = 30002
	OrderDelivered = 30003
	OrderWaitingForRefund = 40001
	OrderRefundInProcess = 40002
	OrderRefunded = 40003
	OrderCancelled = 50000
	OrderPaymentError = 50001
	OrderDeliveryError = 50002
	OrderRefundError = 50003
	OrderExpired = 50004
	OrderComplete = 90000
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
