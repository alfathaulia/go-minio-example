package upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func controllers(base *gin.RouterGroup) {
	api := base.Group("/attachment")

	api.POST("/", func(c *gin.Context) {
		file, err := uploader.Upload(c)
		if err != nil {
			abort(c, err)
		} else {
			c.JSON(http.StatusOK, file)
		}
	})

	api.GET("/:id", func(c *gin.Context) {
		file, err := uploader.Get(c.Param("id"))
		if err != nil {
			abort(c, err)
			return
		}
		c.JSON(http.StatusOK, file)
	})

	api.DELETE("/:id", func(c *gin.Context) {
		file, err := uploader.Delete(c.Param("id"))
		if err != nil {
			abort(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"name":    file.Filename,
		})
	})
}
