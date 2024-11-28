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
	msLessonTeacherServiceMapItemFieldNames          = builder.RawFieldNames(&MsLessonTeacherServiceMapItem{})
	msLessonTeacherServiceMapItemRows                = strings.Join(msLessonTeacherServiceMapItemFieldNames, ",")
	msLessonTeacherServiceMapItemRowsExpectAutoSet   = strings.Join(stringx.Remove(msLessonTeacherServiceMapItemFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	msLessonTeacherServiceMapItemRowsWithPlaceHolder = strings.Join(stringx.Remove(msLessonTeacherServiceMapItemFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	msLessonTeacherServiceMapItemModel interface {
		Insert(ctx context.Context, data *MsLessonTeacherServiceMapItem) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MsLessonTeacherServiceMapItem, error)
		Update(ctx context.Context, data *MsLessonTeacherServiceMapItem) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMsLessonTeacherServiceMapItemModel struct {
		conn  sqlx.SqlConn
		table string
	}

	MsLessonTeacherServiceMapItem struct {
		Id                  int64           `db:"id"`
		LessonTeacherId     sql.NullInt64   `db:"lesson_teacher_id"`      // 对应lesson_teacher表里的id
		TeacherServiceMapId sql.NullInt64   `db:"teacher_service_map_id"` // 对应teacher_service_map的id
		ServiceType         sql.NullInt64   `db:"service_type"`
		ServiceItemId       sql.NullInt64   `db:"service_item_id"`
		ServiceItemChildId  sql.NullInt64   `db:"service_item_child_id"`
		Price               sql.NullFloat64 `db:"price"`           // 按课时收费
		PriceType           sql.NullString  `db:"price_type"`      // 1 按课时 2 按项目一笔 3 按项目两笔 4 按项目三笔
		PriceScale1         sql.NullInt64   `db:"price_scale1"`    // 按项目收费该字段 第一笔
		PriceScale2         sql.NullInt64   `db:"price_scale2"`    // 按项目收费该字段 第二笔
		PriceScale3         sql.NullInt64   `db:"price_scale3"`    // 按项目收费该字段 第三笔
		SettlementNode      sql.NullInt64   `db:"settlement_node"` // 1 项目开启 2 项目结束
		CreateName          sql.NullString  `db:"create_name"`     // 创建人
		AddTime             sql.NullTime    `db:"add_time"`
		UpdatedTime         time.Time       `db:"updated_time"`
	}
)

func newMsLessonTeacherServiceMapItemModel(conn sqlx.SqlConn) *defaultMsLessonTeacherServiceMapItemModel {
	return &defaultMsLessonTeacherServiceMapItemModel{
		conn:  conn,
		table: "`ms_lesson_teacher_service_map_item`",
	}
}

func (m *defaultMsLessonTeacherServiceMapItemModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMsLessonTeacherServiceMapItemModel) FindOne(ctx context.Context, id int64) (*MsLessonTeacherServiceMapItem, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", msLessonTeacherServiceMapItemRows, m.table)
	var resp MsLessonTeacherServiceMapItem
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

func (m *defaultMsLessonTeacherServiceMapItemModel) Insert(ctx context.Context, data *MsLessonTeacherServiceMapItem) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, msLessonTeacherServiceMapItemRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.LessonTeacherId, data.TeacherServiceMapId, data.ServiceType, data.ServiceItemId, data.ServiceItemChildId, data.Price, data.PriceType, data.PriceScale1, data.PriceScale2, data.PriceScale3, data.SettlementNode, data.CreateName, data.AddTime, data.UpdatedTime)
	return ret, err
}

func (m *defaultMsLessonTeacherServiceMapItemModel) Update(ctx context.Context, data *MsLessonTeacherServiceMapItem) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, msLessonTeacherServiceMapItemRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.LessonTeacherId, data.TeacherServiceMapId, data.ServiceType, data.ServiceItemId, data.ServiceItemChildId, data.Price, data.PriceType, data.PriceScale1, data.PriceScale2, data.PriceScale3, data.SettlementNode, data.CreateName, data.AddTime, data.UpdatedTime, data.Id)
	return err
}

func (m *defaultMsLessonTeacherServiceMapItemModel) tableName() string {
	return m.table
}
