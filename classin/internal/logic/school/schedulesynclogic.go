package school

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/model/school"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
	"sync"
	"sync/atomic"
)

type ScheduleSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewScheduleSyncLogic 校校课表数据同步落库
func NewScheduleSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleSyncLogic {
	return &ScheduleSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ScheduleSync 获取远程接口请求条件并存放到生产者队列中
func (l *ScheduleSyncLogic) ScheduleSync(req *types.ScheduleSyncRequest) (resp *types.ScheduleSyncResponse, err error) {
	params := &school.ScheduleViewCond{
		ViewType:   req.ViewType,
		PageID:     req.PageID,
		PageSize:   req.PageSize,
		Dfrom:      req.Dfrom,
		Dto:        req.Dto,
		Search:     req.Search,
		CampusIDs:  req.CampusIDs,
		ClassIDs:   req.ClassIDs,
		CourseIDs:  req.CourseIDs,
		ClzRooms:   req.ClzRooms,
		TeacherIDs: req.TeacherIds,
		ZhuJiaoID:  req.ZhuJiaoID,
		Status:     req.Status,
		ShowCancel: req.ShowCancel,
		ExceptNull: req.ExceptNull,
		Asc:        req.Asc,
	}
	remoteScheduleList, err := l.svcCtx.SScheduleModel.GetScheduleList(l.ctx, params)
	if err != nil {
		return nil, err
	}
	if remoteScheduleList.Data == nil {
		return nil, errors.New("获取远程课表暂无数据")
	}
	resp = &types.ScheduleSyncResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "添加远程课表任务成功"
	resp.Data.Dto = req.Dto
	resp.Data.Dfrom = req.Dfrom
	resp.Data.Search = req.Search
	resp.Data.PageID = req.PageID
	resp.Data.PageSize = req.PageSize
	resp.Data.Status = req.Status
	resp.Data.ViewType = req.ViewType
	resp.Data.ZhuJiaoID = req.ZhuJiaoID
	resp.Data.TeacherIds = req.TeacherIds
	resp.Data.ClassIDs = req.ClassIDs
	resp.Data.ClzRooms = req.ClzRooms
	resp.Data.CourseIDs = req.CourseIDs
	resp.Data.CampusIDs = req.CampusIDs
	resp.Data.Total = remoteScheduleList.Total
	if remoteScheduleList.Total > 0 {
		mqParams := &types.SyncData{}
		mqParams.Search = req.Search
		mqParams.Status = req.Status
		mqParams.Dto = req.Dto
		mqParams.Dfrom = req.Dfrom
		mqParams.PageID = req.PageID
		mqParams.PageSize = req.PageSize
		mqParams.ClassIDs = req.ClassIDs
		mqParams.ClzRooms = req.ClzRooms
		mqParams.CourseIDs = req.CourseIDs
		mqParams.CampusIDs = req.CampusIDs
		mqParams.ZhuJiaoID = req.ZhuJiaoID
		mqParams.TeacherIds = req.TeacherIds
		mqParams.ViewType = req.ViewType
		mqParams.ShowCancel = req.ShowCancel
		mqParams.ExceptNull = req.ExceptNull
		mqParams.Asc = req.Asc
		mqParams.Total = len(remoteScheduleList.Data) //待同步数据总数 remoteScheduleList.Total
		//将任务发到生产者MQ队列
		msg, _ := json.Marshal(mqParams)
		err = l.svcCtx.RabbitMqProducer.SendClassInMsgToMq(l.svcCtx.Config.RabbitMq.ScheduleJobSyncTopic, msg)
	}
	return resp, nil
}

// ScheduleDataSync 根据条件获取远程接口数据并存放到生产者队列中
func (l *ScheduleSyncLogic) ScheduleDataSync(scheduleListSyncReq []byte) error {
	req := &types.SyncData{}
	err := json.Unmarshal(scheduleListSyncReq, req)
	if err != nil {
		l.Error(err.Error() + "校校源数据：" + string(scheduleListSyncReq))
	}
	var jobMetric int64
	//页码大小
	pageSize := req.PageSize
	goroutineTotal := 20
	totalPages := (req.Total / pageSize) + 1
	wg := &sync.WaitGroup{}
	ch := make(chan int64)
	for i := 0; i < goroutineTotal; i++ {
		wg.Add(1)
		go func(req *types.SyncData, wg *sync.WaitGroup, ch chan int64) {
			defer wg.Done()
			for page := range ch {
				params := &school.ScheduleViewCond{
					Dfrom:      req.Dfrom,
					Dto:        req.Dto,
					Search:     req.Search,
					PageID:     int(page),
					PageSize:   pageSize,
					Status:     req.Status,
					ViewType:   req.ViewType,
					ZhuJiaoID:  req.ZhuJiaoID,
					ClassIDs:   req.ClassIDs,
					CourseIDs:  req.CourseIDs,
					ClzRooms:   req.ClzRooms,
					TeacherIDs: req.TeacherIds,
					ShowCancel: req.ShowCancel,
					ExceptNull: req.ExceptNull,
					Asc:        req.Asc,
				}
				scheduleLists, _ := l.svcCtx.SScheduleModel.GetScheduleList(l.ctx, params)
				if err == nil {
					atomic.AddInt64(&jobMetric, int64(len(scheduleLists.Data)))
					if len(scheduleLists.Data) > 0 {
						//将课表列表记录同步到MQ
						scheduleSyncReq := &types.ViewsData{}
						for _, schedule := range scheduleLists.Data {
							scheduleSyncReq.ID = schedule.ID
							scheduleSyncReq.ClzID = schedule.ClzID
							scheduleSyncReq.CourseID = schedule.CourseID
							scheduleSyncReq.CampusID = schedule.CampusID
							scheduleSyncReq.OrganizationID = schedule.OrganizationID
							scheduleSyncReq.Way = schedule.Way
							scheduleSyncReq.CourseName = schedule.CourseName
							scheduleSyncReq.CourseType = schedule.CourseType
							scheduleSyncReq.StartTime = schedule.StartTime
							scheduleSyncReq.EndTime = schedule.EndTime
							scheduleSyncReq.ClzName = schedule.ClzName
							scheduleSyncReq.ClzZhujiao = schedule.ClzZhujiao
							scheduleSyncReq.Content = schedule.Content
							scheduleSyncReq.Place = schedule.Place
							scheduleSyncReq.Duration = schedule.Duration
							scheduleSyncReq.DurationMinutes = schedule.DurationMinutes
							scheduleSyncReq.Note = schedule.Note
							scheduleSyncReq.Category = schedule.Category
							scheduleSyncReq.ClassType = schedule.ClassType
							scheduleSyncReq.CourseType = schedule.CourseType
							scheduleSyncReq.CourseCategory = schedule.CourseCategory
							scheduleSyncReq.Status = schedule.Status
							scheduleSyncReq.SCountClz = schedule.SCountClz
							scheduleSyncReq.SCountJoin = schedule.SCountJoin
							scheduleSyncReq.SCountLeave = schedule.SCountLeave
							scheduleSyncReq.StartTimeStr = schedule.StartTimeStr
							scheduleSyncReq.EndTimeStr = schedule.EndTimeStr
							scheduleSyncReq.Teacher.ID = schedule.Teacher.ID
							scheduleSyncReq.Teacher.Name = schedule.Teacher.Name
							scheduleSyncReq.Teacher.Mobile = schedule.Teacher.Mobile
							scheduleSyncReq.Teacher.ClassInUID = schedule.Teacher.ClassInUID
							msg, _ := json.Marshal(scheduleSyncReq)
							err = l.svcCtx.RabbitMqProducer.SendClassInMsgToMq(l.svcCtx.Config.RabbitMq.ScheduleDataSyncTopic, msg)
						}
					}
				}
			}

		}(req, wg, ch)
	}
	for i := int64(1); i <= int64(totalPages); i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	l.Info(jobMetric)
	return nil
}
