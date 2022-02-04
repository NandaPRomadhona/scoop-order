package repository

import (
	"context"
	"database/sql"
	"fmt"
	"scoop-order/internal/databases"
	"time"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO core_orders(created, modified, order_number, total_amount, final_amount, user_id, client_id, partner_id,
                        is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id,
                        order_status, remote_order_number, is_renewal)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
`

type CreateOrderParams struct {
	Created           time.Time      `json:"created"`
	Modified          time.Time      `json:"modified"`
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

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (databases.CoreOrder, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.Created,
		arg.Modified,
		arg.OrderNumber,
		arg.TotalAmount,
		arg.FinalAmount,
		arg.UserID,
		arg.ClientID,
		arg.PartnerID,
		arg.IsActive,
		arg.PointReward,
		arg.CurrencyCode,
		arg.PaymentgatewayID,
		arg.TierCode,
		arg.PlatformID,
		arg.TemporderID,
		arg.OrderStatus,
		arg.RemoteOrderNumber,
		arg.IsRenewal,
	)
	var i databases.CoreOrder
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderNumber,
		&i.TotalAmount,
		&i.FinalAmount,
		&i.UserID,
		&i.ClientID,
		&i.PartnerID,
		&i.IsActive,
		&i.PointReward,
		&i.CurrencyCode,
		&i.PaymentgatewayID,
		&i.TierCode,
		&i.PlatformID,
		&i.TemporderID,
		&i.OrderStatus,
		&i.RemoteOrderNumber,
		&i.IsRenewal,
	)
	return i, err
}

const createOrderlines = `-- name: CreateOrderlines :one
INSERT INTO core_orderlines(created, modified, name, offer_id, is_active, is_free, is_discount, user_id, campaign_id, order_id, quantity,
                            orderline_status, currency_code, price, final_price, localized_currency_code,
                            localized_final_price, is_trial)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING created, modified, id, name, offer_id, is_active, is_free, is_discount, user_id, campaign_id, order_id, quantity, orderline_status, currency_code, price, final_price, localized_currency_code, localized_final_price, is_trial
`

type CreateOrderlinesParams struct {
	Created 	time.Time `json:"created"`
	Modified time.Time `json:"modified"`
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

func (q *Queries) CreateOrderlines(ctx context.Context, arg CreateOrderlinesParams) (databases.CoreOrderline, error) {
	row := q.db.QueryRowContext(ctx, createOrderlines,
		arg.Created,
		arg.Modified,
		arg.Name,
		arg.OfferID,
		arg.IsActive,
		arg.IsFree,
		arg.IsDiscount,
		arg.UserID,
		arg.CampaignID,
		arg.OrderID,
		arg.Quantity,
		arg.OrderlineStatus,
		arg.CurrencyCode,
		arg.Price,
		arg.FinalPrice,
		arg.LocalizedCurrencyCode,
		arg.LocalizedFinalPrice,
		arg.IsTrial,
	)
	var i databases.CoreOrderline
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.Name,
		&i.OfferID,
		&i.IsActive,
		&i.IsFree,
		&i.IsDiscount,
		&i.UserID,
		&i.CampaignID,
		&i.OrderID,
		&i.Quantity,
		&i.OrderlineStatus,
		&i.CurrencyCode,
		&i.Price,
		&i.FinalPrice,
		&i.LocalizedCurrencyCode,
		&i.LocalizedFinalPrice,
		&i.IsTrial,
	)
	return i, err
}

const selectDetailOrderByID = `-- name: SelectDetailOrderByID :many
Select o.created, o.modified, o.id, order_number, total_amount, final_amount, o.user_id, client_id, partner_id, o.is_active, point_reward, o.currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal, ol.created, ol.modified, ol.id, name, offer_id, ol.is_active, is_free, is_discount, ol.user_id, campaign_id, order_id, quantity, orderline_status, ol.currency_code, price, final_price, localized_currency_code, localized_final_price, is_trial
from core_orders o
         JOIN core_orderlines ol
              ON o.id = ol.order_id
WHERE o.id = $1
`

type SelectDetailOrderByIDRow struct {
	Created               time.Time       `json:"created"`
	Modified              time.Time       `json:"modified"`
	ID                    int32           `json:"id"`
	OrderNumber           int64           `json:"order_number"`
	TotalAmount           float64         `json:"total_amount"`
	FinalAmount           float64         `json:"final_amount"`
	UserID                int32           `json:"user_id"`
	ClientID              sql.NullInt32   `json:"client_id"`
	PartnerID             sql.NullInt32   `json:"partner_id"`
	IsActive              sql.NullBool    `json:"is_active"`
	PointReward           sql.NullInt32   `json:"point_reward"`
	CurrencyCode          string          `json:"currency_code"`
	PaymentgatewayID      int32           `json:"paymentgateway_id"`
	TierCode              sql.NullString  `json:"tier_code"`
	PlatformID            sql.NullInt32   `json:"platform_id"`
	TemporderID           sql.NullInt32   `json:"temporder_id"`
	OrderStatus           sql.NullInt32   `json:"order_status"`
	RemoteOrderNumber     sql.NullString  `json:"remote_order_number"`
	IsRenewal             sql.NullBool    `json:"is_renewal"`
	Created_2             time.Time       `json:"created_2"`
	Modified_2            time.Time       `json:"modified_2"`
	ID_2                  int32           `json:"id_2"`
	Name                  sql.NullString  `json:"name"`
	OfferID               int32           `json:"offer_id"`
	IsActive_2            sql.NullBool    `json:"is_active_2"`
	IsFree                sql.NullBool    `json:"is_free"`
	IsDiscount            sql.NullBool    `json:"is_discount"`
	UserID_2              sql.NullInt32   `json:"user_id_2"`
	CampaignID            sql.NullInt32   `json:"campaign_id"`
	OrderID               sql.NullInt32   `json:"order_id"`
	Quantity              sql.NullInt32   `json:"quantity"`
	OrderlineStatus       sql.NullInt32   `json:"orderline_status"`
	CurrencyCode_2        sql.NullString  `json:"currency_code_2"`
	Price                 sql.NullFloat64 `json:"price"`
	FinalPrice            sql.NullFloat64 `json:"final_price"`
	LocalizedCurrencyCode sql.NullString  `json:"localized_currency_code"`
	LocalizedFinalPrice   sql.NullFloat64 `json:"localized_final_price"`
	IsTrial               sql.NullBool    `json:"is_trial"`
}

func (q *Queries) SelectDetailOrderByID(ctx context.Context, id int32) ([]SelectDetailOrderByIDRow, error) {
	rows, err := q.db.QueryContext(ctx, selectDetailOrderByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectDetailOrderByIDRow
	for rows.Next() {
		var i SelectDetailOrderByIDRow
		if err := rows.Scan(
			&i.Created,
			&i.Modified,
			&i.ID,
			&i.OrderNumber,
			&i.TotalAmount,
			&i.FinalAmount,
			&i.UserID,
			&i.ClientID,
			&i.PartnerID,
			&i.IsActive,
			&i.PointReward,
			&i.CurrencyCode,
			&i.PaymentgatewayID,
			&i.TierCode,
			&i.PlatformID,
			&i.TemporderID,
			&i.OrderStatus,
			&i.RemoteOrderNumber,
			&i.IsRenewal,
			&i.Created_2,
			&i.Modified_2,
			&i.ID_2,
			&i.Name,
			&i.OfferID,
			&i.IsActive_2,
			&i.IsFree,
			&i.IsDiscount,
			&i.UserID_2,
			&i.CampaignID,
			&i.OrderID,
			&i.Quantity,
			&i.OrderlineStatus,
			&i.CurrencyCode_2,
			&i.Price,
			&i.FinalPrice,
			&i.LocalizedCurrencyCode,
			&i.LocalizedFinalPrice,
			&i.IsTrial,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectPendingOrder = `-- name: SelectOrder :many
SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
FROM core_orders WHERE order_status = 20001
ORDER BY id DESC LIMIT $1
OFFSET $2
`

type SelectPendingOrderParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (q *Queries) SelectPendingOrder(ctx context.Context, arg SelectPendingOrderParams) ([]databases.CoreOrder, error) {
	rows, err := q.db.QueryContext(ctx, selectPendingOrder, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []databases.CoreOrder
	for rows.Next() {
		var i databases.CoreOrder
		if err := rows.Scan(
			&i.Created,
			&i.Modified,
			&i.ID,
			&i.OrderNumber,
			&i.TotalAmount,
			&i.FinalAmount,
			&i.UserID,
			&i.ClientID,
			&i.PartnerID,
			&i.IsActive,
			&i.PointReward,
			&i.CurrencyCode,
			&i.PaymentgatewayID,
			&i.TierCode,
			&i.PlatformID,
			&i.TemporderID,
			&i.OrderStatus,
			&i.RemoteOrderNumber,
			&i.IsRenewal,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectOrder = `-- name: SelectOrder :many
SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
FROM core_orders
ORDER BY id DESC LIMIT $1
OFFSET $2
`

type SelectOrderParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (q *Queries) SelectOrder(ctx context.Context, arg SelectOrderParams) ([]databases.CoreOrder, error) {
	rows, err := q.db.QueryContext(ctx, selectOrder, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []databases.CoreOrder
	for rows.Next() {
		var i databases.CoreOrder
		if err := rows.Scan(
			&i.Created,
			&i.Modified,
			&i.ID,
			&i.OrderNumber,
			&i.TotalAmount,
			&i.FinalAmount,
			&i.UserID,
			&i.ClientID,
			&i.PartnerID,
			&i.IsActive,
			&i.PointReward,
			&i.CurrencyCode,
			&i.PaymentgatewayID,
			&i.TierCode,
			&i.PlatformID,
			&i.TemporderID,
			&i.OrderStatus,
			&i.RemoteOrderNumber,
			&i.IsRenewal,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectOrderByID = `-- name: SelectOrderByID :one
SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
FROM core_orders
WHERE id = $1
`

func (q *Queries) SelectOrderByID(ctx context.Context, id int32) (databases.CoreOrder, error) {
	row := q.db.QueryRowContext(ctx, selectOrderByID, id)
	var i databases.CoreOrder
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderNumber,
		&i.TotalAmount,
		&i.FinalAmount,
		&i.UserID,
		&i.ClientID,
		&i.PartnerID,
		&i.IsActive,
		&i.PointReward,
		&i.CurrencyCode,
		&i.PaymentgatewayID,
		&i.TierCode,
		&i.PlatformID,
		&i.TemporderID,
		&i.OrderStatus,
		&i.RemoteOrderNumber,
		&i.IsRenewal,
	)
	return i, err
}

const selectOrderByOrderNumber = `-- name: SelectOrderByOrderNumber :one
SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
FROM core_orders
WHERE order_number = $1
`

func (q *Queries) SelectOrderByOrderNumber(ctx context.Context, orderNumber int64) (databases.CoreOrder, error) {
	fmt.Println("orderNUmber: ", orderNumber)
	row := q.db.QueryRowContext(ctx, selectOrderByOrderNumber, orderNumber)
	var i databases.CoreOrder
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderNumber,
		&i.TotalAmount,
		&i.FinalAmount,
		&i.UserID,
		&i.ClientID,
		&i.PartnerID,
		&i.IsActive,
		&i.PointReward,
		&i.CurrencyCode,
		&i.PaymentgatewayID,
		&i.TierCode,
		&i.PlatformID,
		&i.TemporderID,
		&i.OrderStatus,
		&i.RemoteOrderNumber,
		&i.IsRenewal,
	)
	return i, err
}

const selectOrderByUserID = `-- name: SelectOrderByUserID :many
SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
FROM core_orders
WHERE user_id = $1
`

func (q *Queries) SelectOrderByUserID(ctx context.Context, userID int32) ([]databases.CoreOrder, error) {
	rows, err := q.db.QueryContext(ctx, selectOrderByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []databases.CoreOrder
	for rows.Next() {
		var i databases.CoreOrder
		if err := rows.Scan(
			&i.Created,
			&i.Modified,
			&i.ID,
			&i.OrderNumber,
			&i.TotalAmount,
			&i.FinalAmount,
			&i.UserID,
			&i.ClientID,
			&i.PartnerID,
			&i.IsActive,
			&i.PointReward,
			&i.CurrencyCode,
			&i.PaymentgatewayID,
			&i.TierCode,
			&i.PlatformID,
			&i.TemporderID,
			&i.OrderStatus,
			&i.RemoteOrderNumber,
			&i.IsRenewal,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrder = `-- name: UpdateOrder :one
UPDATE core_orders
SET modified            = $2,
    order_status        = $3
WHERE id = $1 RETURNING created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
`

type UpdateOrderParams struct {
	ID          int32         `json:"id"`
	Modified    time.Time     `json:"modified"`
	OrderStatus sql.NullInt32 `json:"order_status"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (databases.CoreOrder, error) {
	row := q.db.QueryRowContext(ctx, updateOrder,
		arg.ID,
		arg.Modified,
		arg.OrderStatus,
	)
	var i databases.CoreOrder
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderNumber,
		&i.TotalAmount,
		&i.FinalAmount,
		&i.UserID,
		&i.ClientID,
		&i.PartnerID,
		&i.IsActive,
		&i.PointReward,
		&i.CurrencyCode,
		&i.PaymentgatewayID,
		&i.TierCode,
		&i.PlatformID,
		&i.TemporderID,
		&i.OrderStatus,
		&i.RemoteOrderNumber,
		&i.IsRenewal,
	)
	return i, err
}

