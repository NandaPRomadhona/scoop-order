package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"scoop-order/internal/databases"
)

const selectOfferByID = `-- name: SelectOfferByID :one
select created, modified, id, name, offer_status, sort_priority, is_active, offer_type_id, exclusive_clients, is_free, offer_code, item_code, price_usd, price_idr, price_point, discount_id, discount_tag, discount_name, discount_price_usd, discount_price_idr, discount_price_point, is_discount, image_highres, image_normal, vendor_price_usd, vendor_price_idr, vendor_price_point, long_name from core_offers where id = $1
`

func (q *Queries) SelectOfferByID(ctx context.Context, id int32) (databases.CoreOffer, error) {
	row := q.db.QueryRowContext(ctx, selectOfferByID, id)
	var i databases.CoreOffer
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.Name,
		&i.OfferStatus,
		&i.SortPriority,
		&i.IsActive,
		&i.OfferTypeID,
		pq.Array(&i.ExclusiveClients),
		&i.IsFree,
		&i.OfferCode,
		&i.ItemCode,
		&i.PriceUsd,
		&i.PriceIdr,
		&i.PricePoint,
		pq.Array(&i.DiscountID),
		&i.DiscountTag,
		&i.DiscountName,
		&i.DiscountPriceUsd,
		&i.DiscountPriceIdr,
		&i.DiscountPricePoint,
		&i.IsDiscount,
		&i.ImageHighres,
		&i.ImageNormal,
		&i.VendorPriceUsd,
		&i.VendorPriceIdr,
		&i.VendorPricePoint,
		&i.LongName,
	)
	return i, err
}

const selectItemOfSingleOffer = `-- name: SelectItemOfSingleOffer :one
Select i.name, i.countries, i.item_distribution_country_group_id from core_offers_items oi 
join core_offers o on oi.offer_id = o.id 
join core_items i on oi.item_id = i.id
where o.id = $1`

type SelectItemOfSingleOfferRows struct {
	Name                           string        `json:"name"`
	Countries                      []uint8       `json:"countries"`
	ItemDistributionCountryGroupID sql.NullInt32 `json:"item_distribution_country_group_id"`
}

func (q *Queries) SelectItemOfSingleOffer(ctx context.Context, offerID int32) (SelectItemOfSingleOfferRows, error) {
	row := q.db.QueryRowContext(ctx, selectItemOfSingleOffer, offerID)
	var i SelectItemOfSingleOfferRows
	err := row.Scan(
		&i.Name,
		&i.Countries,
		&i.ItemDistributionCountryGroupID,
	)
	return i, err
}

const selectRestrictCountriesByOffer = `-- name: selectRestrictCountriesByItem :one
Select o.id, i.name, i.item_distribution_country_group_id, dc.countries, dc.group_type as restrict_type from core_offers_items oi 
join core_offers o on oi.offer_id = o.id 
join core_items i on oi.item_id = i.id
join core_distributioncountries dc on i.item_distribution_country_group_id = dc.id
where o.id = $1`

type SelectRestrictCountriesByOfferRows struct {
	OfferID                        int32    `json:"offer_id"`
	Name                           string   `json:"name"`
	ItemDistributionCountryGroupID int32    `json:"item_distribution_country_group_id"`
	Countries                      []string `json:"countries"`
	RestrictType                   int32    `json:"restrict_type"`
}

func (q *Queries) SelectRestrictCountriesByOffer(ctx context.Context, offerID int32) (SelectRestrictCountriesByOfferRows, error) {
	row := q.db.QueryRowContext(ctx, selectRestrictCountriesByOffer, offerID)
	var i SelectRestrictCountriesByOfferRows
	err := row.Scan(
		&i.OfferID,
		&i.Name,
		&i.ItemDistributionCountryGroupID,
		pq.Array(&i.Countries),
		&i.RestrictType,
	)
	return i, err
}

const selectPlatformOffer = `-- name: selectPlatformOffer :one
select po.offer_id, o.name, po.tier_id, po.tier_code, po.currency, po.price_usd, po.price_idr, po.price_point,
po.discount_tier_id, po.discount_tier_code, po.discount_tag, po.discount_name,
po.discount_tier_price, po.discount_price_usd, po.discount_price_idr, po.discount_id 
from core_platforms_offers po 
join core_offers o on o.id = po.offer_id
where po.offer_id = $1 and po.platform_id = $2`

func (q *Queries) SelectPlatformOffer(ctx context.Context, offerID int32, platformID int32) (databases.CoreOffer, error) {
	row := q.db.QueryRowContext(ctx, selectPlatformOffer, offerID, platformID)
	var i databases.CoreOffer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.TierID,
		&i.TierCode,
		&i.Currency,
		&i.PriceUsd,
		&i.PriceIdr,
		&i.PricePoint,
		&i.DiscountTierID,
		&i.DiscountTierCode,
		&i.DiscountTag,
		&i.DiscountName,
		&i.DiscountTierPrice,
		&i.DiscountPriceUsd,
		&i.DiscountPriceIdr,
		pq.Array(&i.DiscountID),
	)
	return i, err
}