package mongodb

import (
	"dicio-scrapper/config"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewConnection(cfg *config.EnvConfig) *mongo.Database {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to mongo database")
	}
	return client.Database("dicio")
}
