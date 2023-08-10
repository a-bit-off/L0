package postgres

import (
	"L0/internal/storage"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
	"log"
	"strings"
)

type Storage struct {
	db *sql.DB
}

func New(connectionString string) (*Storage, error) {
	const op = "L0/internal/storage/postgres.New"

	// Строка подключения к базе данных
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	// Проверяем, что соединение действительно работает
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	log.Println("Init storage successful!")

	// Создаем таблицу
	stmt, err := db.Prepare(storage.CreateTableOrders)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Println("Create table orders successful!")

	return &Storage{db: db}, nil
}

func (s Storage) GetById(id string) ([]byte, error) {
	const op = "storage.postgres.GetById"

	var jsonB []byte
	err := s.db.QueryRow(storage.GetByIdFromOrders, id).Scan(&jsonB)
	if err != nil {
		if err != sql.ErrNoRows {
			return jsonB, fmt.Errorf("%s: %s", op, err)
		}
	}

	return jsonB, nil
}

func (s Storage) AddOrder(id, order string) error {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare(storage.InsertIntoOrders)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, order)
	if err != nil {
		if strings.Contains(err.Error(), "invalid input syntax for type json") {
			return fmt.Errorf("Order not a json!")
		} else if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("Not a unique ID!")
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
