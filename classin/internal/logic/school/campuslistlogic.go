package school

import (
	"context"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CampusListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCampusListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CampusListLogic {
	return &CampusListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CampusListLogic) CampusList() (resp *types.CampusListResponse, err error) {
	campuses, err := l.svcCtx.SCampusModel.GetAllCampuses(l.ctx)
	if err != nil {
		l.Errorf("获取校区列表错误，err:%v", err)
		return nil, err
	}
	resp = &types.CampusListResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取成功"
	if len(campuses) > 0 {
		campus := types.Campus{}
		for _, c := range campuses {
			campus.Id = c.Id
			campus.Name = c.Name
			campus.Type = c.Type
			resp.CampusList = append(resp.CampusList, campus)
		}
	}
	return resp, nil
}
