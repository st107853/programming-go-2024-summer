package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Anonymously import the driver package
)

type PostgresTransactionLogger struct {
	events chan<- Event // Write-only channel for sending events
	errors <-chan error // Read-only channels for receving errors
	db     *sql.DB      // The database access interface
}

type PostgresDbParams struct {
	dbName   string
	host     string
	user     string
	password string
}

func (l *PostgresTransactionLogger) WritePut(id uint64, title, artist, prise string) {
	l.events <- Event{EventType: EventPut, Id: id, Title: title, Artist: artist, Prise: prise}
}

func (l *PostgresTransactionLogger) WriteDelete(id uint64) {
	l.events <- Event{EventType: EventDelete, Id: id}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func NewPostgresTransactionLogger(config PostgresDbParams) (TransactionLogger, error) {

	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		config.host, config.dbName, config.user, config.password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping() // Test the database connection
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresTransactionLogger{db: db}

	exists, err := logger.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}

func (l *PostgresTransactionLogger) createTable() error {
	var err error

	createQuery := `CREATE TABLE transalbum (
						event_type SMALLINT,
						id INT NOT NULL,
						title VARCHAR(128) NOT NULL,
						artist VARCHAR(128) NOT NULL,
						price DECIMAL(5, 2) NOT NULL,
					);`

	_, err = l.db.Exec(createQuery)
	if err != nil {
		return err
	}

	return nil
}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	const table = "transalbum"

	var result string

	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() && result != table {
		rows.Scan(&result)
	}

	return result == table, rows.Err()
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() {
		query := `INSERT INTO transalbum
				(event_type, id, title, artist, price)
				VALUES ($1, $2, $3, $4, $5)`
		for e := range events { // Retrieve the next Event
			_, err := l.db.Exec( // Execute the INSERT query
				query, e.EventType, e.Id, e.Title, e.Artist, e.Prise)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)    // An unbuffered events channel
	outError := make(chan error, 1) // A buffered errors channel

	query := `SELECT event_type , id , title , artist , price FROM transalbum`

	go func() {
		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		rows, err := l.db.Query(query) // Run query; get result set
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close() // This is important!

		var e Event // Create an empty Event

		for rows.Next() { // Iterate over the rows

			err = rows.Scan( // Read the values from the
				&e.EventType, // row into the Event.
				&e.Id, &e.Title, &e.Artist, &e.Prise)

			if err != nil {
				outError <- err
				return
			}

			outEvent <- e // Send e to the channel
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}
