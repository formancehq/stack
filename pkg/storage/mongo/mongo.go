package mongo

import (
	"context"
	"time"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/formancehq/webhooks/cmd/flag"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type Store struct {
	uri                string
	client             *mongo.Client
	configsCollection  *mongo.Collection
	attemptsCollection *mongo.Collection
}

var _ storage.Store = &Store{}

func NewStore() (storage.Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDBUri := viper.GetString(flag.StorageMongoConnString)
	sharedlogging.Infof("connecting to mongoDB URI: %s", mongoDBUri)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBUri).SetMonitor(otelmongo.NewMonitor()))
	if err != nil {
		return Store{}, errors.Wrap(err, "mongo.Connect")
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return Store{}, errors.Wrap(err, "mongo.Client.Ping")
	}

	return Store{
		uri:    mongoDBUri,
		client: client,
		configsCollection: client.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionConfigs),
		attemptsCollection: client.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts),
	}, nil
}

func (s Store) FindManyConfigs(ctx context.Context, filter map[string]any) ([]webhooks.Config, error) {
	res := []webhooks.Config{}
	opts := options.Find().SetSort(bson.M{webhooks.KeyUpdatedAt: -1})
	cur, err := s.configsCollection.Find(ctx, filter, opts)
	if err != nil {
		return res, errors.Wrap(err, "mongo.Collection.Find")
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &res); err != nil {
		return []webhooks.Config{}, errors.Wrap(err, "mongo.Cursor.All")
	}

	return res, nil
}

func (s Store) InsertOneConfig(ctx context.Context, cfgUser webhooks.ConfigUser) (webhooks.Config, error) {
	cfg := webhooks.NewConfig(cfgUser)
	_, err := s.configsCollection.InsertOne(ctx, cfg)
	if err != nil {
		return webhooks.Config{}, errors.Wrap(err, "store.Collection.InsertOne")
	}

	return cfg, nil
}

func (s Store) DeleteOneConfig(ctx context.Context, id string) error {
	res, err := s.configsCollection.DeleteOne(ctx, bson.D{
		{Key: webhooks.KeyID, Value: id},
	})
	if err != nil {
		return errors.Wrap(err, "momgo.Collection.DeleteOne")
	}

	if res.DeletedCount == 0 {
		return storage.ErrConfigNotFound
	}

	return nil
}

func (s Store) UpdateOneConfigActivation(ctx context.Context, id string, active bool) (webhooks.Config, error) {
	cfg := webhooks.Config{}
	filter := bson.D{{Key: webhooks.KeyID, Value: id}}
	if err := s.configsCollection.FindOne(ctx, filter).Decode(&cfg); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return webhooks.Config{}, storage.ErrConfigNotFound
		}
		return webhooks.Config{}, errors.Wrap(err, "decode config")
	}
	if cfg.Active == active {
		return webhooks.Config{}, storage.ErrConfigNotModified
	}

	filter = bson.D{
		{Key: webhooks.KeyID, Value: id},
		{Key: webhooks.KeyActive, Value: !active},
	}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: webhooks.KeyActive, Value: active},
		{Key: webhooks.KeyUpdatedAt, Value: time.Now().UTC()},
	}}}
	if _, err := s.configsCollection.UpdateOne(ctx, filter, update); err != nil {
		return webhooks.Config{}, errors.Wrap(err, "mongo.Collection.UpdateOne")
	}

	cfg.Active = active
	return cfg, nil
}

func (s Store) UpdateOneConfigSecret(ctx context.Context, id, secret string) (webhooks.Config, error) {
	cfg := webhooks.Config{}
	filter := bson.D{{Key: webhooks.KeyID, Value: id}}
	if err := s.configsCollection.FindOne(ctx, filter).Decode(&cfg); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return webhooks.Config{}, storage.ErrConfigNotFound
		}
		return webhooks.Config{}, errors.Wrap(err, "decode updated config")
	}
	if cfg.Secret == secret {
		return webhooks.Config{}, storage.ErrConfigNotModified
	}

	filter = bson.D{
		{Key: webhooks.KeyID, Value: id},
		{Key: webhooks.KeySecret, Value: bson.D{
			{Key: "$ne", Value: secret},
		}},
	}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: webhooks.KeySecret, Value: secret},
		{Key: webhooks.KeyUpdatedAt, Value: time.Now().UTC()},
	}}}
	if _, err := s.configsCollection.UpdateOne(ctx, filter, update); err != nil {
		return webhooks.Config{}, errors.Wrap(err, "mongo.Collection.UpdateOne")
	}

	cfg.Secret = secret
	return cfg, nil
}

func (s Store) FindManyAttempts(ctx context.Context, filter map[string]any) ([]webhooks.Attempt, error) {
	res := []webhooks.Attempt{}
	opts := options.Find().SetSort(bson.M{webhooks.KeyID: -1})
	cur, err := s.attemptsCollection.Find(ctx, filter, opts)
	if err != nil {
		return res, errors.Wrap(err, "mongo.Collection.Find")
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &res); err != nil {
		return res, errors.Wrap(err, "mongo.Cursor.All")
	}

	return res, nil
}

func (s Store) FindDistinctWebhookIDs(ctx context.Context, filter map[string]any) ([]string, error) {
	dis, err := s.attemptsCollection.Distinct(ctx, webhooks.KeyWebhookID, filter, nil)
	if err != nil {
		return nil, errors.Wrap(err, "mongo.Collection.Distinct")
	}

	res := make([]string, len(dis))
	for i, d := range dis {
		res[i] = d.(string)
	}

	return res, nil
}

func (s Store) UpdateManyAttemptsStatus(ctx context.Context, webhookID, status string) ([]webhooks.Attempt, error) {
	atts, err := s.FindManyAttempts(ctx, map[string]any{webhooks.KeyWebhookID: webhookID})
	if err != nil {
		return []webhooks.Attempt{}, errors.Wrap(err, "mongo.Collection.UpdateMany")
	}
	if len(atts) == 0 {
		return []webhooks.Attempt{}, storage.ErrAttemptIDNotFound
	}

	filter := bson.D{
		{Key: webhooks.KeyWebhookID, Value: webhookID},
		{Key: webhooks.KeyStatus, Value: bson.D{
			{Key: "$ne", Value: status},
		}},
	}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: webhooks.KeyStatus, Value: status},
	}}}

	res, err := s.attemptsCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return []webhooks.Attempt{}, errors.Wrap(err, "mongo.Collection.UpdateMany")
	}
	if res.ModifiedCount == 0 {
		return []webhooks.Attempt{}, storage.ErrAttemptNotModified
	}

	for i := range atts {
		atts[i].Status = status
	}
	return atts, nil
}

func (s Store) InsertOneAttempt(ctx context.Context, att webhooks.Attempt) error {
	_, err := s.attemptsCollection.InsertOne(ctx, att)
	if err != nil {
		return errors.Wrap(err, "store.Collection.InsertOne")
	}

	return nil
}

func (s Store) Close(ctx context.Context) error {
	if s.client == nil {
		return nil
	}

	if err := s.client.Disconnect(ctx); err != nil {
		return errors.Wrap(err, "mongo.Client.Disconnect")
	}

	return nil
}
