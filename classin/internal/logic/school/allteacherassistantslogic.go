package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllTeacherAssistantsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllTeacherAssistantsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllTeacherAssistantsLogic {
	return &AllTeacherAssistantsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllTeacherAssistantsLogic) AllTeacherAssistants(req *types.TeacherReq) (resp *types.TeacherAssistantsResponse, err error) {
	teacherAssistants, total, err := l.svcCtx.STeacherModel.GetTeacherAssistants(l.ctx, req.PageId, req.PageSize)
	if err != nil {
		return nil, err
	}
	resp = &types.TeacherAssistantsResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取成功"
	resp.Pages.PageID = req.PageId
	resp.Pages.PageSize = req.PageSize
	resp.Pages.Total = total
	resp.Pages.Pages = total/req.PageSize + 1
	if total > 0 {
		for _, teacher := range teacherAssistants {
			t := types.Teacher{}
			_ = utils.ConvertStruct(teacher, &t)
			resp.AssistantList = append(resp.AssistantList, t)
		}
	}
	return resp, nil
}