const updateOrderlines = `-- name: UpdateOrderlines :one
UPDATE core_orderlines
SET modified=$2,
    orderline_status=$3
WHERE order_id = $1 RETURNING created, modified, id, name, offer_id, is_active, is_free, is_discount, user_id, campaign_id, order_id, quantity, orderline_status, currency_code, price, final_price, localized_currency_code, localized_final_price, is_trial
`

type UpdateOrderlinesParams struct {
	OrderID         int32         `json:"order_id"`
	Modified        time.Time     `json:"modified"`
	OrderlineStatus sql.NullInt32 `json:"orderline_status"`
}

func (q *Queries) UpdateOrderlines(ctx context.Context, arg UpdateOrderlinesParams) (databases.CoreOrderline, error) {
	row := q.db.QueryRowContext(ctx, updateOrderlines,
		arg.OrderID,
		arg.Modified,
		arg.OrderlineStatus,
	)
	var i databases.CoreOrderline
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.Name,
		&i.OfferID,
		&i.IsActive,
		&i.IsFree,
		&i.IsDiscount,
		&i.UserID,
		&i.CampaignID,
		&i.OrderID,
		&i.Quantity,
		&i.OrderlineStatus,
		&i.CurrencyCode,
		&i.Price,
		&i.FinalPrice,
		&i.LocalizedCurrencyCode,
		&i.LocalizedFinalPrice,
		&i.IsTrial,
	)
	return i, err
}

