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

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB database.User
}

func NewUserHandler(db database.User) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// AuthenticateUser godoc
//
//	@Summary      Authenticate user
//	@Description  Authenticate a user
//	@Tags         sessions
//	@Accept       json
//	@Produce      json
//	@Param        request body dto.CreateSessionDTO true "user credencials"
//	@Success      200  {object}  dto.CreateAccessTokenDTO
//	@Failure      404  {object}  Error
//	@Failure      500  {object}  Error
//	@Router       /sessions [post]
func (userHandler *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// Get the JWT
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExperesIn := r.Context().Value("JWTExperesIn").(int)

	var body dto.CreateSessionDTO

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := userHandler.UserDB.FindByEmail(body.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	if !user.ValidatePassword(body.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, token, _ := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExperesIn)).Unix(),
	})

	accessToken := dto.CreateAccessTokenDTO{AccessToken: token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// CreateUser godoc
//
//	@Summary      Create user
//	@Description  Create a user
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Param        request body dto.CreateUserDTO true "user request"
//	@Success      201
//	@Failure      500  {object}  Error
//	@Router       /users [post]
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
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = userHandler.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
