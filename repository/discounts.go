package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"scoop-order/internal/databases"
)

const selectDiscountByID = `-- name: SelectDiscountByID :one
select created, modified, id, name, tag_name, description, campaign_id, valid_to, valid_from, discount_rule, discount_type, discount_status, discount_schedule_type, is_active, discount_usd, discount_idr, discount_point, min_usd_order_price, max_usd_order_price, min_idr_order_price, max_idr_order_price, predefined_group, vendor_participation, partner_participation, sales_recognition, bin_codes from core_discounts where id = $1`

func (q *Queries) SelectDiscountByID(ctx context.Context, id int32) (databases.CoreDiscount, error) {
	row := q.db.QueryRowContext(ctx, selectDiscountByID, id)
	var i databases.CoreDiscount
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.Name,
		&i.TagName,
		&i.Description,
		&i.CampaignID,
		&i.ValidTo,
		&i.ValidFrom,
		&i.DiscountRule,
		&i.DiscountType,
		&i.DiscountStatus,
		&i.DiscountScheduleType,
		&i.IsActive,
		&i.DiscountUsd,
		&i.DiscountIdr,
		&i.DiscountPoint,
		&i.MinUsdOrderPrice,
		&i.MaxUsdOrderPrice,
		&i.MinIdrOrderPrice,
		&i.MaxIdrOrderPrice,
		&i.PredefinedGroup,
		&i.VendorParticipation,
		&i.PartnerParticipation,
		&i.SalesRecognition,
		pq.Array(&i.BinCodes),
		//&i.TrialTime,
	)
	return i, err
}

const selectDiscountByIDs = `-- name: SelectDiscountByIDs :many
SELECT created, modified, id, name, tag_name, description, campaign_id, valid_to, valid_from, discount_rule, discount_type,
       discount_status, discount_schedule_type, is_active, discount_usd, discount_idr, discount_point, min_usd_order_price, max_usd_order_price, min_idr_order_price,
       max_idr_order_price, predefined_group, vendor_participation, partner_participation, sales_recognition, bin_codes
FROM core_discounts WHERE id = any($1) ORDER BY valid_to DESC
`

