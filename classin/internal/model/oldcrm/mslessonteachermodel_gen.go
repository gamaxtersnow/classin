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
	msLessonTeacherFieldNames          = builder.RawFieldNames(&MsLessonTeacher{})
	msLessonTeacherRows                = strings.Join(msLessonTeacherFieldNames, ",")
	msLessonTeacherRowsExpectAutoSet   = strings.Join(stringx.Remove(msLessonTeacherFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	msLessonTeacherRowsWithPlaceHolder = strings.Join(stringx.Remove(msLessonTeacherFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	msLessonTeacherModel interface {
		Insert(ctx context.Context, data *MsLessonTeacher) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MsLessonTeacher, error)
		Update(ctx context.Context, data *MsLessonTeacher) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMsLessonTeacherModel struct {
		conn  sqlx.SqlConn
		table string
	}

	MsLessonTeacher struct {
		Id                              int64           `db:"id"`
		Name                            sql.NullString  `db:"name"`              // 老师名称
		Phone                           sql.NullString  `db:"phone"`             // 老师电话
		Email                           sql.NullString  `db:"email"`             // 老师邮箱
		Sex                             sql.NullInt64   `db:"sex"`               // 性别 1 男 2 女
		CountryId                       sql.NullInt64   `db:"country_id"`        // 国籍id  关联country表主键
		ResumeFile                      sql.NullString  `db:"resume_file"`       // 简历文件
		WexinCode                       sql.NullString  `db:"wexin_code"`        // 微信号
		Education                       sql.NullInt64   `db:"education"`         //  最高学历 1 中学 2本科  3硕士  4 博士  5其他
		Hightitle                       sql.NullInt64   `db:"hightitle"`         // 头衔  1 中学 2本科  3硕士  4 博士  5其他
		SchoolInfo                      sql.NullString  `db:"school_info"`       // 学校信息存的json
		CooperationType                 sql.NullInt64   `db:"cooperation_type"`  // 合作类型 1-全职 2-兼职
		Status                          sql.NullInt64   `db:"status"`            // 状态 1-启用 2-禁用
		CreatorId                       sql.NullInt64   `db:"creator_id"`        // 创建人ID
		CreatorName                     sql.NullString  `db:"creator_name"`      // 创建人名称
		ClassHour                       sql.NullInt64   `db:"class_hour"`        // 课时
		BankAccountName                 sql.NullString  `db:"bank_account_name"` // 开户名
		BankCardNum                     sql.NullString  `db:"bank_card_num"`     // 卡号
		IdcardNum                       sql.NullString  `db:"idcard_num"`        // 证件号码
		IdcardType                      sql.NullString  `db:"idcard_type"`       // 卡片类型
		BankId                          sql.NullInt64   `db:"bank_id"`           // 开户网点id关联网点表
		BankName                        sql.NullString  `db:"bank_name"`         // 开户行名称
		BankNote                        sql.NullString  `db:"bank_note"`         // 备注
		OneLessonSalary                 sql.NullFloat64 `db:"one_lesson_salary"` // 单课次薪酬
		CreatedAt                       time.Time       `db:"created_at"`
		UpdatedAt                       time.Time       `db:"updated_at"`
		IdCardNo                        sql.NullString  `db:"id_card_no"`                           // 身份证号
		TeacherLevelFilesUrl            sql.NullString  `db:"teacher_level_files_url"`              // 教师入职定级表文件URL
		TeacherCertificateFilesUrl      sql.NullString  `db:"teacher_certificate_files_url"`        // 教师资格证文件URL
		TeacherIdcardFilesUrl           sql.NullString  `db:"teacher_idcard_files_url"`             // 教师身份证文件URL
		IsSelfPayee                     sql.NullInt64   `db:"is_self_payee"`                        // 是否本人收款
		Bank                            sql.NullString  `db:"bank"`                                 // 开户银行
		ApproveStatus                   sql.NullInt64   `db:"approve_status"`                       // 审批状态:0未提交审批,1审批中,2审批通过,3审批驳回
		NeedApproveTeacherServiceMapIds sql.NullString  `db:"need_approve_teacher_service_map_ids"` // 当前需要审批的老师的服务业务id(ms_lesson_teacher_service_map.id)
	}
)

func newMsLessonTeacherModel(conn sqlx.SqlConn) *defaultMsLessonTeacherModel {
	return &defaultMsLessonTeacherModel{
		conn:  conn,
		table: "`ms_lesson_teacher`",
	}
}

func (m *defaultMsLessonTeacherModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMsLessonTeacherModel) FindOne(ctx context.Context, id int64) (*MsLessonTeacher, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", msLessonTeacherRows, m.table)
	var resp MsLessonTeacher
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

func (m *defaultMsLessonTeacherModel) Insert(ctx context.Context, data *MsLessonTeacher) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, msLessonTeacherRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.Phone, data.Email, data.Sex, data.CountryId, data.ResumeFile, data.WexinCode, data.Education, data.Hightitle, data.SchoolInfo, data.CooperationType, data.Status, data.CreatorId, data.CreatorName, data.ClassHour, data.BankAccountName, data.BankCardNum, data.IdcardNum, data.IdcardType, data.BankId, data.BankName, data.BankNote, data.OneLessonSalary, data.IdCardNo, data.TeacherLevelFilesUrl, data.TeacherCertificateFilesUrl, data.TeacherIdcardFilesUrl, data.IsSelfPayee, data.Bank, data.ApproveStatus, data.NeedApproveTeacherServiceMapIds)
	return ret, err
}

func (m *defaultMsLessonTeacherModel) Update(ctx context.Context, data *MsLessonTeacher) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, msLessonTeacherRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.Phone, data.Email, data.Sex, data.CountryId, data.ResumeFile, data.WexinCode, data.Education, data.Hightitle, data.SchoolInfo, data.CooperationType, data.Status, data.CreatorId, data.CreatorName, data.ClassHour, data.BankAccountName, data.BankCardNum, data.IdcardNum, data.IdcardType, data.BankId, data.BankName, data.BankNote, data.OneLessonSalary, data.IdCardNo, data.TeacherLevelFilesUrl, data.TeacherCertificateFilesUrl, data.TeacherIdcardFilesUrl, data.IsSelfPayee, data.Bank, data.ApproveStatus, data.NeedApproveTeacherServiceMapIds, data.Id)
	return err
}

func (m *defaultMsLessonTeacherModel) tableName() string {
	return m.table
}
