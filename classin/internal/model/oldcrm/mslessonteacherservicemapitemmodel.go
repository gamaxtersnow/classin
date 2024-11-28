package oldcrm

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MsLessonTeacherServiceMapItemModel = (*customMsLessonTeacherServiceMapItemModel)(nil)

type (
	// MsLessonTeacherServiceMapItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMsLessonTeacherServiceMapItemModel.
	MsLessonTeacherServiceMapItemModel interface {
		msLessonTeacherServiceMapItemModel
		FindListByTeacherServiceMapIds(ctx context.Context, serviceMapIds []int64) ([]MsLessonTeacherServiceMapItem, error)
		withSession(session sqlx.Session) MsLessonTeacherServiceMapItemModel
	}

	customMsLessonTeacherServiceMapItemModel struct {
		*defaultMsLessonTeacherServiceMapItemModel
	}
)

// NewMsLessonTeacherServiceMapItemModel returns a model for the database table.
func NewMsLessonTeacherServiceMapItemModel(conn sqlx.SqlConn) MsLessonTeacherServiceMapItemModel {
	return &customMsLessonTeacherServiceMapItemModel{
		defaultMsLessonTeacherServiceMapItemModel: newMsLessonTeacherServiceMapItemModel(conn),
	}
}

func (m *customMsLessonTeacherServiceMapItemModel) withSession(session sqlx.Session) MsLessonTeacherServiceMapItemModel {
	return NewMsLessonTeacherServiceMapItemModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultMsLessonTeacherServiceMapItemModel) FindListByTeacherServiceMapIds(ctx context.Context, serviceMapIds []int64) ([]MsLessonTeacherServiceMapItem, error) {
	if len(serviceMapIds) == 0 {
		return nil, nil
	}

	var resp []MsLessonTeacherServiceMapItem
	query := fmt.Sprintf("select %s from %s where `teacher_service_map_id` IN (%s)", msLessonTeacherServiceMapItemRows, m.table, utils.Int64SliceToCommaSeparatedString(serviceMapIds))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
