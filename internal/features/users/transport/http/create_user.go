package users_transport_http

import (
	"net/http"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164" example:"+79991234567"`
}

type CreateUserResponse UserDTOResponse

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// CreateUser 	godoc
// @Summary 	Создать пользователя
// @Description Создание нового пользователя в системе
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		request body CreateUserRequest true "CreateUser тело запроса"
// @Success 	201 {object} CreateUserResponse "Успешно созданный пользователь"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users [post]
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
