package school

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/school"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
	"net/http"
)

// ScheduleExportHandler 课表导出
func ScheduleExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ScheduleExportRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := school.NewScheduleExportLogic(r.Context(), svcCtx)
		resp, err := l.ScheduleExport(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 设置响应头
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
