package school

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/school"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

func ScheduleViewsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ScheduleViewsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := school.NewScheduleViewsLogic(r.Context(), svcCtx)
		resp, err := l.ScheduleViews(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
