package db

import (
	"context"
	"denet/models"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	p *pgxpool.Pool
}

func New(connString string) *DB {
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		return nil
	}

	return &DB{p: db}
}

func (d *DB) Close() {
	d.p.Close()
}

func (d *DB) GetUserByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	return user, d.p.QueryRow(ctx, "SELECT uid, username, refer_uid, points FROM users WHERE uid = $1", id).Scan(&user.ID, &user.Login, &user.ReferUID, &user.Points)
}

func (d *DB) GetUserByLoginAndPassword(ctx context.Context, username, password string) (models.User, error) {
	var user models.User

	// Ищем пользователя по логину
	err := d.p.QueryRow(ctx, `
		SELECT uid, username, password_hash, refer_uid, points
		FROM users
		WHERE username = $1
	`, username).Scan(&user.ID, &user.Login, &user.Password, &user.ReferUID, &user.Points)
	if err != nil {
		return models.User{}, fmt.Errorf("user not found %w", err)
	}

	if !checkPasswordHash(password, user.Password) {
		return models.User{}, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (d *DB) CreateUser(ctx context.Context, passHash, username string) error {
	_, err := d.p.Exec(ctx, `
    INSERT INTO users (username, password_hash, refer_uid, points)
    VALUES ($1, $2, $3, $4)
`, username, passHash, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (d *DB) GetLeaderboard(ctx context.Context, count int) ([]models.User, error) {
	rows, err := d.p.Query(ctx, "SELECT uid, username, refer_uid, points FROM users ORDER BY points DESC LIMIT $1", count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0, count)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Login, &user.ReferUID, &user.Points)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
func (d *DB) IncrementPoints(ctx context.Context, tx pgx.Tx, userID int) error {
	var err error
	if tx == nil {
		// Начинаем транзакцию
		tx, err = d.p.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer func() {
			if err != nil {
				tx.Rollback(ctx)
				return
			}
			tx.Commit(ctx)
		}()
	}
	// Выполняем SQL-запрос для увеличения points на 1
	result, err := tx.Exec(ctx, `
		UPDATE users
		SET points = points + 1
		WHERE uid = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("failed to increment points: %w", err)
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (d *DB) SetRefferer(ctx context.Context, userID, referID int) error {
	// Проверяем, существует ли refer_uid
	_, err := d.GetUserByID(ctx, referID)
	if err != nil {
		return fmt.Errorf("refer not found: %w", err)
	}

	// Начинаем транзакцию
	tx, err := d.p.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Обновляем refer_uid, только если он равен 0
	result, err := tx.Exec(ctx, `
		UPDATE users
		SET refer_uid = $1
		WHERE uid = $2 AND refer_uid = 0
	`, referID, userID)
	if err != nil {
		return fmt.Errorf("failed to update refer_uid: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("refer_uid is already set or user not found")
	}

	err = d.IncrementPoints(ctx, tx, userID)
	if err != nil {
		return fmt.Errorf("failed to increment refer points: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