type OrderPublished struct {
	OrderID           int32   `json:"order_id"`
	PartnerID         int32   `json:"partner_id"`
	IsActive          bool    `json:"is_active"`
	PointReward       int32   `json:"point_reward"`
	OrderStatus       int32   `json:"order_status"`
	RemoteOrderNumber string  `json:"remote_order_number"`
	Currency          string  `json:"currency"`
	TierCode          string  `json:"tier_code"`
	TotalAmount       float64 `json:"total_amount"`
	FinalaAmount      float64 `json:"finala_amount"`
	OrderNumber       int     `json:"order_number"`
}

func DestructorsOrders(orders []databases.CoreOrder) []OrderPublished {
	var Orders []OrderPublished

	for _, hit := range orders {
		id := hit.ID
		partnerID := hit.PartnerID
		isActive := hit.IsActive
		pointReward := hit.PointReward
		orderStatus := hit.OrderStatus
		remoteOrderNumber := hit.RemoteOrderNumber
		currency := hit.CurrencyCode
		tierCode := hit.TierCode
		totalAmount := hit.TotalAmount
		finalAmount := hit.FinalAmount
		orderNumber := hit.OrderNumber

		Orders = append(Orders, OrderPublished{
			OrderID:           id,
			PartnerID:         partnerID.Int32,
			IsActive:          isActive.Bool,
			PointReward:       pointReward.Int32,
			OrderStatus:       orderStatus.Int32,
			RemoteOrderNumber: remoteOrderNumber.String,
			Currency:          currency,
			TierCode:          tierCode.String,
			TotalAmount:       totalAmount,
			OrderNumber:       orderNumber,
			FinalaAmount:      finalAmount,
		})
	}
	return Orders
}

