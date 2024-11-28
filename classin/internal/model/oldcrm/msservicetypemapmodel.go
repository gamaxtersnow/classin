package oldcrm

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MsServicetypeMapModel = (*customMsServicetypeMapModel)(nil)

type (
	// MsServicetypeMapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMsServicetypeMapModel.
	MsServicetypeMapModel interface {
		msServicetypeMapModel
		FindListByIds(ctx context.Context, ids []int64) ([]MsServicetypeMap, error)
		withSession(session sqlx.Session) MsServicetypeMapModel
	}

	customMsServicetypeMapModel struct {
		*defaultMsServicetypeMapModel
	}
)

// NewMsServicetypeMapModel returns a model for the database table.
func NewMsServicetypeMapModel(conn sqlx.SqlConn) MsServicetypeMapModel {
	return &customMsServicetypeMapModel{
		defaultMsServicetypeMapModel: newMsServicetypeMapModel(conn),
	}
}

func (m *customMsServicetypeMapModel) withSession(session sqlx.Session) MsServicetypeMapModel {
	return NewMsServicetypeMapModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultMsServicetypeMapModel) FindListByIds(ctx context.Context, ids []int64) ([]MsServicetypeMap, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var resp []MsServicetypeMap
	query := fmt.Sprintf("select %s from %s where `id` IN (%s)", msServicetypeMapRows, m.table, utils.Int64SliceToCommaSeparatedString(ids))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
