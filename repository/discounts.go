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
