package main

import (
	"context"
	"flag"
	"fmt"
	"meishiedu.com/classin/internal/logic"
	"meishiedu.com/classin/mq"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/config"
	"meishiedu.com/classin/internal/constant"
	"meishiedu.com/classin/internal/handler"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
)

var configFile = flag.String("f", "etc/classin-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	server := rest.MustNewServer(c.RestConf, rest.WithCors(), rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
		httpx.WriteJson(w, http.StatusOK, types.TokenErrorInfo{
			ErrorCode: constant.TokenExpired,
			ErrorMsg:  err.Error(),
		})
	}))
	defer server.Stop()
	httpx.SetErrorHandler(func(err error) (int, any) {
		responseBody := types.ErrorResponse{}
		responseBody.ErrorInfo.ErrorCode = 102
		responseBody.ErrorInfo.ErrorMsg = err.Error()
		return http.StatusOK, responseBody
	})
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	consumerCtx := context.Background()
	lessonClipsLogic := logic.NewLessonClipsLogic(consumerCtx, ctx)
	//课节同步生产者
	go mq.NewRabbitMqConsumer(consumerCtx, c.RabbitMq, c.RabbitMq.LessonSyncJobTopic).Consumer(logic.NewSyncJobLogic(consumerCtx, ctx).SyncLesson)
	//课程课节消费者
	go mq.NewRabbitMqConsumer(consumerCtx, c.RabbitMq, c.RabbitMq.ClassInCourseTopic).Consumer(logic.NewCourseListLogic(consumerCtx, ctx).AddCourses)
	//课节视频消费者
	go mq.NewRabbitMqConsumer(consumerCtx, c.RabbitMq, c.RabbitMq.ClassInLessonClipTopic).Consumer(lessonClipsLogic.AddLessonClipToMq)
	//视频落库消费者
	go mq.NewRabbitMqConsumer(consumerCtx, c.RabbitMq, c.RabbitMq.OSSUploadLessonClipTopic).Consumer(lessonClipsLogic.AddLessonClip)
	//视频上传消费者
	go lessonClipsLogic.UploadLessonClipToOss()

	//发送消息通知消费者
	go mq.NewRabbitMqConsumer(consumerCtx, c.RabbitMq, c.RabbitMq.SendMsgTopic).Consumer(ctx.SNotificationService.SendMsg)
	//go mq.NewRabbitMqConsumer(consumerCtx, c.RabbitMq, c.RabbitMq.ScheduleJobSyncTopic).Consumer(school.NewScheduleSyncLogic(consumerCtx, ctx).ScheduleDataSync)
	//go mq.NewRabbitMqConsumer(consumerCtx, c.RabbitMq, c.RabbitMq.ScheduleDataSyncTopic).Consumer(school.NewScheduleViewsLogic(consumerCtx, ctx).AddScheduleToTable)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
