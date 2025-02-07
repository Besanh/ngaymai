package service

import (
	"context"
	"net/http"
	"ngaymai/common/cache"
	"ngaymai/model"
)

type (
	IVideo interface {
		PutVideoRanking(ctx context.Context, id string, jsonBody model.VideoActionRequest) (code int, result any)
		GetVideoRanking(ctx context.Context, request model.VideoRankingRequest) (code int, result any)
	}
	Video struct{}
)

func NewVideo() IVideo {
	return &Video{}
}

func (s *Video) PutVideoRanking(ctx context.Context, id string, jsonBody model.VideoActionRequest) (code int, result any) {
	var increment float64
	switch jsonBody.Action {
	case "like":
		increment = 10
	case "dislike":
		increment = -5
	case "view":
		increment = 1
	default:
		return http.StatusBadRequest, map[string]any{
			"error": "action is not valid",
		}
	}

	score, err := cache.RCache.ZIncrBy("video_ranking", increment, id)
	if err != nil {
		return http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		}
	}

	// TODO: I can update sql here

	return http.StatusOK, map[string]any{
		"score": score,
	}
}

func (s *Video) GetVideoRanking(ctx context.Context, request model.VideoRankingRequest) (code int, result any) {
	topVideos, err := cache.RCache.ZRevRangeWithScores("video_ranking", int64(request.Limit), int64(request.Offset))
	if err != nil {
		return http.StatusInternalServerError, map[string]string{"error": "Failed to get top videos"}
	}

	var rs []map[string]any
	for _, v := range topVideos {
		rs = append(rs, map[string]any{
			"video_id": v.Member,
			"score":    v.Score,
		})
	}
	return http.StatusOK, rs
}
