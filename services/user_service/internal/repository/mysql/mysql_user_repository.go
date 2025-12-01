package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pmii/user-service/internal/domain/entity"
	"github.com/pmii/user-service/internal/domain/repository"
)

type mysql_user_repository struct {
	db *sql.DB
}

// NewMysqlUserRepository creates new MySQL user repository
func NewMysqlUserRepository(db *sql.DB) repository.UserRepository {
	return &mysql_user_repository{
		db: db,
	}
}

func (r *mysql_user_repository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO tbl_user (user_name, user_email, user_password, user_level, user_status, user_photo, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		user.UserName,
		user.UserEmail,
		user.UserPassword,
		user.UserLevel,
		user.UserStatus,
		user.UserPhoto,
		now,
		now,
	)
	
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	user.UserId = int(id)
	user.CreatedAt = now
	user.UpdatedAt = now
	
	return nil
}

func (r *mysql_user_repository) FindById(ctx context.Context, userId int) (*entity.User, error) {
	query := `
		SELECT user_id, user_name, user_email, user_password, user_level, user_status, user_photo, created_at, updated_at
		FROM tbl_user
		WHERE user_id = ?
	`
	
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		&user.UserId,
		&user.UserName,
		&user.UserEmail,
		&user.UserPassword,
		&user.UserLevel,
		&user.UserStatus,
		&user.UserPhoto,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return user, nil
}

func (r *mysql_user_repository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT user_id, user_name, user_email, user_password, user_level, user_status, user_photo, created_at, updated_at
		FROM tbl_user
		WHERE user_email = ?
	`
	
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserId,
		&user.UserName,
		&user.UserEmail,
		&user.UserPassword,
		&user.UserLevel,
		&user.UserStatus,
		&user.UserPhoto,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return user, nil
}

func (r *mysql_user_repository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `
		SELECT user_id, user_name, user_email, user_password, user_level, user_status, user_photo, created_at, updated_at
		FROM tbl_user
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserEmail,
			&user.UserPassword,
			&user.UserLevel,
			&user.UserStatus,
			&user.UserPhoto,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (r *mysql_user_repository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE tbl_user
		SET user_name = ?, user_email = ?, user_level = ?, user_status = ?, user_photo = ?, updated_at = ?
		WHERE user_id = ?
	`
	
	_, err := r.db.ExecContext(ctx, query,
		user.UserName,
		user.UserEmail,
		user.UserLevel,
		user.UserStatus,
		user.UserPhoto,
		time.Now(),
		user.UserId,
	)
	
	return err
}

func (r *mysql_user_repository) Delete(ctx context.Context, userId int) error {
	query := `DELETE FROM tbl_user WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userId)
	return err
}

func (r *mysql_user_repository) CountAll(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM tbl_user`
	
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}
