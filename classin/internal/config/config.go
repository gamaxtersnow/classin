package config

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/crmsdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/notify"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/osssdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/titansdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	FrontendDomain       string
	ResourceCenterDomain string
	Auth                 struct {
		AccessSecret string
		AccessExpire int64
	}
	OfficeAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	Mysql struct {
		DataSource       string
		OldcrmDataSource string
	}
	Cache           cache.CacheConf
	RabbitMq        RabbitMqConf
	Oss             osssdk.OssConf
	XiaoxiaoApiConf xiaoxiaosdk.XiaoxiaoApiConf
	CrmApiConf      crmsdk.CrmApiConf
	TiTanApiConf    titansdk.TiTanApiConf
	NotifyConf      notify.Conf
}

type RabbitMqConf struct {
	Uri                      string
	Exchange                 string
	ExchangeType             string
	LessonSyncJobTopic       Topic
	ClassInCourseTopic       Topic
	ClassInLessonClipTopic   Topic
	OSSUploadLessonClipTopic Topic
	SendMsgTopic             Topic
	ScheduleJobSyncTopic     Topic
	ScheduleDataSyncTopic    Topic
}
type Topic struct {
	RoutingKey string
	Queue      string
	ConnName   string
	CTag       string
}
type SchoolApiConf struct {
	BaseUrl string
}
type MsCrmApiConf struct {
	UserId   string
	Email    string
	Rate     int `json:",default=300"`
	LoginUrl string
	ApiUrl   string
}
type Api struct {
	Name string
	Uri  string
}
