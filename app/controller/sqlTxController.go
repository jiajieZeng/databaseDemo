package controller

import (
	model2 "databaseDemo/app/model"
	"databaseDemo/global"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"
)

func TxBegin(ctx *gin.Context) {
	tx, err := global.App.RDB.Begin()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
		return
	}
	global.App.Tx = tx
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Tx begin",
	})
}

func TxCommit(ctx *gin.Context) {
	tx := global.App.Tx
	err := tx.Commit()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Tx commited",
	})
}

func TxRollBack(ctx *gin.Context) {
	tx := global.App.Tx
	err := tx.Rollback()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Tx rolled back",
	})
}

func TxRaw(ctx *gin.Context) {
	var requestBody model2.RequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}
	sql := requestBody.SQL
	tx := global.App.Tx
	if tx == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "please begin a Tx first",
		})
		return
	}
	if sql[0] == 's' || sql[0] == 'S' {
		rows, err := tx.Query(sql)
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
		if _, err := tx.Exec(sql); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	}
}
