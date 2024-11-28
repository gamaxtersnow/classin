package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/notify"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/notify/template/schedule"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"meishiedu.com/classin/internal/model/school"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
	"strings"
)

type ScheduleExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewScheduleExportLogic 课表导出
func NewScheduleExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleExportLogic {
	return &ScheduleExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScheduleExportLogic) ScheduleExport(req *types.ScheduleExportRequest) (resp *types.ScheduleExportResponse, err error) {
	// 1.  获取老师数据总数
	params := &school.ScheduleViewCond{}
	params.Search = req.Search
	params.Status = req.Status
	params.CampusIDs = req.CampusIDs
	params.ClassIDs = req.ClassIDs
	params.TeacherIDs = req.TeacherIds
	params.CourseIDs = req.CourseIDs
	params.ZhuJiaoID = req.ZhuJiaoID
	params.ClzRooms = req.ClzRooms
	params.Dfrom = req.Dfrom
	params.Dto = req.Dto
	params.PageID = req.PageID
	params.PageSize = req.PageSize
	params.ViewType = req.ViewType
	params.ShowCancel = req.ShowCancel
	params.ExceptNull = req.ExceptNull
	params.Asc = req.Asc
	if req.Dlall == 1 {
		params.Dlall = req.Dlall
		params.PageID = 0
		params.PageSize = 0
	}
	allData, err := l.svcCtx.SScheduleModel.GetScheduleList(l.ctx, params)
	if err != nil {
		return nil, err
	}
	if allData.Total == 0 {
		return nil, errors.New("获取远程课表暂无数据")
	}
	logx.Infof("获取远程课表数据总数：%d", allData.Total)
	// 2. 获取所有数据
	allTeacherData := make([]school.Schedule, 0, allData.Total)

	// 获取校区名字列表映射
	campusMap, err := l.getCampusNames()
	if err != nil {
		return nil, err
	}
	// 处理校区名称
	for i := range allData.Data {
		if campus, exists := campusMap[allData.Data[i].CampusID]; exists {
			allData.Data[i].CampusNames = campus.Name
		}
		allTeacherData = append(allTeacherData, allData.Data[i])
	}
	if len(allTeacherData) == 0 {
		return nil, errors.New("获取远程课表数据为空")
	}
	f := excelize.NewFile()
	defer func(f *excelize.File) {
		_ = f.Close()
	}(f)
	defaultSheet := f.GetSheetName(0)
	_ = f.SetSheetName(defaultSheet, "")
	f.SetActiveSheet(0)
	sheetName := f.GetSheetName(0)
	//循环获取老师兼职数据-待提供crm新接口
	// 写入数据
	if err := l.writeExcelData(f, sheetName, allTeacherData); err != nil {
		return nil, err
	}
	//上传文件到oss
	ossUrl, err := l.uploadFile(f)
	if err != nil {
		return nil, err
	}
	resp = &types.ScheduleExportResponse{
		ErrorInfo: types.ErrorInfo{
			ErrorCode: 1,
			ErrorMsg:  "success",
		},
		Pages: types.Pages{
			Total:    allData.Total,
			PageID:   req.PageID,
			PageSize: req.PageSize,
		},
		Url: ossUrl.Url,
	}
	return resp, nil
}
func (l *ScheduleExportLogic) GenerateObjectKey() (uuid string, name string, objectKey string, fileType string) {
	uuid = utils.GetUUID()
	name = "日签"
	objectKey = "upload/schedule/" + name + "_" + uuid + ".xlsx"
	fileType = "xlsx"
	return
}

// 获取助教姓名列表
func (l *ScheduleExportLogic) getZjNames(clzZj []school.Person) string {
	if len(clzZj) == 0 {
		return ""
	}
	names := make([]string, 0, len(clzZj))
	for _, zj := range clzZj {
		if zj.Name != "" {
			names = append(names, zj.Name)
		}
	}
	return strings.Join(names, ",")
}
func (l *ScheduleExportLogic) setHeaderStyle(f *excelize.File, sheetName string) error {
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#CCCCCC"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}
	return f.SetRowStyle(sheetName, 1, 1, style)
}

