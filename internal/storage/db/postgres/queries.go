package postgres

const (
	userInsert        = `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id;`
	userSelectByLogin = `SELECT id, login, password FROM users WHERE login = $1;`
)
