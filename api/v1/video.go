package v1

import (
	"net/http"
	"ngaymai/model"
	"ngaymai/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	video service.IVideo
}

func NewVideoHandler(e *gin.Engine, videoService service.IVideo) {
	handler := VideoHandler{
		video: videoService,
	}

	group := e.Group("video/v1")
	{
		group.PUT(":video_id", handler.PutVideoRanking)
		group.GET("ranking", handler.GetVideoRanking)
	}
}

// Put Video Ranking godoc
// @Summary put action view, like, comment, share video
// @Schemes
// @Tags video
// @Accept json
// @Produce json
// @Success 200 {string} success
// @Router /video/v1/:video_id [put]
// @Param data body model.VideoActionRequest true "data"
func (h *VideoHandler) PutVideoRanking(c *gin.Context) {
	jsonBody := model.VideoActionRequest{}
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id := c.Param("video_id")
	if len(id) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "video_id is required",
		})
		return
	}

	code, result := h.video.PutVideoRanking(c, id, jsonBody)
	c.JSON(code, result)
}

// Get Video Ranking godoc
// @Summary Get Video Ranking
// @Schemes
// @Tags video
// @Accept json
// @Produce json
// @Success 200 {string} success
// @Router /video/v1/ranking [get]
// @Param limit query int false "limit"
// @Param offset query int false "offset"
func (h *VideoHandler) GetVideoRanking(c *gin.Context) {
	limitParam := c.Query("limit")
	offsetParam := c.Query("offset")
	if len(limitParam) < 1 {
		limitParam = "10"
	}
	if len(offsetParam) < 1 {
		offsetParam = "0"
	}

	limit, _ := strconv.ParseInt(limitParam, 10, 64)
	offset, _ := strconv.ParseInt(offsetParam, 10, 64)

	code, result := h.video.GetVideoRanking(c, model.VideoRankingRequest{
		Limit:  int(limit),
		Offset: int(offset),
	})
	c.JSON(code, result)
}
