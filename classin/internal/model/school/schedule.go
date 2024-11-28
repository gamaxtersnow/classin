package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk/xapi"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type ScheduleViewCond struct {
	ViewType   int64  // 视图类型
	Search     string // 搜索关键字
	CourseIDs  string // 课程 ID 列表
	CampusIDs  string // 校区 ID 列表
	TeacherIDs string // 教师 ID 列表
	ZhuJiaoID  string // 助教 ID
	ClassIDs   string // 班级 ID 列表
	Status     string // 状态
	PageID     int    // 当前页码
	PageSize   int    // 每页大小
	ClzRooms   string // 教室列表
	ShowRange  int64  // 显示范围
	ExceptNull int64  // 排除空值
	ShowCancel int64  // 显示取消
	Asc        int64  // 升序或降序
	Dfrom      string //日期起始时间
	Dto        string //日期结束时间
	Dlall      int    // 导出全部
}
type ScheduleViews struct {
	Total         int        // 总数
	ConflictCount int        // 冲突计数
	Data          []Schedule // 课程安排列表
}

type Schedule struct {
	ID                int64       // 课程安排 ID
	LockVersion       int64       // 锁版本
	ClzID             int64       // 班级 ID
	OrganizationID    int64       // 组织 ID
	CourseID          int64       // 课程 ID
	CourseName        string      // 课程名称
	CoopLevelValue    int         //兼职分级数值
	CoopLevel         string      // 兼职分级名称
	CoopType          string      //兼职类型
	Teacher           Person      // 教师信息
	StartTime         int64       // 开始时间（毫秒）
	EndTime           int64       // 结束时间（毫秒）
	Way               int64       // 方式
	Content           string      // 内容
	Place             string      // 地点
	Duration          string      // 持续时间（格式化）
	DurationMinutes   int64       // 持续时间（分钟）
	FromSmart         bool        // 是否来自智能
	TDay              int64       // 天
	Category          int64       // 类别
	CourseType        int64       // 课程类型
	TagContent        string      // 标签内容
	TagCIndex         int64       // 标签索引
	Ended             bool        // 是否结束
	IsConflict        bool        // 是否冲突
	PreoccupyConflict bool        // 预占冲突
	StudentConflict   bool        // 学生冲突
	IFromSmart        bool        // 是否来自智能
	Status            int64       // 状态
	EdStatus          bool        // 编辑状态
	Attendanced       bool        // 是否考勤
	Attendance        int64       // 考勤
	SCountClz         int64       // 班级计数
	SCountJoin        int64       // 加入计数
	SCountLeave       int64       // 离开计数
	StartTimeStr      string      // 开始时间（字符串）
	EndTimeStr        string      // 结束时间（字符串）
	Diseditable       bool        // 是否可编辑
	PeriodOfClz       float64     // 班级周期
	PeriodOfCur       float64     // 当前周期
	ClassPlatform     int64       // 班级平台
	CourseCategory    string      // 课程类别
	ClzName           string      // 班级名称
	ScheduleCount     int64       // 日程计数
	TheSmart          bool        // 智能
	IngStatus         int64       // 进行状态
	ClassType         int64       // 班级类型
	CampusID          int64       // 校区 ID
	CampusNames       string      // 校区名称
	ClzZhujiao        []Person    // 班级助教
	Note              interface{} //备注
}
type Person struct {
	ID             int64       // 人员 ID
	OrganizationID int64       // 组织 ID
	Name           string      // 姓名
	Mobile         string      // 手机号
	Role           int64       // 角色
	Teaching       interface{} // 教学
	PhotoURL       interface{} // 照片 URL
	Wechat         interface{} // 微信
	Gender         interface{} // 性别
	Introduce      interface{} // 介绍
	ClassInUID     int64       // ClassIn 用户 ID
	Token          interface{} // 令牌
	Gzhvc          interface{} // 公众号验证
	FeedBack       interface{} // 反馈
	AccPriosList   interface{} // 账户优先级列表
}

var _ ScheduleModel = (*customScheduleModel)(nil)

type (
	ScheduleModel interface {
		GetScheduleList(ctx context.Context, req *ScheduleViewCond) (*ScheduleViews, error)
	}
	customScheduleModel struct {
		model xapi.ScheduleModel
		//detail mapi.TeacherModel
		cache cache.Cache
	}
)

func NewScheduleModel(xClient *xiaoxiaosdk.HttpClient, cache cache.Cache) ScheduleModel {
	return &customScheduleModel{
		model: xapi.NewScheduleModel(xClient),
		cache: cache,
	}
}

// GetScheduleList 获取排课列表
func (s *customScheduleModel) GetScheduleList(ctx context.Context, req *ScheduleViewCond) (*ScheduleViews, error) {
	param := &xapi.ScheduleViewReq{
		ViewType:   req.ViewType,
		ExceptNull: req.ExceptNull,
		ShowCancel: req.ShowCancel,
		PageID:     req.PageID,
		PageSize:   req.PageSize,
		Asc:        req.Asc,
		Status:     req.Status,
		Search:     req.Search,
		ShowRange:  req.ShowRange,
		ClassIDs:   req.ClassIDs,
		CampusIDs:  req.CampusIDs,
		ClzRooms:   req.ClzRooms,
		CourseIDs:  req.CourseIDs,
		ZhuJiaoID:  req.ZhuJiaoID,
		TeacherIDs: req.TeacherIDs,
		Dfrom:      req.Dfrom,
		Dto:        req.Dto,
	}
	schedules := &ScheduleViews{}
	scheduleViewRes, err := s.model.GetScheduleList(ctx, param)
	if err != nil {
		return schedules, err
	}
	if !scheduleViewRes.Success {
		return schedules, errors.New(scheduleViewRes.Messages)
	}
	schedules.Total = scheduleViewRes.Total
	schedules.ConflictCount = scheduleViewRes.ConflictCount
	for _, schedule := range scheduleViewRes.Data {
		sch := Schedule{}
		_ = utils.ConvertStruct(schedule, &sch)
		schedules.Data = append(schedules.Data, sch)
	}
	return schedules, nil
}
