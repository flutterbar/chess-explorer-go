package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/flutterbar/chess-explorer-go/internal/pgntodb"
	"github.com/notnil/chess"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type searchFENResult struct {
	game   pgntodb.Game
	moveId int
}

func searchFentHandler(w http.ResponseWriter, r *http.Request) {
	defer timeTrack(time.Now(), "searchFentHandler")

	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
	default:
		fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		return
	}

	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// create game filter
	filter := gameFilterFromRequest(r)
	gameFilterBson := bsonFromGameFilter(filter)

	fen := strings.TrimSpace(r.FormValue("fen"))
	maxMoves, _ := strconv.Atoi(r.FormValue("maxMoves"))

	go searchFEN(fen, maxMoves, gameFilterBson) // launch background job and return immediately
}

func searchFEN(fen string, maxMoves int, gameFilterBson primitive.M) {
	log.Println("Searching for FEN: " + fen)
	log.Println("Maximum", maxMoves, "moves per games")

	// start a ticker
	ticker := time.NewTicker(15000 * time.Millisecond)
	tickerChannel := make(chan bool)
	go func() {
		for {
			select {
			case <-tickerChannel:
				return
			case <-ticker.C:
				log.Println("Searching for FEN ...")
			}
		}
	}()

	// start the log accumulator
	logChannel := make(chan *searchFENResult)
	go func() {
		var logs []*searchFENResult
		for {
			item := <-logChannel
			if item != nil {
				logs = append(logs, item)
			} else {
				log.Println(strconv.Itoa(len(logs)) + " hits")
				winWins, blackWins, draw := 0, 0, 0
				for _, logItem := range logs {
					log.Println("move " + strconv.Itoa(logItem.moveId) + " in game " + logItem.game.Link + " " + logItem.game.Result)
					switch logItem.game.Result {
					case "1-0":
						winWins = winWins + 1
					case "0-1":
						blackWins = blackWins + 1
					default:
						draw = draw + 1
					}
				}
				log.Println("1-0: " + strconv.Itoa(winWins) + ", 0-1: " + strconv.Itoa(blackWins) + ", 1/2-1/2: " + strconv.Itoa(draw))
				return
			}
		}
	}()

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

	gamesCollection := client.Database(viper.GetString("mongo-db-name")).Collection("games")

	cur, error := gamesCollection.Find(ctx, gameFilterBson)
	if error != nil {
		log.Fatal(err)
	}

	concurrency := 20
	concurrencyChannel := make(chan bool, concurrency)

	count := 0
	for cur.Next(context.TODO()) {
		var gameHolder pgntodb.Game
		err := cur.Decode(&gameHolder)

		concurrencyChannel <- true // take a slot
		go replay(gameHolder, fen, maxMoves, concurrencyChannel, logChannel)

		if err != nil {
			log.Fatal(err)
		}
		count++
	}

	// wait for everything to be finished
	for i := 0; i < cap(concurrencyChannel); i++ {
		concurrencyChannel <- true
	}

	log.Printf("replayed " + strconv.Itoa(count) + " games")

	// stop the ticker
	ticker.Stop()
	tickerChannel <- true

	// dump the logs
	logChannel <- nil
}

func replay(game pgntodb.Game, fen string, maxMoves int, concurrencyChannel chan bool, logChannel chan *searchFENResult) {

	defer func() { <-concurrencyChannel }() // release the slot when finished

	// Process game.PGN (remove "1." etc)
	var pgnMoves []string
	if len(game.PGN) > 0 {
		pgnMoves = strings.Split(game.PGN, " ")
	}

	i := 0 // output index
	for _, x := range pgnMoves {
		if !strings.HasSuffix(x, ".") {
			// copy and increment index
			pgnMoves[i] = x
			i++
		}
	}
	pgnMoves = pgnMoves[:i] // strip final result

	// Replay game
	chessGame := chess.NewGame()
	iMove := 0
	for _, move := range pgnMoves {
		chessGame.MoveStr(move)

		// Compare
		if chessGame.Position().String() == fen {
			iMove++
			logChannel <- &searchFENResult{game: game, moveId: iMove}
			break
		}

		iMove++
		if iMove == maxMoves {
			break
		}
	}
}
