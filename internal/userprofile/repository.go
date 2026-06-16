package userprofile

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userProfileRepository struct {
	db *pgxpool.Pool
}

func NewUserProfileRepository(db *pgxpool.Pool) UserProfileRepository {
	return &userProfileRepository{db: db}
}

type UserProfileRepository interface {
	Create(ctx context.Context, request UserProfile) (UserProfile, error)
	List(ctx context.Context) ([]UserProfile, error)
	GetByID(ctx context.Context, userID string) (UserProfile, error)
	Update(ctx context.Context, userID string, request UserProfile) (UserProfile, error)
	Delete(ctx context.Context, userID string) error
}

func (r *userProfileRepository) Create(ctx context.Context, request UserProfile) (UserProfile, error) {
	const query = `
		INSERT INTO userprofile (user_id,password,first_name, last_name, address, birthdate, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id,password, first_name, last_name, address, birthdate, email, created_at, updated_at`

	rows, err := r.db.Query(ctx, query,
		request.UserID,
		request.Password,
		request.FirstName,
		request.LastName,
		request.Address,
		request.Birthdate,
		request.Email)

	if err != nil {
		return UserProfile{}, err
	}

	defer rows.Close()
	// Decode เข้า Struct ได้ทันทีด้วย RowToStructByName
	insertedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserProfile])
	if err != nil {
		return UserProfile{}, err
	}
	return insertedUser, nil
}

func (r *userProfileRepository) List(ctx context.Context) ([]UserProfile, error) {
	const query = `
		SELECT user_id, first_name, password, last_name, address, birthdate, email, created_at, updated_at
		FROM userprofile
		ORDER BY created_at DESC, user_id DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listProfile, err := pgx.CollectRows(rows, pgx.RowToStructByName[UserProfile])

	return listProfile, nil
}

func (r *userProfileRepository) GetByID(ctx context.Context, userID string) (UserProfile, error) {
	const query = `
		SELECT user_id, password, first_name, last_name, address, birthdate, email, created_at, updated_at
		FROM userprofile
		WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)

	if err != nil {
		return UserProfile{}, err
	}

	defer rows.Close()
	// Decode เข้า Struct ได้ทันทีด้วย RowToStructByName
	userProfile, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserProfile])
	if err != nil {
		return UserProfile{}, err
	}
	return userProfile, nil
}

func (r *userProfileRepository) Update(ctx context.Context, userID string, request UserProfile) (UserProfile, error) {
	const query = `
		UPDATE userprofile
		SET password = $2,
			first_name = $3,
			last_name = $4,
			address = $5,
			birthdate = $6,
			email = $7,
			updated_at = now()
		WHERE user_id = $1
		RETURNING user_id,password, first_name, last_name, address, birthdate, email, created_at, updated_at`

	rows, err := r.db.Query(ctx, query,
		userID,
		request.Password,
		request.FirstName,
		request.LastName,
		request.Address,
		request.Birthdate,
		request.Email)

	if err != nil {
		return UserProfile{}, err
	}

	defer rows.Close()
	// Decode เข้า Struct ได้ทันทีด้วย RowToStructByName
	updatedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserProfile])
	if err != nil {
		return UserProfile{}, err
	}
	return updatedUser, nil
}

func (r *userProfileRepository) Delete(ctx context.Context, userID string) error {
	const query = `DELETE FROM userprofile WHERE user_id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
