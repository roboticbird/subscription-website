package tasks

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func OrderUpdate(database *sql.DB) {
	log.Printf("ORDER_UPDATE: Starting nightly update")
	type Sub struct {
		ID         int
		ProductID  int
		UserID     int
		StartDate  time.Time
		Status     string
		NextStatus string
	}

	// all subscriptions that need to be renewed
	rows, err := database.Query("SELECT sub.id, sub.product_id, sub.user_id, " +
		"sub.start_date, sub.status, sub.nextStatus " +
		"FROM subscriptions as sub " +
		"LEFT JOIN products as prd on prd.id = sub.product_id " +
		"WHERE sub.start_date < DATE_ADD(NOW(), INTERVAL prd.duration DAY) " +
		"AND sub.status != 'CANCELLED' " +
		"AND sub.status != 'EXPIRED'")

	if err != nil {
		log.Printf("ORDER_UPDATE: Error reading database: %s", err)
		return
	}
	defer rows.Close()
	subscriptions := []Sub{}
	for rows.Next() {
		sub := Sub{}

		err := rows.Scan(&sub.ID, &sub.ProductID, &sub.UserID,
			&sub.StartDate, &sub.Status, &sub.NextStatus)
		if err != nil {
			log.Printf("ORDER_UPDATE: Error reading database: %s", err)
			return
		}
		subscriptions = append(subscriptions, sub)
	}
	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		log.Printf("ORDER_UPDATE: Error reading database: %s", err)
		return
	}

	query := ""
	for _, sub := range subscriptions {
		newStatus := ""
		if sub.NextStatus == "QUEUED" {
			newStatus = "ACTIVE"
		} else {
			newStatus = "PAUSED"
		}
		query += fmt.Sprintf("UPDATE subscriptions SET status = 'EXPIRED' WHERE id = %d;", sub.ID)
		query += fmt.Sprintf("INSERT INTO subscriptions (user_id, product_id, start_date, "+
			"status, nextStatus) VALUES(%d, %d, NOW(), '%s', '%s');",
			sub.UserID, sub.ProductID, newStatus, sub.NextStatus)
	}
	_, err = database.Exec(query)
	if err != nil {
		log.Printf("ORDER_UPDATE: Error writing to database: %s\n", err)
		return
	}
	return
}
