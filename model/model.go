package model

import (
	"time"

	"github.com/uptrace/bun"
)

type (
	User struct {
		bun.BaseModel `bun:"table:users"`

		UserID   int64      `bun:",pk,autoincrement"`
		Username string     `bun:",notnull"`
		Email    string     `bun:",unique,notnull"`
		UserType string     `bun:",notnull" json:"user_type"`
		Channels []*Channel `bun:"rel:has-many,join:user_id=creator_id"`
	}

	Channel struct {
		bun.BaseModel `bun:"table:channels"`

		ChannelID int64     `bun:",pk,autoincrement"`
		CreatorID int64     `bun:",notnull"`
		Name      string    `bun:",notnull"`
		CreatedAt time.Time `bun:",default:current_timestamp"`

		// Relationship
		Creator *User    `bun:"rel:belongs-to,join:creator_id=user_id"`
		Videos  []*Video `bun:"rel:has-many,join:channel_id=channel_id"`
	}

	Video struct {
		bun.BaseModel `bun:"table:videos"`

		VideoID    int64     `bun:",pk,autoincrement"`
		ChannelID  int64     `bun:",notnull"`
		Title      string    `bun:",notnull"`
		Score      float64   `bun:",default:0"`
		UploadTime time.Time `bun:",default:current_timestamp"`

		// Relationship
		Channel      *Channel       `bun:"rel:belongs-to,join:channel_id=channel_id"`
		Interactions []*Interaction `bun:"rel:has-many,join:video_id=video_id"`
	}

	Interaction struct {
		bun.BaseModel `bun:"table:interactions"`

		InteractionID    int64     `bun:",pk,autoincrement"`
		VideoID          int64     `bun:",notnull"`
		UserID           int64     `bun:",notnull"`
		InteractionType  string    `bun:",notnull"`
		InteractionValue float64   `bun:",default:1"`
		CreatedAt        time.Time `bun:",default:current_timestamp"`

		// Relationship
		Video *Video `bun:"rel:belongs-to,join:video_id=video_id"`
		User  *User  `bun:"rel:belongs-to,join:user_id=user_id"`
	}

	VideoRanking struct {
		bun.BaseModel `bun:"table:video_rankings"`

		RankID       int64     `bun:",pk,autoincrement"`
		VideoID      int64     `bun:",notnull"`
		Rank         int       `bun:",notnull"`
		Score        float64   `bun:",notnull"`
		CalculatedAt time.Time `bun:",default:current_timestamp"`

		// Relationship
		Video *Video `bun:"rel:belongs-to,join:video_id=video_id"`
	}
)
