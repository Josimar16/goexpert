package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/josimar16/goexpert/apis/internal/dto"
	entity "github.com/josimar16/goexpert/apis/internal/entities"
	"github.com/josimar16/goexpert/apis/internal/infra/database"
)

type UserHandler struct {
	UserDB       database.User
	JWT          *jwtauth.JWTAuth
	JWTExperesIn int
}

func NewUserHandler(
	db database.User,
	jwt *jwtauth.JWTAuth,
	JWTExperesIn int,
) *UserHandler {
	return &UserHandler{
		UserDB:       db,
		JWT:          jwt,
		JWTExperesIn: JWTExperesIn,
	}
}

func (userHandler *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// Get the JWT
	var body dto.CreateSessionDTO

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := userHandler.UserDB.FindByEmail(body.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !user.ValidatePassword(body.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, token, _ := userHandler.JWT.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(userHandler.JWTExperesIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (userHandler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Create a user
	var body dto.CreateUserDTO

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := entity.NewUser(body.Name, body.Email, body.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = userHandler.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
