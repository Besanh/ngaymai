package repository

import (
	"context"
	"ngaymai/common/sqlclient"
	"ngaymai/model"

	"github.com/uptrace/bun"
)

type (
	IVideo interface {
		GetVideoById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.Video, error)
		UpdateVideoScoreInDB(ctx context.Context, tx bun.Tx, id string, increment float64) error
		GetVideoRanking(ctx context.Context, db sqlclient.ISqlClientConn, request model.VideoRankingRequest) (result []model.Video, err error)
	}
	Video struct{}
)

func NewVideo() IVideo {
	return &Video{}
}

func (repo *Video) GetVideoById(ctx context.Context, db sqlclient.ISqlClientConn, id string) (*model.Video, error) {
	var video model.Video
	err := db.GetDB().NewSelect().Model(&video).Where("id = ?", id).Scan(ctx)
	return &video, err
}

func (repo *Video) UpdateVideoScoreInDB(ctx context.Context, tx bun.Tx, id string, increment float64) error {
	_, err := tx.NewUpdate().
		Model((*model.Video)(nil)).
		Set("score = score + ?", increment).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (repo *Video) GetVideoRanking(ctx context.Context, db sqlclient.ISqlClientConn, request model.VideoRankingRequest) (result []model.Video, err error) {
	err = db.GetDB().NewSelect().
		Model((*model.Video)(nil)).
		Column("id", "score").
		OrderExpr("score DESC").
		Limit(request.Limit).
		Offset(request.Offset).
		Scan(ctx, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
