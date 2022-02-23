package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	//"os"
	"parse-log/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	input := []byte(`
    {
        "request":{
           "method":"GET",
           "uri":"\/",
           "url":"http:\/\/yost.com",
           "size":174,
           "querystring":[
              
           ],
           "headers":{
              "accept":"*\/*",
              "host":"yost.com",
              "user-agent":"curl\/7.37.1"
           }
        },
        "upstream_uri":"\/",
        "response":{
           "status":500,
           "size":878,
           "headers":{
              "Content-Length":"197",
              "via":"kong\/1.3.0",
              "Connection":"close",
              "access-control-allow-credentials":"true",
              "Content-Type":"application\/json",
              "server":"nginx",
              "access-control-allow-origin":"*"
           }
        },
        "authenticated_entity":{
           "consumer_id":{
              "uuid":"72b34d31-4c14-3bae-9cc6-516a0939c9d6"
           }
        },
        "route":{
           "created_at":1564823899,
           "hosts":"miller.com",
           "id":"0636a119-b7ee-3828-ae83-5f7ebbb99831",
           "methods":[
              "GET",
              "POST",
              "PUT",
              "DELETE",
              "PATCH",
              "OPTIONS",
              "HEAD"
           ],
           "paths":[
              "\/"
           ],
           "preserve_host":false,
           "protocols":[
              "http",
              "https"
           ],
           "regex_priority":0,
           "service":{
              "id":"c3e86413-648a-3552-90c3-b13491ee07d6"
           },
           "strip_path":true,
           "updated_at":1564823899
        },
        "service":{
           "connect_timeout":60000,
           "created_at":1563589483,
           "host":"ritchie.com",
           "id":"c3e86413-648a-3552-90c3-b13491ee07d6",
           "name":"ritchie",
           "path":"\/",
           "port":80,
           "protocol":"http",
           "read_timeout":60000,
           "retries":5,
           "updated_at":1563589483,
           "write_timeout":60000
        },
        "latencies":{
           "proxy":1836,
           "kong":8,
           "request":1058
        },
        "client_ip":"75.241.168.121",
        "started_at":1566660387
     }
    `)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to MongoDB")
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logrus.WithError(err).Fatal("Failed to ping MongoDB")
	}

	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	fmt.Println("here")

	//db := mongo.NewStore(client.Database("main"))

	gateway := models.Gateway{}
	if err := json.Unmarshal(input, &gateway); err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(gateway, "", "  ")
	fmt.Println(string(b))

	db := client.Database("main")

	/*r, err := db.Collection("gateway").InsertOne(ctx, gateway)
	fmt.Println(err)
	fmt.Println(r)*/

	var test models.Gateway

	err = db.Collection("gateway").FindOne(ctx, bson.M{"clientip": "75.241.168.121"}).Decode(&test)
	fmt.Println(err)
	fmt.Println(test)
}
