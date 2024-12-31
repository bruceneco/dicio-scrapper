package mongodb

import (
	"context"
	"dicio-scrapper/internal/adapters/db/models"
	"dicio-scrapper/internal/domain/core"
	"dicio-scrapper/internal/ports/wordports"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ wordports.Repo = (*Word)(nil)

type (
	Word struct {
		collection *mongo.Collection
	}
)

func NewWord(db *mongo.Database) *Word {
	return &Word{
		collection: db.Collection("words"),
	}
}

func (w *Word) Insert(ctx context.Context, word core.Word) (core.Word, error) {
	var model models.Word

	model.FromCore(word)
	model.BeforeInsert()

	inserted, err := w.collection.InsertOne(ctx, &model)
	if err != nil {
		return core.Word{}, err
	}

	id, ok := inserted.InsertedID.(bson.ObjectID)
	if !ok {
		return core.Word{}, errors.New("inserted ID is not ObjectID")
	}
	model.ID = id
	return model.ToCore(), nil
}

func (w *Word) FindByContent(ctx context.Context, content string) (core.Word, error) {
	regex := bson.M{"$regex": bson.Regex{Pattern: content, Options: "i"}}
	result := w.collection.FindOne(ctx, bson.M{"content": regex})

	if err := result.Err(); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return core.Word{}, wordports.ErrWordNotFound
		default:
			return core.Word{}, err
		}
	}

	var model models.Word
	if err := result.Decode(&model); err != nil {
		return core.Word{}, err
	}

	return model.ToCore(), nil
}
