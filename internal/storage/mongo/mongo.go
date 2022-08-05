package mongo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/numary/webhooks-cloud/internal/storage"
	"github.com/numary/webhooks-cloud/pkg/model"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Store struct {
	uri        string
	client     *mongo.Client
	collection *mongo.Collection
}

func NewConfigStore() (storage.Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDBUri := viper.GetString(constants.StorageMongoConnStringFlag)
	if mongoDBUri == "" {
		mongoDBUri = constants.DefaultMongoConnString
	}
	sharedlogging.Infof("connecting to mongoDB URI: %s", mongoDBUri)
	sharedlogging.Infof("env: %+v", os.Environ())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBUri))
	if err != nil {
		return Store{}, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return Store{}, err
	}

	return Store{
		uri:        mongoDBUri,
		client:     client,
		collection: client.Database("webhooks").Collection("configs"),
	}, nil
}

func (s Store) FindAllConfigs(ctx context.Context) (sharedapi.Cursor[model.ConfigInserted], error) {
	opts := options.Find().SetSort(bson.M{"insertedAt": -1})
	cur, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return sharedapi.Cursor[model.ConfigInserted]{}, fmt.Errorf("mongo.Collection.Find: %w", err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		if err := cur.Close(ctx); err != nil {
			sharedlogging.GetLogger(ctx).Errorf("mongo.Cursor.Close: %s", err)
		}
	}(cur, ctx)

	var results []model.ConfigInserted
	if err := cur.All(ctx, &results); err != nil {
		return sharedapi.Cursor[model.ConfigInserted]{}, fmt.Errorf("mongo.Cursor.All: %w", err)
	}

	return sharedapi.Cursor[model.ConfigInserted]{
		Data: results,
	}, nil
}

func (s Store) FindLastConfig(ctx context.Context) (*model.ConfigInserted, error) {
	res := model.ConfigInserted{}
	opts := options.FindOne().SetSort(bson.M{"insertedAt": -1})
	if err := s.collection.FindOne(ctx, bson.M{}, opts).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &res, nil
}

func (s Store) InsertOneConfig(ctx context.Context, config model.Config) (string, error) {
	configInserted := model.ConfigInserted{
		Config:     config,
		ID:         uuid.New().String(),
		InsertedAt: int(time.Now().UnixNano()),
	}

	res, err := s.collection.InsertOne(ctx, configInserted)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(string), nil
}

func (s Store) DropConfigsCollection(ctx context.Context) error {
	if err := s.collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (s Store) Close(ctx context.Context) error {
	if s.client == nil {
		return nil
	}

	return s.client.Disconnect(ctx)
}
