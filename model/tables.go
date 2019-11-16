package model

import "time"

func init() {
	table["user"] = &User{}
	table["appointment"] = &Appointment{}
}

type WashOption uint

const (
	WashOptionNone WashOption = iota
	WashOut
	WashInAndOut
	WashDelux
)

type Car struct {
	Car     CarType `json:"-"`
	CarType string  `gorm:"-" json:",omitempty"`
}

type User struct {
	IgnoreTimeModel
	Password string `json:",omitempty"`
	password string

	Admin bool `binding:"isdefault" gorm:"default:false" json:",omitempty"`

	FirstName string `binding:"required"`
	LastName  string `binding:"required"`
	Email     string `gorm:"unique;not null" binding:"required,email" json:",omitempty"`

	Mobile  string `json:",omitempty"`
	Home    string `json:",omitempty"`
	Work    string `json:",omitempty"`
	Address string `json:",omitempty"`

	CarInfo string `json:",omitempty"`
	Car
}

type Appointment struct {
	IgnoreTimeModel
	UserID uint `json:",omitempty"`
	//DateTime string    `gorm:"-"`
	//Time     time.Time `json:"-" gorm:"unique;not null"`
	Time    time.Time  `binding:"required" gorm:"unique" json:",omitempty"`
	Option  WashOption `json:",omitempty" binding:"min=1,max=3"`
	Comment string     `json:",omitempty"`
	Car
	User *User `json:",omitempty" binding:"isdefault"`
}
