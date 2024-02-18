package domain

// This is the interface that that will define the behavior to interact with the database

type UserStorage interface {
	CheckUser(string) (bool, error)
	CreateUser(string) error
}

type IUserService interface {
	CheckUser(string) (bool, error)
	CreateUser(string) error
}

// User struct is used to hold the user data received from the database
type User struct {
	Id        int    `json:"id"`
	CognitoId string `json:"cognitoId"`
	Teams     []Team `json:"teams"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

// CreateUserRequest struct is used to hold the request data for creating a new user
type CreateUserRequest struct {
	CognitoId string `json:"cognitoId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

func NewCreateUserRequest(username, email, cognitoId string) *CreateUserRequest {
	return &CreateUserRequest{
		Username:  username,
		Email:     email,
		CognitoId: cognitoId,
	}
}
