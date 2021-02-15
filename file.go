package upload

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//File the file stored in db
type File struct {
	ID           uuid.UUID `gorm:"type:varchar(36);primary_key;"`
	Filename     string
	Ext          string
	Path         string
	Size         int64
	MIMEType     string
	URL          string `gorm:"-"`
	CloudStorage string
	CreatedAt    time.Time
}

//BeforeCreate initial uuid for event
func (file *File) BeforeCreate(tx *gorm.DB) (err error) {
	if file.ID == uuid.Nil {
		file.ID = uuid.New()
	}
	return
}
