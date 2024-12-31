package core

import (
	"time"
)

type Meaning struct {
	Tag, Content string
}
type Phrase struct {
	Content, By string
}
type Word struct {
	ID          OID
	Content     string
	Meanings    []Meaning
	Synonyms    []string
	Etymologies []string
	Phrases     []Phrase
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
