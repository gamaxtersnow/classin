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
	msLessonTeacherServiceMapFieldNames          = builder.RawFieldNames(&MsLessonTeacherServiceMap{})
	msLessonTeacherServiceMapRows                = strings.Join(msLessonTeacherServiceMapFieldNames, ",")
	msLessonTeacherServiceMapRowsExpectAutoSet   = strings.Join(stringx.Remove(msLessonTeacherServiceMapFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	msLessonTeacherServiceMapRowsWithPlaceHolder = strings.Join(stringx.Remove(msLessonTeacherServiceMapFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	msLessonTeacherServiceMapModel interface {
		Insert(ctx context.Context, data *MsLessonTeacherServiceMap) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MsLessonTeacherServiceMap, error)
		Update(ctx context.Context, data *MsLessonTeacherServiceMap) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMsLessonTeacherServiceMapModel struct {
		conn  sqlx.SqlConn
		table string
	}

	MsLessonTeacherServiceMap struct {
		Id              int64          `db:"id"`
		LessonTeacherId sql.NullInt64  `db:"lesson_teacher_id"` // 老师ID
		ServiceType     sql.NullInt64  `db:"service_type"`      // 服务类型ID
		ClassModality   sql.NullString `db:"class_modality"`    // 上课形式1线上 2线下
		ClassType       sql.NullString `db:"class_type"`        // 课程类型 1.1V1  2.1V2  3.1V3  4.1V4   5.1V6  6.1V8  7.1V10  8.1V12
		Jingyingfuwu    sql.NullString `db:"jingyingfuwu"`      // 关联菁英服务 1 关联 0 不关联
		Lavel           sql.NullString `db:"lavel"`             // 导师等级 vip1 vip2 vip3
		ServiceGooddesc sql.NullString `db:"service_gooddesc"`  // 服务亮点
		SignDate        sql.NullTime   `db:"sign_date"`         // 签约日期
		EndDate         sql.NullTime   `db:"end_date"`          // 到期日期
		CoopType        sql.NullInt64  `db:"coop_type"`         //  合作形式 1 个人 2 机构
		CoopLavel       sql.NullInt64  `db:"coop_lavel"`        // 合作性质 1 全职 2 兼职
		Company         sql.NullString `db:"company"`           // 任职机构
		Year            sql.NullInt64  `db:"year"`              // 服务年限
		Exp             sql.NullString `db:"exp"`               // 经验
		Chengguo        sql.NullString `db:"chengguo"`          // 成果
		Lingyudalei     sql.NullString `db:"lingyudalei"`       // 领域大类
		Yanjiulingyu    sql.NullString `db:"yanjiulingyu"`      // 研究领域
		Shanchang       sql.NullString `db:"shanchang"`         // 擅长
		BillFile        sql.NullString `db:"bill_file"`         // 海报宣传
		ContractFile    sql.NullString `db:"contract_file"`     // 合同文件
		CertFile        sql.NullString `db:"cert_file"`         // 行业证书
		CreatedAt       time.Time      `db:"created_at"`
		UpdatedAt       time.Time      `db:"updated_at"`
		JuralCount      sql.NullInt64  `db:"jural_count"`  // 义务课时数
		PartnerId       sql.NullInt64  `db:"partner_id"`   // 所属合作方 关联 ms_late_product_partner表id
		Status          sql.NullString `db:"status"`       // 合同状态值
		OperateType     sql.NullInt64  `db:"operate_type"` // 操作类型:0添加，1续签，2追加服务业务，3修改
		IsDeleted       sql.NullInt64  `db:"is_deleted"`   // 是否删除
		RelateId        sql.NullInt64  `db:"relate_id"`    // 关联的id
		RelateIds       sql.NullString `db:"relate_ids"`   // 关联的多个id
	}
)

func newMsLessonTeacherServiceMapModel(conn sqlx.SqlConn) *defaultMsLessonTeacherServiceMapModel {
	return &defaultMsLessonTeacherServiceMapModel{
		conn:  conn,
		table: "`ms_lesson_teacher_service_map`",
	}
}

func (m *defaultMsLessonTeacherServiceMapModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMsLessonTeacherServiceMapModel) FindOne(ctx context.Context, id int64) (*MsLessonTeacherServiceMap, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", msLessonTeacherServiceMapRows, m.table)
	var resp MsLessonTeacherServiceMap
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

func (m *defaultMsLessonTeacherServiceMapModel) Insert(ctx context.Context, data *MsLessonTeacherServiceMap) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, msLessonTeacherServiceMapRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.LessonTeacherId, data.ServiceType, data.ClassModality, data.ClassType, data.Jingyingfuwu, data.Lavel, data.ServiceGooddesc, data.SignDate, data.EndDate, data.CoopType, data.CoopLavel, data.Company, data.Year, data.Exp, data.Chengguo, data.Lingyudalei, data.Yanjiulingyu, data.Shanchang, data.BillFile, data.ContractFile, data.CertFile, data.JuralCount, data.PartnerId, data.Status, data.OperateType, data.IsDeleted, data.RelateId, data.RelateIds)
	return ret, err
}

func (m *defaultMsLessonTeacherServiceMapModel) Update(ctx context.Context, data *MsLessonTeacherServiceMap) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, msLessonTeacherServiceMapRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.LessonTeacherId, data.ServiceType, data.ClassModality, data.ClassType, data.Jingyingfuwu, data.Lavel, data.ServiceGooddesc, data.SignDate, data.EndDate, data.CoopType, data.CoopLavel, data.Company, data.Year, data.Exp, data.Chengguo, data.Lingyudalei, data.Yanjiulingyu, data.Shanchang, data.BillFile, data.ContractFile, data.CertFile, data.JuralCount, data.PartnerId, data.Status, data.OperateType, data.IsDeleted, data.RelateId, data.RelateIds, data.Id)
	return err
}

func (m *defaultMsLessonTeacherServiceMapModel) tableName() string {
	return m.table
}
