package dto

type (
	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	LoginResponse struct {
		Token string      `json:"token"`
		User  UserProfile `json:"user"`
	}

	UserProfile struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		RoleID   string `json:"roleId"`
	}
)
