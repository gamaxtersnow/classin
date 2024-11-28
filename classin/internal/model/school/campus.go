package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk/xapi"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"time"
)

var (
	_                     CampusModel = (*conCampusModel)(nil)
	cacheCampusAllKey                 = "cache:school:campus:all"
	cacheCampusRoomAllKey             = "cache:school:campus:room:all"
	cacheCampusExpiration             = 2 * time.Hour
)

type Campus struct {
	Id   int64
	Type int64
	Name string
}

type CampusInfo struct {
	ID       int         // 校区的唯一标识符
	Name     string      // 校区名称
	RoomList []Classroom // 校区中的教室列表
}

type Classroom struct {
	ID   int    // 教室的唯一标识符
	Name string // 教室名称
}
type (
	CampusModel interface {
		GetAllCampuses(ctx context.Context) ([]Campus, error)
		GetCampusRoomList(ctx context.Context) ([]CampusInfo, error)
	}
)
type (
	conCampusModel struct {
		model xapi.CampusModel
		cache cache.Cache
	}
)

func NewCampusModel(client *xiaoxiaosdk.HttpClient, cache cache.Cache) CampusModel {
	return &conCampusModel{
		model: xapi.NewCampusModel(client),
		cache: cache,
	}
}

// GetAllCampuses 获取校区列表
func (c *conCampusModel) GetAllCampuses(ctx context.Context) ([]Campus, error) {
	var campusList []Campus
	if err := c.cache.GetCtx(ctx, cacheCampusAllKey, &campusList); err == nil {
		return campusList, nil
	}
	campuses, err := c.model.GetAllCampuses(ctx)
	if err != nil {
		return campusList, err
	}
	if !campuses.Success {
		return campusList, errors.New(campuses.Messages)
	}
	campus := Campus{}
	for _, ca := range campuses.CampusList {
		campus.Id = ca.Id
		campus.Name = ca.Name
		campus.Type = ca.Type
		campusList = append(campusList, campus)
	}
	_ = c.cache.SetWithExpireCtx(ctx, cacheCampusAllKey, campusList, cacheCampusExpiration)
	return campusList, nil
}

// GetCampusRoomList 获取校区教室列表
func (c *conCampusModel) GetCampusRoomList(ctx context.Context) ([]CampusInfo, error) {
	var campusRoomList []CampusInfo
	if err := c.cache.GetCtx(ctx, cacheCampusRoomAllKey, &campusRoomList); err == nil {
		return campusRoomList, nil
	}
	campusRoomListResponse, err := c.model.GetCampusRoomList(ctx)
	if err != nil {
		return campusRoomList, err
	}
	if !campusRoomListResponse.Success {
		return campusRoomList, errors.New(campusRoomListResponse.Messages)
	}
	for _, cr := range campusRoomListResponse.Data {
		campus := CampusInfo{}
		campus.ID = cr.ID
		campus.Name = cr.Name
		for _, room := range cr.OutClzRoomList {
			campus.RoomList = append(campus.RoomList, Classroom(room))
		}
		campusRoomList = append(campusRoomList, campus)
	}
	_ = c.cache.SetWithExpireCtx(ctx, cacheCampusRoomAllKey, campusRoomList, cacheCampusExpiration)
	return campusRoomList, nil
}
