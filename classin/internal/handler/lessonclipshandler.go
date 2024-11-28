package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

func LessonClipsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LesssonClipRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLessonClipsLogic(r.Context(), svcCtx)
		resp, err := l.LessonClips(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
