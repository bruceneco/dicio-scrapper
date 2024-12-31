package wordports

import (
	"context"
	"dicio-scrapper/internal/domain/core"
)

type (
	Repo interface {
		Insert(ctx context.Context, word core.Word) (core.Word, error)
		FindByContent(ctx context.Context, content string) (core.Word, error)
	}
)
