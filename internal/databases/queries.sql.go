package databases
//
//import (
//	"context"
//	"database/sql"
//	"github.com/lib/pq"
//	"scoop-order/cmd/pkg/repository"
//	"time"
//)
//
//const createOrder = `-- name: CreateOrder :one
//INSERT INTO core_orders(created, modified, order_number, total_amount, final_amount, user_id, client_id, partner_id,
//                        is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id,
//                        order_status, remote_order_number, is_renewal)
//VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
//`
//
//type CreateOrderParams struct {
//	Created           time.Time      `json:"created"`
//	Modified          time.Time      `json:"modified"`
//	OrderNumber       int            `json:"order_number"`
//	TotalAmount       float64        `json:"total_amount"`
//	FinalAmount       float64        `json:"final_amount"`
//	UserID            int32          `json:"user_id"`
//	ClientID          sql.NullInt32  `json:"client_id"`
//	PartnerID         sql.NullInt32  `json:"partner_id"`
//	IsActive          sql.NullBool   `json:"is_active"`
//	PointReward       sql.NullInt32  `json:"point_reward"`
//	CurrencyCode      string         `json:"currency_code"`
//	PaymentgatewayID  int32          `json:"paymentgateway_id"`
//	TierCode          sql.NullString `json:"tier_code"`
//	PlatformID        sql.NullInt32  `json:"platform_id"`
//	TemporderID       sql.NullInt32  `json:"temporder_id"`
//	OrderStatus       sql.NullInt32  `json:"order_status"`
//	RemoteOrderNumber sql.NullString `json:"remote_order_number"`
//	IsRenewal         sql.NullBool   `json:"is_renewal"`
//}
//
//func (q *repository.Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (CoreOrder, error) {
//	row := q.db.QueryRowContext(ctx, createOrder,
//		arg.Created,
//		arg.Modified,
//		arg.OrderNumber,
//		arg.TotalAmount,
//		arg.FinalAmount,
//		arg.UserID,
//		arg.ClientID,
//		arg.PartnerID,
//		arg.IsActive,
//		arg.PointReward,
//		arg.CurrencyCode,
//		arg.PaymentgatewayID,
//		arg.TierCode,
//		arg.PlatformID,
//		arg.TemporderID,
//		arg.OrderStatus,
//		arg.RemoteOrderNumber,
//		arg.IsRenewal,
//	)
//	var i CoreOrder
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.OrderNumber,
//		&i.TotalAmount,
//		&i.FinalAmount,
//		&i.UserID,
//		&i.ClientID,
//		&i.PartnerID,
//		&i.IsActive,
//		&i.PointReward,
//		&i.CurrencyCode,
//		&i.PaymentgatewayID,
//		&i.TierCode,
//		&i.PlatformID,
//		&i.TemporderID,
//		&i.OrderStatus,
//		&i.RemoteOrderNumber,
//		&i.IsRenewal,
//	)
//	return i, err
//}
//
//const createOrderlines = `-- name: CreateOrderlines :one
//INSERT INTO core_orderlines(name, offer_id, is_active, is_free, is_discount, user_id, campaign_id, order_id, quantity,
//                            orderline_status, currency_code, price, final_price, localized_currency_code,
//                            localized_final_price, is_trial)
//VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING created, modified, id, name, offer_id, is_active, is_free, is_discount, user_id, campaign_id, order_id, quantity, orderline_status, currency_code, price, final_price, localized_currency_code, localized_final_price, is_trial
//`
//
//type CreateOrderlinesParams struct {
//	Name                  sql.NullString  `json:"name"`
//	OfferID               int32           `json:"offer_id"`
//	IsActive              sql.NullBool    `json:"is_active"`
//	IsFree                sql.NullBool    `json:"is_free"`
//	IsDiscount            sql.NullBool    `json:"is_discount"`
//	UserID                sql.NullInt32   `json:"user_id"`
//	CampaignID            sql.NullInt32   `json:"campaign_id"`
//	OrderID               sql.NullInt32   `json:"order_id"`
//	Quantity              sql.NullInt32   `json:"quantity"`
//	OrderlineStatus       sql.NullInt32   `json:"orderline_status"`
//	CurrencyCode          sql.NullString  `json:"currency_code"`
//	Price                 sql.NullFloat64 `json:"price"`
//	FinalPrice            sql.NullFloat64 `json:"final_price"`
//	LocalizedCurrencyCode sql.NullString  `json:"localized_currency_code"`
//	LocalizedFinalPrice   sql.NullFloat64 `json:"localized_final_price"`
//	IsTrial               sql.NullBool    `json:"is_trial"`
//}
//
//func (q *repository.Queries) CreateOrderlines(ctx context.Context, arg CreateOrderlinesParams) (CoreOrderline, error) {
//	row := q.db.QueryRowContext(ctx, createOrderlines,
//		arg.Name,
//		arg.OfferID,
//		arg.IsActive,
//		arg.IsFree,
//		arg.IsDiscount,
//		arg.UserID,
//		arg.CampaignID,
//		arg.OrderID,
//		arg.Quantity,
//		arg.OrderlineStatus,
//		arg.CurrencyCode,
//		arg.Price,
//		arg.FinalPrice,
//		arg.LocalizedCurrencyCode,
//		arg.LocalizedFinalPrice,
//		arg.IsTrial,
//	)
//	var i CoreOrderline
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.Name,
//		&i.OfferID,
//		&i.IsActive,
//		&i.IsFree,
//		&i.IsDiscount,
//		&i.UserID,
//		&i.CampaignID,
//		&i.OrderID,
//		&i.Quantity,
//		&i.OrderlineStatus,
//		&i.CurrencyCode,
//		&i.Price,
//		&i.FinalPrice,
//		&i.LocalizedCurrencyCode,
//		&i.LocalizedFinalPrice,
//		&i.IsTrial,
//	)
//	return i, err
//}
//
//const selectDetailOrderByID = `-- name: SelectDetailOrderByID :many
//Select o.created, o.modified, o.id, order_number, total_amount, final_amount, o.user_id, client_id, partner_id, o.is_active, point_reward, o.currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal, ol.created, ol.modified, ol.id, name, offer_id, ol.is_active, is_free, is_discount, ol.user_id, campaign_id, order_id, quantity, orderline_status, ol.currency_code, price, final_price, localized_currency_code, localized_final_price, is_trial
//from core_orders o
//         JOIN core_orderlines ol
//              ON o.id = ol.order_id
//WHERE o.id = $1
//`
//
//type SelectDetailOrderByIDRow struct {
//	Created               time.Time       `json:"created"`
//	Modified              time.Time       `json:"modified"`
//	ID                    int32           `json:"id"`
//	OrderNumber           int64           `json:"order_number"`
//	TotalAmount           float64         `json:"total_amount"`
//	FinalAmount           float64         `json:"final_amount"`
//	UserID                int32           `json:"user_id"`
//	ClientID              sql.NullInt32   `json:"client_id"`
//	PartnerID             sql.NullInt32   `json:"partner_id"`
//	IsActive              sql.NullBool    `json:"is_active"`
//	PointReward           sql.NullInt32   `json:"point_reward"`
//	CurrencyCode          string          `json:"currency_code"`
//	PaymentgatewayID      int32           `json:"paymentgateway_id"`
//	TierCode              sql.NullString  `json:"tier_code"`
//	PlatformID            sql.NullInt32   `json:"platform_id"`
//	TemporderID           sql.NullInt32   `json:"temporder_id"`
//	OrderStatus           sql.NullInt32   `json:"order_status"`
//	RemoteOrderNumber     sql.NullString  `json:"remote_order_number"`
//	IsRenewal             sql.NullBool    `json:"is_renewal"`
//	Created_2             time.Time       `json:"created_2"`
//	Modified_2            time.Time       `json:"modified_2"`
//	ID_2                  int32           `json:"id_2"`
//	Name                  sql.NullString  `json:"name"`
//	OfferID               int32           `json:"offer_id"`
//	IsActive_2            sql.NullBool    `json:"is_active_2"`
//	IsFree                sql.NullBool    `json:"is_free"`
//	IsDiscount            sql.NullBool    `json:"is_discount"`
//	UserID_2              sql.NullInt32   `json:"user_id_2"`
//	CampaignID            sql.NullInt32   `json:"campaign_id"`
//	OrderID               sql.NullInt32   `json:"order_id"`
//	Quantity              sql.NullInt32   `json:"quantity"`
//	OrderlineStatus       sql.NullInt32   `json:"orderline_status"`
//	CurrencyCode_2        sql.NullString  `json:"currency_code_2"`
//	Price                 sql.NullFloat64 `json:"price"`
//	FinalPrice            sql.NullFloat64 `json:"final_price"`
//	LocalizedCurrencyCode sql.NullString  `json:"localized_currency_code"`
//	LocalizedFinalPrice   sql.NullFloat64 `json:"localized_final_price"`
//	IsTrial               sql.NullBool    `json:"is_trial"`
//}
//
//func (q *repository.Queries) SelectDetailOrderByID(ctx context.Context, id int32) ([]SelectDetailOrderByIDRow, error) {
//	rows, err := q.db.QueryContext(ctx, selectDetailOrderByID, id)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	var items []SelectDetailOrderByIDRow
//	for rows.Next() {
//		var i SelectDetailOrderByIDRow
//		if err := rows.Scan(
//			&i.Created,
//			&i.Modified,
//			&i.ID,
//			&i.OrderNumber,
//			&i.TotalAmount,
//			&i.FinalAmount,
//			&i.UserID,
//			&i.ClientID,
//			&i.PartnerID,
//			&i.IsActive,
//			&i.PointReward,
//			&i.CurrencyCode,
//			&i.PaymentgatewayID,
//			&i.TierCode,
//			&i.PlatformID,
//			&i.TemporderID,
//			&i.OrderStatus,
//			&i.RemoteOrderNumber,
//			&i.IsRenewal,
//			&i.Created_2,
//			&i.Modified_2,
//			&i.ID_2,
//			&i.Name,
//			&i.OfferID,
//			&i.IsActive_2,
//			&i.IsFree,
//			&i.IsDiscount,
//			&i.UserID_2,
//			&i.CampaignID,
//			&i.OrderID,
//			&i.Quantity,
//			&i.OrderlineStatus,
//			&i.CurrencyCode_2,
//			&i.Price,
//			&i.FinalPrice,
//			&i.LocalizedCurrencyCode,
//			&i.LocalizedFinalPrice,
//			&i.IsTrial,
//		); err != nil {
//			return nil, err
//		}
//		items = append(items, i)
//	}
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//	return items, nil
//}
//
//const selectPendingOrder = `-- name: SelectOrder :many
//SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
//FROM core_orders WHERE order_status = 20001
//ORDER BY id DESC LIMIT $1
//OFFSET $2
//`
//
//type SelectPendingOrderParams struct {
//	Limit  int `json:"limit"`
//	Offset int `json:"offset"`
//}
//
//func (q *repository.Queries) SelectPendingOrder(ctx context.Context, arg SelectPendingOrderParams) ([]CoreOrder, error) {
//	rows, err := q.db.QueryContext(ctx, selectPendingOrder, arg.Limit, arg.Offset)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	var items []CoreOrder
//	for rows.Next() {
//		var i CoreOrder
//		if err := rows.Scan(
//			&i.Created,
//			&i.Modified,
//			&i.ID,
//			&i.OrderNumber,
//			&i.TotalAmount,
//			&i.FinalAmount,
//			&i.UserID,
//			&i.ClientID,
//			&i.PartnerID,
//			&i.IsActive,
//			&i.PointReward,
//			&i.CurrencyCode,
//			&i.PaymentgatewayID,
//			&i.TierCode,
//			&i.PlatformID,
//			&i.TemporderID,
//			&i.OrderStatus,
//			&i.RemoteOrderNumber,
//			&i.IsRenewal,
//		); err != nil {
//			return nil, err
//		}
//		items = append(items, i)
//	}
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//	return items, nil
//}
//const selectOrderByID = `-- name: SelectOrderByID :one
//SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
//FROM core_orders
//WHERE id = $1
//`
//
//func (q *repository.Queries) SelectOrderByID(ctx context.Context, id int32) (CoreOrder, error) {
//	row := q.db.QueryRowContext(ctx, selectOrderByID, id)
//	var i CoreOrder
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.OrderNumber,
//		&i.TotalAmount,
//		&i.FinalAmount,
//		&i.UserID,
//		&i.ClientID,
//		&i.PartnerID,
//		&i.IsActive,
//		&i.PointReward,
//		&i.CurrencyCode,
//		&i.PaymentgatewayID,
//		&i.TierCode,
//		&i.PlatformID,
//		&i.TemporderID,
//		&i.OrderStatus,
//		&i.RemoteOrderNumber,
//		&i.IsRenewal,
//	)
//	return i, err
//}
//
//const selectOrderByOrderNumber = `-- name: SelectOrderByOrderNumber :one
//SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
//FROM core_orders
//WHERE order_number = $1
//`
//
//func (q *repository.Queries) SelectOrderByOrderNumber(ctx context.Context, orderNumber int64) (CoreOrder, error) {
//	row := q.db.QueryRowContext(ctx, selectOrderByOrderNumber, orderNumber)
//	var i CoreOrder
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.OrderNumber,
//		&i.TotalAmount,
//		&i.FinalAmount,
//		&i.UserID,
//		&i.ClientID,
//		&i.PartnerID,
//		&i.IsActive,
//		&i.PointReward,
//		&i.CurrencyCode,
//		&i.PaymentgatewayID,
//		&i.TierCode,
//		&i.PlatformID,
//		&i.TemporderID,
//		&i.OrderStatus,
//		&i.RemoteOrderNumber,
//		&i.IsRenewal,
//	)
//	return i, err
//}
//
//const selectOrderByUserID = `-- name: SelectOrderByUserID :many
//SELECT created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
//FROM core_orders
//WHERE user_id = $1
//`
//
//func (q *repository.Queries) SelectOrderByUserID(ctx context.Context, userID int32) ([]CoreOrder, error) {
//	rows, err := q.db.QueryContext(ctx, selectOrderByUserID, userID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	var items []CoreOrder
//	for rows.Next() {
//		var i CoreOrder
//		if err := rows.Scan(
//			&i.Created,
//			&i.Modified,
//			&i.ID,
//			&i.OrderNumber,
//			&i.TotalAmount,
//			&i.FinalAmount,
//			&i.UserID,
//			&i.ClientID,
//			&i.PartnerID,
//			&i.IsActive,
//			&i.PointReward,
//			&i.CurrencyCode,
//			&i.PaymentgatewayID,
//			&i.TierCode,
//			&i.PlatformID,
//			&i.TemporderID,
//			&i.OrderStatus,
//			&i.RemoteOrderNumber,
//			&i.IsRenewal,
//		); err != nil {
//			return nil, err
//		}
//		items = append(items, i)
//	}
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//	return items, nil
//}
//
//const updateOrder = `-- name: UpdateOrder :one
//UPDATE core_orders
//SET modified            = $2,
//    order_status        = $3
//WHERE id = $1 RETURNING created, modified, id, order_number, total_amount, final_amount, user_id, client_id, partner_id, is_active, point_reward, currency_code, paymentgateway_id, tier_code, platform_id, temporder_id, order_status, remote_order_number, is_renewal
//`
//
//type UpdateOrderParams struct {
//	ID          int32         `json:"id"`
//	Modified    time.Time     `json:"modified"`
//	OrderStatus sql.NullInt32 `json:"order_status"`
//}
//
//func (q *repository.Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (CoreOrder, error) {
//	row := q.db.QueryRowContext(ctx, updateOrder,
//		arg.ID,
//		arg.Modified,
//		arg.OrderStatus,
//	)
//	var i CoreOrder
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.OrderNumber,
//		&i.TotalAmount,
//		&i.FinalAmount,
//		&i.UserID,
//		&i.ClientID,
//		&i.PartnerID,
//		&i.IsActive,
//		&i.PointReward,
//		&i.CurrencyCode,
//		&i.PaymentgatewayID,
//		&i.TierCode,
//		&i.PlatformID,
//		&i.TemporderID,
//		&i.OrderStatus,
//		&i.RemoteOrderNumber,
//		&i.IsRenewal,
//	)
//	return i, err
//}
//
//const updateOrderlines = `-- name: UpdateOrderlines :one
//UPDATE core_orderlines
//SET modified=$2,
//    orderline_status=$3
//WHERE order_id = $1 RETURNING created, modified, id, name, offer_id, is_active, is_free, is_discount, user_id, campaign_id, order_id, quantity, orderline_status, currency_code, price, final_price, localized_currency_code, localized_final_price, is_trial
//`
//
//type UpdateOrderlinesParams struct {
//	OrderID         int32         `json:"order_id"`
//	Modified        time.Time     `json:"modified"`
//	OrderlineStatus sql.NullInt32 `json:"orderline_status"`
//}
//
//func (q *repository.Queries) UpdateOrderlines(ctx context.Context, arg UpdateOrderlinesParams) (CoreOrderline, error) {
//	row := q.db.QueryRowContext(ctx, updateOrderlines,
//		arg.OrderID,
//		arg.Modified,
//		arg.OrderlineStatus,
//	)
//	var i CoreOrderline
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.Name,
//		&i.OfferID,
//		&i.IsActive,
//		&i.IsFree,
//		&i.IsDiscount,
//		&i.UserID,
//		&i.CampaignID,
//		&i.OrderID,
//		&i.Quantity,
//		&i.OrderlineStatus,
//		&i.CurrencyCode,
//		&i.Price,
//		&i.FinalPrice,
//		&i.LocalizedCurrencyCode,
//		&i.LocalizedFinalPrice,
//		&i.IsTrial,
//	)
//	return i, err
//}
//
//const selectOfferByID = `-- name: SelectOfferByID :one
//select created, modified, id, name, offer_status, sort_priority, is_active, offer_type_id, exclusive_clients, is_free, offer_code, item_code, price_usd, price_idr, price_point, discount_id, discount_tag, discount_name, discount_price_usd, discount_price_idr, discount_price_point, is_discount, image_highres, image_normal, vendor_price_usd, vendor_price_idr, vendor_price_point, long_name from core_offers where id = $1
//`
//
//func (q *repository.Queries) SelectOfferByID(ctx context.Context, id int32) (CoreOffer, error) {
//	row := q.db.QueryRowContext(ctx, selectOfferByID, id)
//	var i CoreOffer
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.Name,
//		&i.OfferStatus,
//		&i.SortPriority,
//		&i.IsActive,
//		&i.OfferTypeID,
//		pq.Array(&i.ExclusiveClients),
//		&i.IsFree,
//		&i.OfferCode,
//		&i.ItemCode,
//		&i.PriceUsd,
//		&i.PriceIdr,
//		&i.PricePoint,
//		pq.Array(&i.DiscountID),
//		&i.DiscountTag,
//		&i.DiscountName,
//		&i.DiscountPriceUsd,
//		&i.DiscountPriceIdr,
//		&i.DiscountPricePoint,
//		&i.IsDiscount,
//		&i.ImageHighres,
//		&i.ImageNormal,
//		&i.VendorPriceUsd,
//		&i.VendorPriceIdr,
//		&i.VendorPricePoint,
//		&i.LongName,
//	)
//	return i, err
//}
//
//const selectDiscountByID = `-- name: SelectDiscountByID :one
//select created, modified, id, name, tag_name, description, campaign_id, valid_to, valid_from, discount_rule, discount_type, discount_status, discount_schedule_type, is_active, discount_usd, discount_idr, discount_point, min_usd_order_price, max_usd_order_price, min_idr_order_price, max_idr_order_price, predefined_group, vendor_participation, partner_participation, sales_recognition, bin_codes, trial_time from core_discounts where id = $1`
//
//func (q *repository.Queries) SelectDiscountByID(ctx context.Context, id int32) (CoreDiscount, error) {
//	row := q.db.QueryRowContext(ctx, selectDiscountByID, id)
//	var i CoreDiscount
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.Name,
//		&i.TagName,
//		&i.Description,
//		&i.CampaignID,
//		&i.ValidTo,
//		&i.ValidFrom,
//		&i.DiscountRule,
//		&i.DiscountType,
//		&i.DiscountStatus,
//		&i.DiscountScheduleType,
//		&i.IsActive,
//		&i.DiscountUsd,
//		&i.DiscountIdr,
//		&i.DiscountPoint,
//		&i.MinUsdOrderPrice,
//		&i.MaxUsdOrderPrice,
//		&i.MinIdrOrderPrice,
//		&i.MaxIdrOrderPrice,
//		&i.PredefinedGroup,
//		&i.VendorParticipation,
//		&i.PartnerParticipation,
//		&i.SalesRecognition,
//		pq.Array(&i.BinCodes),
//		&i.TrialTime,
//	)
//	return i, err
//}
//
//const selectDiscountByIDs = `-- name: SelectDiscountByIDs :many
//SELECT created, modified, id, name, tag_name, description, campaign_id, valid_to, valid_from, discount_rule, discount_type,
//       discount_status, discount_schedule_type, is_active, discount_usd, discount_idr, discount_point, min_usd_order_price, max_usd_order_price, min_idr_order_price,
//       max_idr_order_price, predefined_group, vendor_participation, partner_participation, sales_recognition, bin_codes, trial_time
//FROM core_discounts WHERE id = any($1) ORDER BY valid_to DESC
//`
//
//func (q *repository.Queries) SelectDiscountByIDs(ctx context.Context, dollar_1 []int32) ([]CoreDiscount, error) {
//	rows, err := q.db.QueryContext(ctx, selectDiscountByIDs, pq.Array(dollar_1))
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	var items []CoreDiscount
//	for rows.Next() {
//		var i CoreDiscount
//		if err := rows.Scan(
//			&i.Created,
//			&i.Modified,
//			&i.ID,
//			&i.Name,
//			&i.TagName,
//			&i.Description,
//			&i.CampaignID,
//			&i.ValidTo,
//			&i.ValidFrom,
//			&i.DiscountRule,
//			&i.DiscountType,
//			&i.DiscountStatus,
//			&i.DiscountScheduleType,
//			&i.IsActive,
//			&i.DiscountUsd,
//			&i.DiscountIdr,
//			&i.DiscountPoint,
//			&i.MinUsdOrderPrice,
//			&i.MaxUsdOrderPrice,
//			&i.MinIdrOrderPrice,
//			&i.MaxIdrOrderPrice,
//			&i.PredefinedGroup,
//			&i.VendorParticipation,
//			&i.PartnerParticipation,
//			&i.SalesRecognition,
//			pq.Array(&i.BinCodes),
//			&i.TrialTime,
//		); err != nil {
//			return nil, err
//		}
//		items = append(items, i)
//	}
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//	return items, nil
//}
//
//const selectDiscountCodeByCode = `-- name: SelectDiscountCodeByCode :one
//select created, modified, id, code, max_uses, current_uses, is_active, discount_type, discount_id, is_for_new_user from core_discountcodes where code = $1`
//
//func (q *repository.Queries) SelectDiscountCodeByCode(ctx context.Context, code string) (CoreDiscountcode, error) {
//	row := q.db.QueryRowContext(ctx, selectDiscountCodeByCode, code)
//	var i CoreDiscountcode
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.Code,
//		&i.MaxUses,
//		&i.CurrentUses,
//		&i.IsActive,
//		&i.DiscountType,
//		&i.DiscountID,
//		&i.IsForNewUser,
//	)
//	return i, err
//}
//
//const updateDiscountCode = `-- name: UpdateDiscountCode :one
//UPDATE core_discountcodes
//SET
//    modified=$2, current_uses=$3
//WHERE id = $1 RETURNING created, modified, id, code, max_uses, current_uses, is_active, discount_type, discount_id, is_for_new_user
//`
//
//type UpdateDiscountCodeParams struct {
//	ID          int32         `json:"id"`
//	Modified    sql.NullTime  `json:"modified"`
//	CurrentUses sql.NullInt32 `json:"current_uses"`
//}
//
//func (q *repository.Queries) UpdateDiscountCode(ctx context.Context, arg UpdateDiscountCodeParams) (CoreDiscountcode, error) {
//	row := q.db.QueryRowContext(ctx, updateDiscountCode, arg.ID, arg.Modified, arg.CurrentUses)
//	var i CoreDiscountcode
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.Code,
//		&i.MaxUses,
//		&i.CurrentUses,
//		&i.IsActive,
//		&i.DiscountType,
//		&i.DiscountID,
//		&i.IsForNewUser,
//	)
//	return i, err
//}
//
//const selectItemOfSingleOffer = `-- name: SelectItemOfSingleOffer :one
//Select i.name, i.countries, i.item_distribution_country_group_id from core_offers_items oi
//join core_offers o on oi.offer_id = o.id
//join core_items i on oi.item_id = i.id
//where o.id = $1`
//
//type SelectItemOfSingleOfferRows struct {
//	Name                           string        `json:"name"`
//	Countries                      []uint8       `json:"countries"`
//	ItemDistributionCountryGroupID sql.NullInt32 `json:"item_distribution_country_group_id"`
//}
//
//func (q *repository.Queries) SelectItemOfSingleOffer(ctx context.Context, offerID int32) (SelectItemOfSingleOfferRows, error) {
//	row := q.db.QueryRowContext(ctx, selectItemOfSingleOffer, offerID)
//	var i SelectItemOfSingleOfferRows
//	err := row.Scan(
//		&i.Name,
//		&i.Countries,
//		&i.ItemDistributionCountryGroupID,
//	)
//	return i, err
//}
//
//const selectRestrictCountriesByOffer = `-- name: selectRestrictCountriesByItem :one
//Select o.id, i.name, i.item_distribution_country_group_id, dc.countries, dc.group_type from core_offers_items oi
//join core_offers o on oi.offer_id = o.id
//join core_items i on oi.item_id = i.id
//join core_distributioncountries dc on i.item_distribution_country_group_id = dc.id
//where o.id = $1`
//
//type SelectRestrictCountriesByOfferRows struct {
//	OfferID                        int32    `json:"offer_id"`
//	Name                           string   `json:"name"`
//	ItemDistributionCountryGroupID int32    `json:"item_distribution_country_group_id"`
//	Countries                      []string `json:"countries"`
//	RestrictType                   int32    `json:"restrict_type"`
//}
//
//func (q *repository.Queries) SelectRestrictCountriesByOffer(ctx context.Context, offerID int32) (SelectRestrictCountriesByOfferRows, error) {
//	row := q.db.QueryRowContext(ctx, selectRestrictCountriesByOffer, offerID)
//	var i SelectRestrictCountriesByOfferRows
//	err := row.Scan(
//		&i.OfferID,
//		&i.Name,
//		&i.ItemDistributionCountryGroupID,
//		pq.Array(&i.Countries),
//		&i.RestrictType,
//	)
//	return i, err
//}
//
//const selectPlatformOffer = `-- name: selectPlatformOffer :one
//select po.offer_id, o.name, po.tier_id, po.tier_code, po.currency, po.price_usd, po.price_idr, po.price_point,
//po.discount_tier_id, po.discount_tier_code, po.discount_tag, po.discount_name,
//po.discount_tier_price, po.discount_price_usd, po.discount_price_idr, po.discount_id
//from core_platforms_offers po
//join core_offers o on o.id = po.offer_id
//where po.offer_id = $1 and po.platform_id = $2`
//
//func (q *repository.Queries) SelectPlatformOffer(ctx context.Context, offerID int32, platformID int32) (CoreOffer, error) {
//	row := q.db.QueryRowContext(ctx, selectPlatformOffer, offerID, platformID)
//	var i CoreOffer
//	err := row.Scan(
//		&i.ID,
//		&i.Name,
//		&i.TierID,
//		&i.TierCode,
//		&i.Currency,
//		&i.PriceUsd,
//		&i.PriceIdr,
//		&i.PricePoint,
//		&i.DiscountTierID,
//		&i.DiscountTierCode,
//		&i.DiscountTag,
//		&i.DiscountName,
//		&i.DiscountTierPrice,
//		&i.DiscountPriceUsd,
//		&i.DiscountPriceIdr,
//		pq.Array(&i.DiscountID),
//	)
//	return i, err
//}
//
//const selectPaymentGateways = `-- name: selectPaymentGateways :one
// SELECT id, name, is_active, base_currency_id, minimal_amount, is_renewal, payment_group FROM core_paymentgateways WHERE id = $1`
//
//type SelectPaymentGatewaysRows struct {
//	PaymentGatewayID int32           `json:"payment_gateway_id"`
//	Name             string          `json:"name"`
//	IsActive         bool            `json:"is_active"`
//	BaseCurrencyID   int32           `json:"base_currency_id"`
//	MinimalAmount    sql.NullFloat64 `json:"minimal_amount"`
//	IsRenewal        bool            `json:"is_renewal"`
//	PaymentGroup     string          `json:"payment_group"`
//}
//
//func (q *repository.Queries) SelectPaymentGateways(ctx context.Context, paymentGatewayID int32) (SelectPaymentGatewaysRows, error) {
//	row := q.db.QueryRowContext(ctx, selectPaymentGateways, paymentGatewayID)
//	var i SelectPaymentGatewaysRows
//	err := row.Scan(
//		&i.PaymentGatewayID,
//		&i.Name,
//		&i.IsActive,
//		&i.BaseCurrencyID,
//		&i.MinimalAmount,
//		&i.IsRenewal,
//		&i.PaymentGroup,
//	)
//	return i, err
//}
//
//const selectUser = `-- name: selectUser :one
// SELECT id, username, email FROM cas_users WHERE id = $1`
//
//type SelectUserRows struct {
//	UserID   int32  `json:"user_id"`
//	UserName string `json:"user_name"`
//	Email    string `json:"email"`
//}
//
//func (q *repository.Queries) SelectUser(ctx context.Context, userID int32) (SelectUserRows, error) {
//	row := q.db.QueryRowContext(ctx, selectUser, userID)
//	var i SelectUserRows
//	err := row.Scan(
//		&i.UserID,
//		&i.UserName,
//		&i.Email,
//	)
//	return i, err
//}
//
//const createPayment = `-- name: CreatePayment :one
//INSERT INTO core_payments(created, modified, order_id, user_id, paymentgateway_id, currency_code, amount, payment_status, is_active, is_test_payment, payment_datetime, financial_archive_date, is_trial)
//VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
//RETURNING created, modified, id, order_id, user_id, paymentgateway_id, currency_code, amount, payment_status, is_active, is_test_payment, payment_datetime, financial_archive_date, is_trial, merchant_params
//`
//
//type CreatePaymentParams struct {
//	Created              time.Time `json:"created"`
//	Modified             time.Time `json:"modified"`
//	OrderID              int32     `json:"order_id"`
//	UserID               int32     `json:"user_id"`
//	PaymentgatewayID     int32     `json:"paymentgateway_id"`
//	CurrencyCode         string    `json:"currency_code"`
//	Amount               float64   `json:"amount"`
//	PaymentStatus        int32     `json:"payment_status"`
//	IsActive             bool      `json:"is_active"`
//	IsTestPayment        bool      `json:"is_test_payment"`
//	PaymentDatetime      time.Time `json:"payment_datetime"`
//	FinancialArchiveDate time.Time `json:"financial_archive_date"`
//	IsTrial              bool      `json:"is_trial"`
//	//MerchantParams       json.RawMessage `json:"merchant_params"`
//}
//
//func (q *repository.Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (CorePayment, error) {
//	row := q.db.QueryRowContext(ctx, createPayment,
//		arg.Created,
//		arg.Modified,
//		arg.OrderID,
//		arg.UserID,
//		arg.PaymentgatewayID,
//		arg.CurrencyCode,
//		arg.Amount,
//		arg.PaymentStatus,
//		arg.IsActive,
//		arg.IsTestPayment,
//		arg.PaymentDatetime,
//		arg.FinancialArchiveDate,
//		arg.IsTrial,
//		//arg.MerchantParams,
//	)
//	var i CorePayment
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.OrderID,
//		&i.UserID,
//		&i.PaymentgatewayID,
//		&i.CurrencyCode,
//		&i.Amount,
//		&i.PaymentStatus,
//		&i.IsActive,
//		&i.IsTestPayment,
//		&i.PaymentDatetime,
//		&i.FinancialArchiveDate,
//		&i.IsTrial,
//		&i.MerchantParams,
//	)
//	return i, err
//}
//
//const updatePayment = `-- name: UpdateOrder :one
//UPDATE core_orders
//SET modified            = $2,
//    order_status        = $3
//WHERE order_id = $1
//RETURNING created, modified, id, order_id, user_id, paymentgateway_id, currency_code, amount, payment_status, is_active, is_test_payment, payment_datetime, financial_archive_date, is_trial, merchant_params
//`
//
//type UpdatePaymentParams struct {
//	OrderID       int32         `json:"id"`
//	Modified      time.Time     `json:"modified"`
//	PaymentStatus sql.NullInt32 `json:"payment_status"`
//}
//
//func (q *repository.Queries) UpdatePaymentByOrder(ctx context.Context, arg UpdatePaymentParams) (CorePayment, error) {
//	row := q.db.QueryRowContext(ctx, updatePayment,
//		arg.OrderID,
//		arg.Modified,
//		arg.PaymentStatus,
//	)
//	var i CorePayment
//	err := row.Scan(
//		&i.Created,
//		&i.Modified,
//		&i.ID,
//		&i.OrderID,
//		&i.UserID,
//		&i.PaymentgatewayID,
//		&i.CurrencyCode,
//		&i.Amount,
//		&i.PaymentStatus,
//		&i.IsActive,
//		&i.IsTestPayment,
//		&i.PaymentDatetime,
//		&i.FinancialArchiveDate,
//		&i.IsTrial,
//		&i.MerchantParams,
//	)
//	return i, err
//}