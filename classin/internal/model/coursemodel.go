package model

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseModel = (*customCourseModel)(nil)

const SyncStatusInit int64 = 0     //未同步
const SyncStatusSync int64 = 1     //同步中
const SyncStatusComplete int64 = 2 //已完成
const SourceTypeManual int64 = 1   //手动抓取
const SourceTypePush int64 = 2     //手动抓取
type (
	// CourseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseModel.
	CourseModel interface {
		courseModel
		FindByUid(ctx context.Context, uid string) (*Course, error)
		CountLesson(ctx context.Context, params map[string]interface{}) (int64, error)
		GetLessonList(ctx context.Context, params map[string]interface{}, offset int, pageSize int) ([]Course, error)
		SetSyncStatusSyncByUniqueId(ctx context.Context, uniqueId string) error
		SetSyncStatusCompleteByUniqueId(ctx context.Context, uniqueId string) error
		FindNotSyncRow(ctx context.Context, id int64, limit int) (*Course, error)
		UpdateSyncStatusToComplete(ctx context.Context) error
	}

	customCourseModel struct {
		*defaultCourseModel
	}
)

// NewCourseModel returns a model for the database table.
func NewCourseModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseModel {
	return &customCourseModel{
		defaultCourseModel: newCourseModel(conn, c, opts...),
	}
}

func (c *defaultCourseModel) FindByUid(ctx context.Context, uid string) (*Course, error) {
	var resp Course
	query := fmt.Sprintf("select %s from %s where `uniqueId` = ? limit 1", courseRows, c.table)
	err := c.QueryRowNoCacheCtx(ctx, &resp, query, uid)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *defaultCourseModel) CountLesson(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	whereClause, args, _ := utils.BuildQuery(c.table, params)
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", c.table, whereClause)
	logx.Info(query, args)
	err := c.QueryRowNoCacheCtx(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *defaultCourseModel) GetLessonList(ctx context.Context, params map[string]interface{}, offset int, pageSize int) ([]Course, error) {
	var resp []Course
	whereClause, args, _ := utils.BuildQuery(c.table, params)
	query := fmt.Sprintf("SELECT %s FROM %s %s LIMIT ? , ?", courseRows, c.table, whereClause)
	args = append(args, offset)
	args = append(args, pageSize)
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (c *defaultCourseModel) SetSyncStatusCompleteByUniqueId(ctx context.Context, uniqueId string) error {
	query := fmt.Sprintf("update %s set `syncStatus`= ? where `uniqueId` = ?", c.table)
	_, err := c.ExecNoCacheCtx(ctx, query, SyncStatusComplete, uniqueId)
	if err != nil {
		return err
	}
	return nil
}
func (c *defaultCourseModel) SetSyncStatusSyncByUniqueId(ctx context.Context, uniqueId string) error {
	query := fmt.Sprintf("update %s set `syncStatus`= ? where `uniqueId` = ?", c.table)
	_, err := c.ExecNoCacheCtx(ctx, query, SyncStatusSync, uniqueId)
	if err != nil {
		return err
	}
	return nil
}
func (c *defaultCourseModel) UpdateSyncStatusToComplete(ctx context.Context) error {
	query := fmt.Sprintf("update %s set `syncStatus`= ? where `syncStatus` <> 2", c.table)
	_, err := c.ExecNoCacheCtx(ctx, query, SyncStatusComplete)
	if err != nil {
		return err
	}
	return nil
}

func (c *defaultCourseModel) FindNotSyncRow(ctx context.Context, id int64, limit int) (*Course, error) {
	var resp Course
	query := fmt.Sprintf("select %s from %s where id > ? and `syncStatus` = ?  order by id asc limit ?", courseRows, c.table)
	err := c.QueryRowNoCacheCtx(ctx, &resp, query, id, SyncStatusInit, limit)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
