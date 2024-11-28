package school

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/school"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

// ScheduleSyncHandler 校校课表数据同步落库
func ScheduleSyncHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ScheduleSyncRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := school.NewScheduleSyncLogic(r.Context(), svcCtx)
		resp, err := l.ScheduleSync(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
