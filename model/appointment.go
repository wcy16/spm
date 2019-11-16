package model

import (
	"time"
)

// marshall time and check valid
// valid hour should be [8, 12), [13, 17)
// deprecated
func (a *Appointment) MarshallTime() (err error) {
	//a.Time, _ = time.Parse(a.DateTime, "2006-01-02T15")
	//hour := a.Time.Hour()
	//if hour < 8 || hour == 12 || hour >= 17 {
	//	err = errors.New("not the correct time slot")
	//}
	return
}

// deprecated
func (a *Appointment) UnmarshallTime() {
	//a.DateTime = a.Time.Format("2006-01-02T15:04:05")
}

func (a *Appointment) Create(UserID uint) error {
	a.ID = 0
	a.UserID = UserID
	return conn.Create(a).Error
}

func (a *Appointment) Update(UserID uint) {
	if a.ID == 0 {
		return
	}
	a.UserID = UserID
	conn.Where("user_id = ?", UserID).Save(a)
}

// query and delete
func (a *Appointment) Delete(UserID uint) (NotFound bool) {
	if a.ID == 0 {
		return
	}
	NotFound = conn.First(a, a.ID).RecordNotFound()
	conn.Model(a).Where("user_id = ?", UserID).Unscoped().Delete(a)
	return
}

func GetAppointments(UserID uint, a *[]Appointment) {
	conn.Where("user_id = ?", UserID).Find(a)
}

func AllAppointments(from, to time.Time, userid uint, appointments *[]Appointment) (err error) {
	from_ := from.Format("2006-01-02")

	to_ := to.Format("2006-01-02")

	// check if admin
	user := User{}
	conn.Select("admin").Find(&user, userid)

	tmp := conn.Where("DATE(`time`) BETWEEN ? AND ?", from_, to_)
	if user.Admin == false {
		tmp = tmp.Select("`time`")
	}

	err = tmp.Find(appointments).Error
	if err != nil {
		return
	}

	// add user
	if user.Admin {
		var users []User
		conn.Find(&users)

		for id, _ := range *appointments {
			for uid, _ := range users {
				if users[uid].ID == (*appointments)[id].UserID {
					(*appointments)[id].User = &users[uid]
				}
			}
		}
	}
	return
}

func GetPrice(option WashOption) string {
	var price string
	switch option {
	case WashOut:
		price = "wash outside only $15"
	case WashInAndOut:
		price = "wash inside and outside $25"
	case WashDelux:
		price = "deluxe wash $30"
	default:
		price = ""
	}

	return price
}
