package db

import "database/sql"

type Expression struct {
	ID     int
	UserID int
	Expr   string
	Status string
	Result sql.NullFloat64
}
