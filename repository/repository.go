package repository

import (
	"context"
	"log"
	"ngaymai/common/sql"
	"ngaymai/model"
)

var DBConn sql.ISqlClientConn

func CreateTable(ctx context.Context, db sql.ISqlClientConn, entity any) (err error) {
	_, err = db.GetDB().NewCreateTable().Model(entity).
		IfNotExists().
		Exec(ctx)
	return
}

func InitRepositories() {

}

func InitTables(ctx context.Context, dbConn sql.ISqlClientConn) {
	if err := CreateTable(ctx, dbConn, (*model.User)(nil)); err != nil {
		log.Fatal(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.Channel)(nil)); err != nil {
		log.Fatal(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.Video)(nil)); err != nil {
		log.Fatal(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.Interaction)(nil)); err != nil {
		log.Fatal(err)
	}
	if err := CreateTable(ctx, dbConn, (*model.VideoRanking)(nil)); err != nil {
		log.Fatal(err)
	}
}
