package customer

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/logic/customer"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

// 查询客户基本信息
func GetCustomerBasicDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CustomerRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := customer.NewGetCustomerBasicDataLogic(r.Context(), svcCtx)
		resp, err := l.GetCustomerBasicData(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
