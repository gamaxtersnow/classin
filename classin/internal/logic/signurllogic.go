package logic

import (
	"context"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignUrlLogic {
	return &SignUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignUrlLogic) SignUrl(req *types.SignUrlRequest) (resp *types.SignUrlResponse, err error) {
	signUrl, err := l.svcCtx.MsOssModel.GetSignUrl(l.ctx, req.Objectkey, 3600)
	l.Info(signUrl)
	if err != nil {
		return nil, err
	}
	resp = &types.SignUrlResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "操作成功"
	resp.Data.Url = signUrl
	return
}
