package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/file"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

// 根据uuid获取文件预览信息
func GetFileViewByUUIDHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileViewByUUIDReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewGetFileViewByUUIDLogic(r.Context(), svcCtx)
		resp, err := l.GetFileViewByUUID(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
