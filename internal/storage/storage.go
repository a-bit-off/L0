package storage

const (
	CreateTableOrders = `
	CREATE TABLE IF NOT EXISTS orders (
    	id text PRIMARY KEY,
    	order_data JSONB);`
	InsertIntoOrders = `INSERT INTO orders(id, order_data) VALUES(?, ?);`
	DeleteFromOrders = `DELETE FROM orders WHERE id == ?;`
)
