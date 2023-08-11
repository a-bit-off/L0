package storage

const (
	CreateTableOrders = `
	CREATE TABLE IF NOT EXISTS orders (
    	id text PRIMARY KEY,
    	order_data JSONB);`
	InsertIntoOrders  = `INSERT INTO orders(id, order_data) VALUES($1, $2);`
	GetByIdFromOrders = `SELECT order_data FROM orders WHERE id = $1;`
	GetAllFromOrders  = `SELECT * FROM orders`
)
