package logic

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/constant"
	"meishiedu.com/classin/internal/model"
	"meishiedu.com/classin/internal/model/classin"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
)

type SyncJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncJobLogic {
	return &SyncJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncJobLogic) SyncJob(req *types.SyncLessonRequest) (resp *types.SyncLessonResponse, err error) {
	params := l.getParams(req.CourseName, req.ClassName, req.CourseStartDate, req.CourseEndDate, 1, 1)
	courses := &classin.CourseResponse{}
	courses, err = l.svcCtx.ClassInCourseModel.GetCourseList(l.ctx, req.Cookie, params)
	if err != nil {
		return nil, err
	}
	if courses.ErrorInfo.ErrorCode != 1 {
		return nil, errors.New("查询出错:[" + courses.ErrorInfo.ErrorMsg + "],请重试！")
	}
	resp = &types.SyncLessonResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "任务添加成功"
	resp.Data.CourseName = req.CourseName
	resp.Data.ClassName = req.ClassName
	resp.Data.CourseStartDate = req.CourseStartDate
	resp.Data.CourseEndDate = req.CourseEndDate
	resp.Data.Total = courses.Data.TotalClassNum
	if courses.Data.TotalClassNum > 0 {
		//将任务发送到消息队列
		LessonSyncJobReq := &types.LessonSyncJobReq{}
		LessonSyncJobReq.CourseName = req.CourseName
		LessonSyncJobReq.ClassName = req.ClassName
		LessonSyncJobReq.CourseStartDate = req.CourseStartDate
		LessonSyncJobReq.CourseEndDate = req.CourseEndDate
		LessonSyncJobReq.Cookie = req.Cookie
		LessonSyncJobReq.UserAgent = l.ctx.Value("User-Agent").(string)
		LessonSyncJobReq.Total = courses.Data.TotalClassNum
		LessonSyncJobReq.SourceType = model.SourceTypeManual

		msg, _ := json.Marshal(LessonSyncJobReq)
		err = l.svcCtx.RabbitMqProducer.SendClassInMsgToMq(l.svcCtx.Config.RabbitMq.LessonSyncJobTopic, msg)
		//创建消息通知人员缓存
		qwUserId := l.ctx.Value("wid").(string)
		if qwUserId != "" {
			mu := sync.RWMutex{}
			mu.Lock()
			defer mu.Unlock()
			qwUserIds := ""
			_ = l.svcCtx.Cache.GetCtx(l.ctx, constant.ClassinJobUsersCacheKey, &qwUserIds)
			if qwUserIds == "" {
				qwUserIds = qwUserId
			} else {
				qwUserIds += "|" + qwUserId
			}
			_ = l.svcCtx.Cache.SetWithExpireCtx(l.ctx, constant.ClassinJobUsersCacheKey, qwUserIds, constant.ClassinJobUsersCacheExpiry)
			delayExpire := utils.GetCurrentTimestamp() + 120
			_ = l.svcCtx.Cache.SetWithExpireCtx(l.ctx, constant.ClassinJobUsersCacheDelayKey, delayExpire, constant.ClassinJobUsersCacheExpiry)
		}
	}
	return resp, nil
}
func (l *SyncJobLogic) SyncLesson(lessonClipSyncReq []byte) error {
	req := &types.LessonSyncJobReq{}
	err := json.Unmarshal(lessonClipSyncReq, req)
	if err != nil {
		l.Error(err.Error() + "源数据：" + string(lessonClipSyncReq))
		return nil
	}
	if req.UserAgent == "" || req.Cookie == "" || req.CourseStartDate <= 0 || req.CourseEndDate <= 0 {
		l.Error("消息中参数非法，获取数据为：" + string(lessonClipSyncReq))
		return nil
	}
	var jobMetric int64
	//页码大小
	pageSize := int64(100)
	//goroutineTotal
	goroutineTotal := 20
	totalPages := (req.Total / pageSize) + 1
	wg := &sync.WaitGroup{}
	ch := make(chan int64)
	for i := 0; i < goroutineTotal; i++ {
		wg.Add(1)
		go func(req *types.LessonSyncJobReq, wg *sync.WaitGroup, ch chan int64) {
			defer wg.Done()
			for page := range ch {
				params := l.getParams(req.CourseName, req.ClassName, req.CourseStartDate, req.CourseEndDate, page, pageSize)
				l.ctx = context.WithValue(l.ctx, "User-Agent", req.UserAgent)
				lessons, _ := l.svcCtx.ClassInCourseModel.GetCourseList(l.ctx, req.Cookie, params)
				if err == nil {
					atomic.AddInt64(&jobMetric, int64(len(lessons.Data.ClassList)))
					if len(lessons.Data.ClassList) > 0 {
						//将课节信息同步到
						lessonSyncReq := &types.LessonSyncReq{}
						for _, lesson := range lessons.Data.ClassList {
							lessonSyncReq.UserAgent = req.UserAgent
							lessonSyncReq.Cookie = req.Cookie
							lessonSyncReq.ClassId = lesson.ClassId
							lessonSyncReq.CourseId = lesson.CourseId
							lessonSyncReq.ClassName = lesson.ClassName
							lessonSyncReq.CourseName = lesson.CourseName
							lessonSyncReq.TeacherName = lesson.TeacherInfo.Name
							lessonSyncReq.TeacherMobile = lesson.TeacherInfo.Mobile
							lessonSyncReq.ClassStartTime = lesson.ClassStartTime
							lessonSyncReq.ClassEndTime = lesson.ClassEndTime
							lessonSyncReq.SourceType = model.SourceTypeManual
							lessonSyncReq.SeatNum = lesson.SeatNum
							lessonSyncReq.IsHd = lesson.IsHd
							lessonSyncReq.IsDc = lesson.IsDc
							msg, _ := json.Marshal(lessonSyncReq)
							err = l.svcCtx.RabbitMqProducer.SendClassInMsgToMq(l.svcCtx.Config.RabbitMq.ClassInCourseTopic, msg)
						}
					}
				}
			}

		}(req, wg, ch)
	}
	for i := int64(1); i <= totalPages; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	l.Info(jobMetric)
	return nil
}

func (l *SyncJobLogic) getParams(courseName string, className string, courseStartDate int64, courseEndDate int64, page int64, pageSize int64) url.Values {
	params := url.Values{}
	params.Set("withAuth", "1")      //默认值
	params.Set("classStatus", "1,3") //未开始与已结课
	if courseName != "" {
		params.Set("courseName", courseName)
	}
	if className != "" {
		searchKey := &types.ClassSearchKey{}
		searchKey.KeyName = "className"
		searchKey.KeyValue = className
		sk, _ := json.Marshal(searchKey)
		params.Set("searchKey", string(sk))
		sort := &types.Sort{
			SortName:  "classBtime",
			SortValue: 1,
		}
		s, _ := json.Marshal(sort)
		params.Set("sort", string(s))
	}
	if courseStartDate > 0 && courseEndDate > 0 && courseStartDate <= courseEndDate {
		timeRange := &types.TimeRange{}
		timeRange.StartTime = courseStartDate
		timeRange.EndTime = courseEndDate
		tr, _ := json.Marshal(timeRange)
		params.Set("timeRange", string(tr))
	}
	params.Set("page", strconv.FormatInt(page, 10))
	params.Set("perpage", strconv.FormatInt(pageSize, 10))
	return params
}
