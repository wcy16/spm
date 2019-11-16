package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"spm/model"
	"spm/service"
	"spm/util"
	"strconv"
	"time"
)

const (
	TokenLength = 32
)

func SignUp(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err == nil {
		if err := user.Create(); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		} else {
			//go service.EmailSignUp(user.Email, user.FirstName)
			c.JSON(http.StatusOK, SuccessResponse(nil))
		}
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}
}

func SignIn(c *gin.Context) {
	type SignIn struct {
		Email    string `binding:"required"`
		Password string `binding:"required"`
	}

	signIn := SignIn{}

	if err := c.ShouldBindJSON(&signIn); err == nil {
		user := model.User{}
		if user.Find(signIn.Email) {
			if user.Auth(signIn.Password) {
				id := strconv.Itoa(int(user.ID))
				token := util.RandomString(TokenLength)
				service.Cache.Set(id, token, 0)
				service.SetCookies(c, map[string]string{
					"id":    id,
					"token": token,
				}, service.Cookies10days)
				c.JSON(http.StatusOK, SuccessResponse(gin.H{"Admin": user.Admin}))
			} else {
				c.JSON(http.StatusForbidden, ErrorResponse("wrong email or password"))
			}
		} else {
			c.JSON(http.StatusForbidden, ErrorResponse("wrong email or password"))
		}
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}
}

func SignOut(c *gin.Context) {
	id := c.Keys["id"].(uint)
	service.Cache.Del(fmt.Sprint(id))
	service.SetCookies(c, map[string]string{
		"id":    "",
		"token": "",
	}, -1)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func GetProfile(c *gin.Context) {
	id := c.Keys["id"].(uint)
	user := model.User{}
	user.Retrieve(id)
	user.EnumToStr()
	c.JSON(http.StatusOK, SuccessResponse(user))
}

func EditProfile(c *gin.Context) {
	type Edit struct {
		FirstName string `binding:"required"`
		LastName  string `binding:"required"`

		Mobile  string `json:",omitempty"`
		Home    string `json:",omitempty"`
		Work    string `json:",omitempty"`
		Address string `json:",omitempty"`

		CarInfo string `json:",omitempty"`
		CarType string `json:",omitempty"`
	}

	id := c.Keys["id"].(uint)
	edit := Edit{}
	if err := c.ShouldBindJSON(&edit); err == nil {
		user := model.User{}
		user.Retrieve(id)
		{
			user.FirstName = edit.FirstName
			user.LastName = edit.LastName
			user.Mobile = edit.Mobile
			user.Home = edit.Home
			user.Work = edit.Work
			user.Address = edit.Address
			user.CarInfo = edit.CarInfo
			user.CarType = edit.CarType
		}
		user.Update()
		c.JSON(http.StatusOK, SuccessResponse(nil))
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}
}

func MakeAppointment(c *gin.Context) {
	id := c.Keys["id"].(uint)
	appointment := model.Appointment{}
	if err := c.ShouldBindJSON(&appointment); err == nil {
		if err = appointment.MarshallTime(); err == nil {
			appointment.StrToEnum()
			if err = appointment.Create(id); err == nil {
				go service.EmailConfirm(id, appointment, "creating")

				c.JSON(http.StatusOK, SuccessResponse(nil))
			} else {
				c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
			}
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		}
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}
}

func EditAppointment(c *gin.Context) {
	id := c.Keys["id"].(uint)
	appointment := model.Appointment{}
	if err := c.ShouldBindJSON(&appointment); err == nil {
		if err = appointment.MarshallTime(); err == nil {
			appointment.StrToEnum()
			appointment.Update(id)

			go service.EmailConfirm(id, appointment, "rescheduling")

			c.JSON(http.StatusOK, SuccessResponse(nil))
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		}
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}
}

func DeleteAppointment(c *gin.Context) {
	type Req struct {
		ID uint
	}

	id := c.Keys["id"].(uint)
	appointment := model.Appointment{}
	req := Req{}
	if err := c.ShouldBindJSON(&req); err == nil {
		appointment.ID = req.ID
		NotFound := appointment.Delete(id)
		if !NotFound {
			go service.EmailConfirm(id, appointment, "canceling")
		}
		c.JSON(http.StatusOK, SuccessResponse(nil))
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}
}

func GetAppointments(c *gin.Context) {
	id := c.Keys["id"].(uint)
	var appointments []model.Appointment
	model.GetAppointments(id, &appointments)
	for i := 0; i != len(appointments); i++ {
		appointments[i].EnumToStr()
	}
	c.JSON(http.StatusOK, SuccessResponse(appointments))
}

func GetAllAppointments(c *gin.Context) {
	type Req struct {
		From *time.Time
		To   *time.Time
	}

	var appointments []model.Appointment
	var from, to time.Time

	id := c.Keys["id"].(uint)
	req := Req{}
	if err := c.ShouldBindJSON(&req); err == nil {
		if req.From == nil {
			from = time.Now()
		} else {
			from = *req.From
		}

		if req.To == nil {
			to = from.Add(240 * time.Hour)
		} else {
			to = *req.To
		}
	} else {
		from = time.Now()
		to = from.Add(240 * time.Hour)
	}

	err := model.AllAppointments(from, to, id, &appointments)
	if err == nil {
		for i := 0; i != len(appointments); i++ {
			appointments[i].EnumToStr()
		}
		c.JSON(http.StatusOK, SuccessResponse(appointments))
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

}

func GetAnyProfile(c *gin.Context) {
	id := c.Keys["id"].(uint)
	user := model.User{}
	user.Retrieve(id)
	if user.Admin {
		u := model.User{}
		sid := c.Param("id")
		uid, err := strconv.Atoi(sid)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		} else {
			u.Retrieve(uint(uid))
			u.EnumToStr()
			c.JSON(http.StatusOK, SuccessResponse(u))
		}
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse("not allowed"))
	}
}
