package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/file"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

// GetFileByUUIDHandler 根据uuid获取文件信息
func GetFileByUUIDHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileByUUIDReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewGetFileByUUIDLogic(r.Context(), svcCtx)
		resp, err := l.GetFileByUUID(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
