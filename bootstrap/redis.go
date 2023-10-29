package bootstrap

import (
    "context"
    "github.com/go-redis/redis/v9"
    "go.uber.org/zap"
    "databaseDemo/global"
    "fmt"
)

func InitializeRedis() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     global.App.Config.Redis.Host + ":" + fmt.Sprintf("%d", global.App.Config.Redis.Port),
        Password: global.App.Config.Redis.Password, // no password set
        DB:       global.App.Config.Redis.DB,       // use default DB
    })
    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        global.App.Log.Error("Redis connect ping failed, err:", zap.Any("err", err))
        return nil
    }
    return client
}