package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

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

	token, err := auth.GetNewOneTimeToken(email, time.Now().Add(1*time.Hour))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Can't Create New Token!"))
		return
	}

	if err := h.Db.Create(&token).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Can't Create New Token!"))
		return
	}

	w.Write([]byte(token.Token))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	token := r.FormValue("token")
	if password == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Please Set Password!"))
		return
	}

	// check OneTimeScript is valid
	var t = models.OneTimeScript{}
	if err := h.Db.Where("token = ?", token).First(&t).Error; gorm.IsRecordNotFoundError(err) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - token is not valid!"))
		return
	}

	// hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Password cannot use!"))
		return
	}

	var u = models.User{
		Email:    t.Email,
		Password: string(hashedPassword),
	}
	if err := h.Db.Create(&u).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Can't Create New User!"))
		return
	}

	jwtToken, err := auth.GetNewJWTToken(u.ID, u.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Can't Create New JWT Token!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwtToken))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	userId, err := strconv.ParseInt(fmt.Sprintf("%.f", user.(*jwt.Token).Claims.(jwt.MapClaims)["id"]), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - UserID is Invalid!"))
		return
	}
	var u = models.User{}
	if err := h.Db.First(&u, userId).Error; gorm.IsRecordNotFoundError(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("500 - This User doesn't exist!"))
		return
	}

	res, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Failed to Marshal JSON!"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
