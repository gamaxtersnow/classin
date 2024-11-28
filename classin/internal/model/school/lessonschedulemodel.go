package school

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LessonScheduleModel = (*customLessonScheduleModel)(nil)

type (
	// LessonScheduleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonScheduleModel.
	LessonScheduleModel interface {
		lessonScheduleModel
	}

	customLessonScheduleModel struct {
		*defaultLessonScheduleModel
	}
)

// NewLessonScheduleModel returns a model for the database table.
func NewLessonScheduleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LessonScheduleModel {
	return &customLessonScheduleModel{
		defaultLessonScheduleModel: newLessonScheduleModel(conn, c, opts...),
	}
}
