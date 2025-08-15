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
	mux.Route("/users", func(r chi.Router) {
		r.Post("/", middleware.JSON(ua.signup))
		r.Post("/login", middleware.JSON(ua.login))

		r.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth)
			r.Get("/", middleware.JSON(ua.getUsers))
			r.Get("/", middleware.JSON(ua.getUser))
		})
	})
	mux.Route("/admin", func(r chi.Router) {
		r.Post("/", middleware.JSON(ua.adminSignup))
		r.Post("/login", middleware.JSON(ua.adminLogin))

		r.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth)
			r.Get("/", middleware.JSON(ua.getUsers))
			r.Get("/", middleware.JSON(ua.getUser))
		})
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

func (ua *UserAPIs) getUser(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Starting to get user")
	var input GetUserInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}
	user, err := ua.repo.GetUser(input.Id)
	if err != nil {
		logger.Logger.Errorf("Ran into an error %v", err)
		return nil, http.StatusBadRequest
	}
	return user, http.StatusOK
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
	if hash != *user.Password {
		return map[string]string{"error": "Incorrect Password"}, http.StatusBadRequest
	}
	token, err := auth.GenerateJWT(*user.Email)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}

	return map[string]string{"status": "logged_in", "token": token}, http.StatusCreated
}

func (ua *UserAPIs) adminLogin(w http.ResponseWriter, r *http.Request) (interface{}, int) {
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
	if hash != *user.Password {
		return map[string]string{"error": "Incorrect Password"}, http.StatusBadRequest
	}
	if *user.Role != "ADMIN" {
		return map[string]string{"error": "No admin found with this email id"}, http.StatusBadRequest
	}
	token, err := auth.GenerateJWT(*user.Email)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}

	return map[string]string{"status": "logged_in", "token": token}, http.StatusCreated
}

func (ua *UserAPIs) adminSignup(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Signing up for admin")
	var input SignupInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}

	adminRole := ConstRoleAdmin
	createUserInput := User{
		Name:     &input.Name,
		Email:    &input.Email,
		Password: &hash,
		Role:     &adminRole,
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
	userRole := ConstRoleShopper
	createUserInput := User{
		Name:     &input.Name,
		Email:    &input.Email,
		Password: &hash,
		Role:     &userRole,
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
