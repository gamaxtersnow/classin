package school

import (
	"context"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MakeLessonGetTeacherPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewMakeLessonGetTeacherPriceLogic 新增排课获取老师和服务价格数据
func NewMakeLessonGetTeacherPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeLessonGetTeacherPriceLogic {
	return &MakeLessonGetTeacherPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MakeLessonGetTeacherPriceLogic) MakeLessonGetTeacherPrice(req *types.TeacherServiceMapRequest) (resp *types.TeacherServiceMapResponse, err error) {
	_, _ = l.svcCtx.STeacherServiceModel.GetTeacherServicePriceByPhone(l.ctx, req.Phone)
	return resp, nil
}
