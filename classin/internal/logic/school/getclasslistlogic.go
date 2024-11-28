package school

import (
	"context"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetClassListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClassListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClassListLogic {
	return &GetClassListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClassListLogic) GetClassList(req *types.ClassListRequest) (resp *types.ClassListResponse, err error) {
	classes, err := l.svcCtx.SClassModel.GetClassListByRole(l.ctx, req.PageID, req.PageSize, req.CampusIDs, req.CourseID, req.Search, req.Status)
	if err != nil {
		return nil, err
	}
	resp = &types.ClassListResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取成功"
	resp.Pages = types.Pages(classes.Pages)
	for _, class := range classes.Data {
		resp.ClassList = append(resp.ClassList, types.Class(class))
	}

	return resp, nil
}
