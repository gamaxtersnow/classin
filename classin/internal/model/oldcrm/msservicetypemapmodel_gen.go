package oldcrm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	msServicetypeMapFieldNames          = builder.RawFieldNames(&MsServicetypeMap{})
	msServicetypeMapRows                = strings.Join(msServicetypeMapFieldNames, ",")
	msServicetypeMapRowsExpectAutoSet   = strings.Join(stringx.Remove(msServicetypeMapFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	msServicetypeMapRowsWithPlaceHolder = strings.Join(stringx.Remove(msServicetypeMapFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	msServicetypeMapModel interface {
		Insert(ctx context.Context, data *MsServicetypeMap) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MsServicetypeMap, error)
		Update(ctx context.Context, data *MsServicetypeMap) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMsServicetypeMapModel struct {
		conn  sqlx.SqlConn
		table string
	}

	MsServicetypeMap struct {
		Id          int64          `db:"id"`
		ServiceType sql.NullInt64  `db:"service_type"` // 类型id 1留学 2 语培 ... 9 移民
		Pid         sql.NullInt64  `db:"pid"`          // 父节点id 对应data_distionary里的id
		Name        sql.NullString `db:"name"`
	}
)

func newMsServicetypeMapModel(conn sqlx.SqlConn) *defaultMsServicetypeMapModel {
	return &defaultMsServicetypeMapModel{
		conn:  conn,
		table: "`ms_servicetype_map`",
	}
}

func (m *defaultMsServicetypeMapModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMsServicetypeMapModel) FindOne(ctx context.Context, id int64) (*MsServicetypeMap, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", msServicetypeMapRows, m.table)
	var resp MsServicetypeMap
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMsServicetypeMapModel) Insert(ctx context.Context, data *MsServicetypeMap) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, msServicetypeMapRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ServiceType, data.Pid, data.Name)
	return ret, err
}

func (m *defaultMsServicetypeMapModel) Update(ctx context.Context, data *MsServicetypeMap) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, msServicetypeMapRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.ServiceType, data.Pid, data.Name, data.Id)
	return err
}

func (m *defaultMsServicetypeMapModel) tableName() string {
	return m.table
}
