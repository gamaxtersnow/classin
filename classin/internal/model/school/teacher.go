package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/crmsdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/crmsdk/mapi"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk/xapi"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Teacher struct {
	ID             int    // 用户 ID
	OrganizationID int    // 组织 ID
	Password       string // 密码
	Name           string // 姓名
	Role           int    // 角色
	Mobile         string // 手机号
	Teaching       string // 教学信息（可能为空）
	Username       string // 用户名
	SysUse         bool   // 系统使用标志
	PhotoURL       string // 照片 URL（可能为空）
	Wechat         string // 微信（可能为空）
	Gender         string // 性别（可能为空）
	Introduce      string // 介绍（可能为空）
	TeachLength    string // 教学时长（可能为空）
	TeachStyle     string // 教学风格（可能为空）
	City           string // 城市（可能为空）
	PhotoPro       string // 照片专业（可能为空）
	VideoPro       string // 视频专业（可能为空）
	ShowYangyu     bool   // 是否显示杨语
	ShowExpYangyu  bool   // 是否显示杨语经验
	ClassInUID     int    // ClassIn 用户 ID
	BDel           bool   // 删除标志
	RAccountID     int    // 账户 ID（可能为空）
	TeaDuration    string // 教学持续时间
	Token          string // 令牌（可能为空）
	Gzhvc          string // 公众号验证（可能为空）
	AccPriosList   string // 账户优先级列表（可能为空）
	Disabled       bool   // 是否禁用
	JustRecovery   bool   // 是否刚恢复
	Type           int    // 老师类型
	TypeName       string // 老师类型名称
}

var _ TeacherModel = (*customTeacherModel)(nil)

type (
	TeacherModel interface {
		GetTeachers(ctx context.Context, pageId int, pageSize int) ([]Teacher, int, error)
		GetTeacherAssistants(ctx context.Context, pageId int, pageSize int) ([]Teacher, int, error)
		GetTeacherType(ctx context.Context, phones string, pageSize int) ([]Teacher, error)
	}
	customTeacherModel struct {
		model             xapi.TeacherModel
		cache             cache.Cache
		mTeacherTypeModel mapi.TeacherModel
	}
)

func NewTeacherModel(client *xiaoxiaosdk.HttpClient, mClient *crmsdk.HttpClient, cache cache.Cache) TeacherModel {
	return &customTeacherModel{
		model:             xapi.NewTeacherModel(client),
		cache:             cache,
		mTeacherTypeModel: mapi.NewTeacherModel(mClient),
	}
}

// GetTeachers 获取老师列表
func (t *customTeacherModel) GetTeachers(ctx context.Context, pageId int, pageSize int) ([]Teacher, int, error) {
	return t.getTeacherList(ctx, pageId, pageSize, 3)
}

// GetTeacherAssistants 获取助教列表
func (t *customTeacherModel) GetTeacherAssistants(ctx context.Context, pageId int, pageSize int) ([]Teacher, int, error) {
	return t.getTeacherList(ctx, pageId, pageSize, 4)
}

func (t *customTeacherModel) getTeacherList(ctx context.Context, pageId int, pageSize int, role int) ([]Teacher, int, error) {
	var teacherList []Teacher
	req := xapi.TeacherListReq{
		Role:     role,
		PageSize: pageSize,
		PageID:   pageId,
	}
	teachers, err := t.model.GetTeacherList(ctx, req)
	if err != nil {
		return teacherList, 0, err
	}
	if !teachers.Success {
		return teacherList, 0, errors.New(teachers.Messages)
	}
	for _, teacher := range teachers.Data {
		teacherStruct := Teacher{}
		_ = utils.ConvertStruct(teacher, &teacherStruct)
		teacherList = append(teacherList, teacherStruct)
	}
	return teacherList, teachers.Total, nil
}

func (t *customTeacherModel) GetTeacherType(ctx context.Context, phones string, pageSize int) ([]Teacher, error) {
	req := mapi.CrmTeacherRequest{
		Page:       1,
		PageSize:   pageSize,
		MultiPhone: phones,
	}
	teacherTypeList, err := t.mTeacherTypeModel.GetTeacherDetailList(ctx, req)
	if err != nil {
		return nil, err
	}
	var teacherTypeDataList []Teacher
	for _, teacherType := range teacherTypeList.Data.List {
		teacherTypeDataList = append(teacherTypeDataList, Teacher{
			ID:       teacherType.ID,
			Mobile:   teacherType.Phone,
			TypeName: teacherType.CoopType,
			Type:     teacherType.CoopLevelValue,
		})
	}
	return teacherTypeDataList, nil
}
