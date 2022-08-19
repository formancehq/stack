package mongo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/numary/webhooks/internal/model"
	"github.com/numary/webhooks/internal/storage"
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
		uri:    mongoDBUri,
		client: client,
		collection: client.Database(
			viper.GetString(constants.StorageMongoDatabaseNameFlag)).
			Collection("configs"),
	}, nil
}

func (s Store) FindAllConfigs(ctx context.Context) (sharedapi.Cursor[model.ConfigInserted], error) {
	opts := options.Find().SetSort(bson.M{"updatedAt": -1})
	cur, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return sharedapi.Cursor[model.ConfigInserted]{}, fmt.Errorf("mongo.Collection.Find: %w", err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			sharedlogging.GetLogger(ctx).Errorf("mongo.Cursor.Close: %s", err)
		}
	}()

	var results []model.ConfigInserted
	if err := cur.All(ctx, &results); err != nil {
		return sharedapi.Cursor[model.ConfigInserted]{}, fmt.Errorf("mongo.Cursor.All: %w", err)
	}

	return sharedapi.Cursor[model.ConfigInserted]{
		Data: results,
	}, nil
}

func (s Store) InsertOneConfig(ctx context.Context, cfg model.Config) (string, error) {
	configInserted := model.ConfigInserted{
		Config:    cfg,
		ID:        uuid.New().String(),
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	res, err := s.collection.InsertOne(ctx, configInserted)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(string), nil
}

func (s Store) DeleteOneConfig(ctx context.Context, id string) (int64, error) {
	res, err := s.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (s Store) UpdateOneConfigActive(ctx context.Context, id string, active bool) (model.ConfigInserted, int64, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	resFind := s.collection.FindOne(ctx, filter)
	if err := resFind.Err(); err != nil {
		return model.ConfigInserted{}, 0, fmt.Errorf("mongo.Collection.FindOne: %w", err)
	}

	var cfg model.ConfigInserted
	if err := resFind.Decode(&cfg); err != nil {
		return model.ConfigInserted{}, 0, fmt.Errorf("mongo.SingleResult.Decode: %w", err)
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: active}}}}
	resUpdate, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.ConfigInserted{}, 0, fmt.Errorf("mongo.Collection.UpdateOne: %w", err)
	}

	return cfg, resUpdate.ModifiedCount, nil
}

func (s Store) UpdateOneConfigSecret(ctx context.Context, id, secret string) (int64, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "secret", Value: secret}}}}
	resUpdate, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("mongo.Collection.UpdateOne: %w", err)
	}

	return resUpdate.ModifiedCount, nil
}

func (s Store) FindEventType(ctx context.Context, eventType string) (bool, error) {
	filter := bson.D{{Key: "eventTypes", Value: eventType}}
	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("mongo.Collection.Find: %w", err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			sharedlogging.GetLogger(ctx).Errorf("mongo.Cursor.Close: %s", err)
		}
	}()

	var results []model.ConfigInserted
	if err := cur.All(ctx, &results); err != nil {
		return false, fmt.Errorf("mongo.Cursor.All: %w", err)
	}

	return len(results) > 0, nil
}

func (s Store) Close(ctx context.Context) error {
	if s.client == nil {
		return nil
	}

	return s.client.Disconnect(ctx)
}
