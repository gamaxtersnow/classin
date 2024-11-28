package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk/xapi"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"strconv"
)

type Class struct {
	Pages ClassListByRolePages // 分页信息
	Data  []ClassCourse        // 课程数据列表
}

type ClassListByRolePages struct {
	Total    int // 总记录数
	PageSize int // 每页大小
	PageID   int // 当前页码
	Pages    int // 总页数
}

type ClassCourse struct {
	ID         int    // 班级ID
	Name       string // 班级名称
	CampusID   int    // 校区 ID
	CampusName string // 校区名称
	CourseID   int    // 课程 ID
	CourseName string // 课程名称
	ClassType  int    // 课程类型
	Category   string // 课程类别
	Status     int    // 课程状态
}

var _ ClassModel = (*customClassModel)(nil)

type (
	ClassModel interface {
		GetClassListByRole(ctx context.Context, pageId int, pageSize int, search string, courseId string, campusIds string, status int) (*Class, error)
	}
	customClassModel struct {
		model xapi.ClassModel
		cache cache.Cache
	}
)

func NewClassModel(client *xiaoxiaosdk.HttpClient, cache cache.Cache) ClassModel {
	return &customClassModel{
		model: xapi.NewClassModel(client),
		cache: cache,
	}
}

// GetClassListByRole 通过角色获取班级列表
func (t *customClassModel) GetClassListByRole(ctx context.Context, pageId int, pageSize int, search string, courseId string, campusIds string, status int) (*Class, error) {
	class := &Class{}
	req := xapi.ClassListByRoleReq{
		Search:    search,
		CourseID:  courseId,
		CampusIDs: campusIds,
		PageID:    pageId,
		PageSize:  pageSize,
		Status:    "",
	}
	if status != -1 {
		req.Status = strconv.Itoa(status)
	}
	classes, err := t.model.GetClassListByRole(ctx, req)
	if err != nil {
		return class, err
	}
	if !classes.Success {
		return class, errors.New(classes.Messages)
	}
	class.Pages = ClassListByRolePages(classes.Pages)
	for _, c := range classes.Data {
		class.Data = append(class.Data, ClassCourse(c))
	}
	return class, nil
}
