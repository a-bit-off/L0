package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL

	"L0/internal/storage"
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
