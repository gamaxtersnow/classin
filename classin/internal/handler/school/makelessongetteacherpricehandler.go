package school

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/school"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

// 新增排课获取老师和服务价格数据
func MakeLessonGetTeacherPriceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TeacherServiceMapRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := school.NewMakeLessonGetTeacherPriceLogic(r.Context(), svcCtx)
		resp, err := l.MakeLessonGetTeacherPrice(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
