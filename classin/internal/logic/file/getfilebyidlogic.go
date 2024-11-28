package file

import (
	"context"
	"errors"
	"meishiedu.com/classin/internal/model"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据id获取文件信息
func NewGetFileByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileByIdLogic {
	return &GetFileByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileByIdLogic) GetFileById(req *types.FileByIdReq) (resp *types.FileInfoResp, err error) {
	file, err := l.svcCtx.SFilesModel.FindOne(l.ctx, req.Id)
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
