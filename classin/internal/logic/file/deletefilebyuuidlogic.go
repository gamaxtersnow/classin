package file

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"meishiedu.com/classin/internal/model/school"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileByUUIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据uuid删除文件
func NewDeleteFileByUUIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileByUUIDLogic {
	return &DeleteFileByUUIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileByUUIDLogic) DeleteFileByUUID(req *types.DeleteFileByUUIDReq) (resp *types.DeleteFileResp, err error) {
	file, err := l.svcCtx.SFilesModel.FindOneByUuid(l.ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	file.Status = school.FileStatusDeleted
	file.UpdateTime = utils.GetCurrentTimestamp()
	err = l.svcCtx.SFilesModel.Update(l.ctx, file)
	if err != nil {
		return nil, err
	}
	resp = &types.DeleteFileResp{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "删除成功"
	return resp, nil
}
