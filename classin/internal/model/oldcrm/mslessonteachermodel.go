package oldcrm

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MsLessonTeacherModel = (*customMsLessonTeacherModel)(nil)

type (
	// MsLessonTeacherModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMsLessonTeacherModel.
	MsLessonTeacherModel interface {
		msLessonTeacherModel
		FindInfoByPhone(ctx context.Context, phone string) (*MsLessonTeacher, error)
		withSession(session sqlx.Session) MsLessonTeacherModel
	}

	customMsLessonTeacherModel struct {
		*defaultMsLessonTeacherModel
	}
)

// NewMsLessonTeacherModel returns a model for the database table.
func NewMsLessonTeacherModel(conn sqlx.SqlConn) MsLessonTeacherModel {
	return &customMsLessonTeacherModel{
		defaultMsLessonTeacherModel: newMsLessonTeacherModel(conn),
	}
}

func (m *customMsLessonTeacherModel) withSession(session sqlx.Session) MsLessonTeacherModel {
	return NewMsLessonTeacherModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultMsLessonTeacherModel) FindInfoByPhone(ctx context.Context, phone string) (*MsLessonTeacher, error) {
	var resp MsLessonTeacher
	query := fmt.Sprintf("select %s from %s where `phone` = ?", msLessonTeacherRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
