package core_http_utils

import (
	"context"
	"net/http"

	core_logger "github.com/wasstend/todoapp-golang/internal/core/logger"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
)

func GetCtxLogResp(
	w http.ResponseWriter,
	r *http.Request,
) (
	context.Context,
	core_logger.Logger,
	*core_http_response.HTTPResponseHandler,
) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	return ctx, logger, responseHandler

}
