package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/odmishien/go-auth-playground/auth"
	"github.com/odmishien/go-auth-playground/models"
)

type UserHandler struct {
	Db *gorm.DB
}

func (h *UserHandler) PreCreateUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Please Set Email!"))
		return
	}
	var u = models.User{Email: email, Activated: false}

	h.Db.Create(&u)

	token, err := auth.GetNewToken(u.ID, u.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Can't Create New Token!"))
		return
	}
	w.Write([]byte(token))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	if password == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Please Set Password!"))
		return
	}
	user := r.Context().Value("user")
	id, err := strconv.ParseInt(fmt.Sprintf("%.f", user.(*jwt.Token).Claims.(jwt.MapClaims)["id"]), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - token is not valid!"))
		return
	}
	var u = models.User{}
	if err := h.Db.First(&u, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - user is not found!"))
		return
	}
	fmt.Printf("%#v\n", u)
	// TODO: 平文で保存してるのでダメ
	u.Password = password
	u.Activated = true
	h.Db.Save(&u)
	fmt.Fprintf(w, "You are logged in as %s", u.Email)
}
