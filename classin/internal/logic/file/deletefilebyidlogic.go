package file

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"meishiedu.com/classin/internal/model/school"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据id删除文件
func NewDeleteFileByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileByIdLogic {
	return &DeleteFileByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileByIdLogic) DeleteFileById(req *types.DeleteFileByIdReq) (resp *types.DeleteFileResp, err error) {
	file, err := l.svcCtx.SFilesModel.FindOne(l.ctx, req.Id)
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
