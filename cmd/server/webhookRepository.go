package main

import (
	"database/sql"
	"time"
)

type webhookRepository interface {
	recordWebhook(watchID string, attemptedAt time.Time, status int) error
	close() error
}

////////////////////////////////////////////////////////
// generic implementation for most
// databases (SQLite, MySQL)
////////////////////////////////////////////////////////
type genericSQLWebhookRepository struct {
	db *sql.DB
}

func (r *genericSQLWebhookRepository) close() error {
	return r.db.Close()
}

func (r *genericSQLWebhookRepository) recordWebhook(watchID string, attemptedAt time.Time, status int) error {
	query := `insert into webhook_stats (watch_id, attempted_at, status) values (?, ?, ?);`
	stmt, err := r.db.Prepare(query)
	return insertWebhook(watchID, attemptedAt, status, err, stmt)
}

////////////////////////////////////////////////////////
// postgres implementation
////////////////////////////////////////////////////////
type postgresWebhookRepository struct {
	db *sql.DB
}

func (r *postgresWebhookRepository) close() error {
	return r.db.Close()
}

func (r *postgresWebhookRepository) recordWebhook(watchID string, attemptedAt time.Time, status int) error {
	query := `insert into webhook_stats (watch_id, attempted_at, status) values ($1, $2, $3);`
	stmt, err := r.db.Prepare(query)
	return insertWebhook(watchID, attemptedAt, status, err, stmt)
}

// Common method across all databases.

// This function will return a webhookRepository for a specific database that requires specific handling of
// queries such as Postgres and Oracle. Other databases such as SQLite and MySQL will get a generic repository.
func getWebhookRepo(dbType string, db *sql.DB) webhookRepository {
	switch dbType {
	case "postgres":
		return &postgresWebhookRepository{db}
	default:
		return &genericSQLWebhookRepository{db}
	}
}

func insertWebhook(watchID string, attemptedAt time.Time, status int, err error, stmt *sql.Stmt) error {
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(watchID, attemptedAt, status)
	return err
}
