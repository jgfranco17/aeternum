package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jgfranco17/aeternum/api/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	exec "github.com/jgfranco17/aeternum/execution"
)

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoClient(ctx context.Context, uri string, username string, token string) (*MongoClient, error) {
	log := logging.FromContext(ctx)

	serverAPI := mongo_options.ServerAPI(mongo_options.ServerAPIVersion1)
	appliedUri := fmt.Sprintf("mongodb+srv://%s:%s@%s", username, token, uri)
	opts := mongo_options.Client().ApplyURI(appliedUri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping MongoDB: %w", err)
	}
	log.Debugf("Connection to database secured")
	return &MongoClient{Client: client, Database: client.Database("test")}, nil
}

func (m *MongoClient) Disconnect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return m.Client.Disconnect(ctx)
}

func (m *MongoClient) GetResult(ctx context.Context, id string) (*exec.CheckResponse, error) {
	var result exec.CheckResponse
	database := m.Client.Database("tests")
	if database == nil {
		return nil, fmt.Errorf("No test database found")
	}
	collection := database.Collection("results")
	if collection == nil {
		return nil, fmt.Errorf("No results collection in test database found")
	}
	err := collection.FindOne(ctx, bson.D{{"id", id}}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("Failed to fetch test results: %w", err)
	}
	return &result, nil
}
