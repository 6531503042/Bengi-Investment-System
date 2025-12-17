package dto

type (
	RegisterRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		FullName string `json:"fullName" validate:"required,min=2"`
		Phone    string `json:"phone,omitempty"`
	}

	RegisterResponse struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
	}
)
