package orm

import (
	"time"
)

type Replies struct {
	ReplyID int64  `gorm:"index, primaryKey"`
	Author  string `gorm:"index"`
	Text    string `gorm:"index"`
	Url     string `gorm:"index"`
	Visited bool   `gorm:"index"`

	PositiveProb float64
	NegativeProb float64
	CreateAt     time.Time `gorm:"autoCreateTime"`
}

type Users struct {
	ChatID   int64 `gorm:"index"`
	UserName string
	CreateAt time.Time `gorm:"autoCreateTime"`
}

type Rules struct {
	ID       int       `gorm:"primaryKey"`
	Type     string    `gorm:"index"`
	Content  string    `gorm:"index"`
	CreateAt time.Time `gorm:"autoCreateTime"`
}

type KV struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

func AutoCreateTable() {
	err := Init()
	if err != nil {
		panic(err)
	}

	db := GetConn()
	err = db.AutoMigrate(&Users{}, &Replies{}, &Rules{}, &KV{})
	if err != nil {
		panic(err)
	}
}
