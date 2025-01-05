package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nfb/goblackjack/pkg"
)

type PlayerView struct {
	PlayerHand []pkg.ACard
	DealerHand []pkg.ACard
	State      int
}

func newGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	round := pkg.NewRound()
	round.Play(-1)
	err := round.SaveToRedis()
	if err != nil {
		slog.Error("Failed to save to redis", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	roundBytes, err := json.Marshal(round.CurrentViewableState())
	if err != nil {
		slog.Error("Failed to encode viewable round", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error:" + err.Error()))
	} else {
		w.Write(roundBytes)
	}
}

func showGame(w http.ResponseWriter, r *http.Request) {
	round, err := pkg.LoadFromRedis()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Failed Load from Redis, returing 500")
		return
	}
	roundBytes, err := json.Marshal(round.CurrentViewableState())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Failed encoding round, returing 500")
		return
	}
	w.Write(roundBytes)
}

func playGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	playerCommand, err := strconv.Atoi(string(body))
	if err != nil {
		slog.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Info("Recieved command", playerCommand)

	round, err := pkg.LoadFromRedis()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Failed Load from Redis, returing 500")
		return
	}
	round.Play(playerCommand)

	roundBytes, err := json.Marshal(round.CurrentViewableState())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Failed encoding round, returing 500")
		return
	}
	w.Write(roundBytes)

}

func StartAPI() {
	http.HandleFunc("/new", newGame)
	http.HandleFunc("/show", showGame)
	http.HandleFunc("/play", playGame)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/home/odd/go/src/github.com/nfb/goblackjack/static"))))

	slog.Info("Starting API Server...")
	http.ListenAndServe(":8080", nil)
}
