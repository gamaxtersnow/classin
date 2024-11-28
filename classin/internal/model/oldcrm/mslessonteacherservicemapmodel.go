package oldcrm

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MsLessonTeacherServiceMapModel = (*customMsLessonTeacherServiceMapModel)(nil)

type (
	// MsLessonTeacherServiceMapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMsLessonTeacherServiceMapModel.
	MsLessonTeacherServiceMapModel interface {
		msLessonTeacherServiceMapModel
		FindListByTeacherId(ctx context.Context, teacherId int64) ([]MsLessonTeacherServiceMap, error)
		withSession(session sqlx.Session) MsLessonTeacherServiceMapModel
	}

	customMsLessonTeacherServiceMapModel struct {
		*defaultMsLessonTeacherServiceMapModel
	}
)

// NewMsLessonTeacherServiceMapModel returns a model for the database table.
func NewMsLessonTeacherServiceMapModel(conn sqlx.SqlConn) MsLessonTeacherServiceMapModel {
	return &customMsLessonTeacherServiceMapModel{
		defaultMsLessonTeacherServiceMapModel: newMsLessonTeacherServiceMapModel(conn),
	}
}

func (m *customMsLessonTeacherServiceMapModel) withSession(session sqlx.Session) MsLessonTeacherServiceMapModel {
	return NewMsLessonTeacherServiceMapModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultMsLessonTeacherServiceMapModel) FindListByTeacherId(ctx context.Context, teacherId int64) ([]MsLessonTeacherServiceMap, error) {
	var resp []MsLessonTeacherServiceMap
	query := fmt.Sprintf("select %s from %s where `lesson_teacher_id` = ? AND `status`='enable' AND is_deleted=0", msLessonTeacherServiceMapRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, teacherId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
