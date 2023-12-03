package controller

import (
	"databaseDemo/app/model"
	model2 "databaseDemo/app/model"
	"databaseDemo/global"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/goccy/go-json"
)

const Ex06RankKey = "ex06_rank_zset"

type Ex06ItemScore struct {
	ItemNam string
	Score   float64
}

type Value struct {
	Str string
	Num string
}

// Ex06 排行榜
// go run main.go init // 初始化积分
// go run main.go Ex06 rev_order // 输出完整榜单
// go run main.go  Ex06 order_page 0 // 逆序分页输出，offset=1
// go run main.go  Ex06 get_rank user2 // 获取user2的排名
// go run main.go  Ex06 get_score user2 // 获取user2的分数
// go run main.go  Ex06 add_user_score user2 10 // 为user2设置为10分
// zadd ex06_rank_zset 15 andy
// zincrby ex06_rank_zset -9 andy // andy 扣9分，排名掉到最后一名
func Zset(ctx *gin.Context) {
	var requestBody model2.RedisRequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		panic("CheckIn: ctx.ShouldBindJSON failed\n")
	}
	arg1 := requestBody.Command
	switch arg1 {
	case "init":
		Ex06InitUserScore(ctx)
	case "rev_order":
		GetRevOrderAllList(ctx, 0, -1)
	case "get_rank":
		user := requestBody.ItemNam
		GetUserRankByName(ctx, user)
	case "get_score":
		user := requestBody.ItemNam
		GetUserScoreByName(ctx, user)
	case "add_user_score":
		// if len(args) < 3 {
		// 	fmt.Printf("参数错误，可能是缺少需要增加的分值。eg：go run main.go  Ex06 add_user_score user2 10\n")
		// 	return
		// }
		// score, err := strconv.ParseFloat(args[2], 64)
		s := requestBody.Score
		score, err := strconv.ParseFloat(s, 64)
		if err != nil {
			ctx.JSON(403, gin.H{
				"status": err.Error(),
			})
		}
		user := requestBody.ItemNam
		AddUserScore(ctx, user, score)
	case "add":
		var user model.ZsetUser
		user.Name = requestBody.ItemNam
		user.Score, _ = strconv.ParseFloat(requestBody.Score, 64)
		Add(ctx, user)
	}
	return
}

func Ex06InitUserScore(ctx *gin.Context) {
	members := make([]string, 100)
	for i := 1; i <= 100; i++ {
		Name := fmt.Sprintf("user%d", i)
		members[i-1] = Name
	}
	scores := make([]float64, 100)
	rand.Seed(time.Now().UnixNano())
	min := 0.0   // 最小值
	max := 500.0 // 最大值
	for i := 1; i <= 100; i++ {
		scores[i-1] = min + rand.Float64()*(max-min)
	}
	// users := make(map[string]string
	var users []map[string]string
	for i := 0; i < 100; i++ {
		// users[members[i]] = fmt.Sprintf("%f", scores[i])
		rowData := make(map[string]string)
		rowData["user"] = members[i]
		rowData["score"] = fmt.Sprintf("%f", scores[i])
		users = append(users, rowData)
	}
	initList := make([]redis.Z, 100)
	for i := 1; i <= 100; i++ {
		initList[i-1] = redis.Z{Member: members[i-1], Score: scores[i-1]}
	}
	// 清空榜单
	if err := global.App.Redis.Del(ctx, Ex06RankKey).Err(); err != nil {
		panic(err)
	}

	_, err := global.App.Redis.ZAdd(ctx, Ex06RankKey, initList...).Result()
	if err != nil {
		panic(err)
	}

	// fmt.Printf("初始化榜单Item数量:%d\n", nums)
	jsonData, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	ctx.Header("Content-Type", "application/json")
	ctx.Status(http.StatusOK)
	ctx.Writer.Write([]byte(jsonData))
}

// 榜单逆序输出
// ZRANGE ex06_rank_zset +inf -inf BYSCORE  rev WITHSCORES
// 正序输出
// ZRANGE ex06_rank_zset 0 -1 WITHSCORES
func GetRevOrderAllList(ctx *gin.Context, limit, offset int64) {
	resList, err := global.App.Redis.ZRevRangeWithScores(ctx, Ex06RankKey, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	// list := make(map[string]Value)
	var result []map[string]string
	// fmt.Printf("\n榜单:\n")
	for i, z := range resList {
		rowData := make(map[string]string)
		// fmt.Printf("第%d名 %s\t%f\n", i+1, z.Member, z.Score)
		// s := fmt.Sprintf("%d", i + 1);
		// zeros := strings.Repeat("0", 15 - len(s))
		rowData["rank"] = fmt.Sprintf("%d", i+1)
		rowData["name"] = z.Member.(string)
		rowData["score"] = fmt.Sprintf("%f", z.Score)
		// list[fmt.Sprintf("%s%d", zeros, i + 1)] = Value{z.Member.(string), fmt.Sprintf("%f", z.Score)}
		result = append(result, rowData)
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	ctx.Header("Content-Type", "application/json")
	ctx.Status(http.StatusOK)
	ctx.Writer.Write([]byte(jsonData))
}

// GetUserRankByName 获取用户排名
func GetUserRankByName(ctx *gin.Context, name string) {
	rank, err := global.App.Redis.ZRevRank(ctx, Ex06RankKey, name).Result()
	if err != nil {
		// fmt.Errorf("error getting name=%s, err=%v", name, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": fmt.Sprintf("error getting name=%s, err=%v", name, err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"name": name,
		"rank": fmt.Sprintf("%d", rank+1),
	})
	// fmt.Printf("name=%s, 排名=%d\n", name, rank+1)
}

// GetUserScoreByName 获取用户分值
func GetUserScoreByName(ctx *gin.Context, name string) {
	score, err := global.App.Redis.ZScore(ctx, Ex06RankKey, name).Result()
	if err != nil {
		// fmt.Errorf("error getting name=%s, err=%v", name, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": fmt.Sprintf("error getting name=%s, err=%v", name, err),
		})
		return
	}
	// fmt.Println(time.Now().UnixMilli())
	// fmt.Printf("name=%s, 分数=%f\n", name, score)
	ctx.JSON(http.StatusOK, gin.H{
		"name": name,
		"rank": fmt.Sprintf("%f", score),
	})
}

// AddUserScore 排名用户
func AddUserScore(ctx *gin.Context, name string, score float64) {
	num, err := global.App.Redis.ZIncrBy(ctx, Ex06RankKey, score, name).Result()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("name=%s, add_score=%f, score=%f\n", name, score, num)
	ctx.JSON(http.StatusOK, gin.H{
		"name":      name,
		"add_score": fmt.Sprintf("%f", score),
		"score":     fmt.Sprintf("%f", num),
	})
}

func Add(ctx *gin.Context, user model.ZsetUser) {
	redisData := redis.Z{Member: user.Name, Score: user.Score}
	_, err := global.App.Redis.ZAdd(ctx, Ex06RankKey, redisData).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
