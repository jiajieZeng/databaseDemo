package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"
	model2 "databaseDemo/app/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"databaseDemo/global"
)

var ctx = context.Background()

const continuesCheckKey = "cc_uid_%d"

func CheckIn(ctx *gin.Context) {
	var requestBody model2.RedisRequestBody
	RedisClient := global.App.Redis
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		panic("CheckIn: ctx.ShouldBindJSON failed\n")
	}
	userID := requestBody.ID
	key := fmt.Sprintf(continuesCheckKey, userID)

	// 锁一天
	now := time.Now()
	expAt := beginningOfDay().Add(48 * time.Hour)
	secondsUnitlTomorrow := int(expAt.Sub(now).Seconds())
	fmt.Println(time.Duration(secondsUnitlTomorrow * 1000) * time.Millisecond)
	acquired, err := RedisClient.SetNX(ctx, key, 1, time.Duration(secondsUnitlTomorrow * 1000) * time.Millisecond).Result()
	if err != nil {
		// fmt.Errorf("用户[%d]连续签到失败", userID)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": fmt.Sprintf("用户[%d]连续签到失败", userID),
		})
		return
	} 
	if  !acquired {
		ctx.JSON(403, gin.H{
			"status": fmt.Sprintf("用户[%d]今日已签到, 过期时间:%s", userID, expAt.Format("2006-01-02 15:04:05")),
		})		
		return
	}

	// 1. 连续签到数+1
	err = RedisClient.Incr(ctx, key).Err()
	if err != nil {
		// fmt.Errorf("用户[%d]连续签到失败", userID)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": fmt.Sprintf("用户[%d]连续签到失败", userID),
		})
	} else {
		// 2. 设置签到记录在后天的0点到期
		if err := RedisClient.ExpireAt(ctx, key, expAt).Err(); err != nil {
			panic(err)
		} else {
			// 3. 打印用户续签后的连续签到天数
			day, err := getUserCheckInDays(ctx, userID)
			if err != nil {
				panic(err)
			}
			// fmt.Printf("用户[%d]连续签到：%d(天), 过期时间:%s", userID, day, expAt.Format("2006-01-02 15:04:05"))
			ctx.JSON(http.StatusOK, gin.H{
				"status": fmt.Sprintf("用户[%d]连续签到：%d(天), 过期时间:%s", userID, day, expAt.Format("2006-01-02 15:04:05")),
			})
		}
	}
}

// getUserCheckInDays 获取用户连续签到天数
func getUserCheckInDays(ctx context.Context, userID int64) (int64, error) {
	key := fmt.Sprintf(continuesCheckKey, userID)
	RedisClient := global.App.Redis
	days, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if daysInt, err := strconv.ParseInt(days, 10, 64); err != nil {
		panic(err)
	} else {
		return daysInt, nil
	}
}

// beginningOfDay 获取今天0点时间
func beginningOfDay() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}
