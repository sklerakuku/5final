package db

import "database/sql"

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func SaveExpression(userID int, expr string) (int, error) {
	res, err := db.Exec(
		"INSERT INTO expressions (user_id, expression, status) VALUES (?, ?, 'pending')",
		userID, expr,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func GetExpressionByID(id int) (Expression, error) {
	var expr Expression
	err := db.QueryRow(
		"SELECT id, user_id, expression, status, result FROM expressions WHERE id = ?", id,
	).Scan(&expr.ID, &expr.UserID, &expr.Expr, &expr.Status, &expr.Result)
	return expr, err
}

func UpdateExpressionStatus(id int, status string) error {
	_, err := db.Exec(
		"UPDATE expressions SET status = ? WHERE id = ?",
		status, id,
	)
	return err
}

func UpdateExpressionResult(id int, status string, result float64) error {
	_, err := db.Exec(
		"UPDATE expressions SET status = ?, result = ? WHERE id = ?",
		status, result, id,
	)
	return err
}

func GetExpressionsByUser(userID int) ([]Expression, error) {
	rows, err := db.Query(
		"SELECT id, expression, status, result FROM expressions WHERE user_id = ? ORDER BY id DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []Expression
	for rows.Next() {
		var expr Expression
		if err := rows.Scan(&expr.ID, &expr.Expr, &expr.Status, &expr.Result); err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}
