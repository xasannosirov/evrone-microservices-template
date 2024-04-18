package postgresql

import (
	"context"
	"fmt"
	"user-service/internal/entity"
	"user-service/internal/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

const (
	usersTableName      = "users"
	usersServiceName    = "userService"
	usersSpanRepoPrefix = "usersRepo"
)

type usersRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUsersRepo(db *postgres.PostgresDB) *usersRepo {
	return &usersRepo{
		tableName: usersTableName,
		db:        db,
	}
}

func (p *usersRepo) usersSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.
		Select(
			"id",
			"first_name",
			"last_name",
			"username",
			"email",
			"password",
			"bio",
			"website",
			"created_at",
			"updated_at",
		).From(p.tableName)
}

func (p usersRepo) Create(ctx context.Context, news *entity.User) error {
	data := map[string]any{
		"id":         news.GUID,
		"first_name": news.FirstName,
		"last_name":  news.LastName,
		"username":   news.Username,
		"email":      news.Email,
		"password":   news.Password,
		"bio":        news.Bio,
		"website":    news.Website,
		"created_at": news.CreatedAt,
		"updated_at": news.UpdatedAt,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

func (p usersRepo) Update(ctx context.Context, users *entity.User) error {
	clauses := map[string]any{
		"first_name": users.FirstName,
		"last_name":  users.LastName,
		"username":   users.Username,
		"email":      users.Email,
		"password":   users.Password,
		"bio":        users.Bio,
		"website":    users.Website,
		"updated_at": users.UpdatedAt,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", users.GUID)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p usersRepo) Delete(ctx context.Context, guid string) error {
	sqlStr, args, err := p.db.Sq.Builder.
		Delete(p.tableName).
		Where(p.db.Sq.Equal("id", guid)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" delete")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p usersRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	var (
		user entity.User
	)

	queryBuilder := p.usersSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&user.GUID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.Website,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}

	return &user, nil
}

func (p usersRepo) List(ctx context.Context, limit uint64, offset uint64, filter map[string]string) ([]*entity.User, error) {
	var (
		users []*entity.User
	)
	queryBuilder := p.usersSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()
	users = make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		if err = rows.Scan(
			&user.GUID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Bio,
			&user.Website,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}
		users = append(users, &user)
	}

	return users, nil
}
