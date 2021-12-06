package transactions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"scoop-order/internal/schemas"
	"scoop-order/repository"
)

// Transaction defines all function to execute db queries and transaction
type Transaction interface {
	repository.Querier
	PricingTx(ctx context.Context, arg schemas.CheckPricingRequest) (schemas.CheckPricingResponse, error)
	CheckoutTx(ctx context.Context, arg schemas.CheckoutTxParams) (schemas.CheckoutTxResult, error)
	PaymentTx(ctx context.Context, checkout schemas.CheckoutTxResult, paymentStatus int32) error
	CompletePaymentTx(ctx context.Context, arg schemas.PaymentTxParams) (string, error)
}

type SQLTransaction struct {
	*repository.Queries
	db          *sql.DB //require to create a new db transaction
	clientRedis *redis.Client
}

// NewTransaction : create new Transaction
func NewTransaction(db *sql.DB, clientRedis *redis.Client) Transaction {
	return &SQLTransaction{
		db:          db,
		Queries:     repository.New(db),
		clientRedis: clientRedis,
	}
}

// execDBTx is unexported function, only for each specific function
func (transaction *SQLTransaction) execDBTx(ctx context.Context, fn func(*repository.Queries) error) error {
	tx, err := transaction.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := repository.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

