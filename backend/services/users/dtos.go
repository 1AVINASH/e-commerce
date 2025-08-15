package user

type GetUserInputDto struct {
	Id int64 `json:"id"`
}

type LoginInputDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupInputDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