type OrderLinePublished struct {
	ID         int32   `json:"id"`
	IsDiscount bool    `json:"is_discount"`
	OfferID    int32   `json:"offer_id"`
	TotalPrice float64 `json:"total_price"`
	FinalPrice float64 `json:"final_price"`
}

type OrderDetailPublished struct {
	OrderPublished
	OrderLines []OrderLinePublished `json:"order_lines"`
}

func DestructorOrderDetail(order databases.CoreOrder, orderlines []SelectDetailOrderByIDRow) interface{} {
	var detailOrder OrderDetailPublished
	var orderLines []OrderLinePublished

	for _, hit := range orderlines {
		id := hit.ID_2
		isDiscount := hit.IsDiscount
		offerId := hit.OfferID
		totalPrice := hit.Price.Float64
		finalPrice := hit.FinalPrice.Float64

		orderLines = append(orderLines, OrderLinePublished{
			ID:         id,
			IsDiscount: isDiscount.Bool,
			OfferID:    offerId,
			TotalPrice: totalPrice,
			FinalPrice: finalPrice,
		})
	}
	// assign order response
	detailOrder.OrderID = order.ID
	detailOrder.OrderStatus = order.OrderStatus.Int32
	detailOrder.IsActive = order.IsActive.Bool
	detailOrder.RemoteOrderNumber = order.RemoteOrderNumber.String
	detailOrder.OrderLines = orderLines
	detailOrder.OrderNumber = order.OrderNumber
	detailOrder.TotalAmount = order.TotalAmount
	detailOrder.FinalaAmount = order.FinalAmount
	detailOrder.Currency = order.CurrencyCode
	detailOrder.PartnerID = order.PartnerID.Int32

	return detailOrder
}

