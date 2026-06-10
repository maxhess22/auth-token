package repositories

import (
	"database/sql"
	"errors"
	"max/auth/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	// 1. Cambiamos $1, $2, $3 por ? y removemos el RETURNING id
	query := `INSERT INTO users (name, email, password, created_at) VALUES (?, ?, ?, NOW())`

	// 2. Usamos Exec en lugar de QueryRow porque MySQL no retorna filas en un INSERT convencional
	result, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	// 3. Obtenemos el ID autogenerado e insertado mediante LastInsertId
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Asignamos el ID al struct del usuario
	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Cambiamos $1 por ? para cumplir con el estándar de MySQL
	query := `SELECT id, name, email, password FROM users WHERE email = ?`

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return &user, nil
}
