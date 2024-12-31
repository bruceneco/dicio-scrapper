package models

import (
	"dicio-scrapper/internal/domain/core"
	"dicio-scrapper/internal/ports/conversions"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var _ conversions.CoreConverter[core.Word] = (*Word)(nil)

type Meaning struct {
	Tag     string `bson:"tag,omitempty"`
	Content string `bson:"content,omitempty"`
}
type Phrase struct {
	Content string `bson:"content,omitempty"`
	By      string `bson:"by,omitempty"`
}

type Word struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Content     string        `bson:"content,omitempty"`
	Meanings    []Meaning     `bson:"meanings,omitempty"`
	Synonyms    []string      `bson:"synonyms,omitempty"`
	Etymologies []string      `bson:"etymologies,omitempty"`
	Phrases     []Phrase      `bson:"phrases,omitempty"`
	CreatedAt   time.Time     `bson:"created_at,omitempty"`
	UpdatedAt   time.Time     `bson:"updated_at,omitempty"`
}

func (w *Word) BeforeInsert() {
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now()
	}
	w.UpdatedAt = time.Now()
}

func (w *Word) FromCore(word core.Word) {
	*w = Word{
		ID:          bson.ObjectID(word.ID),
		Content:     word.Content,
		Synonyms:    word.Synonyms,
		Etymologies: word.Etymologies,
		CreatedAt:   word.CreatedAt,
		UpdatedAt:   word.UpdatedAt,
	}
	for _, m := range word.Meanings {
		w.Meanings = append(w.Meanings, Meaning{
			Tag:     m.Tag,
			Content: m.Content,
		})
	}
	for _, p := range word.Phrases {
		w.Phrases = append(w.Phrases, Phrase{
			By:      p.By,
			Content: p.Content,
		})
	}
}

func (w *Word) ToCore() core.Word {
	cw := core.Word{
		ID:          core.OID(w.ID),
		Content:     w.Content,
		Synonyms:    w.Synonyms,
		Etymologies: w.Etymologies,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
	for _, m := range w.Meanings {
		cw.Meanings = append(cw.Meanings, core.Meaning{
			Tag:     m.Tag,
			Content: m.Content,
		})
	}
	for _, p := range w.Phrases {
		cw.Phrases = append(cw.Phrases, core.Phrase{
			By:      p.By,
			Content: p.Content,
		})
	}
	return cw
}