const createOrderDetail = `-- name: CreateOrderDetail :one
INSERT INTO public.core_orderdetails(
    created, modified, order_id, user_id, user_email, user_name, user_street_address, user_city, user_zipcode, user_state, user_country, latitude, longitude, note, ip_address, os_version, client_version, device_model, temporder_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
RETURNING created, modified, id, order_id, user_id, user_email, user_name, user_street_address, user_city, user_zipcode, user_state, user_country, latitude, longitude, note, ip_address, os_version, client_version, device_model, temporder_id
`

type CreateOrderDetailParams struct {
	Created           time.Time      `json:"created"`
	Modified          time.Time      `json:"modified"`
	OrderID           int32          `json:"order_id"`
	UserID            sql.NullInt32  `json:"user_id"`
	UserEmail         sql.NullString `json:"user_email"`
	UserName          sql.NullString `json:"user_name"`
	UserStreetAddress sql.NullString `json:"user_street_address"`
	UserCity          sql.NullString `json:"user_city"`
	UserZipcode       sql.NullString `json:"user_zipcode"`
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

func (q *Queries) CreateOrderDetail(ctx context.Context, arg CreateOrderDetailParams) (databases.CoreOrderDetail, error) {
	row := q.db.QueryRowContext(ctx, createOrderDetail,
		arg.Created,
		arg.Modified,
		arg.OrderID,
		arg.UserID,
		arg.UserEmail,
		arg.UserName,
		arg.UserStreetAddress,
		arg.UserCity,
		arg.UserZipcode,
		arg.UserState,
		arg.UserCountry,
		arg.Latitude,
		arg.Longitude,
		arg.Note,
		arg.IpAddress,
		arg.OsVersion,
		arg.ClientVersion,
		arg.DeviceModel,
		arg.TemporderID,
	)
	var i databases.CoreOrderDetail
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.UserEmail,
		&i.UserName,
		&i.UserStreetAddress,
		&i.UserCity,
		&i.UserZipCode,
		&i.UserState,
		&i.UserCountry,
		&i.Latitude,
		&i.Longitude,
		&i.Note,
		&i.IpAddress,
		&i.OsVersion,
		&i.ClientVersion,
		&i.DeviceModel,
		&i.TemporderID,
	)
	return i, err
}

