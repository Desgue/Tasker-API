package main

// User struct is used to hold the user data received from the database
type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CognitoId string `json:"cognitoId"`
	CreatedAt string `json:"createdAt"`
}

// CreateUserRequest struct is used to hold the request data for creating a new user
type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CognitoId string `json:"cognitoId"`
}

func NewCreateUserRequest(username, email, cognitoId string) *CreateUserRequest {
	return &CreateUserRequest{
		Username:  username,
		Email:     email,
		CognitoId: cognitoId,
	}
}
