package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk/xapi"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"time"
)

type Course struct {
	ID                  int         // 课程 ID
	OrganizationID      int         // 组织 ID
	Name                string      // 课程名称
	Classification      string      // 课程分类
	Category            int         // 课程类别
	SwPrice             bool        // 价格开关标志
	FeeForm             string      // 费用形式
	Period              string      // 期间
	PriceRaw            float64     // 原始价格
	Price               float64     // 价格
	Description         string      // 描述
	Pics                string      // 图片
	ContactWeixin       string      // 联系微信
	ContactWeixinV2Code string      // 联系微信 V2 代码
	CtContents          []CtContent // 课程内容
	IPkg                bool        // IPkg 标志
	PkgItems            []string    // 包含项目
	Kind                int         // 类型
	PricePerHour        float64     // 每小时价格
	JoinedClz           string      // 加入的班级
	JoinedClzList       string      // 加入的班级列表
	PkgIncJson          string      // 包含 JSON
	CotPkgCourse        int         // 包课程数量
	CotPkgSundry        int         // 包杂项数量
	IngForAdd           bool        // 添加标志
}

type CtContent struct {
	ID             int    // 内容 ID
	OrganizationID int    // 组织 ID
	CourseID       int    // 课程 ID
	Content        string // 内容描述
}

var (
	_                     CourseModel = (*customCourseModel)(nil)
	cacheCourseAllKey                 = "cache:school:course:all"
	cacheCourseExpiration             = 2 * time.Hour
)

type (
	CourseModel interface {
		GetAllCourses(ctx context.Context) ([]Course, error)
	}
	customCourseModel struct {
		model xapi.CourseModel
		cache cache.Cache
	}
)

func NewCourseModel(client *xiaoxiaosdk.HttpClient, cache cache.Cache) CourseModel {
	return &customCourseModel{
		model: xapi.NewCourseModel(client),
		cache: cache,
	}
}

// GetAllCourses 获取所有课程
func (c *customCourseModel) GetAllCourses(ctx context.Context) ([]Course, error) {
	var courseList []Course
	if err := c.cache.GetCtx(ctx, cacheCourseAllKey, &courseList); err == nil {
		return courseList, nil
	}
	courses, err := c.model.GetAllCourses(ctx)
	if err != nil {
		return nil, err
	}
	if !courses.Success {
		return courseList, errors.New(courses.Messages)
	}
	if courses.Total > 0 {
		for _, courseInfo := range courses.Data {
			course := Course{}
			course.ID = courseInfo.ID
			course.Name = courseInfo.Name
			course.Category = courseInfo.Category
			course.Classification = courseInfo.Classification
			course.Pics = courseInfo.Pics
			course.Description = courseInfo.Description
			courseList = append(courseList, course)
		}
	}
	_ = c.cache.SetWithExpireCtx(ctx, cacheCourseAllKey, courseList, cacheCourseExpiration)
	return courseList, nil
}
