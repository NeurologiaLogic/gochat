package database

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//harus di set dulu
var connectionString = ""

type DB struct{
	client *mongo.Client
}

var db DB

func initializeDB() (*DB,error){
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil { return nil,err }
	err = client.Ping(ctx, readpref.Primary())
	if err != nil { return nil,err }
	db.client = client
	return &db,nil
}

func GetDB() (*DB){
	if db.client == nil{
		_,err := initializeDB()
		if err != nil { panic(err) }
	}
	return &db
}
