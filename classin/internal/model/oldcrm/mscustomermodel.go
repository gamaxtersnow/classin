package oldcrm

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MsCustomerModel = (*customMsCustomerModel)(nil)

type (
	// MsCustomerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMsCustomerModel.
	MsCustomerModel interface {
		msCustomerModel
		FindListByIds(ctx context.Context, ids []int64) ([]MsCustomer, error)
		withSession(session sqlx.Session) MsCustomerModel
	}

	customMsCustomerModel struct {
		*defaultMsCustomerModel
	}
)

// NewMsCustomerModel returns a model for the database table.
func NewMsCustomerModel(conn sqlx.SqlConn) MsCustomerModel {
	return &customMsCustomerModel{
		defaultMsCustomerModel: newMsCustomerModel(conn),
	}
}

func (m *customMsCustomerModel) withSession(session sqlx.Session) MsCustomerModel {
	return NewMsCustomerModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultMsCustomerModel) FindListByIds(ctx context.Context, ids []int64) ([]MsCustomer, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var resp []MsCustomer
	query := fmt.Sprintf("select %s from %s where `id` IN (%s)", msCustomerRows, m.table, utils.Int64SliceToCommaSeparatedString(ids))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
