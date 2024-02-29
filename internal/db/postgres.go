package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func RunDB(addr string) {
	var err error
	DB, err = sql.Open("postgres", addr)
	if err != nil {
		panic(fmt.Sprintf("Не удалось подключиться к БД: %v", err))
	}

	createTableQuery :=
		`CREATE TABLE IF NOT EXISTS users (
                      id SERIAL PRIMARY KEY,
                      name VARCHAR(255) NOT NULL UNIQUE,
   					  email VARCHAR(255),
    				  address VARCHAR(255),
    				  password VARCHAR(255) NOT NULL,
                      loyalty INT);
		CREATE TABLE IF NOT EXISTS orders (
                      id SERIAL PRIMARY KEY,
                      number BIGINT NOT NULL,
   					  user_id INT NOT NULL,
    				  status VARCHAR(255) NOT NULL,
    				  accrual INT,
    				  uploaded_at VARCHAR(255));
		CREATE TABLE IF NOT EXISTS withdrawals (
                      id SERIAL PRIMARY KEY,
                      order_id VARCHAR(255) NOT NULL,
   					  user_id INT NOT NULL,
    				  withdrawal INT,
				      processed_at VARCHAR(255) NOT NULL);`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}
}
