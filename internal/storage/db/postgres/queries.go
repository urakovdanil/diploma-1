package postgres

const (
	userInsert        = `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id;`
	userSelectByLogin = `SELECT id, login, password FROM users WHERE login = $1;`

	orderInsert      = `INSERT INTO orders (number, status, accrual, user_id) VALUES ($1, $2, $3, $4) RETURNING id, number, status, accrual, user_id;`
	orderSelect      = `SELECT id, number, status, accrual, user_id FROM orders WHERE number = $1;`
	ordersListByUser = `SELECT id, number, status, accrual, user_id, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC;`
)
