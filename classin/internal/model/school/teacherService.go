package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"meishiedu.com/classin/internal/model/oldcrm"
	"meishiedu.com/classin/internal/types"
	"strconv"
)

type TeacherService struct {
}

var _ TeacherServiceModel = (*customTeacherServiceModel)(nil)

type (
	TeacherServiceModel interface {
		GetTeacherServicePriceByPhone(ctx context.Context, phone string) (*TeacherService, error)
		GetTeacherServicePriceByPhones(ctx context.Context, phones []string) ([]*TeacherService, error)
	}
	customTeacherServiceModel struct {
		oldcrmLessonTeacherModel               oldcrm.MsLessonTeacherModel
		oldcrmLessonTeacherServiceMapModel     oldcrm.MsLessonTeacherServiceMapModel
		oldcrmLessonTeacherServiceMapItemModel oldcrm.MsLessonTeacherServiceMapItemModel
		oldcrmDataDictionaryModel              oldcrm.MsDataDictionaryModel
		oldcrmServicetypeMapModel              oldcrm.MsServicetypeMapModel
		cache                                  cache.Cache
	}
)

func NewCustomTeacherServiceModel(conn sqlx.SqlConn, cache cache.Cache) TeacherServiceModel {
	return &customTeacherServiceModel{
		oldcrmLessonTeacherModel:               oldcrm.NewMsLessonTeacherModel(conn),
		oldcrmLessonTeacherServiceMapModel:     oldcrm.NewMsLessonTeacherServiceMapModel(conn),
		oldcrmLessonTeacherServiceMapItemModel: oldcrm.NewMsLessonTeacherServiceMapItemModel(conn),
		oldcrmDataDictionaryModel:              oldcrm.NewMsDataDictionaryModel(conn),
		oldcrmServicetypeMapModel:              oldcrm.NewMsServicetypeMapModel(conn),
		cache:                                  cache,
	}
}

// GetTeachers 获取老师列表

func (t *customTeacherServiceModel) GetTeacherServicePriceByPhone(ctx context.Context, phone string) (*TeacherService, error) {
	//查询 ms_lesson_teacher 的数据
	teacherInfo, err := t.oldcrmLessonTeacherModel.FindInfoByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	//查询 ms_lesson_teacher_service_map 的数据
	teacherServiceMapList, err := t.oldcrmLessonTeacherServiceMapModel.FindListByTeacherId(ctx, teacherInfo.Id)
	if err != nil {
		return nil, err
	}

	teacherServiceMapListData := make(map[int64]map[string]string, 0)
	for _, serviceMap := range teacherServiceMapList {
		teacherServiceMapListData[serviceMap.Id] = map[string]string{
			"id":                    fmt.Sprintf("%d", serviceMap.Id),
			"cooperation_type_id":   strconv.Itoa(int(serviceMap.CoopLavel.Int64)),
			"cooperation_type_name": oldcrm.TeacherCoopLavelMap[int(serviceMap.CoopLavel.Int64)],
			"service_type_id":       strconv.Itoa(int(serviceMap.ServiceType.Int64)),
			"service_type_name":     oldcrm.ServiceTypeMap[int(serviceMap.ServiceType.Int64)],
		}
	}
	//查询 ms_lesson_teacher_service_map_item 的数据
	var serviceMapIds []int64
	for _, serviceMap := range teacherServiceMapList {
		serviceMapIds = append(serviceMapIds, serviceMap.Id)
	}
	teacherServiceMapItemList, err := t.oldcrmLessonTeacherServiceMapItemModel.FindListByTeacherServiceMapIds(ctx, serviceMapIds)
	if err != nil {
		return nil, err
	}

	//查询 ms_data_dictionary 的数据
	var ids []int64
	for _, item := range teacherServiceMapItemList {
		ids = append(ids, item.ServiceItemId.Int64)
	}
	dataDictionaryList, err := t.oldcrmDataDictionaryModel.FindListByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	dataDictionaryListData := make(map[int64]string, 0)
	for _, dataDictionary := range dataDictionaryList {
		dataDictionaryListData[dataDictionary.Id] = dataDictionary.Name.String
	}

	//查询 ms_servicetype_map 的数据
	ids = []int64{}
	for _, item := range teacherServiceMapItemList {
		ids = append(ids, item.ServiceItemChildId.Int64)
	}
	servicetypeMapList, err := t.oldcrmServicetypeMapModel.FindListByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	servicetypeMapListData := make(map[int64]string, 0)
	for _, servicetypeMap := range servicetypeMapList {
		servicetypeMapListData[servicetypeMap.Id] = servicetypeMap.Name.String
	}

	//组合数据结果
	var respData []types.TeacherServiceMapData
	for _, serviceMapItem := range teacherServiceMapItemList {
		oneTeacherServiceMapId := serviceMapItem.TeacherServiceMapId.Int64
		respData = append(respData, types.TeacherServiceMapData{
			TeacherId:           int(teacherInfo.Id),
			TeacherName:         teacherInfo.Name.String,
			TeacherPhone:        phone,
			CooperationTypeId:   utils.StringToInt(teacherServiceMapListData[oneTeacherServiceMapId]["cooperation_type_id"]),
			CooperationTypeName: teacherServiceMapListData[oneTeacherServiceMapId]["cooperation_type_name"],
			ServiceTypeId:       utils.StringToInt(teacherServiceMapListData[oneTeacherServiceMapId]["service_type_id"]),
			ServiceTypeName:     teacherServiceMapListData[oneTeacherServiceMapId]["service_type_name"],
			CategoryId:          int(serviceMapItem.ServiceItemId.Int64),
			CategoryName:        dataDictionaryListData[serviceMapItem.ServiceItemId.Int64],
			ItemId:              int(serviceMapItem.ServiceItemChildId.Int64),
			ItemName:            servicetypeMapListData[serviceMapItem.ServiceItemChildId.Int64],
			Price:               serviceMapItem.Price.Float64,
		})
	}

	return nil, nil
}
func (t *customTeacherServiceModel) GetTeacherServicePriceByPhones(ctx context.Context, phones []string) ([]*TeacherService, error) {
	return nil, nil
}
