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
	webhooks "github.com/numary/webhooks/pkg"
	"github.com/numary/webhooks/pkg/storage"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Store struct {
	uri                string
	client             *mongo.Client
	configsCollection  *mongo.Collection
	requestsCollection *mongo.Collection
}

var _ storage.Store = &Store{}

func NewStore() (storage.Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDBUri := viper.GetString(constants.StorageMongoConnStringFlag)
	sharedlogging.Infof("connecting to mongoDB URI: %s", mongoDBUri)
	sharedlogging.Infof("env: %+v", os.Environ())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBUri))
	if err != nil {
		return Store{}, fmt.Errorf("mongo.Connect: %w", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return Store{}, fmt.Errorf("mongo.Client.Ping: %w", err)
	}

	return Store{
		uri:    mongoDBUri,
		client: client,
		configsCollection: client.Database(
			viper.GetString(constants.StorageMongoDatabaseNameFlag)).
			Collection(constants.MongoCollectionConfigs),
		requestsCollection: client.Database(
			viper.GetString(constants.StorageMongoDatabaseNameFlag)).
			Collection(constants.MongoCollectionRequests),
	}, nil
}

func (s Store) FindManyConfigs(ctx context.Context, filter map[string]any) (sharedapi.Cursor[webhooks.Config], error) {
	opts := options.Find().SetSort(bson.M{webhooks.KeyUpdatedAt: -1})
	cur, err := s.configsCollection.Find(ctx, filter, opts)
	if err != nil {
		return sharedapi.Cursor[webhooks.Config]{},
			fmt.Errorf("mongo.Collection.Find: %w", err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			sharedlogging.GetLogger(ctx).Errorf("mongo.Cursor.Close: %s", err)
		}
	}()

	var res []webhooks.Config
	if err := cur.All(ctx, &res); err != nil {
		return sharedapi.Cursor[webhooks.Config]{},
			fmt.Errorf("mongo.Cursor.All: %w", err)
	}

	return sharedapi.Cursor[webhooks.Config]{
		Data: res,
	}, nil
}

func (s Store) InsertOneConfig(ctx context.Context, cfgUser webhooks.ConfigUser) (string, error) {
	cfg := webhooks.Config{
		ConfigUser: cfgUser,
		ID:         uuid.NewString(),
		Active:     true,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	res, err := s.configsCollection.InsertOne(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("store.Collection.InsertOne: %w", err)
	}

	return res.InsertedID.(string), nil
}

func (s Store) DeleteOneConfig(ctx context.Context, id string) (int64, error) {
	res, err := s.configsCollection.DeleteOne(ctx, bson.D{
		{Key: webhooks.KeyID, Value: id},
	})
	if err != nil {
		return 0, fmt.Errorf("momgo.Collection.DeleteOne: %w", err)
	}

	return res.DeletedCount, nil
}

func (s Store) UpdateOneConfigActivation(ctx context.Context, id string, active bool) (*webhooks.Config, int64, error) {
	filter := bson.D{{Key: webhooks.KeyID, Value: id}}
	resFind := s.configsCollection.FindOne(ctx, filter)
	if err := resFind.Err(); err != nil {
		return nil, 0, fmt.Errorf("mongo.Collection.FindOne: %w", err)
	}

	var cfg webhooks.Config
	if err := resFind.Decode(&cfg); err != nil {
		return nil, 0, fmt.Errorf("mongo.SingleResult.Decode: %w", err)
	}

	if cfg.Active == active {
		return &cfg, 0, nil
	}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: webhooks.KeyActive, Value: active},
		{Key: webhooks.KeyUpdatedAt, Value: time.Now().UTC()},
	}}}
	resUpdate, err := s.configsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, 0, fmt.Errorf("mongo.Collection.UpdateOne: %w", err)
	}

	return &cfg, resUpdate.ModifiedCount, nil
}

func (s Store) UpdateOneConfigSecret(ctx context.Context, id, secret string) (int64, error) {
	filter := bson.D{{Key: webhooks.KeyID, Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: webhooks.KeySecret, Value: secret},
		{Key: webhooks.KeyUpdatedAt, Value: time.Now().UTC()},
	}}}
	resUpdate, err := s.configsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("mongo.Collection.UpdateOne: %w", err)
	}

	return resUpdate.ModifiedCount, nil
}

func (s Store) InsertOneRequest(ctx context.Context, req webhooks.Request) (primitive.ObjectID, error) {
	res, err := s.requestsCollection.InsertOne(ctx, req)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("store.Collection.InsertOne: %w", err)
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (s Store) Close(ctx context.Context) error {
	if s.client == nil {
		return nil
	}

	if err := s.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("mongo.Client.Disconnect: %w", err)
	}
	return nil
}
