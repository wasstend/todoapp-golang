package users_transport_http

import (
	"net/http"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164"`
}

type CreateUserResponse UserDTOResponse

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke CreateUser Handler")

	var request CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	userDomain := domainFromDTO(request)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JsonResponse(response, http.StatusCreated)

}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
