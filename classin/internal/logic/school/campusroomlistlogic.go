package school

import (
	"context"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CampusRoomListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCampusRoomListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CampusRoomListLogic {
	return &CampusRoomListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CampusRoomListLogic) CampusRoomList() (resp *types.CampusRoomListResponse, err error) {
	campusRooms, err := l.svcCtx.SCampusModel.GetCampusRoomList(l.ctx)
	if err != nil {
		return nil, err
	}
	resp = &types.CampusRoomListResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取成功"
	if len(campusRooms) > 0 {
		for _, c := range campusRooms {
			campus := types.CampusInfo{}
			campus.ID = c.ID
			campus.Name = c.Name
			for _, room := range c.RoomList {
				campus.OutClzRoomList = append(campus.OutClzRoomList, types.Classroom(room))
			}
			resp.CampusInfo = append(resp.CampusInfo, campus)
		}
	}
	return resp, nil
}
