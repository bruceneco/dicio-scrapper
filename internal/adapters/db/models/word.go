package models

import (
	"dicio-scrapper/internal/domain/core"
	"dicio-scrapper/internal/ports/conversions"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var _ conversions.CoreConverter[core.Word] = (*Word)(nil)

type Word struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Content   string        `bson:"content,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
}

func (w *Word) BeforeInsert() {
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now()
	}
	w.UpdatedAt = time.Now()
}

func (w *Word) FromCore(word core.Word) {
	*w = Word{
		ID:        bson.ObjectID(word.ID),
		Content:   word.Content,
		CreatedAt: word.CreatedAt,
		UpdatedAt: word.UpdatedAt,
	}
}

func (w *Word) ToCore() core.Word {
	return core.Word{
		ID:        core.OID(w.ID),
		Content:   w.Content,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}