const selectOrderDetail = `-- name: SelectOrderDetail :one
SELECT created, modified, id, order_id, user_id, user_email, user_name, user_street_address, user_city, user_zipcode, user_state, user_country, latitude, longitude, note, ip_address, os_version, client_version, device_model, temporder_id
FROM public.core_orderdetails WHERE order_id = $1
`

func (q *Queries) SelectOrderDetail(ctx context.Context, orderID sql.NullInt32) (databases.CoreOrderDetail, error) {
	row := q.db.QueryRowContext(ctx, selectOrderDetail, orderID)
	var i databases.CoreOrderDetail
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.UserEmail,
		&i.UserName,
		&i.UserStreetAddress,
		&i.UserCity,
		&i.UserZipCode,
		&i.UserState,
		&i.UserCountry,
		&i.Latitude,
		&i.Longitude,
		&i.Note,
		&i.IpAddress,
		&i.OsVersion,
		&i.ClientVersion,
		&i.DeviceModel,
		&i.TemporderID,
	)
	return i, err
}

const createOrderlineDiscount = `-- name: createOrderlineDiscount :one
INSERT INTO public.core_orderlinediscounts(
	created, modified, order_id, orderline_id, discount_id, discount_name, currency_code, discount_code, discount_type, discount_value, raw_price, final_price)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING created, modified, id, order_id, orderline_id, discount_id, discount_name, currency_code, discount_code, discount_type, discount_value, raw_price, final_price
`

func (q *Queries) CreateOrderlineDiscount(ctx context.Context, arg databases.CoreOrderlineDiscount) (databases.CoreOrderlineDiscount, error) {
	row := q.db.QueryRowContext(ctx, createOrderlineDiscount,
		arg.Created,
		arg.Modified,
		arg.OrderID,
		arg.OrderlineID,
		arg.DiscountID,
		arg.DiscountName,
		arg.CurrencyCode,
		arg.DiscountCode,
		arg.DiscountType,
		arg.DiscountValue,
		arg.RawPrice,
		arg.FinalPrice,
	)
	var i databases.CoreOrderlineDiscount
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderID,
		&i.OrderlineID,
		&i.DiscountID,
		&i.DiscountName,
		&i.CurrencyCode,
		&i.DiscountCode,
		&i.DiscountType,
		&i.DiscountValue,
		&i.RawPrice,
		&i.FinalPrice,
	)
	return i, err
}