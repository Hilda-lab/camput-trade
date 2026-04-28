package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func PurchaseItem(db *sql.DB, orderID, itemID, buyerID string) error {
	if db == nil {
		return errors.New("database not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var status int
	row := tx.QueryRowContext(ctx, "SELECT status FROM item WHERE item_id = ? FOR UPDATE", itemID)
	if err := row.Scan(&status); err != nil {
		return fmt.Errorf("item not found: %w", err)
	}
	if status == 1 {
		return errors.New("item already sold")
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO orders (order_id, buyer_id, item_id, order_date) VALUES (?, ?, ?, NOW())", orderID, buyerID, itemID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "UPDATE item SET status = 1 WHERE item_id = ?", itemID); err != nil {
		return err
	}

	return tx.Commit()
}
