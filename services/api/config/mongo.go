package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo เชื่อมต่อ MongoDB และ ping เพื่อยืนยันว่าต่อติดจริง
// คืน *mongo.Database พร้อมใช้งาน ถ้าต่อไม่ได้จะ log.Fatal ทันที
func ConnectMongo(uri, dbName string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("mongo connect failed: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping failed: %v", err)
	}

	log.Printf("connected to MongoDB: %s (db=%s)", uri, dbName)
	return client.Database(dbName)
}
