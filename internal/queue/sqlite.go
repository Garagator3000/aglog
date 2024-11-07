package queue

import (
	"aglog/internal/log"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteQueue struct {
	db     *sql.DB
	logger log.Logger
}

const (
	CreateTable = `CREATE TABLE IF NOT EXISTS queue (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	message TEXT,
    	timestamp INTEGER,
    	handled INTEGER NOT NULL DEFAULT 0);`

	SaveMessage       = `INSERT INTO queue (message, timestamp) VALUES (?, ?)`
	ReadOldestMessage = `SELECT id, message, timestamp FROM queue WHERE handled = 0 ORDER BY timestamp ASC LIMIT 1;`
	MarkAsHandled     = `UPDATE queue SET handled = 1 WHERE id = ?`
)

func NewSqliteQueue(dbPath string, logger log.Logger) *SqliteQueue {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(fmt.Errorf("failed to open sqlite db: %w", err))
	}

	_, err = db.Exec(CreateTable)
	if err != nil {
		panic(fmt.Errorf("failed to create sqlite table: %w", err))
	}

	return &SqliteQueue{
		db:     db,
		logger: logger,
	}
}

func (q *SqliteQueue) Close() {
	err := q.db.Close()
	if err != nil {
		q.logger.Error("failed to close sqlite db", log.Error(err))
	}
}

func (q *SqliteQueue) Enqueue(message string) {
	_, err := q.db.Exec(SaveMessage, message, time.Now().UnixNano())
	if err != nil {
		q.logger.Error("failed to save message to sqlite", log.Error(err))
	}

	q.logger.Debug("enqueued message", log.String("message", message))
}

func (q *SqliteQueue) Dequeue() (int64, string) {
	var id int
	var message string
	var timestamp int64

	row := q.db.QueryRow(ReadOldestMessage)
	err := row.Scan(&id, &message, &timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ""
		}
		q.logger.Error("failed to read oldest message from sqlite", log.Error(err))
	}

	_, err = q.db.Exec(MarkAsHandled, id)
	if err != nil {
		q.logger.Error("failed to mark message as handled in sqlite", log.Error(err))
	}

	q.logger.Debug("dequeued message", log.String("message", message))
	return timestamp, message
}
