package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/flutterbar/chess-explorer-go/internal/pgntodb"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func gameHandler(w http.ResponseWriter, r *http.Request) {

	type gameResponse struct {
		Error string       `json:"error"`
		Data  pgntodb.Game `json:"data"`
	}

	defer timeTrack(time.Now(), "gameHandler")

	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	gameID := strings.TrimSpace(r.FormValue("gameId"))

	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo-url")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping MongoDB
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Cannot connect to DB " + viper.GetString("mongo-url"))
	}

	games := client.Database(viper.GetString("mongo-db-name")).Collection("games")

	result := games.FindOne(ctx, bson.M{"_id": gameID})

	var game pgntodb.Game

	if result != nil {
		result.Decode(&game)
	}

	response := gameResponse{}
	response.Data = game
	json.NewEncoder(w).Encode(response)

}
