package file

import (
	"context"
	"errors"
	"meishiedu.com/classin/internal/model"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileByUUIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据uuid获取文件信息
func NewGetFileByUUIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileByUUIDLogic {
	return &GetFileByUUIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileByUUIDLogic) GetFileByUUID(req *types.FileByUUIDReq) (resp *types.FileInfoResp, err error) {
	file, err := l.svcCtx.SFilesModel.FindOneByUuid(l.ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("文件不存在")
		}
		return nil, err
	}
	resp = &types.FileInfoResp{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "查询成功"
	resp.FileInfo = types.FileInfo{
		Uuid: file.Uuid,
		Name: file.Name,
	}
	return resp, nil
}
