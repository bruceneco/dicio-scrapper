package wordports

import "context"

type Publisher interface {
	ExtractWord(ctx context.Context, word string) error
}
