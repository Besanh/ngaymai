package service

import (
	"context"
	"log"
	"net/http"
	"ngaymai/common/cache"
	"ngaymai/common/sqlclient"
	"ngaymai/model"
	"ngaymai/repository"
)

type (
	IVideo interface {
		PutVideoRanking(ctx context.Context, id string, jsonBody model.VideoActionRequest) (code int, result any)
		GetVideoRanking(ctx context.Context, request model.VideoRankingRequest) (code int, result any)
	}
	Video struct {
		redisCache cache.IRedisCache
		dbClient   sqlclient.ISqlClientConn
	}
)

func NewVideo(redisCache cache.IRedisCache, dbClient sqlclient.ISqlClientConn) *Video {
	return &Video{
		redisCache: redisCache,
		dbClient:   dbClient,
	}
}

/*
 * Update video ranking with specified action
 * Use transaction to update both redis and database
 * @return http code and result(json)
 */
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

	// Start PostgreSQL transaction
	tx, err := repository.DBConn.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, map[string]any{
			"error": "Failed to start transaction",
		}
	}

	// Rollback transaction in case of error
	defer tx.Rollback()

	// Get video information
	video, err := repository.NewVideo().GetVideoById(ctx, repository.DBConn, id)
	if err != nil {
		log.Fatalf("Failed to get video: %v", err)
		return http.StatusInternalServerError, map[string]any{
			"error": "Failed to get video: " + err.Error(),
		}
	}
	if video == nil {
		log.Fatal("Video not found")
		return http.StatusNotFound, map[string]any{
			"error": "Video not found",
		}
	}

	// Update score in database
	err = repository.NewVideo().UpdateVideoScoreInDB(ctx, tx, id, increment)
	if err != nil {
		log.Fatalf("Failed to update database: %v", err)
		return http.StatusInternalServerError, map[string]any{
			"error": "Failed to update database: " + err.Error(),
		}
	}

	// Update score in Redis
	score, err := cache.RCache.ZIncrBy("video_ranking", increment, id)
	if err != nil {
		log.Fatalf("Failed to update Redis: %v", err)
		return http.StatusInternalServerError, map[string]any{
			"error": "Failed to update Redis: " + err.Error(),
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
		return http.StatusInternalServerError, map[string]any{
			"error": "Failed to commit transaction: " + err.Error(),
		}
	}

	return http.StatusOK, map[string]any{
		"score": score,
	}
}

/*
 * Get video ranking, prioritize from Redis if available
 * @return http code and result(json)
 */
func (s *Video) GetVideoRanking(ctx context.Context, request model.VideoRankingRequest) (code int, result any) {
	// Get from Redis
	topVideos, err := cache.RCache.ZRevRangeWithScores("video_ranking", int64(request.Limit), int64(request.Offset))
	if err != nil {
		log.Fatalf("Failed to get top videos: %v", err)
		return http.StatusInternalServerError, map[string]string{"error": "Failed to get top videos"}
	}

	var rs []map[string]any
	if len(topVideos) > 0 {
		for _, v := range topVideos {
			rs = append(rs, map[string]any{
				"video_id": v.Member,
				"score":    v.Score,
			})
		}
		return http.StatusOK, rs
	}

	// Get from database
	videos, err := repository.NewVideo().GetVideoRanking(ctx, repository.DBConn, request)
	if err != nil {
		log.Fatalf("Failed to get top videos: %v", err)
		return http.StatusInternalServerError, map[string]string{"error": "Failed to get top videos"}
	}

	if len(videos) > 0 {
		for _, v := range videos {
			rs = append(rs, map[string]any{
				"video_id": v.VideoID,
				"score":    v.Score,
			})
		}
		return http.StatusOK, rs
	}

	return http.StatusOK, videos
}
