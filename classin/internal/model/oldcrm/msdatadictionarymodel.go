package oldcrm

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MsDataDictionaryModel = (*customMsDataDictionaryModel)(nil)

type (
	// MsDataDictionaryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMsDataDictionaryModel.
	MsDataDictionaryModel interface {
		msDataDictionaryModel
		FindListByIds(ctx context.Context, ids []int64) ([]MsDataDictionary, error)
		withSession(session sqlx.Session) MsDataDictionaryModel
	}

	customMsDataDictionaryModel struct {
		*defaultMsDataDictionaryModel
	}
)

// NewMsDataDictionaryModel returns a model for the database table.
func NewMsDataDictionaryModel(conn sqlx.SqlConn) MsDataDictionaryModel {
	return &customMsDataDictionaryModel{
		defaultMsDataDictionaryModel: newMsDataDictionaryModel(conn),
	}
}

func (m *customMsDataDictionaryModel) withSession(session sqlx.Session) MsDataDictionaryModel {
	return NewMsDataDictionaryModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultMsDataDictionaryModel) FindListByIds(ctx context.Context, ids []int64) ([]MsDataDictionary, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var resp []MsDataDictionary
	query := fmt.Sprintf("select %s from %s where `id` IN (%s)", msDataDictionaryRows, m.table, utils.Int64SliceToCommaSeparatedString(ids))
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
