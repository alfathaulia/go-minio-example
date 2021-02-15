package upload

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func mkdir(args ...interface{}) {
	var path string
	var mode os.FileMode
	mode = 0755
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			path = filepath.Join(path, arg)
		case os.FileMode:
			mode = arg
		}
	}
	os.MkdirAll(path, mode)
}

func url(path ...string) string {
	return filepath.ToSlash(filepath.Join(path...))
}

//abort with optional code, error. default code: 500; error "error"
func abort(c *gin.Context, args ...interface{}) {
	code := http.StatusInternalServerError
	err := errors.New("error")
	for _, arg := range args {
		switch arg := arg.(type) {
		case error:
			err = arg
		case int:
			code = arg
		}
	}
	c.JSON(code, gin.H{
		"message": err.Error(),
	})
}
