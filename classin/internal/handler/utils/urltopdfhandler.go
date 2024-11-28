package utils

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/utils"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

// url转pdf
func UrlToPdfHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UrlToPdfRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := utils.NewUrlToPdfLogic(r.Context(), svcCtx)
		resp, err := l.UrlToPdf(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
