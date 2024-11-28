package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/model"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

type CourseListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseListLogic {
	return &CourseListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseListLogic) CourseList(req *types.CourseRequest) (resp *types.CourseResponse, err error) {

	params := make(map[string]interface{})

	if req.CourseName != "" {
		params["courseName LIKE"] = "%" + req.CourseName + "%"
	}
	if req.ClassName != "" {
		params["className LIKE"] = "%" + req.ClassName + "%"
	}
	if req.TeacherName != "" {
		params["teacherName LIKE"] = "%" + req.TeacherName + "%"
	}
	if req.TeacherMobile > 0 {
		params["teacherMobile"] = req.TeacherMobile
	}
	if req.SeatNum == 2 {
		params["seatNum"] = req.SeatNum
	} else {
		params["seatNum >="] = req.SeatNum
	}
	if req.IsHd > 0 {
		params["isHd"] = req.IsHd
	}
	if req.CourseStartDate > 0 && req.CourseEndDate > 0 && req.CourseStartDate <= req.CourseEndDate {
		params["courseStartTime >="] = req.CourseStartDate
		params["courseStartTime <="] = req.CourseEndDate
	}
	page := 1
	pageSize := 10
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	courseList := &types.CourseResponse{}
	count, err := l.svcCtx.CourseModel.CountLesson(l.ctx, params)
	if err != nil {
		l.Error(err, params)
		return nil, errors.New("未获取查询数据")
	}

	if count > 0 {
		offset := (page - 1) * pageSize
		var lessonList []model.Course
		lessonList, err = l.svcCtx.CourseModel.GetLessonList(l.ctx, params, offset, pageSize)
		if err != nil {
			l.Error(err, params)
			return nil, errors.New("未获取查询数据")
		}
		courseList.ErrorInfo.ErrorCode = 1
		courseList.ErrorInfo.ErrorMsg = "获取成功"
		courseList.Data.TotalClassNum = count
		for _, class := range lessonList {
			responseCourse := types.CourseResponseCourse{}
			responseCourse.CourseId = class.CourseId
			responseCourse.CourseName = class.CourseName
			responseCourse.ClassId = class.ClassId
			responseCourse.ClassName = class.ClassName
			responseCourse.ClassStartTime = class.CourseStartTime
			responseCourse.ClassEndTime = class.CourseEndTime
			responseCourse.TeacherInfo.Name = class.TeacherName
			responseCourse.TeacherInfo.Mobile = class.TeacherMobile
			responseCourse.SyncStatus = class.SyncStatus
			responseCourse.SeatNum = class.SeatNum
			responseCourse.IsHd = class.IsHd
			responseCourse.SyncStatusName = l.getLessonSyncStatusName(class.SyncStatus)
			courseList.Data.ClassList = append(courseList.Data.ClassList, responseCourse)
		}
	}
	return courseList, nil
}
func (l *CourseListLogic) getLessonSyncStatusName(syncStatus int64) string {
	fileStatusMap := make(map[int64]string)
	fileStatusMap[model.SyncStatusInit] = "未同步"
	fileStatusMap[model.SyncStatusSync] = "同步中"
	fileStatusMap[model.SyncStatusComplete] = "已完成"
	if _, ok := fileStatusMap[syncStatus]; ok {
		return fileStatusMap[syncStatus]
	}
	return "未知"
}
func (l *CourseListLogic) AddCourses(LessonSyncReq []byte) error {
	course := &types.LessonSyncReq{}
	err := json.Unmarshal(LessonSyncReq, course)
	if err != nil {
		l.Error(err.Error() + "源数据：" + string(LessonSyncReq))
		return nil
	}
	if course.CourseId <= 0 || course.ClassId <= 0 || course.UserAgent == "" || course.Cookie == "" {
		l.Error("消息非法，源数据：" + string(LessonSyncReq))
		return nil
	}
	uniqueId := strconv.FormatInt(course.CourseId, 10) + "-" + strconv.FormatInt(course.ClassId, 10)
	res, _ := l.svcCtx.CourseModel.FindByUid(l.ctx, uniqueId)
	courseModel := &model.Course{}
	courseModel.UniqueId = uniqueId
	courseModel.CourseId = course.CourseId
	courseModel.CourseName = course.CourseName
	courseModel.ClassId = course.ClassId
	courseModel.ClassName = course.ClassName
	courseModel.CourseStartTime = course.ClassStartTime
	courseModel.CourseEndTime = course.ClassEndTime
	courseModel.TeacherName = course.TeacherName
	courseModel.TeacherMobile = course.TeacherMobile
	courseModel.SyncStatus = model.FileStatusInit
	courseModel.SourceType = model.SourceTypeManual
	courseModel.AddTime = time.Now().Unix()
	courseModel.SeatNum = course.SeatNum
	courseModel.IsDc = course.IsDc
	courseModel.IsHd = course.IsHd
	if res == nil {
		//增加课程数据
		_, err = l.svcCtx.CourseModel.Insert(l.ctx, courseModel)
	} else {
		courseModel.Id = res.Id
		courseModel.AddTime = res.AddTime
		_ = l.svcCtx.CourseModel.Update(l.ctx, courseModel)
	}
	if course.SourceType == model.SourceTypeManual {
		if courseModel.SyncStatus == model.FileStatusInit {
			//将消息推送到到课节视频队列
			lessonClipSyncReq := types.LessonClipSyncReq{
				UserAgent:  course.UserAgent,
				CourseId:   course.CourseId,
				ClassId:    course.ClassId,
				Cookie:     course.Cookie,
				SourceType: model.SourceTypeManual,
			}
			msg, _ := json.Marshal(lessonClipSyncReq)
			err = l.svcCtx.RabbitMqProducer.SendClassInMsgToMq(l.svcCtx.Config.RabbitMq.ClassInLessonClipTopic, msg)
		}
	}
	return nil
}
