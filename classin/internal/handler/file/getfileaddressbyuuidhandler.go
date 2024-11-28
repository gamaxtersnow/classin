package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/file"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

// 根据uuid获取文件下载地址
func GetFileAddressByUUIDHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileAddressByUUIDReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewGetFileAddressByUUIDLogic(r.Context(), svcCtx)
		resp, err := l.GetFileAddressByUUID(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
