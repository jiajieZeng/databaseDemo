package controller

import (
	"context"
	model2 "databaseDemo/app/model"
	"databaseDemo/global"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

var ctx = context.Background()

const continuesCheckKey = "cc_uid_%s"
const bitmapKey = "%s_b"

func CheckIn(ctx *gin.Context) {
	var requestBody model2.RedisRequestBody
	RedisClient := global.App.Redis
	if err := ctx.ShouldBindBodyWith(&requestBody, binding.JSON); err != nil {
		// panic("CheckIn: ctx.ShouldBindJSON failed\n")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
		return
	}
	userID := requestBody.ID
	key := fmt.Sprintf(continuesCheckKey, userID)
	bKey := fmt.Sprintf(bitmapKey, key)
	flag, err := checkSign(ctx, bKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
		return
	}
	if flag == 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": fmt.Sprintf("用户[%s]今天已经签过到了", userID),
		})
		return
	}
	// 1. 连续签到数+1
	err = RedisClient.Incr(ctx, key).Err()
	if err != nil {
		// fmt.Errorf("用户[%d]连续签到失败", userID)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": fmt.Sprintf("用户[%s]连续签到失败", userID),
		})
		global.App.Log.Error("用户连续签到失败", zap.Any("err", err))
		return
	} else {
		expAt := beginningOfDay().Add(2 * time.Minute)
		// 2. 设置签到记录在后天的0点到期
		if err := RedisClient.ExpireAt(ctx, key, expAt).Err(); err != nil {
			panic(err)
		} else {
			// 3. 打印用户续签后的连续签到天数
			id, err := strconv.ParseInt(userID, 10, 64)
			if err != nil {
				ctx.JSON(403, gin.H{
					"status": err.Error(),
				})
				return
			}
			day, err := getUserCheckInDays(ctx, id)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status": err.Error(),
				})
				return
			}
			var offset int = time.Now().Local().Minute() - 1
			_, err = RedisClient.SetBit(ctx, bKey, int64(offset), 1).Result()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"status": fmt.Sprintf("用户[%s]连续签到：%d(天), 过期时间:%s", userID, day, expAt.Format("2006-01-02 15:04:05")),
			})
			return
		}
	}
}

// getUserCheckInDays 获取用户连续签到天数
func getUserCheckInDays(ctx context.Context, userID int64) (int64, error) {
	key := fmt.Sprintf("cc_uid_%d", userID)
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

func checkSign(ctx context.Context, key string) (int64, error) {
	var offset int = time.Now().Local().Minute() - 1
	RedisClient := global.App.Redis
	return RedisClient.GetBit(ctx, key, int64(offset)).Result()
}

// beginningOfDay 获取今天0点时间
func beginningOfDay() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	h := now.Hour()
	min := now.Minute()
	return time.Date(y, m, d, h, min, 0, 0, time.Local)
}
