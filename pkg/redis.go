package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func LoadFromRedis() (*BlackJackRound, error) {
	var round BlackJackRound
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()
	resp := client.Get(ctx, "game")
	roundBytes, err := resp.Bytes()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	err = json.Unmarshal(roundBytes, &round)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &round, nil
}

func (round BlackJackRound) SaveToRedis() error {
	roundBytes, err := json.Marshal(round)
	if err != nil {
		return errors.New("Failed Encoding Round")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()

	client.Set(ctx, "game", string(roundBytes), 0)
	return nil
}
