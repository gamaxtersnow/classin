package oldcrm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	msDataDictionaryFieldNames          = builder.RawFieldNames(&MsDataDictionary{})
	msDataDictionaryRows                = strings.Join(msDataDictionaryFieldNames, ",")
	msDataDictionaryRowsExpectAutoSet   = strings.Join(stringx.Remove(msDataDictionaryFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	msDataDictionaryRowsWithPlaceHolder = strings.Join(stringx.Remove(msDataDictionaryFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	msDataDictionaryModel interface {
		Insert(ctx context.Context, data *MsDataDictionary) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MsDataDictionary, error)
		Update(ctx context.Context, data *MsDataDictionary) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMsDataDictionaryModel struct {
		conn  sqlx.SqlConn
		table string
	}

	MsDataDictionary struct {
		Id          int64          `db:"id"`
		ParentId    sql.NullInt64  `db:"parent_id"` // 父级id
		EngName     sql.NullString `db:"eng_name"`  // 英文名称
		Name        sql.NullString `db:"name"`      // 名称
		ValueOne    sql.NullString `db:"value_one"` // 扩展值1
		ValueTwo    sql.NullString `db:"value_two"` // 扩展值2
		IsJapan     sql.NullInt64  `db:"is_japan"`  // 是否是日本
		IsActive    sql.NullInt64  `db:"is_active"` // 是否老的配置只在详情展示用, 新添加时候下拉菜单不显示   1 active 激活 0 未激活
		Status      sql.NullInt64  `db:"status"`    // 状态 1启用  2禁用
		Note        sql.NullString `db:"note"`      // 备注
		Sort        sql.NullInt64  `db:"sort"`      // 排序 从小到大
		AddTime     time.Time      `db:"add_time"`
		UpdateTime  time.Time      `db:"update_time"`
		DeleteTime  sql.NullTime   `db:"delete_time"`
		CreatorId   sql.NullInt64  `db:"creator_id"`   // 创建人ID
		CreatorName sql.NullString `db:"creator_name"` // 创建人姓名
	}
)

func newMsDataDictionaryModel(conn sqlx.SqlConn) *defaultMsDataDictionaryModel {
	return &defaultMsDataDictionaryModel{
		conn:  conn,
		table: "`ms_data_dictionary`",
	}
}

func (m *defaultMsDataDictionaryModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMsDataDictionaryModel) FindOne(ctx context.Context, id int64) (*MsDataDictionary, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", msDataDictionaryRows, m.table)
	var resp MsDataDictionary
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

func (m *defaultMsDataDictionaryModel) Insert(ctx context.Context, data *MsDataDictionary) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, msDataDictionaryRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ParentId, data.EngName, data.Name, data.ValueOne, data.ValueTwo, data.IsJapan, data.IsActive, data.Status, data.Note, data.Sort, data.AddTime, data.DeleteTime, data.CreatorId, data.CreatorName)
	return ret, err
}

func (m *defaultMsDataDictionaryModel) Update(ctx context.Context, data *MsDataDictionary) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, msDataDictionaryRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.ParentId, data.EngName, data.Name, data.ValueOne, data.ValueTwo, data.IsJapan, data.IsActive, data.Status, data.Note, data.Sort, data.AddTime, data.DeleteTime, data.CreatorId, data.CreatorName, data.Id)
	return err
}

func (m *defaultMsDataDictionaryModel) tableName() string {
	return m.table
}
