package web_transport_http

import (
	"net/http"

	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

func (h *WebHTTPHandler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	_, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke GetMainPage handler")

	html, err := h.webService.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get main page",
		)
		return
	}

	responseHandler.HTMLResponse(html)
}
