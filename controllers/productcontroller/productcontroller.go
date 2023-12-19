package productcontroller

import (
	"encoding/json"
	"go-gonic-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context){
	
	var products []models.Product

	models.DB.Find(&products)
	c.JSON(http.StatusOK, gin.H{"products" : products})
}

func Show(c *gin.Context){

	var product models.Product

	id := c.Param("id")
	if err := models.DB.First(&product, id).Error; err != nil{
		switch err{
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message":"Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"product":product})
}

func Create(c *gin.Context){

	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message":"Data tidak ditemukan"})
		return
	}

	models.DB.Create(&product)
	c.JSON(http.StatusOK, gin.H{"product":product})
}

func Update(c *gin.Context){
	
	var product models.Product

	id := c.Param("id")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message":"Data tidak ditemukan"})
		return
	}

	if models.DB.Model(&product).Where("id = ?", id).Updates(&product).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message":"Data gagal diupdate"})
		return
	}

	product_id, _ := strconv.ParseInt(id, 10, 64)
	product.Id = product_id
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui","result": product})

}

func Delete(c *gin.Context){

	var product models.Product

	var input struct {
		Id json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message":"Data tidak ditemukan"})
		return
	}

	id, _ := input.Id.Int64()
	if models.DB.Delete(&product, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message":"Data gagal dihapus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})

}
