package customer

import (
	"context"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCustomerBasicDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询客户基本信息
func NewGetCustomerBasicDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerBasicDataLogic {
	return &GetCustomerBasicDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCustomerBasicDataLogic) GetCustomerBasicData(req *types.CustomerRequest) (resp *types.CustomerBasicResponse, err error) {
	customerData, err := l.svcCtx.SCustomerModel.GetCustomerBasicData(l.ctx, req.Ids)
	resp = &types.CustomerBasicResponse{}
	resp.Data = customerData
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取数据成功"
	return resp, nil
}
