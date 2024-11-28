package school

import (
	"context"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllCoursesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllCoursesLogic {
	return &AllCoursesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllCoursesLogic) AllCourses() (resp *types.CoursesResponse, err error) {
	courses, err := l.svcCtx.SCourseModel.GetAllCourses(l.ctx)
	if err != nil {
		return nil, err
	}
	resp = &types.CoursesResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取成功"
	if len(courses) > 0 {
		for _, course := range courses {
			c := types.Course{}
			c.ID = course.ID
			c.Name = course.Name
			c.Category = course.Category
			c.Classification = course.Classification
			c.Pics = course.Pics
			c.Description = course.Description
			resp.CourseList = append(resp.CourseList, c)
		}
	}
	return resp, nil
}
