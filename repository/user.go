package repository

import (
	"context"
	"ngaymai/common/sql"
	"ngaymai/model"
)

type (
	IUser interface {
		GetUserById(ctx context.Context, db sql.ISqlClientConn, id string) (*model.User, error)
	}
	User struct{}
)

func NewUser() IUser {
	return &User{}
}

func (u *User) GetUserById(ctx context.Context, db sql.ISqlClientConn, id string) (*model.User, error) {
	result := new(model.User)
	query := db.GetDB().NewSelect().Model(&result).
		Where("id = ?", id)
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return result, nil
}
