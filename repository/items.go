package repository

import (
	"context"
	"scoop-order/internal/databases"
)

const selectItemByID = `-- name: selectItemByID :one
SELECT id, name, edition_code, is_featured, is_extra, is_active, brand_id, item_type, content_type, item_status, meta, parent_item_id, created, modified FROM core_items WHERE id = $1
`

func (q *Queries) SelectItemByID(ctx context.Context, id int32) (databases.CoreItem, error) {
	row := q.db.QueryRowContext(ctx, selectItemByID, id)
	var item databases.CoreItem
	err := row.Scan(
		&item.Id, &item.Name, &item.EditionCode, &item.IsFeatured, &item.IsExtra, &item.IsActive,
		&item.BrandId, &item.ItemType, &item.ContentType, &item.ItemStatus, &item.Meta,
		&item.ParentItemId, &item.Created, &item.Modified,
	)
	return item, err
}

const selectItemByBrandID = `-- name: selectItemByBrandID :many
SELECT id, name, edition_code, is_featured, is_extra, is_active, brand_id, item_type, content_type, item_status, meta, parent_item_id, created, modified FROM core_items WHERE brand_id = $1
`

func (q *Queries) SelectItemByBrandID(ctx context.Context, brandID int32) ([]databases.CoreItem, error) {
	rows, err := q.db.QueryContext(ctx, selectItemByBrandID, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []databases.CoreItem
	for rows.Next() {
		var item databases.CoreItem
		if err := rows.Scan(
				&item.Id, &item.Name, &item.EditionCode, &item.IsFeatured, &item.IsExtra, &item.IsActive,
				&item.BrandId, &item.ItemType, &item.ContentType, &item.ItemStatus, &item.Meta,
				&item.ParentItemId, &item.Created, &item.Modified,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectItemsByOfferID = `-- name: selectItemsByOfferID :many
SELECT id, name, edition_code, is_featured, is_extra, is_active, brand_id, item_type, content_type, item_status, meta, parent_item_id, created, modified 
FROM core_items WHERE id = any(select item_id from core_offers_items where offer_id = $1)
`

func (q *Queries) SelectItemByOfferID(ctx context.Context, offerID int32) ([]databases.CoreItem, error) {
	rows, err := q.db.QueryContext(ctx, selectItemsByOfferID, offerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []databases.CoreItem
	for rows.Next() {
		var item databases.CoreItem
		if err := rows.Scan(
			&item.Id, &item.Name, &item.EditionCode, &item.IsFeatured, &item.IsExtra, &item.IsActive,
			&item.BrandId, &item.ItemType, &item.ContentType, &item.ItemStatus, &item.Meta,
			&item.ParentItemId, &item.Created, &item.Modified,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectItemsBrandByOfferID = `-- name: selectItemsBrandByOfferID :many
SELECT id, name, edition_code, is_featured, is_extra, is_active, brand_id, item_type, content_type, item_status, meta, parent_item_id, created, modified 
FROM core_items WHERE brand_id = any(select brand_id from core_offers_brands where offer_id = $1)
`

func (q *Queries) SelectItemBrandByOfferID(ctx context.Context, offerID int32) ([]databases.CoreItem, error) {
	rows, err := q.db.QueryContext(ctx, selectItemsBrandByOfferID, offerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []databases.CoreItem
	for rows.Next() {
		var item databases.CoreItem
		if err := rows.Scan(
			&item.Id, &item.Name, &item.EditionCode, &item.IsFeatured, &item.IsExtra, &item.IsActive,
			&item.BrandId, &item.ItemType, &item.ContentType, &item.ItemStatus, &item.Meta,
			&item.ParentItemId, &item.Created, &item.Modified,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

