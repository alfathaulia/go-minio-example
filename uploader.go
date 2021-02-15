package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//Uploader the uploader instance
type Uploader struct {
	StaticRoot   string
	UploadFolder string
	DB           *gorm.DB
}

var uploader *Uploader

//New File Uploader instance
func New() *Uploader {
	return &Uploader{
		StaticRoot:   "attachments",
		UploadFolder: "attachments",
	}
}

//Register the file service
func Register(base *gin.RouterGroup, u *Uploader) {
	if u.DB == nil {
		mkdir("data")
		DB, err := gorm.Open(sqlite.Open("data/file.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		u.DB = DB
	}
	u.DB.AutoMigrate(&File{})
	uploader = u
	controllers(base)
}

//Upload file
func (u *Uploader) Upload(c *gin.Context) (file *File, err error) {
	// endpoint := "storage.mojavstudio.com"
	// accessKeyID := "@bc12345"
	// secretAccessKey := "@bc123456"
	// useSSL := true

	// // Initialize minio client object.
	// minioClient, err := minio.New(endpoint, &minio.Options{
	// 	Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	// 	Secure: useSSL,
	// })
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Printf("%s\n", minioClient.EndpointURL())

	f, err := c.FormFile("file")
	if err != nil {
		return

	}
	id := uuid.New()
	ext := filepath.Ext(f.Filename)
	storedFileName := strings.Replace(id.String(), "-", "", -1) + ext
	t := time.Now()
	dir := filepath.Join(
		strconv.Itoa(t.Year()),
		strconv.Itoa(int(t.Month())),
	)
	dest := filepath.Join(
		u.UploadFolder,
		dir,
	)
	mkdir(dest)
	dest = filepath.Join(dest, storedFileName)
	if err = c.SaveUploadedFile(f, dest); err != nil {
		return
	}
	path := filepath.Join(dir, storedFileName)

	// userMetaData := map[string]string{"x-amz-acl": "public-read"}
	// cacheControl := "max-age=31536000"
	// buffer := make([]byte, f.Size)
	// fileBytes := bytes.NewReader(buffer)
	// n, err := minioClient.PutObject(c, "hipam-app", f.Filename, fileBytes, f.Size, minio.PutObjectOptions{ContentType: f.Header.Get("Content-Type"), CacheControl: cacheControl, UserMetadata: userMetaData})
	// // Set request parameters
	// reqParams := make(urs.Values)
	// reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")
	// object, err := minioClient.PresignedGetObject(c, "hipam-app", f.Filename, time.Second*24*60*60, reqParams)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Successfully get bytes: ", object)
	// fmt.Println("Successfully uploaded bytes: ", n)
	file = &File{
		ID:       id,
		Filename: f.Filename,
		Size:     f.Size,
		MIMEType: f.Header.Get("Content-Type"),
		Path:     path,
		Ext:      ext,
		URL:      url("/", u.StaticRoot, path),
		// CloudStorage: object.Path,
	}
	err = u.DB.Create(file).Error
	if err != nil {
		return
	}

	return
}

//Get file
func (u *Uploader) Get(id interface{}) (file *File, err error) {
	var uid uuid.UUID
	switch id := id.(type) {
	case uuid.UUID:
		uid = id
	case string:
		uid, err = uuid.Parse(id)
	default:
		err = &Error{
			Code:    ParserError,
			Message: fmt.Sprintf("Unknown argument %T", id),
		}
	}
	if err != nil {
		return
	}
	file = &File{}
	err = u.DB.First(file, uid).Error
	if err != nil {
		return
	}
	file.URL = url("/", u.StaticRoot, file.Path)
	return
}

//Delete file, accepts string or uuid
func (u *Uploader) Delete(id interface{}) (file *File, err error) {
	file, err = u.Get(id)
	if err != nil {
		return
	}
	err = os.Remove(filepath.Join(u.UploadFolder, file.Path))
	if err != nil {
		u.DB.Delete(file)
		return
	}
	err = u.DB.Delete(file).Error
	if err != nil {
		return
	}
	return
}

//GetURL file
func (u *Uploader) GetURL(id interface{}) (URL string, err error) {
	file, err := u.Get(id)
	if err != nil {
		return
	}
	URL = url("/", u.StaticRoot, file.Path)
	return
}
