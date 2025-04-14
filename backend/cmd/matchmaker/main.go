package matchmaker

import (
	"context"

	"github.com/MommusWinner/MicroDurak/services/matchmaker/config"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func run(ctx context.Context) error {
	config, err := config.Load()
	if err != nil {
		return err
	}

	opt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return err
	}

	client := redis.NewClient(opt)

	res := client.ZRangeByScore(ctx, "matchmaking:queue", &redis.ZRangeBy{Min: "1400", Max: "1600"})
	if res.Err() != nil {
		return res.Err()
	}

	err = client.ZAdd(ctx, "matchmaking:queue", redis.Z{Score: 1500, Member: "hello"}).Err()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	e := echo.New()
	ctx := context.Background()
	if err := run(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