func (q *Queries) SelectDiscountByIDs(ctx context.Context, dollar_1 []int32) ([]databases.CoreDiscount, error) {
	rows, err := q.db.QueryContext(ctx, selectDiscountByIDs, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var items []databases.CoreDiscount
	for rows.Next() {
		var i databases.CoreDiscount
		if err := rows.Scan(
			&i.Created,
			&i.Modified,
			&i.ID,
			&i.Name,
			&i.TagName,
			&i.Description,
			&i.CampaignID,
			&i.ValidTo,
			&i.ValidFrom,
			&i.DiscountRule,
			&i.DiscountType,
			&i.DiscountStatus,
			&i.DiscountScheduleType,
			&i.IsActive,
			&i.DiscountUsd,
			&i.DiscountIdr,
			&i.DiscountPoint,
			&i.MinUsdOrderPrice,
			&i.MaxUsdOrderPrice,
			&i.MinIdrOrderPrice,
			&i.MaxIdrOrderPrice,
			&i.PredefinedGroup,
			&i.VendorParticipation,
			&i.PartnerParticipation,
			&i.SalesRecognition,
			pq.Array(&i.BinCodes),
			//&i.TrialTime,
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

const selectDiscountCodeByCode = `-- name: SelectDiscountCodeByCode :one
select created, modified, id, code, max_uses, current_uses, is_active, discount_type, discount_id, is_for_new_user from core_discountcodes where code = $1`

func (q *Queries) SelectDiscountCodeByCode(ctx context.Context, code string) (databases.CoreDiscountcode, error) {
	row := q.db.QueryRowContext(ctx, selectDiscountCodeByCode, code)
	var i databases.CoreDiscountcode
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.Code,
		&i.MaxUses,
		&i.CurrentUses,
		&i.IsActive,
		&i.DiscountType,
		&i.DiscountID,
		&i.IsForNewUser,
	)
	return i, err
}

const updateDiscountCode = `-- name: UpdateDiscountCode :one
UPDATE core_discountcodes
SET
    modified=$2, current_uses=$3
WHERE id = $1 RETURNING created, modified, id, code, max_uses, current_uses, is_active, discount_type, discount_id, is_for_new_user
`

type UpdateDiscountCodeParams struct {
	ID          int32         `json:"id"`
	Modified    sql.NullTime  `json:"modified"`
	CurrentUses sql.NullInt32 `json:"current_uses"`
}

func (q *Queries) UpdateDiscountCode(ctx context.Context, arg UpdateDiscountCodeParams) (databases.CoreDiscountcode, error) {
	row := q.db.QueryRowContext(ctx, updateDiscountCode, arg.ID, arg.Modified, arg.CurrentUses)
	var i databases.CoreDiscountcode
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.Code,
		&i.MaxUses,
		&i.CurrentUses,
		&i.IsActive,
		&i.DiscountType,
		&i.DiscountID,
		&i.IsForNewUser,
	)
	return i, err
}

const selectDiscountByPredefinedGroups = `-- name: selectDiscountByPredefinedGroups :one
SELECT d.created, d.modified, d.id, d.name, d.tag_name, d.description, d.campaign_id, d.valid_to, d.valid_from, d.discount_rule, d.discount_type,
       d.discount_status, d.discount_schedule_type, d.is_active, d.discount_usd, d.discount_idr, d.discount_point, d.min_usd_order_price, d.max_usd_order_price, d.min_idr_order_price,
       d.max_idr_order_price, d.predefined_group, d.vendor_participation, d.partner_participation, d.sales_recognition, d.bin_codes
FROM core_discounts d 
JOIN core_discounts_platforms dp ON dp.discount_id = d.id
JOIN core_discounts_paymentgateways dpg ON dpg.discount_id = d.id 
WHERE d.predefined_group = any($1) AND d.discount_type = 1 AND d.is_active = True AND platform_id = $2 AND paymentgateway_id = $3 AND valid_to >= $4 ORDER BY id DESC LIMIT 1
`

func (q *Queries) SelectDiscountByPredefinedGroups(ctx context.Context, predefinedGroups []int32, platformID int32, paymentgatewayID int32, today string) ([]databases.CoreDiscount, error) {
	rows, err := q.db.QueryContext(ctx, selectDiscountByPredefinedGroups, pq.Array(predefinedGroups), platformID, paymentgatewayID, today)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var items []databases.CoreDiscount
	for rows.Next() {
		var i databases.CoreDiscount
		if err := rows.Scan(
			&i.Created,
			&i.Modified,
			&i.ID,
			&i.Name,
			&i.TagName,
			&i.Description,
			&i.CampaignID,
			&i.ValidTo,
			&i.ValidFrom,
			&i.DiscountRule,
			&i.DiscountType,
			&i.DiscountStatus,
			&i.DiscountScheduleType,
			&i.IsActive,
			&i.DiscountUsd,
			&i.DiscountIdr,
			&i.DiscountPoint,
			&i.MinUsdOrderPrice,
			&i.MaxUsdOrderPrice,
			&i.MinIdrOrderPrice,
			&i.MaxIdrOrderPrice,
			&i.PredefinedGroup,
			&i.VendorParticipation,
			&i.PartnerParticipation,
			&i.SalesRecognition,
			pq.Array(&i.BinCodes),
			//&i.TrialTime,
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

const selectAllowedPG = `-- name: selectAllowedPG :one
select discount_id, paymentgateway_id from core_discounts_paymentgateways where discount_id = $1 and paymentgateway_id = $2`

func (q *Queries) SelectAllowedPG(ctx context.Context, discountId int32, paymentgatewayId int32) (int32, error) {
	row := q.db.QueryRowContext(ctx, selectAllowedPG, discountId, paymentgatewayId)
	type coreDiscountsPaymentgateways struct {
		discountID       int32 `json:"discount_id"`
		paymentGatewayID int32 `json:"payment_gateway_id"`
	}
	var i coreDiscountsPaymentgateways
	err := row.Scan(
		&i.discountID,
		&i.paymentGatewayID,
	)
	return i.paymentGatewayID, err
}

const selectAllowedPlatform = `-- name: selectAllowedPlatform :one
select discount_id, platform_id from core_discounts_platforms where discount_id = $1 and platform_id = $2`

func (q *Queries) SelectAllowedPlatform(ctx context.Context, discountId int32, platformId int32) (int32, error) {
	row := q.db.QueryRowContext(ctx, selectAllowedPlatform, discountId, platformId)
	type coreDiscountsPlatforms struct {
		discountID int32 `json:"discount_id"`
		platformId int32 `json:"platform_id"`
	}
	var i coreDiscountsPlatforms
	err := row.Scan(
		&i.discountID,
		&i.platformId,
	)
	return i.platformId, err
}

const selectAllowedOffers = `-- name: selectAllowedOffers :many
SELECT discount_id, offer_id FROM core_discounts_offers WHERE order_id = $1`

func (q *Queries) SelectAllowedOffers(ctx context.Context, discountID int32) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, selectAllowedOffers, discountID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	type coreDiscountsOffers struct {
		discountID int32 `json:"discount_id"`
		offerID    int32 `json:"offer_id"`
	}
	var items []int32
	for rows.Next() {
		var i coreDiscountsOffers
		if err := rows.Scan(
			&i.discountID,
			&i.offerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i.offerID)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const selectDiscountsPaymentGateways = `-- name: selectDiscountsPaymentGateways :one
SELECT d.created, d.modified, d.id, d.name, d.tag_name, d.description, d.campaign_id, d.valid_to, d.valid_from, d.discount_rule, d.discount_type,
       d.discount_status, d.discount_schedule_type, d.is_active, d.discount_usd, d.discount_idr, d.discount_point, d.min_usd_order_price, d.max_usd_order_price, d.min_idr_order_price,
       d.max_idr_order_price, d.predefined_group, d.vendor_participation, d.partner_participation, d.sales_recognition, d.bin_codes
FROM core_discounts d 
JOIN core_discounts_paymentgateways dpg ON dpg.discount_id = d.id 
WHERE d.discount_type = 3 AND d.is_active = True AND paymentgateway_id = $1 AND valid_to >= $2 ORDER BY id DESC LIMIT 1
`

func (q *Queries) SelectDiscountsPaymentGateways(ctx context.Context, paymentgatewayID int32, today string) (databases.CoreDiscount, error) {
	row := q.db.QueryRowContext(ctx, selectDiscountsPaymentGateways, paymentgatewayID, today)
	var i databases.CoreDiscount
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.Name,
		&i.TagName,
		&i.Description,
		&i.CampaignID,
		&i.ValidTo,
		&i.ValidFrom,
		&i.DiscountRule,
		&i.DiscountType,
		&i.DiscountStatus,
		&i.DiscountScheduleType,
		&i.IsActive,
		&i.DiscountUsd,
		&i.DiscountIdr,
		&i.DiscountPoint,
		&i.MinUsdOrderPrice,
		&i.MaxUsdOrderPrice,
		&i.MinIdrOrderPrice,
		&i.MaxIdrOrderPrice,
		&i.PredefinedGroup,
		&i.VendorParticipation,
		&i.PartnerParticipation,
		&i.SalesRecognition,
		pq.Array(&i.BinCodes),
		//&i.TrialTime,
	)
	return i, err
}
