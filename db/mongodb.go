package db

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB connection parameters
const (
    DbName           = "product_management"
    UserCollection   = "users"
    ProductCollection = "products"
)

type MongoDB struct {
	Client *mongo.Client
}


// NewMongoDB creates a new MongoDB connection
func NewMongoDB(connectionString string) (*MongoDB, error) {
   
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()
    
    // log.Println("Using MongoDB connection string:", connectionString)

    clientOptions := options.Client().ApplyURI(connectionString)

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
		client.Disconnect(ctx)
		return nil, err
	}

	log.Println("Connected to MongoDB!")

    return &MongoDB{Client: client}, nil
}


func (db *MongoDB) Close() {
    if db.Client != nil {
        db.Client.Disconnect(context.Background())
    }
}
