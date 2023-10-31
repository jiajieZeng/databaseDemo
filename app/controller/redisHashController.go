package controller

import (
	"fmt"
	"strconv"
	model2 "databaseDemo/app/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/goccy/go-json"
	"databaseDemo/global"
	"github.com/go-redis/redis/v9"
)

const Ex05UserCountKey = "ex05_user_count"

// Ex05 hash数据结果的运用（参考掘金应用）
// go run main.go init 初始化用户计数值
// go run main.go get 1556564194374926  // 打印用户(1556564194374926)的所有计数值
// go run main.go incr_like 1556564194374926 // 点赞数+1
// go run main.go incr_collect 1556564194374926 // 点赞数+1
// go run main.go decr_like 1556564194374926 // 点赞数-1
// go run main.go decr_collect 1556564194374926 // 点赞数-1
func HashData(ctx *gin.Context) {
	var requestBody model2.RedisRequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		panic("CheckIn: ctx.ShouldBindJSON failed\n")
	}
	command := requestBody.Command
	switch command {
	case "init":
		Ex06InitUserCounter(ctx)
	case "get":
		id := requestBody.ID
		userID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": err.Error(),
			})
			return
		}
		GetUserCounter(ctx, userID)
	case "incr_like":
		id := requestBody.ID
		userID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": err.Error(),
			})
			return
		}
		IncrByUserLike(ctx, userID)
	case "incr_collect":
		id := requestBody.ID
		userID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": err.Error(),
			})
			return
		}
		IncrByUserCollect(ctx, userID)
	case "decr_like":
		id := requestBody.ID
		userID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": err.Error(),
			})
			return
		}
		DecrByUserLike(ctx, userID)
	case "decr_collect":
		id := requestBody.ID
		userID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": err.Error(),
			})
			return
		}
		DecrByUserCollect(ctx, userID)
	}
}


func GetUserCounterKey(userID int64) string {
	return fmt.Sprintf("%s_%d", Ex05UserCountKey, userID)
}

func GetUserCounter(ctx *gin.Context, userID int64) {
	pipe := global.App.Redis.Pipeline()
	GetUserCounterKey(userID)
	pipe.HGetAll(ctx, GetUserCounterKey(userID))
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
	responseData := make(map[string]string)  
	for _, cmder := range cmders {  
		counterMap, err := cmder.(*redis.MapStringStringCmd).Result()  
		if err != nil {  
			panic(err)  
		}  
		for field, value := range counterMap {  
			// 将counterMap的数据添加到responseData中  
			responseData[field] = value  
		}  
	} 
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}
	ctx.Header("Content-Type", "application/json")
	ctx.Status(http.StatusOK)
	ctx.Writer.Write([]byte(jsonData))
	// for _, cmder := range cmders {
	// 	counterMap, err := cmder.(*redis.MapStringStringCmd).Result()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for field, value := range counterMap {
	// 		fmt.Printf("%s: %s\n", field, value)
	// 	}
	// }
}

// IncrByUserLike 点赞数+1
func IncrByUserLike(ctx *gin.Context, userID int64) {
	incrByUserField(ctx, userID, "got_digg_count")
}

// IncrByUserCollect 收藏数+1
func IncrByUserCollect(ctx *gin.Context, userID int64) {
	incrByUserField(ctx, userID, "follow_collect_set_count")
}

// DecrByUserLike 点赞数-1
func DecrByUserLike(ctx *gin.Context, userID int64) {
	decrByUserField(ctx, userID, "got_digg_count")
}

// DecrByUserCollect 收藏数-1
func DecrByUserCollect(ctx *gin.Context, userID int64) {
	decrByUserField(ctx, userID, "follow_collect_set_count")
}

func incrByUserField(ctx *gin.Context, userID int64, field string) {
	change(ctx, userID, field, 1)
}

func decrByUserField(ctx *gin.Context, userID int64, field string) {
	change(ctx, userID, field, -1)
}

func change(ctx *gin.Context, userID int64, field string, incr int64) {
	redisKey := GetUserCounterKey(userID)
	before, err := global.App.Redis.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	beforeInt, err := strconv.ParseInt(before, 10, 64)
	if err != nil {
		panic(err)
	}
	if beforeInt+incr < 0 {
		// fmt.Printf("禁止变更计数，计数变更后小于0. %d + (%d) = %d\n", beforeInt, incr, beforeInt+incr)
		ctx.JSON(403, gin.H{
			"status": fmt.Sprintf("禁止变更计数，计数变更后小于0. %d + (%d) = %d\n", beforeInt, incr, beforeInt+incr),
		})
		return
	}
	beforeField := field
	beforeB := before
	// fmt.Printf("user_id: %d\n更新前\n%s = %s\n--------\n", userID, field, before)
	_, err = global.App.Redis.HIncrBy(ctx, redisKey, field, incr).Result()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("更新记录[%d]:%d\n", userID, num)
	count, err := global.App.Redis.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("user_id: %d\n更新后\n%s = %s\n--------\n", userID, field, count)
	ctx.JSON(200, gin.H{
		"userID": strconv.FormatInt(userID, 10),
		"before": fmt.Sprintf("%s = %s", beforeField, beforeB),
		"after": fmt.Sprintf("%s = %s", field, count),
	})
}

