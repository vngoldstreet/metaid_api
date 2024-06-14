package controllers

import (
	"log"
	"net/http"
	"vietvd/mql-api/entity"
	"vietvd/mql-api/service"

	"github.com/gin-gonic/gin"
)

func HandleGetNewMetaID(c *gin.Context) {
	vps_name := c.Query("vps_name")
	data, err := service.GetMongoDB("")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": entity.IDGenerate{},
			"err":  err,
		})
		c.Done()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})

	err = service.UpdateMongoByID(data.ID.Hex(), vps_name)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleGetMetaID(c *gin.Context) {
	vps_name := c.Query("vps_name")
	datas, err := service.GetAllMongoData(vps_name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": entity.IDGenerate{},
		})
		c.Done()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": datas,
	})
}

func HandleDeleteMetaID(c *gin.Context) {
	vps_name := c.Query("vps_name")
	err := service.DeleteAllDataByName(vps_name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"err":  err,
		})
		c.Done()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"err":  nil,
	})
}
