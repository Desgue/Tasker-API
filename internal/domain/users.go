package domain

// User struct is used to hold the user data received from the database
type User struct {
	CognitoId string `json:"cognitoId"`
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
