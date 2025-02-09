package service_test

import (
	"context"
	"ngaymai/mock"
	"ngaymai/model"
	"ngaymai/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*gomock.Controller, *mock.MockIRedisCache, *mock.MockISqlClientConn, *service.Video) {
	ctrl := gomock.NewController(t)
	mockRedis := mock.NewMockIRedisCache(ctrl)
	mockDB := mock.NewMockISqlClientConn(ctrl)

	videoService := service.NewVideo(mockRedis, mockDB)

	return ctrl, mockRedis, mockDB, videoService
}

// Test GetVideoRanking
func TestGetVideoRankingFromRedis(t *testing.T) {
	ctrl, mockRedis, _, videoService := setupTest(t)
	defer ctrl.Finish()

	// Mock
	mockRedis.EXPECT().
		ZRevRangeWithScores("video_ranking", int64(0), int64(5)).
		Return([]redis.Z{
			{Member: "video_1", Score: 100},
			{Member: "video_2", Score: 95},
		}, nil)

	// Test
	code, rs := videoService.GetVideoRanking(context.Background(), model.VideoRankingRequest{
		Limit:  5,
		Offset: 0,
	})
	var err error
	if code != 200 {
		err = rs.(error)
	}
	result := rs.([]map[string]any)

	// Test result
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "video_1", result[0]["video_id"])
	assert.Equal(t, 100.0, result[0]["score"])
}

// Test GetVideoRanking
func TestGetVideoRankingFromDB(t *testing.T) {
	ctrl, mockRedis, _, videoService := setupTest(t)
	defer ctrl.Finish()

	// Mock Redis return empty data
	mockRedis.EXPECT().
		ZRevRangeWithScores("video_ranking", int64(0), int64(5)).
		Return(nil, nil)

	// Call service to get result
	code, rs := videoService.GetVideoRanking(context.Background(), model.VideoRankingRequest{
		Limit:  5,
		Offset: 0,
	})
	var err error
	if code != 200 {
		err = rs.(error)
	}
	result := rs.([]map[string]any)

	// Test result
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "video_1", result[0]["video_id"])
	assert.Equal(t, 100.0, result[0]["score"])
}

// Test UpdateVideoScore when action is "like"
func TestUpdateVideoScore_Like(t *testing.T) {
	ctrl, mockRedis, _, videoService := setupTest(t)
	defer ctrl.Finish()

	// Mock Redis return score
	mockRedis.EXPECT().
		ZIncrBy("video_ranking", float64(10), "video_1").
		Return(110.0, nil)

	// Test video service
	code, rs := videoService.PutVideoRanking(context.Background(), "video_1", model.VideoActionRequest{
		Action: "like",
	})
	var err error
	if code != 200 {
		err = rs.(error)
	}
	result := rs.(map[string]any)

	// Check result
	assert.Nil(t, err)
	assert.Equal(t, 110.0, result["score"])
}

// Test UpdateVideoScore with action is "dislike"
func TestUpdateVideoScore_Dislike(t *testing.T) {
	ctrl, mockRedis, _, videoService := setupTest(t)
	defer ctrl.Finish()

	// Mock Redis trả về điểm mới sau khi giảm
	mockRedis.EXPECT().
		ZIncrBy("video_ranking", float64(-5), "video_1").
		Return(95.0, nil)

	// Test
	code, rs := videoService.PutVideoRanking(context.Background(), "video_1", model.VideoActionRequest{
		Action: "dislike",
	})
	var err error
	if code != 200 {
		err = rs.(error)
	}
	result := rs.(map[string]any)

	// Kiểm tra kết quả
	assert.Nil(t, err)
	assert.Equal(t, 95.0, result["score"])
}

// Test UpdateVideoScore with invalid action
func TestUpdateVideoScore_InvalidAction(t *testing.T) {
	ctrl, _, _, videoService := setupTest(t)
	defer ctrl.Finish()

	// Test with invalid action
	code, rs := videoService.PutVideoRanking(context.Background(), "video_1", model.VideoActionRequest{
		Action: "invalid",
	})
	var err error
	if code != 200 {
		err = rs.(error)
	}
	result := rs.(map[string]any)

	// Check error
	assert.NotNil(t, err)
	assert.Equal(t, "invalid action", err.Error())
	assert.Equal(t, 0.0, result["score"])
}