func Ex06InitUserCounter(ctx *gin.Context) {
	pipe := global.App.Redis.Pipeline()
	userCounters := []map[string]interface{}{
		{"userID": "114514", "got_digg_count": 10693, "got_view_count": 223, "followee_count": 76, "follower_count": 995, "follow_collect_set_count": 575, "subscribe_tag_count": 95},
		{"userID": "1111", "got_digg_count": 19, "got_view_count": 2238438, "followee_count": 1716, "follower_count": 98895, "follow_collect_set_count": 75, "subscribe_tag_count": 5},
		{"userID": "2222", "got_digg_count": 1238, "got_view_count": 22338, "followee_count": 1176, "follower_count": 85, "follow_collect_set_count": 788, "subscribe_tag_count": 99},
		{"userID": "3333", "got_digg_count": 1238, "got_view_count": 38438, "followee_count": 1786, "follower_count": 779895, "follow_collect_set_count": 878, "subscribe_tag_count": 18},
		{"userID": "4444", "got_digg_count": 19, "got_view_count": 3438, "followee_count": 18, "follower_count": 1000, "follow_collect_set_count": 10, "subscribe_tag_count": 114},
		{"userID": "5555", "got_digg_count": 10693, "got_view_count": 18438, "followee_count": 188, "follower_count": 11495, "follow_collect_set_count": 20, "subscribe_tag_count": 541},
		{"userID": "1919810", "got_digg_count": 10693, "got_view_count": 84138, "followee_count": 0, "follower_count": 2547, "follow_collect_set_count": 30, "subscribe_tag_count": 810},
	}
	for _, counter := range userCounters {
		uid, err := strconv.ParseInt(counter["userID"].(string), 10, 64)
		key := GetUserCounterKey(uid)
		rw, err := pipe.Del(ctx, key).Result()
		if err != nil {
			fmt.Printf("del %s, rw=%d\n", key, rw)
		}
		_, err = pipe.HMSet(ctx, key, counter).Result()
		if err != nil {
			panic(err)
		}

		fmt.Printf("设置 uid=%d, key=%s\n", uid, key)
	}
	// 批量执行上面for循环设置好的hmset命令
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": err.Error(),
			})
			return
		}
	}
	sData := []map[string]interface{}{
		{"userID": "114514", "got_digg_count": "10693", "got_view_count": "223", "followee_count": "76", "follower_count": "995", "follow_collect_set_count": "575", "subscribe_tag_count": "95"},
		{"userID": "1111", "got_digg_count": "19", "got_view_count": "2238438", "followee_count": "1716", "follower_count": "98895", "follow_collect_set_count": "75", "subscribe_tag_count": "5"},
		{"userID": "2222", "got_digg_count": "1238", "got_view_count": "22338", "followee_count": "1176", "follower_count": "85", "follow_collect_set_count": "788", "subscribe_tag_count": "99"},
		{"userID": "3333", "got_digg_count": "1238", "got_view_count": "38438", "followee_count": "1786", "follower_count": "779895", "follow_collect_set_count": "878", "subscribe_tag_count": "18"},
		{"userID": "4444", "got_digg_count": "19", "got_view_count": "3438", "followee_count": "18", "follower_count": "1000", "follow_collect_set_count": "10", "subscribe_tag_count": "114"},
		{"userID": "5555", "got_digg_count": "10693", "got_view_count": "18438", "followee_count": "188", "follower_count": "11495", "follow_collect_set_count": "20", "subscribe_tag_count": "541"},
		{"userID": "1919810", "got_digg_count": "10693", "got_view_count": "84138", "followee_count": "0", "follower_count": "2547", "follow_collect_set_count": "30", "subscribe_tag_count": "810"},
	}
	jsonData, err := json.Marshal(sData)
	if err != nil {
		panic(err)
	}
	ctx.Header("Content-Type", "application/json")
	ctx.Status(http.StatusOK)
	ctx.Writer.Write([]byte(jsonData))
}