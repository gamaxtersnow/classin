package logic

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/notify"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/notify/template/classin"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"encoding/json"
	"errors"
	"meishiedu.com/classin/internal/constant"
	"meishiedu.com/classin/internal/model"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonClipsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLessonClipsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonClipsLogic {
	return &LessonClipsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LessonClipsLogic) LessonClips(req *types.LesssonClipRequest) (resp *types.LessonClipResponse, err error) {
	uniqueId := strconv.FormatInt(req.CourseId, 10) + "-" + strconv.FormatInt(req.ClassId, 10)
	lessonInfo, _ := l.svcCtx.CourseModel.FindByUid(l.ctx, uniqueId)
	count, err := l.svcCtx.LessonClipModel.CountByCourseIdAndClassId(l.ctx, req.CourseId, req.ClassId)
	LessonClipResponse := &types.LessonClipResponse{}
	LessonClipResponse.ErrorInfo.ErrorCode = 1
	LessonClipResponse.ErrorInfo.ErrorMsg = "获取成功"
	LessonClipResponse.VodInfo.AllCount = count
	if count > 0 {
		offset := (req.Page - 1) * req.PageSize
		lessonClips, _ := l.svcCtx.LessonClipModel.FindByCourseIdAndClassId(l.ctx, req.CourseId, req.ClassId, offset, req.PageSize)
		if lessonInfo != nil {
			LessonClipResponse.VodInfo.LessonInfo.CourseName = lessonInfo.CourseName
			LessonClipResponse.VodInfo.LessonInfo.ClassName = lessonInfo.ClassName
		}
		for _, lessonClip := range lessonClips {
			lc := types.LessonClipFile{}
			lc.CourseId = lessonClip.CourseId
			lc.ClassId = lessonClip.ClassId
			lc.OriginUrl = lessonClip.FileOriginUrl
			lc.Objectkey = lessonClip.ObjectKey
			lc.FileStatus = lessonClip.FileStatus
			lc.FileSize = lessonClip.FileSize
			lc.SequenceNumber = lessonClip.SequenceNumber
			lc.FileId = lessonClip.FileId
			lc.FileStatusName = l.getLessonClipFileStatusName(lessonClip.FileStatus)
			LessonClipResponse.VodInfo.FileList = append(LessonClipResponse.VodInfo.FileList, lc)
		}
	}
	return LessonClipResponse, nil
}
func (l *LessonClipsLogic) AddLessonClipToMq(lessonClipSyncReq []byte) error {
	req := &types.LessonClipSyncReq{}
	err := json.Unmarshal(lessonClipSyncReq, req)
	if err != nil {
		l.Error(err.Error() + "源数据：" + string(lessonClipSyncReq))
		return nil
	}

	if req.CourseId <= 0 || req.ClassId <= 0 || req.UserAgent == "" || req.Cookie == "" {
		l.Error("消息非法，源数据：" + string(lessonClipSyncReq))
		return nil
	}
	if req.SourceType == model.FileSourceTypeManual {
		//拉取课节视频数据
		l.ctx = context.WithValue(l.ctx, "User-Agent", req.UserAgent)
		lessonClip, err := l.svcCtx.ClassInLessonModel.GetLessonClipList(l.ctx, req.Cookie, l.getParams(req.CourseId, req.ClassId, 1, 100))
		if err != nil {
			return err
		}
		if lessonClip.ErrorInfo.ErrorCode != 1 {
			return errors.New(lessonClip.ErrorInfo.ErrorMsg)
		}
		if lessonClip.Data.VodInfo.AllCount > 0 {
			for _, file := range lessonClip.Data.VodInfo.FileList {
				sequenceNumber := int64(1)
				for _, clip := range file.PlaySet {
					lessonClipModel := &model.LessonClip{}
					lessonClipModel.CourseId = req.CourseId
					lessonClipModel.ClassId = req.ClassId
					lessonClipModel.SequenceNumber = sequenceNumber
					lessonClipModel.FileOriginUrl = clip.Url
					lessonClipModel.FileSize, _ = strconv.ParseInt(file.FileSize, 10, 64)
					lessonClipModel.FileStatus = model.FileStatusInit
					lessonClipModel.FileId = file.FileId
					lessonClipModel.SourceType = req.SourceType
					sequenceNumber++
					//发送消息到上传对队列
					msg, _ := json.Marshal(lessonClipModel)
					err = l.svcCtx.RabbitMqProducer.SendClassInMsgToMq(l.svcCtx.Config.RabbitMq.OSSUploadLessonClipTopic, msg)
				}
			}
		}
	}
	return nil
}
func (l *LessonClipsLogic) AddLessonClip(lessonClip []byte) error {
	clip := &model.LessonClip{}
	err := json.Unmarshal(lessonClip, clip)
	if err != nil {
		l.Error(err.Error() + "源数据：" + string(lessonClip))
		return nil
	}
	if clip.FileOriginUrl == "" {
		l.Error("消息中url地址为空，获取数据为：" + string(lessonClip))
		return nil
	}
	//获取数据库详情
	lessonClipDetail, err := l.svcCtx.LessonClipModel.FindDetailByFileOriginUrl(l.ctx, clip.FileOriginUrl)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		l.Error("查询出错：" + err.Error() + string(lessonClip))
		return err
	}
	if lessonClipDetail == nil {
		//将数据增加到数据库
		clip.AddTime = time.Now().Unix()
		_, _ = l.svcCtx.LessonClipModel.Insert(l.ctx, clip)
	}

	l.Info("数据拉取完毕：" + string(lessonClip))
	return nil
}
func (l *LessonClipsLogic) getLessonClipObjectKey(courseId int64, classId int64, originUrlPath string) string {
	return strconv.FormatInt(courseId, 10) + "/" + strconv.FormatInt(classId, 10) + "/" + strings.Replace(strings.TrimPrefix(originUrlPath, "/"), "/", "-", -1)
}
func (l *LessonClipsLogic) getLessonClipFileStatusName(fileStatus int64) string {
	fileStatusMap := make(map[int64]string)
	fileStatusMap[model.FileStatusInit] = "未同步"
	fileStatusMap[model.FileStatusSync] = "同步中"
	fileStatusMap[model.FileStatusSyncComplete] = "已完成"
	if _, ok := fileStatusMap[fileStatus]; ok {
		return fileStatusMap[fileStatus]
	}
	return "未知"
}
func (l *LessonClipsLogic) UploadLessonClipToOss() {
	for {
		clipId := int64(0)
		_, err := l.svcCtx.LessonClipModel.FindNotSyncRow(l.ctx, clipId, 1)
		if errors.Is(err, model.ErrNotFound) {
			_ = l.SendSyncCompleteNotification()
			_ = l.svcCtx.CourseModel.UpdateSyncStatusToComplete(l.ctx)
			continue
		}
		for {
			lessonClip, err := l.svcCtx.LessonClipModel.FindNotSyncRow(l.ctx, clipId, 1)
			if err != nil {
				if errors.Is(err, model.ErrNotFound) {
					break
				}
				time.Sleep(time.Second)
				continue
			}
			clipId = lessonClip.Id
			urlPath, err := l.svcCtx.MsOssModel.ParseUrl(l.ctx, lessonClip.FileOriginUrl)
			l.Info("开始解析url：" + urlPath)
			if err != nil {
				l.Info("解析url出错：" + string(lessonClip.FileOriginUrl))
				time.Sleep(time.Second)
				continue
			}
			objectKey := l.getLessonClipObjectKey(lessonClip.CourseId, lessonClip.ClassId, urlPath)
			err = l.svcCtx.MsOssModel.UploadOSSByUrl(l.ctx, objectKey, lessonClip.FileOriginUrl)
			l.Info("开始上传OSS, 路径的键值：" + objectKey)
			if err != nil {
				l.Error(err.Error() + "上传oss的源数据：" + lessonClip.FileOriginUrl)
				time.Sleep(time.Second)
				continue
			}
			l.Info("文件[" + lessonClip.FileOriginUrl + "]同步完成")
			//更新视频状态为已同步
			err = l.svcCtx.LessonClipModel.SetSyncStatusCompleteByFileOriginUrl(l.ctx, lessonClip.FileOriginUrl, objectKey)
			if err != nil {
				l.Error(err.Error() + "更新的源数据：" + lessonClip.FileOriginUrl)
				time.Sleep(time.Second)
				continue
			}
			l.Info("文件[" + lessonClip.FileOriginUrl + "]同步状态更新完成")
		}
	}
}
func (l *LessonClipsLogic) getParams(courseId int64, classId int64, page int64, pageSize int64) url.Values {
	params := url.Values{}
	params.Set("courseId", strconv.FormatInt(courseId, 10))
	params.Set("classId", strconv.FormatInt(classId, 10))
	params.Set("page", strconv.FormatInt(page, 10))
	params.Set("perpage", strconv.FormatInt(pageSize, 10))
	return params
}
func (l *LessonClipsLogic) SendSyncCompleteNotification() error {
	delayFlag := int64(-1)
	err := l.svcCtx.Cache.GetCtx(l.ctx, constant.ClassinJobUsersCacheDelayKey, &delayFlag)
	if err == nil && delayFlag < utils.GetCurrentTimestamp() {
		err = l.svcCtx.Cache.DelCtx(l.ctx, constant.ClassinJobUsersCacheDelayKey)
		toUser := ""
		_ = l.svcCtx.Cache.GetCtx(l.ctx, constant.ClassinJobUsersCacheKey, &toUser)
		err = l.svcCtx.Cache.DelCtx(l.ctx, constant.ClassinJobUsersCacheKey)
		if toUser != "" {
			//发送通知
			notification := notify.Msg{
				FromUser:    "",
				ToUser:      toUser,
				MsgType:     notify.MarkDownMsgType,
				MsgPlatform: notify.TitanMsgPlatform,
			}
			syncLessonVideoTemplateData := map[string]interface{}{
				"Content": "classin课程同步完成",
				"Date":    utils.GetCurrentTime(),
			}
			notification.Content, _ = l.svcCtx.SNotificationService.GenerateNotification(classin.SyncLessonVideoTemplate, syncLessonVideoTemplateData)
			msg, _ := json.Marshal(notification)
			err = l.svcCtx.RabbitMqProducer.SendClassInMsgToMq(l.svcCtx.Config.RabbitMq.SendMsgTopic, msg)
		}
	}
	return nil
}
