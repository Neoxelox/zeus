package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgxutil"
	"github.com/rs/xid"

	"github.com/neoxelox/zeus/internal/database"
	"github.com/neoxelox/zeus/pkg/model"
)

// UserRepository interacts with the user repository.
type UserRepository interface {
	Transaction(ctx context.Context, fn func(UserRepository) error) error
	Create(ctx context.Context, m *model.User) (*model.User, error)
	GetByID(ctx context.Context, ID xid.ID) (*model.User, error)
	List(ctx context.Context, username string) ([]model.User, error)
}

// UserDatabase implements a SQL UserRepository.
type UserDatabase struct {
	db    *pgxpool.Pool
	cn    database.Connection
	table string
}

// NewUserDatabase creates a new UserDatabase instance.
func NewUserDatabase(db *pgxpool.Pool) *UserDatabase {
	return &UserDatabase{
		db:    db,
		cn:    db,
		table: "users",
	}
}

// Transaction returns a UserRepository for transactions.
func (r *UserDatabase) Transaction(ctx context.Context, fn func(UserRepository) error) error {
	tx, err := database.BeginTransaction(ctx, r.db)
	if err != nil {
		return err // nolint
	}

	defer database.WatchTransaction(ctx, tx)()

	err = fn(&UserDatabase{
		db: r.db,
		cn: tx,
	})

	return database.FinishTransaction(ctx, err, tx)
}

// Create creates a new user in the database.
func (r *UserDatabase) Create(ctx context.Context, m *model.User) (*model.User, error) {
	var u model.User

	query := fmt.Sprintf(`INSERT INTO "%s" ("id", "name", "username", "age", "created_at", "updated_at", "deleted_at")
			  			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  			  RETURNING *;`, r.table)

	err := pgxutil.SelectStruct(ctx, r.cn, &u, query,
		m.ID, m.Name, m.Username, m.Age, m.CreatedAt, m.UpdatedAt, m.DeletedAt)
	if err != nil {
		return nil, database.Error(err)
	}

	return &u, nil
}

// GetByID gets an existing user in the database by its ID.
func (r *UserDatabase) GetByID(ctx context.Context, ID xid.ID) (*model.User, error) {
	var u model.User

	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "id" = $1;`, r.table)

	err := pgxutil.SelectStruct(ctx, r.cn, &u, query,
		ID)
	if err != nil {
		return nil, database.Error(err)
	}

	return &u, nil
}

// List gets existing users from the database with a similar username.
func (r *UserDatabase) List(ctx context.Context, username string) ([]model.User, error) {
	var us []model.User

	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE "username" LIKE '%%' || $1 || '%%';`, r.table)

	err := pgxutil.SelectAllStruct(ctx, r.cn, &us, query,
		username)
	if err != nil {
		return nil, database.Error(err)
	}

	return us, nil
}
