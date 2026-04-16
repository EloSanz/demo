package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/elosanz/demo/internal/user"
)

// UserSQLiteRepository implements user.UserRepository using database/sql.
type UserSQLiteRepository struct {
	db *sql.DB
}

// NewUserSQLiteRepository constructs a UserSQLiteRepository backed by the given *sql.DB.
func NewUserSQLiteRepository(db *sql.DB) user.UserRepository {
	return &UserSQLiteRepository{db: db}
}

func (r *UserSQLiteRepository) FindAll(ctx context.Context, page int, size int) ([]user.User, int64, error) {
	offset := (page - 1) * size

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting users: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, email, phone, website, created_at, updated_at FROM users ORDER BY id ASC LIMIT ? OFFSET ?",
		size, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		u, err := scanUser(rows.Scan)
		if err != nil {
			return nil, 0, fmt.Errorf("scanning user row: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating user rows: %w", err)
	}
	if users == nil {
		users = []user.User{}
	}
	return users, total, nil
}

func (r *UserSQLiteRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, name, email, phone, website, created_at, updated_at FROM users WHERE id = ?", id,
	)
	u, err := scanUser(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by id %d: %w", id, err)
	}
	return &u, nil
}

func (r *UserSQLiteRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, name, email, phone, website, created_at, updated_at FROM users WHERE email = ?", email,
	)
	u, err := scanUser(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by email: %w", err)
	}
	return &u, nil
}

func (r *UserSQLiteRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users WHERE email = ?", email,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking email existence: %w", err)
	}
	return count > 0, nil
}

func (r *UserSQLiteRepository) Save(ctx context.Context, u user.User) (user.User, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO users (name, email, phone, website) VALUES (?, ?, ?, ?)",
		u.Name, u.Email, nullString(u.Phone), nullString(u.Website),
	)
	if err != nil {
		return user.User{}, fmt.Errorf("inserting user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return user.User{}, fmt.Errorf("getting last insert id: %w", err)
	}
	// Re-read to get DB-generated created_at / updated_at values.
	saved, err := r.FindByID(ctx, id)
	if err != nil {
		return user.User{}, fmt.Errorf("reading saved user: %w", err)
	}
	return *saved, nil
}

func (r *UserSQLiteRepository) Update(ctx context.Context, u user.User) (user.User, error) {
	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET name = ?, email = ?, phone = ?, website = ?, updated_at = ? WHERE id = ?",
		u.Name, u.Email, nullString(u.Phone), nullString(u.Website), now, u.ID,
	)
	if err != nil {
		return user.User{}, fmt.Errorf("updating user %d: %w", u.ID, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return user.User{}, fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		return user.User{}, user.ErrNotFound
	}
	u.UpdatedAt = now
	return u, nil
}

func (r *UserSQLiteRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting user %d: %w", id, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected on delete: %w", err)
	}
	if rows == 0 {
		return user.ErrNotFound
	}
	return nil
}

// scanUser uses the given Scan function to read a user row (works with both *sql.Row and *sql.Rows).
func scanUser(scan func(dest ...any) error) (user.User, error) {
	var u user.User
	var phone, website sql.NullString
	err := scan(&u.ID, &u.Name, &u.Email, &phone, &website, &u.CreatedAt, &u.UpdatedAt)
	u.Phone = phone.String
	u.Website = website.String
	return u, err
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
