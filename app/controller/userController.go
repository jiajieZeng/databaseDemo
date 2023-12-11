package controller

import (
	model2 "databaseDemo/app/model"
	"databaseDemo/global"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {

	db := global.App.DB

	//获取参数
	//此处使用Bind()函数，可以处理不同格式的前端数据
	var requestUser model2.User
	ctx.Bind(&requestUser)
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	global.App.Log.Info("Register: get info success" + fmt.Sprintf("  %s %s %s", name, telephone, password))
	//数据验证
	if len(name) == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户名不能为空",
		})
		return
	}
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "手机号必须为11位",
		})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码不能少于6位",
		})
		return
	}
	global.App.Log.Info("Register: check success")
	//判断手机号是否存在
	var user model2.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户已存在",
		})
		return
	}
	global.App.Log.Info("Register: check exist success")
	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    500,
			"message": "密码加密错误",
		})
		return
	}
	global.App.Log.Info("Register: create user success")
	newUser := model2.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	db.Create(&newUser)
	global.App.Log.Info("Register: insert into db")
	//返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}

func QueryFirst(ctx *gin.Context) {
	db := global.App.DB
	var user model2.User
	db.First(&user)
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"code":       201,
		"name":       user.Name,
		"tlelephone": user.Telephone,
		"password":   user.Password,
	})
}

func RawSQL(ctx *gin.Context) {

	db := global.App.RDB
	var requestBody model2.RequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}
	sql := requestBody.SQL
	if sql[0] == 's' || sql[0] == 'S' {
		rows, err := db.Query(sql)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		defer rows.Close()
		columns, err := rows.Columns()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}

		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		var result []map[string]string
		for rows.Next() {
			err := rows.Scan(values...)
			if err != nil {
				panic(err)
			}

			// 打印查询结果及其字段名
			rowData := make(map[string]string)
			for i, v := range values {
				fmt.Printf("%s: ", columns[i])
				s, ok := (*v.(*interface{})).([]byte)
				rowData[columns[i]] = string(s)
				if !ok {
					// 处理类型断言失败的情况
					continue
				}
				for _, val := range s {
					fmt.Printf("%c", val)
				}
				fmt.Println()
			}
			result = append(result, rowData)
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		ctx.Header("Content-Type", "application/json")
		ctx.Status(http.StatusOK)
		ctx.Writer.Write([]byte(jsonData))
	} else {
		if _, err := db.Exec(sql); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	}
}

func Login(ctx *gin.Context) {

	db := global.App.DB

	//获取参数
	//此处使用Bind()函数，可以处理不同格式的前端数据
	var requestUser model2.User
	ctx.Bind(&requestUser)
	telephone := requestUser.Telephone
	password := requestUser.Password

	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "手机号必须为11位",
		})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码不能少于6位",
		})
		return
	}

	//判断手机号是否存在
	var user model2.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户不存在",
		})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码错误",
		})
	}

	//返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
	})
}
