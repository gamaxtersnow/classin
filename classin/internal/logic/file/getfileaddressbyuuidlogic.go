package file

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

type GetFileAddressByUUIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据uuid获取文件下载地址
func NewGetFileAddressByUUIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileAddressByUUIDLogic {
	return &GetFileAddressByUUIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileAddressByUUIDLogic) GetFileAddressByUUID(req *types.FileAddressByUUIDReq) (resp *types.FileAddressResp, err error) {
	file, err := l.svcCtx.SFilesModel.FindOneByUuid(l.ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	expires := int64(300)
	if req.ExpireTime > 0 {
		expires = req.ExpireTime
	}
	fileAddress, err := l.svcCtx.MsOssModel.GetSignUrl(l.ctx, file.ObjectKey, expires)
	if err != nil {
		return nil, err
	}
	resp = &types.FileAddressResp{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取成功"
	resp.Address = fileAddress
	return resp, nil
}
