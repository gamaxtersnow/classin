package svc

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/crmsdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/notify"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/osssdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/rest"
	"meishiedu.com/classin/internal/config"
	"meishiedu.com/classin/internal/middleware"
	"meishiedu.com/classin/internal/model"
	"meishiedu.com/classin/internal/model/classin"
	"meishiedu.com/classin/internal/model/school"

	"meishiedu.com/classin/mq"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type ServiceContext struct {
	Config               config.Config
	UserAgentMiddleware  rest.Middleware
	Cache                cache.Cache
	CourseModel          model.CourseModel
	LessonClipModel      model.LessonClipModel
	RabbitMqProducer     *mq.RabbitProducer
	MsOssModel           osssdk.OSSModel
	ClassInCourseModel   *classin.CourseModel
	ClassInLessonModel   *classin.LessonModel
	SettingModel         *classin.SettingModel
	SCampusModel         school.CampusModel
	SScheduleModel       school.ScheduleModel
	STeacherModel        school.TeacherModel
	STeacherServiceModel school.TeacherServiceModel
	SCustomerModel       school.CustomerModel
	SClassModel          school.ClassModel
	SCourseModel         school.CourseModel
	SNotificationService notify.Notification
	SFilesModel          school.FilesModel
	LessonScheduleModel  school.LessonScheduleModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	connOldCrm := sqlx.NewMysql(c.Mysql.OldcrmDataSource)
	errNotFound := errors.New("item not found in school cache")
	cacheNode := cache.New(c.Cache, syncx.NewSingleFlight(), cache.NewStat("school"), errNotFound)
	httpClient := &http.Client{}
	cookieJar, _ := cookiejar.New(nil)
	httpClient.Jar = cookieJar
	httpClient.Timeout = 20 * time.Second
	xHttpClient := xiaoxiaosdk.NewHttpClient(c.XiaoxiaoApiConf)
	mHttpClient := crmsdk.NewHttpClient(c.CrmApiConf)
	c.NotifyConf.TitanConfig = c.TiTanApiConf
	ctx := context.Background()
	return &ServiceContext{
		Config:               c,
		UserAgentMiddleware:  middleware.NewUserAgentMiddleware().Handle,
		Cache:                cacheNode,
		CourseModel:          model.NewCourseModel(conn, c.Cache),
		LessonClipModel:      model.NewLessonClipModel(conn, c.Cache),
		RabbitMqProducer:     mq.NewRabbitProducer(ctx, c.RabbitMq),
		MsOssModel:           osssdk.NewAliYunOssModel(c.Oss),
		ClassInCourseModel:   classin.NewCourseModel(httpClient),
		ClassInLessonModel:   classin.NewLessonModel(httpClient),
		SettingModel:         classin.NewSettingModel(httpClient),
		SCampusModel:         school.NewCampusModel(xHttpClient, cacheNode),
		SScheduleModel:       school.NewScheduleModel(xHttpClient, cacheNode),
		STeacherModel:        school.NewTeacherModel(xHttpClient, mHttpClient, cacheNode),
		STeacherServiceModel: school.NewCustomTeacherServiceModel(connOldCrm, cacheNode),
		SCustomerModel:       school.NewCustomerModel(xHttpClient, connOldCrm, cacheNode),
		SClassModel:          school.NewClassModel(xHttpClient, cacheNode),
		SCourseModel:         school.NewCourseModel(xHttpClient, cacheNode),
		SNotificationService: notify.NewNotificationService(c.NotifyConf),
		SFilesModel:          school.NewFilesModel(conn, c.Cache),
		LessonScheduleModel:  school.NewLessonScheduleModel(conn, c.Cache),
	}
}
