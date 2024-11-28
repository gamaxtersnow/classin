package school

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"meishiedu.com/classin/internal/model"

	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/model/school"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

type ScheduleViewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScheduleViewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleViewsLogic {
	return &ScheduleViewsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScheduleViewsLogic) ScheduleViews(req *types.ScheduleViewsRequest) (resp *types.ScheduleViewsResponse, err error) {
	// 获取远程排课列表
	params := &school.ScheduleViewCond{
		Asc:        req.Asc,
		PageID:     req.PageId,
		PageSize:   req.PageSize,
		ViewType:   req.ViewType,
		Search:     req.Search,
		ShowCancel: req.ShowCancel,
		ShowRange:  req.ShowRange,
		ExceptNull: req.ExceptNull,
		Status:     req.Status,
		ClassIDs:   req.ClassIds,
		CampusIDs:  req.CampusIds,
		ClzRooms:   req.ClzRooms,
		CourseIDs:  req.CourseIds,
		ZhuJiaoID:  req.ZhujiaoId,
		TeacherIDs: req.TeacherIds,
		Dfrom:      req.Dfrom,
		Dto:        req.Dto,
	}
	logx.Info("ScheduleViews==请求参数:", params)
	list, err := l.svcCtx.SScheduleModel.GetScheduleList(l.ctx, params)
	if err != nil {
		return nil, err
	}
	//获取老师服务类型
	phones := make(map[string]bool)
	phonesStr := ""
	pageLen := 0
	for _, teacher := range list.Data {
		if teacher.Teacher.Mobile != "" {
			if _, ok := phones[teacher.Teacher.Mobile]; !ok {
				phones[teacher.Teacher.Mobile] = true
				phonesStr += teacher.Teacher.Mobile + ","
				pageLen++
			}
		}
	}
	logx.Infof("phonesStr: %s", phonesStr)
	teacherTypeMap, err := l.getTeacherTypeData(phonesStr, pageLen)
	if err != nil {
		return nil, err
	}
	//循环列表取数据并赋值到新的结构体中
	resp = &types.ScheduleViewsResponse{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取成功"
	resp.Pages.Total = list.Total
	resp.Pages.PageID = req.PageId
	resp.Pages.PageSize = req.PageSize
	resp.Pages.Pages = list.Total/req.PageSize + 1
	resp.ConflictCount = list.ConflictCount
	for _, v := range list.Data {
		scheduleView := types.ViewsData{}
		_ = utils.ConvertStruct(v, &scheduleView)
		_, ok := teacherTypeMap[v.Teacher.Mobile]
		if ok {
			scheduleView.CoopLevelValue = teacherTypeMap[v.Teacher.Mobile]
		}
		resp.ScheduleViews = append(resp.ScheduleViews, scheduleView)
	}
	return resp, nil
}
func (l *ScheduleViewsLogic) getTeacherTypeData(phones string, pageSize int) (map[string]int, error) {
	//定义一个map，key为手机号，value为老师详情
	teacherTypeMap := make(map[string]int)
	//获取老师类型列表
	typeList, err := l.svcCtx.STeacherModel.GetTeacherType(l.ctx, phones, pageSize)
	if err != nil {
		return nil, err
	}
	for _, typeDetail := range typeList {
		teacherTypeMap[typeDetail.Mobile] = typeDetail.Type
	}
	return teacherTypeMap, nil
}

// AddScheduleToTable 添加排课数据到本地数据库表
func (l *ScheduleViewsLogic) AddScheduleToTable(scheduleListSyncReq []byte) error {
	schedule := &school.Schedule{}
	err := json.Unmarshal(scheduleListSyncReq, schedule)
	if err := json.Unmarshal(scheduleListSyncReq, schedule); err != nil {
		return fmt.Errorf("解析排课数据失败: %w", err)
	}
	if schedule.ID <= 0 {
		return fmt.Errorf("无效的排课ID: %d", schedule.ID)
	}
	// 生成新数据的MD5
	newMD5, err := l.GenerateContentMD5(schedule)
	if err != nil {
		return fmt.Errorf("生成MD5失败: %w", err)
	}
	// 查找现有记录
	detail, err := l.svcCtx.LessonScheduleModel.FindOneBySid(l.ctx, schedule.ID)
	if err != nil && err != model.ErrNotFound {
		return fmt.Errorf("查询排课记录失败: %w", err)
	}
	// 准备新的数据模型
	addScheduleModel := &school.LessonSchedule{}
	teacherPerson := types.Person{}
	addScheduleModel.Sid = schedule.ID
	addScheduleModel.Clzid = schedule.ClzID
	addScheduleModel.Organizationid = schedule.OrganizationID
	addScheduleModel.Campusid = schedule.CampusID
	addScheduleModel.Courseid = schedule.CourseID
	addScheduleModel.ClzName = schedule.ClzName
	addScheduleModel.CourseName = schedule.CourseName
	addScheduleModel.Content = schedule.Content
	addScheduleModel.Starttime = schedule.StartTime
	addScheduleModel.Endtime = schedule.EndTime
	addScheduleModel.StarttimeStr = schedule.StartTimeStr
	addScheduleModel.EndtimeStr = schedule.EndTimeStr
	addScheduleModel.Duration = schedule.Duration
	addScheduleModel.Way = schedule.Way
	addScheduleModel.Place = schedule.Place
	addScheduleModel.Classtype = schedule.ClassType
	addScheduleModel.Coursetype = schedule.CourseType
	addScheduleModel.CourseCategory = schedule.CourseCategory
	addScheduleModel.Category = schedule.Category
	addScheduleModel.Status = schedule.Status
	addScheduleModel.Scountclz = schedule.SCountClz
	addScheduleModel.Scountjoin = schedule.SCountJoin
	addScheduleModel.Scountleave = schedule.SCountLeave
	teacherPerson.ID = schedule.Teacher.ID
	teacherPerson.Name = schedule.Teacher.Name
	teacherPerson.Mobile = schedule.Teacher.Mobile
	addScheduleModel.Teacher = l.toJsonString(teacherPerson)
	addScheduleModel.Note = l.toJsonString(schedule.Note)
	addScheduleModel.ClzZhujiao = l.toJsonString(schedule.ClzZhujiao)
	// 设置MD5值
	addScheduleModel.Md5Str = newMD5
	if detail == nil {
		// 新增记录
		addScheduleModel.Version = 1
		addScheduleModel.AddTime = time.Now().Unix()
		if _, err := l.svcCtx.LessonScheduleModel.Insert(l.ctx, addScheduleModel); err != nil {
			return fmt.Errorf("插入排课记录失败: %w", err)
		}
		return nil
	}
	// 如果MD5相同，说明内容没有变化，不需要更新
	if detail.Md5Str == newMD5 {
		//获取老数据的版本号,删除无用数据
		oldVersion := detail.Version
		logx.Infof("老版本号:%d,数据ID:%d", oldVersion, detail.Id)
		return nil
	}
	// 更新记录
	addScheduleModel.Id = detail.Id
	addScheduleModel.Sid = detail.Sid
	addScheduleModel.AddTime = detail.AddTime
	addScheduleModel.Version = detail.Version + 1
	addScheduleModel.UpdateTime = time.Now().Unix()
	if err := l.svcCtx.LessonScheduleModel.Update(l.ctx, addScheduleModel); err != nil {
		return fmt.Errorf("更新排课记录失败: %w", err)
	}
	return nil
}

// GenerateContentMD5 生成内容的MD5值
func (l *ScheduleViewsLogic) GenerateContentMD5(v interface{}) (string, error) {
	content, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("生成MD5时序列化失败: %w", err)
	}
	hash := md5.Sum(content)
	return fmt.Sprintf("%x", hash), nil
}
func (l *ScheduleViewsLogic) toJsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	if data, err := json.Marshal(v); err == nil && string(data) != "null" {
		return string(data)
	}
	return ""
}