func (l *ScheduleExportLogic) getCampusNames() (map[int64]school.Campus, error) {
	//定义一个map，key为id，value为校区详情
	campusMap := make(map[int64]school.Campus)
	//获取校区名称列表
	campusList, err := l.svcCtx.SCampusModel.GetAllCampuses(l.ctx)
	if err != nil {
		return nil, err
	}
	for _, campusDetail := range campusList {
		campusMap[campusDetail.Id] = campusDetail
	}
	return campusMap, nil
}
func (l *ScheduleExportLogic) uploadFile(f *excelize.File) (*types.ScheduleExportResponse, error) {
	uuid, name, objectKey, fileType := l.GenerateObjectKey()
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.MsOssModel.UploadOssByFile(l.ctx, objectKey, buffer.Bytes()); err != nil {
		l.Error(err.Error() + "上传oss的导出数据数据：" + objectKey)
		return nil, err
	}
	objectUrl, err := l.svcCtx.MsOssModel.GetSignUrl(l.ctx, objectKey, 300)
	if err != nil {
		l.Error(err.Error() + "根据文件对象key,获取oss文件地址：" + objectKey)
		return nil, err
	}
	file := &school.Files{
		Uuid:      uuid,
		ObjectKey: objectKey,
		Name:      name,
		FileType:  fileType,
		AddTime:   utils.GetCurrentTimestamp(),
		Remark:    "日签文件导出",
		Status:    school.FileStatusNormal,
	}
	_, err = l.svcCtx.SFilesModel.Insert(l.ctx, file)
	if err != nil {
		return nil, fmt.Errorf("新增文件失败: %v", err)
	}
	_ = l.sendExportNotification(uuid, name)
	return &types.ScheduleExportResponse{
		ErrorInfo: types.ErrorInfo{ErrorCode: 1},
		Url:       objectUrl, //这个用于现在可以直接给一个临时地址objectUrl
	}, nil
}
func (l *ScheduleExportLogic) writeExcelData(f *excelize.File, sheetName string, data []school.Schedule) error {
	// 设置表头
	headers := []string{"校区名称", "班级名称", "上课状态", "课程类别", "上课日期", "开始时间", "结束时间", "上课时长", "上课内容", "上课老师", "老师合作类型", "上课方式", "上课地点", "出勤率", "课程类型", "本班助教", "备注"}
	for i, header := range headers {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%c1", 'A'+i), header)
	}
	// 写入数据
	for i, item := range data {
		row := i + 2 // 从第二行开始写入数据
		values := []interface{}{
			item.CampusNames,
			item.ClzName,
			utils.GetStatusText(item.Status),
			item.CourseCategory,
		}
		// 时间格式化
		startDate, startTime := utils.GetTimeText(item.StartTime)
		values = append(values, startDate, startTime)
		_, endTime := utils.GetTimeText(item.EndTime)
		values = append(values, endTime, item.Duration, item.Content, item.Teacher.Name, item.CoopLevel,
			utils.GetWayText(item.Way), item.Place,
			utils.GetAttendanceText(item.SCountJoin, item.SCountClz),
			utils.GetCourseTypeText(item.CourseType),
			l.getZjNames(item.ClzZhujiao),
			item.Note,
		)
		if err := l.setRowValuesAndAdjustWidth(f, sheetName, row, values); err != nil {
			return err
		}
	}
	_ = l.setHeaderStyle(f, sheetName)
	return nil
}
func (l *ScheduleExportLogic) setRowValuesAndAdjustWidth(f *excelize.File, sheetName string, row int, values []interface{}) error {
	for j, value := range values {
		cell := fmt.Sprintf("%c%d", 'A'+j, row)
		if err := f.SetCellValue(sheetName, cell, value); err != nil {
			return fmt.Errorf("failed to set value for cell %s: %w", cell, err)
		}

		// 调整列宽
		currentWidth, _ := f.GetColWidth(sheetName, fmt.Sprintf("%c", 'A'+j)) // 获取当前列宽
		newWidth := float64(len(fmt.Sprintf("%v", value))) + 2                // 计算新宽度，增加一些额外空间
		if newWidth > currentWidth {
			_ = f.SetColWidth(sheetName, fmt.Sprintf("%c", 'A'+j), fmt.Sprintf("%c", 'A'+j), newWidth) // 设置新的列宽
		}
	}
	return nil
}

// sendExportNotification 发送导出通知到消息队列
func (l *ScheduleExportLogic) sendExportNotification(uuid string, name string) error {
	//发送通知
	notification := notify.Msg{
		FromUser:    "",
		ToUser:      l.ctx.Value("wid").(string),
		MsgType:     notify.MarkDownMsgType,
		MsgPlatform: notify.TitanMsgPlatform,
	}
	exportScheduleExcelTemplateData := map[string]interface{}{
		"System":      l.svcCtx.Config.Name,
		"IP":          "",
		"Address":     "",
		"Date":        utils.GetCurrentTime(),
		"FileName":    name,
		"FileAddress": l.svcCtx.Config.ResourceCenterDomain + "/#/ssoLogin?fileKey=" + uuid,
	}
	notification.Content, _ = l.svcCtx.SNotificationService.GenerateNotification(schedule.ExportScheduleExcelTemplate, exportScheduleExcelTemplateData)
	msg, _ := json.Marshal(notification)
	_ = l.svcCtx.RabbitMqProducer.SendClassInMsgToMq(l.svcCtx.Config.RabbitMq.SendMsgTopic, msg)
	return nil
}
