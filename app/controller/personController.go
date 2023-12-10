package controller

import (
	"databaseDemo/app/model"
	"databaseDemo/global"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetPersonByID(ctx *gin.Context) {
	db := global.App.DB

	var requestBody model.RequestPerson
	ctx.ShouldBind(&requestBody)
	var person model.Persons
	db.First(&person, requestBody.ID)
	ctx.JSON(http.StatusOK, person)
}

func GetBelongingsByID(ctx *gin.Context) {
	db := global.App.DB

	var requestBody model.RequestPerson
	ctx.ShouldBind(&requestBody)
	var belongings model.Belongings
	db.First(&belongings, requestBody.ID)
	ctx.JSON(http.StatusOK, belongings)
}

func GetInfosByID(ctx *gin.Context) {
	db := global.App.DB

	var requestBody model.RequestPerson
	ctx.ShouldBind(&requestBody)
	var infos model.Infos
	db.First(&infos, requestBody.ID)
	ctx.JSON(http.StatusOK, infos)
}

func GetAllByID(ctx *gin.Context) {
	db := global.App.DB
	var requestBody model.RequestPerson
	ctx.ShouldBind(&requestBody)
	var ret model.RequestPerson
	db.Table("persons").
		Joins("JOIN infos ON persons.id = infos.id").
		Joins("JOIN belongings ON persons.id = belongings.id").
		Where("persons.id = ?", requestBody.ID).
		Select("persons.id, persons.home, persons.background, infos.business, infos.address, belongings.cars, belongings.pets, belongings.clothes_size").
		Scan(&ret)
	ctx.JSON(http.StatusOK, ret)
}

func GetAllWearSameSize(ctx *gin.Context) {
	db := global.App.DB
	var requestBody model.RequestPerson
	ctx.ShouldBind(&requestBody)
	var ret []model.RequestPerson
	db.Table("belongings").
		Joins("JOIN persons ON belongings.id = persons.id").
		Joins("JOIN infos ON belongings.id = infos.id").
		Where("belongings.clothes_size = ?", requestBody.ClothesSize).
		Select("persons.id, persons.home, persons.background, infos.business, infos.address, belongings.cars, belongings.pets, belongings.clothes_size").
		Find(&ret)
	ctx.JSON(http.StatusOK, ret)
}

func GetAllRunSameBusiness(ctx *gin.Context) {
	db := global.App.DB
	var requestBody model.RequestPerson
	ctx.ShouldBind(&requestBody)
	var ret []model.RequestPerson
	db.Table("infos").
		Joins("JOIN persons ON infos.id = persons.id").
		Joins("JOIN belongings ON infos.id = belongings.id").
		Where("infos.business = ?", requestBody.Business).
		Select("persons.id, persons.home, persons.background, infos.business, infos.address, belongings.cars, belongings.pets, belongings.clothes_size").
		Find(&ret)
	ctx.JSON(http.StatusOK, ret)
}
