package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

type SetNumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetNumLogic {
	return &SetNumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetNumLogic) SetNum(req *types.SetNumRequest) (resp *types.SetNumResponse, err error) {
	// 接收列表的参数并判断
	classList := req.CourseList
	logx.Infof("classList: %+v", classList)
	//循环调用ClassSetNum方法修改数据
	for _, classJson := range classList {
		logx.Infof("courseId: %+v", classJson.CourseId)
		logx.Infof("classJsonInfo: %+v", classJson.ClassJson)
		resp, err := l.svcCtx.SettingModel.ClassSetNum(
			l.ctx, req.Cookie, classJson.CourseId, classJson.ClassJson)
		logx.Infof("resp: %+v", resp)
		if err != nil {
			return nil, err
		}
		if resp.ErrorInfo.ErrorCode != 1 {
			return nil, errors.New(resp.ErrorInfo.ErrorMsg)
		}
	}
	return &types.SetNumResponse{
		ErrorInfo: types.ErrorInfo{
			ErrorCode: 1,
			ErrorMsg:  "ok",
		},
	}, nil
}
