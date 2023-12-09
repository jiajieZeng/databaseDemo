package controller

import (
	"databaseDemo/app/model"
	"databaseDemo/global"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/goccy/go-json"
)

const Ex05UserCountKey = "ex05_user_count"

func tran2User(request model.HashRequest) model.HashUser {
	var user model.HashUser
	userID, err := strconv.ParseInt(request.ID, 10, 32)
	if err != nil {
		panic(err.Error())
	}
	user.UserID = int(userID)
	digg, err := strconv.ParseInt(request.GotDiggCount, 10, 32)
	user.GotDiggCount = int(digg)
	view, err := strconv.ParseInt(request.GotViewCount, 10, 32)
	user.GotViewCount = int(view)
	follower, err := strconv.ParseInt(request.FollowerCount, 10, 32)
	user.FollowerCount = int(follower)
	floowee, err := strconv.ParseInt(request.FolloweeCount, 10, 32)
	user.FolloweeCount = int(floowee)
	followCollectSetCount, err := strconv.ParseInt(request.FollowCollectSetCount, 10, 32)
	user.FollowCollectSetCount = int(followCollectSetCount)
	sub, err := strconv.ParseInt(request.SubscribeTagCount, 10, 32)
	user.SubscribeTagCount = int(sub)
	fmt.Printf("%d %d %d %d %d %d %d", user.UserID, user.FollowCollectSetCount, user.FolloweeCount,
		user.FollowerCount, user.GotDiggCount, user.GotViewCount, user.SubscribeTagCount)
	return user
}

// Ex05 hash数据结果的运用（参考掘金应用）
// go run main.go init 初始化用户计数值
// go run main.go get 1556564194374926  // 打印用户(1556564194374926)的所有计数值
// go run main.go incr_like 1556564194374926 // 点赞数+1
// go run main.go incr_collect 1556564194374926 // 点赞数+1
// go run main.go decr_like 1556564194374926 // 点赞数-1
// go run main.go decr_collect 1556564194374926 // 点赞数-1
func HashData(ctx *gin.Context) {
	var requestBody model.HashRequest
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
	case "add":
		userBody := tran2User(requestBody)
		AddUser(ctx, userBody)
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
	if len(responseData) != 0 {
		jsonData, err := json.Marshal(responseData)
		if err != nil {
			panic(err)
		}
		ctx.Header("Content-Type", "application/json")
		ctx.Status(http.StatusOK)
		ctx.Writer.Write([]byte(jsonData))
		return
	}
	var user model.HashUser
	db := global.App.DB
	db.Where("user_id", userID).First(&user)
	if user.UserID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "not found",
		})
		return
	}
	info := make(map[string]interface{})
	info["follow_collect_set_count"] = fmt.Sprintf("%d", user.FollowCollectSetCount)
	info["folowee_count"] = user.FolloweeCount
	info["follower_count"] = user.FollowerCount
	info["got_digg_count"] = user.GotDiggCount
	info["got_view_count"] = user.GotViewCount
	info["subscribe_tag_count"] = user.SubscribeTagCount
	info["userID"] = user.UserID

	key := GetUserCounterKey(int64(user.UserID))
	_, err = pipe.Del(ctx, key).Result()
	if err != nil {
		panic("pipe Del")
	}
	_, err = pipe.HMSet(ctx, key, info).Result()
	pipe.Exec(ctx)

	jsonData, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	ctx.Header("Content-Type", "application/json")
	ctx.Status(http.StatusOK)
	ctx.Writer.Write([]byte(jsonData))
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
	count, err := global.App.Redis.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	ctx.JSON(200, gin.H{
		"userID": strconv.FormatInt(userID, 10),
		"before": fmt.Sprintf("%s = %s", beforeField, beforeB),
		"after":  fmt.Sprintf("%s = %s", field, count),
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
	dbList := []model.HashUser{
		{UserID: 114514, GotDiggCount: 10693, GotViewCount: 223, FolloweeCount: 76, FollowerCount: 995, FollowCollectSetCount: 575, SubscribeTagCount: 95},
		{UserID: 1111, GotDiggCount: 19, GotViewCount: 2238438, FolloweeCount: 1716, FollowerCount: 98895, FollowCollectSetCount: 75, SubscribeTagCount: 5},
		{UserID: 2222, GotDiggCount: 1238, GotViewCount: 22338, FolloweeCount: 1176, FollowerCount: 85, FollowCollectSetCount: 788, SubscribeTagCount: 99},
		{UserID: 3333, GotDiggCount: 1238, GotViewCount: 38438, FolloweeCount: 1786, FollowerCount: 779895, FollowCollectSetCount: 878, SubscribeTagCount: 18},
		{UserID: 4444, GotDiggCount: 19, GotViewCount: 3438, FolloweeCount: 18, FollowerCount: 1000, FollowCollectSetCount: 10, SubscribeTagCount: 114},
		{UserID: 5555, GotDiggCount: 10693, GotViewCount: 18438, FolloweeCount: 188, FollowerCount: 11495, FollowCollectSetCount: 20, SubscribeTagCount: 541},
		{UserID: 1919810, GotDiggCount: 10693, GotViewCount: 84138, FolloweeCount: 0, FollowerCount: 2547, FollowCollectSetCount: 30, SubscribeTagCount: 810},
	}
	db := global.App.DB
	for _, value := range dbList {
		db.Save(&value)
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

func AddUser(ctx *gin.Context, user model.HashUser) {
	db := global.App.DB
	db.Save(&user)
	counter := map[string]interface{}{
		"userID":                   fmt.Sprintf("%d", user.ID),
		"got_digg_count":           user.GotDiggCount,
		"got_view_count":           user.GotViewCount,
		"followee_count":           user.FolloweeCount,
		"follower_count":           user.FollowerCount,
		"follow_collect_set_count": user.FollowCollectSetCount,
		"subscribe_tag_count":      user.SubscribeTagCount,
	}
	uid, err := strconv.ParseInt(counter["userID"].(string), 10, 64)
	key := GetUserCounterKey(uid)
	pipe := global.App.Redis.Pipeline()
	rw, err := pipe.Del(ctx, key).Result()
	if err != nil {
		fmt.Printf("del %s, rw=%d\n", key, rw)
	}
	_, err = pipe.HMSet(ctx, key, counter).Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("设置 uid=%d, key=%s\n", uid, key)
	// 批量执行上面for循环设置好的hmset命令
	_, err = pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
