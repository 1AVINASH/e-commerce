package user

import (
	"encoding/json"
	"gotemplate/utility/auth"
	"gotemplate/utility/logger"
	middleware "gotemplate/utility/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserAPIs struct {
	repo *UserRepository
}

func NewUserApis(mux *chi.Mux, repo *UserRepository) *UserAPIs {
	ua := &UserAPIs{repo: repo}
	ua.registerUserRoutes(mux)
	return ua
}

func (ua *UserAPIs) registerUserRoutes(mux *chi.Mux) {
	mux.Route("/users", func(mux chi.Router) {
		mux.Get("/", middleware.JSON(ua.getUsers))
		mux.Post("/", middleware.JSON(ua.createUser))
	})
}

func (ua *UserAPIs) getUsers(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Starting to get users")
	users, err := ua.repo.GetUsers()
	if err != nil {
		logger.Logger.Errorf("Ran into an error %v", err)
		return nil, http.StatusBadRequest
	}
	return users, http.StatusOK
}

func (ua *UserAPIs) login(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Logging in")
	var input LoginInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}
	user, err := ua.repo.GetUserByEmail(input.Email)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}
	if hash != user.Password {
		return map[string]string{"error": "Incorrect Password"}, http.StatusBadRequest
	}

	return map[string]string{"status": "logged_in"}, http.StatusCreated
}

func (ua *UserAPIs) signup(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Signing up")
	var input SignupInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}
	createUserInput := User{
		Name:     &input.Name,
		Email:    &input.Email,
		Password: &hash,
	}
	user, err := ua.repo.CreateUser(&createUserInput)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return map[string]interface{}{"data": jsonData}, http.StatusCreated
}
