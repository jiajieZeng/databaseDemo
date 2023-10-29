package controller

import (
	"fmt"
	model2 "databaseDemo/app/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/goccy/go-json"
	"databaseDemo/global"
	"github.com/go-redis/redis/v9"
	"math/rand"
	"time"
)

const Ex06RankKey = "ex06_rank_zset"

type Ex06ItemScore struct {
	ItemNam string
	Score   float64
}

type Value struct {  
    Str  string  
    Num  float64  
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
func Ex06(ctx *gin.Context) {
	var requestBody model2.ZsetBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		panic("CheckIn: ctx.ShouldBindJSON failed\n")
	}
	arg1 := requestBody.Command
	switch arg1 {
	case "init":
		Ex06InitUserScore(ctx)
	case "rev_order":
		GetRevOrderAllList(ctx, 0, -1)
	case "order_page":
		pageSize := int64(2)
		offset := requestBody.offset
		GetOrderListByPage(ctx, offset, pageSize)
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
		score := requestBody.Score
		user := requestBody.ItemNam
		AddUserScore(ctx, user, score)
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
		scores[i - 1] = min + rand.Float64()*(max-min)
	}  	
	users := make(map[string]float64) 
	for i := 0; i < 100; i++ {
		users[members[i]] = scores[i]
	}
	initList := make([]redis.Z, 100)
	for i := 1; i <= 100; i++ {
		initList[i - 1] = redis.Z{Member: members[i - 1], Score: scores[i - 1]}
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
	list := make(map[string]Value) 
	// fmt.Printf("\n榜单:\n")
	for i, z := range resList {
		// fmt.Printf("第%d名 %s\t%f\n", i+1, z.Member, z.Score)
		list[fmt.Sprintf("第%d名", i + 1)] = Value{string(z.Member), z.Score}
	}
	jsonData, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}
	ctx.Header("Content-Type", "application/json")
	ctx.Status(http.StatusOK)
	ctx.Writer.Write([]byte(jsonData))
}

func GetOrderListByPage(ctx *gin.Context, offset, pageSize int64) {
	// zrange ex06_rank_zset 300 0 byscore rev limit 1 2 withscores // 取300分到0分之间的排名
	// zrange ex06_rank_zset -inf +inf byscore withscores 正序输出
	// ZRANGE ex06_rank_zset +inf -inf BYSCORE  REV WITHSCORES 逆序输出所有排名
	// zrange ex06_rank_zset +inf -inf byscore rev limit 0 2 withscores 逆序分页输出排名
	zRangeArgs := redis.ZRangeArgs{
		Key:     Ex06RankKey,
		ByScore: true,
		Rev:     true,
		Start:   "-inf",
		Stop:    "+inf",
		Offset:  offset,
		Count:   pageSize,
	}
	resList, err := global.App.Redis.ZRangeArgsWithScores(ctx, zRangeArgs).Result()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("\n榜单(offest=%d, pageSize=%d):\n", offset, pageSize)
	offNum := int(pageSize * offset)
	// for i, z := range resList {
	// 	rank := i + 1 + offNum
	// 	fmt.Printf("第%d名 %s\t%f\n", rank, z.Member, z.Score)
	// }
	// fmt.Println()
	list := make(map[string]Value) 
	// fmt.Printf("\n榜单:\n")
	for i, z := range resList {
		// fmt.Printf("第%d名 %s\t%f\n", i+1, z.Member, z.Score)
		list[fmt.Sprintf("第%d名", i + 1 + offNum)] = Value{string(z.Member), z.Score}
	}
	jsonData, err := json.Marshal(list)
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
		"rank": rank + 1,
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
		"rank": score,
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
		"name": name,
		"add_score": score,
		"score": num,
	})
}