package upload

import "github.com/gin-gonic/gin"

//Default the default register
func Default(base *gin.RouterGroup) {
	u := New()
	Register(base, u)
}

//Upload file
func Upload(c *gin.Context) (file *File, err error) {
	return uploader.Upload(c)
}

//Get file
func Get(id interface{}) (file *File, err error) {
	return uploader.Get(id)
}

//Delete file, accepts string or uuid
func Delete(id interface{}) (file *File, err error) {
	return uploader.Delete(id)
}
