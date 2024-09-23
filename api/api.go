package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nfb/goblackjack/pkg"
	"github.com/redis/go-redis/v9"
)

type PlayerView struct {
	PlayerHand []pkg.ACard
	DealerHand []pkg.ACard
	State      int
}

func newGame(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		round := pkg.NewRound()
		round.Play(-1)
		roundBytes, err := json.Marshal(round)
		if err != nil {
			w.Write([]byte("error:" + err.Error()))
		} else {
			w.Write(roundBytes)
			client := redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "", // no password set
				DB:       0,  // use default DB
			})
			ctx := context.Background()

			client.Set(ctx, "game", string(roundBytes), 0)
		}

	}
}

func showGame(w http.ResponseWriter, r *http.Request) {
	var round pkg.BlackJackRound
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()
	resp := client.Get(ctx, "game")
	roundBytes, err := resp.Bytes()
	if err != nil {
		fmt.Println(err)
	} else {
		err = json.Unmarshal(roundBytes, &round)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(roundBytes)
	}
}

func StartAPI() {
	http.HandleFunc("/new", newGame)
	http.HandleFunc("/show", showGame)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/home/odd/go/src/github.com/nfb/goblackjack/static"))))

	slog.Info("Starting API Server...")
	http.ListenAndServe(":8080", nil)
}
