// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	customer "meishiedu.com/classin/internal/handler/customer"
	file "meishiedu.com/classin/internal/handler/file"
	school "meishiedu.com/classin/internal/handler/school"
	utils "meishiedu.com/classin/internal/handler/utils"
	"meishiedu.com/classin/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.UserAgentMiddleware},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/course/list",
					Handler: CourseListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/lesson/clips",
					Handler: LessonClipsHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/lesson/setNum",
					Handler: SetNumHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/lesson/sync",
					Handler: SyncJobHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/login",
					Handler: LoginHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/oss/signurl",
					Handler: SignUrlHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// 查询客户基本信息
				Method:  http.MethodPost,
				Path:    "/customer/customerBasicData",
				Handler: customer.GetCustomerBasicDataHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/school"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// 根据id删除文件
				Method:  http.MethodPost,
				Path:    "/file/deleteFileById",
				Handler: file.DeleteFileByIdHandler(serverCtx),
			},
			{
				// 根据uuid删除文件
				Method:  http.MethodPost,
				Path:    "/file/deleteFileByUUID",
				Handler: file.DeleteFileByUUIDHandler(serverCtx),
			},
			{
				// 根据uuid获取文件下载地址
				Method:  http.MethodPost,
				Path:    "/file/getFileAddressByUUID",
				Handler: file.GetFileAddressByUUIDHandler(serverCtx),
			},
			{
				// 根据id获取文件信息
				Method:  http.MethodPost,
				Path:    "/file/getFileById",
				Handler: file.GetFileByIdHandler(serverCtx),
			},
			{
				// 根据uuid获取文件信息
				Method:  http.MethodPost,
				Path:    "/file/getFileByUUID",
				Handler: file.GetFileByUUIDHandler(serverCtx),
			},
			{
				// 根据uuid获取文件预览信息
				Method:  http.MethodPost,
				Path:    "/file/getFileViewByUUID",
				Handler: file.GetFileViewByUUIDHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/school"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// 校区列表
				Method:  http.MethodPost,
				Path:    "/campus/list",
				Handler: school.CampusListHandler(serverCtx),
			},
			{
				// 校区教室列表
				Method:  http.MethodPost,
				Path:    "/campus/room/list",
				Handler: school.CampusRoomListHandler(serverCtx),
			},
			{
				// 班级列表
				Method:  http.MethodPost,
				Path:    "/class/list",
				Handler: school.GetClassListHandler(serverCtx),
			},
			{
				// 全部课程
				Method:  http.MethodPost,
				Path:    "/course/all",
				Handler: school.AllCoursesHandler(serverCtx),
			},
			{
				// 查看文件链接
				Method:  http.MethodGet,
				Path:    "/files/sign/url",
				Handler: school.FilesSignUrlHandler(serverCtx),
			},
			{
				// 新增排课获取老师和服务价格数据
				Method:  http.MethodPost,
				Path:    "/makeLesson/teacherPrice",
				Handler: school.MakeLessonGetTeacherPriceHandler(serverCtx),
			},
			{
				// 课表导出
				Method:  http.MethodPost,
				Path:    "/schedule/export",
				Handler: school.ScheduleExportHandler(serverCtx),
			},
			{
				// 校校课表数据同步落库
				Method:  http.MethodPost,
				Path:    "/schedule/sync",
				Handler: school.ScheduleSyncHandler(serverCtx),
			},
			{
				// 课表列表
				Method:  http.MethodPost,
				Path:    "/schedule/views",
				Handler: school.ScheduleViewsHandler(serverCtx),
			},
			{
				// 全部老师
				Method:  http.MethodPost,
				Path:    "/teacher/all",
				Handler: school.AllTeachersHandler(serverCtx),
			},
			{
				// 全部助教
				Method:  http.MethodPost,
				Path:    "/teacher/assistant/all",
				Handler: school.AllTeacherAssistantsHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/school"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// url转pdf
				Method:  http.MethodPost,
				Path:    "/utils/urlToPdf",
				Handler: utils.UrlToPdfHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/school"),
	)
}
