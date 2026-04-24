package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PurchaseItem(pool *pgxpool.Pool, orderID, itemID, buyerID string) error {
	if pool == nil {
		return errors.New("database not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var status int
	row := tx.QueryRow(ctx, "SELECT status FROM item WHERE item_id = $1 FOR UPDATE", itemID)
	if err := row.Scan(&status); err != nil {
		return fmt.Errorf("item not found: %w", err)
	}
	if status == 1 {
		return errors.New("item already sold")
	}

	if _, err := tx.Exec(ctx, "INSERT INTO orders (order_id, buyer_id, item_id, order_date) VALUES ($1, $2, $3, NOW())", orderID, buyerID, itemID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "UPDATE item SET status = 1 WHERE item_id = $1", itemID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
