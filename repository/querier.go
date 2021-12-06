package repository

import (
	"context"
	"scoop-order/internal/databases"
)

type Querier interface{
	SelectOrder(ctx context.Context, arg SelectOrderParams) ([]databases.CoreOrder, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (databases.CoreOrder, error)
	CreateOrderlines(ctx context.Context, arg CreateOrderlinesParams) (databases.CoreOrderline, error)
	SelectDetailOrderByID(ctx context.Context, id int32) ([]SelectDetailOrderByIDRow, error)
	SelectPendingOrder(ctx context.Context, arg SelectPendingOrderParams) ([]databases.CoreOrder, error)
	SelectOrderByID(ctx context.Context, id int32) (databases.CoreOrder, error)
	SelectOrderByOrderNumber(ctx context.Context, orderNumber int64) (databases.CoreOrder, error)
	SelectOrderByUserID(ctx context.Context, userID int32) ([]databases.CoreOrder, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (databases.CoreOrder, error)
	UpdateOrderlines(ctx context.Context, arg UpdateOrderlinesParams) (databases.CoreOrderline, error)
	SelectOfferByID(ctx context.Context, id int32)(databases.CoreOffer, error)
	SelectDiscountByID(ctx context.Context, id int32) (databases.CoreDiscount, error)
	SelectDiscountByIDs(ctx context.Context, id []int32) ([]databases.CoreDiscount, error)
	SelectDiscountCodeByCode(ctx context.Context, code string) (databases.CoreDiscountcode, error)
	UpdateDiscountCode(ctx context.Context, arg UpdateDiscountCodeParams) (databases.CoreDiscountcode, error)
	SelectItemOfSingleOffer(ctx context.Context, offerID int32) (SelectItemOfSingleOfferRows, error)
	SelectRestrictCountriesByOffer(ctx context.Context, offerID int32) (SelectRestrictCountriesByOfferRows, error)
	SelectPlatformOffer(ctx context.Context, offerID int32, platformID int32) (databases.CoreOffer, error)
	SelectPaymentGateway(ctx context.Context, paymentGatewayID int32) (SelectPaymentGatewaysRow, error)
	SelectPaymentGateways(ctx context.Context) ([]SelectPaymentGatewaysRow, error)
	SelectUser(ctx context.Context, userID int32) (SelectUserRows, error)
	//CreatePayment(ctx context.Context, arg databases.CreatePaymentParams) (databases.CorePayment, error)
	//UpdatePaymentByOrder(ctx context.Context, arg databases.UpdatePaymentParams) (databases.CorePayment, error)
}

var _ Querier = (*Queries)(nil)