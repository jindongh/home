package piano

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Song struct
type Song struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;"`
	Title    string    `json:"title"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
}
// Songs struct
type Songs struct {
	Songs []Song `json:"songs"`
}

func (song *Song) BeforeCreate(tx *gorm.DB) (err error) {
	song.ID = uuid.New()
	return nil
}
