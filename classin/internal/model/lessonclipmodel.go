package model

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/net/context"
)

var _ LessonClipModel = (*customLessonClipModel)(nil)

const FileStatusInit int64 = 0         //未同步
const FileStatusSync int64 = 1         //同步中
const FileStatusSyncComplete int64 = 2 //已完成
const FileSourceTypeManual int64 = 1   //手动抓取
const FileSourceTypePush int64 = 2     //系统推送

type (
	// LessonClipModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonClipModel.
	LessonClipModel interface {
		lessonClipModel
		CountByCourseIdAndClassId(ctx context.Context, courseId int64, classId int64) (int64, error)
		FindByCourseIdAndClassId(ctx context.Context, courseId int64, classId int64, page int, offset int) ([]LessonClip, error)
		FindCountByFileOriginUrl(ctx context.Context, fileOriginUrl string) (int64, error)
		FindDetailByFileOriginUrl(ctx context.Context, fileOriginUrl string) (*LessonClip, error)
		SetSyncStatusCompleteByFileOriginUrl(ctx context.Context, fileOriginUrl string, fileUrl string) error
		FindNotSyncRow(ctx context.Context, id int64, limit int) (*LessonClip, error)
	}

	customLessonClipModel struct {
		*defaultLessonClipModel
	}
)

// NewLessonClipModel returns a model for the database table.
func NewLessonClipModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LessonClipModel {
	return &customLessonClipModel{
		defaultLessonClipModel: newLessonClipModel(conn, c, opts...),
	}
}
func (m *defaultLessonClipModel) FindByCourseIdAndClassId(ctx context.Context, courseId int64, classId int64, page int, offset int) ([]LessonClip, error) {
	var lessonClips []LessonClip
	query := fmt.Sprintf("select %s from %s where `courseId` = ? and `classId` = ? limit ?,?", lessonClipRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &lessonClips, query, courseId, classId, page, offset)
	switch {
	case err == nil:
		return lessonClips, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLessonClipModel) CountByCourseIdAndClassId(ctx context.Context, courseId int64, classId int64) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where `courseId` = ? and `classId` = ?", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, courseId, classId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (m *defaultLessonClipModel) FindCountByFileOriginUrl(ctx context.Context, fileOriginUrl string) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where `fileOriginUrl` = ?", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, fileOriginUrl)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (m *defaultLessonClipModel) FindDetailByFileOriginUrl(ctx context.Context, fileOriginUrl string) (*LessonClip, error) {
	var lessonClip LessonClip
	query := fmt.Sprintf("select %s from %s where `fileOriginUrl` = ?", lessonClipRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &lessonClip, query, fileOriginUrl)
	switch {
	case err == nil:
		return &lessonClip, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultLessonClipModel) SetSyncStatusCompleteByFileOriginUrl(ctx context.Context, fileOriginUrl string, fileUrl string) error {
	query := fmt.Sprintf("update %s set `fileStatus`= ?,`objectKey`=? where `fileOriginUrl` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, FileStatusSyncComplete, fileUrl, fileOriginUrl)
	if err != nil {
		return err
	}
	return nil
}
func (m *defaultLessonClipModel) FindNotSyncRow(ctx context.Context, id int64, limit int) (*LessonClip, error) {
	var lessonClip LessonClip
	query := fmt.Sprintf("select %s from %s where id> ?  and `fileStatus`= ?  order by id asc limit ?", lessonClipRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &lessonClip, query, id, FileStatusInit, limit)
	switch {
	case err == nil:
		return &lessonClip, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
