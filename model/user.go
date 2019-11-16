package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"spm/util"
)

// auth users
func (u *User) Auth(pwd string) bool {
	if util.ComparePassword(u.password, pwd) {
		return true
	} else {
		return false
	}
}

// find user given email, return true if found
func (u *User) Find(email string) bool {
	return !conn.Where("email = ?", email).First(u).RecordNotFound()
}

func (u *User) Retrieve(id uint) {
	conn.First(u, id)
}

// must update through this function, fields that cannot be updated will be omitted
func (u *User) Update() {
	if u.ID == 0 {
		return
	}
	u.StrToEnum()
	conn.Model(u).Omit("id", "email", "password").Save(u)
}

////////////////////////////////////////////
// hooks for db operations
///////////////////////////////////////////

// before create, hashing password
func (u *User) BeforeCreate() (err error) {
	if u.Password != "" {
		u.Password = util.Password(u.Password)
	}
	return
}

// after save, mask password
func (u *User) AfterSave(scope *gorm.Scope) (err error) {
	u.Password = ""
	return
}

// after find, mask password
func (u *User) AfterFind() (err error) {
	u.password = u.Password
	u.Password = ""
	return
}

func (u *User) Create() error {
	u.ID = 0
	if err := passwordCheck(&u.Password); err != nil {
		return err
	}
	u.StrToEnum()
	return conn.Create(u).Error
}

func passwordCheck(s *string) error {
	if len(*s) < 6 || len(*s) > 16 {
		return errors.New("password length must between 6 and 16")
	}
	return nil
}

func (User) addForeignKey(db *gorm.DB) {
	db.Model(Appointment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
}

func (User) seed(db *gorm.DB) bool {

	u1 := User{}
	u1.FirstName = "a"
	u1.LastName = "a"
	u1.Email = "a@gmail.com"
	u1.Password = util.DefaultPassword
	u1.Admin = true
	db.Create(&u1)

	u2 := User{}
	u2.FirstName = "b"
	u2.LastName = "b"
	u2.Email = "b@gmail.com"
	u2.Password = util.DefaultPassword
	db.Create(&u2)

	u3 := User{}
	u3.FirstName = "c"
	u3.LastName = "c"
	u3.Email = "c@gmail.com"
	u3.Password = util.DefaultPassword
	db.Create(&u3)

	return true
}
