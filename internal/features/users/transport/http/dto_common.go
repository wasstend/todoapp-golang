package users_transport_http

import "github.com/wasstend/todoapp-golang/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id" example:"1"`
	Version     int     `json:"version" example:"1"`
	FullName    string  `json:"full_name" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" example:"+79991234567"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
