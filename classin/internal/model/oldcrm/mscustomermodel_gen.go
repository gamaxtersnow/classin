package oldcrm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"strings"
)

var (
	msCustomerFieldNames          = builder.RawFieldNames(&MsCustomer{})
	msCustomerRows                = strings.Join(msCustomerFieldNames, ",")
	msCustomerRowsExpectAutoSet   = strings.Join(stringx.Remove(msCustomerFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	msCustomerRowsWithPlaceHolder = strings.Join(stringx.Remove(msCustomerFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	msCustomerModel interface {
		Insert(ctx context.Context, data *MsCustomer) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MsCustomer, error)
		Update(ctx context.Context, data *MsCustomer) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMsCustomerModel struct {
		conn  sqlx.SqlConn
		table string
	}

	MsCustomer struct {
		Id                      int64          `db:"id"`                        // 客户ID
		Type                    sql.NullInt64  `db:"type"`                      // 客户类型  1 有效客户 2 在途客户 3 公海池
		SourceType              sql.NullInt64  `db:"source_type"`               // 1客户快速创建 2客户列表导入 3公海池创建 4公海池导入 5超级小伙伴  6 活动导入 7太平系统 8太保系统
		Sex                     sql.NullInt64  `db:"sex"`                       // 性别 1 男 2 女
		Name                    sql.NullString `db:"name"`                      // 客户名称
		IsCollision             sql.NullInt64  `db:"is_collision"`              // 撞单 0 否 1 是
		ChannelId               sql.NullInt64  `db:"channel_id"`                // 渠道ID
		ChannelName             sql.NullString `db:"channel_name"`              // 渠道名称
		CreatorId               sql.NullInt64  `db:"creator_id"`                // 创建人ID
		CreatorName             sql.NullString `db:"creator_name"`              // 创建人名称
		AllotStatus             sql.NullInt64  `db:"allot_status"`              // 分配顾问状态 0 未分配 1 已分配
		AllotTime               sql.NullTime   `db:"allot_time"`                // 分配时间
		SignStatus              sql.NullInt64  `db:"sign_status"`               // 签约状态 0 未签约 1 已签约
		IsTmk                   sql.NullInt64  `db:"is_tmk"`                    // 是否展示在TMK列表 0 否 1 是
		TmkId                   sql.NullInt64  `db:"tmk_id"`                    // TMK_ID
		TmkName                 sql.NullString `db:"tmk_name"`                  // TMK名称
		Phone                   sql.NullString `db:"phone"`                     // 电话
		OtherPhone              sql.NullString `db:"other_phone"`               // 备用电话
		FamilyPhoneOne          sql.NullString `db:"family_phone_one"`          // 亲属电话1
		FamilyPhoneOneName      sql.NullString `db:"family_phone_one_name"`     // 亲属称呼1
		FamilyPhoneTwo          sql.NullString `db:"family_phone_two"`          // 亲属电话2
		FamilyPhoneTwoName      sql.NullString `db:"family_phone_two_name"`     // 亲属称呼2
		Wechat                  sql.NullString `db:"wechat"`                    // 微信号
		Qq                      sql.NullString `db:"qq"`                        // qq号
		Email                   sql.NullString `db:"email"`                     // 邮箱
		ProvinceId              sql.NullInt64  `db:"province_id"`               // 省ID
		ProvinceName            sql.NullString `db:"province_name"`             // 省名称
		CityId                  sql.NullInt64  `db:"city_id"`                   // 市
		CityName                sql.NullString `db:"city_name"`                 // 市名称
		CourseId                sql.NullInt64  `db:"course_id"`                 // 课程体系id
		CourseName              sql.NullString `db:"course_name"`               // 课程体系名称
		SchoolId                sql.NullInt64  `db:"school_id"`                 // 在读学校
		ClassId                 sql.NullInt64  `db:"class_id"`                  // 在读年级
		MajorId                 sql.NullInt64  `db:"major_id"`                  // 在读专业
		Note                    sql.NullString `db:"note"`                      // 备注
		CompanyId               sql.NullInt64  `db:"company_id"`                // 分公司ID
		CompanyName             sql.NullString `db:"company_name"`              // 分公司名称
		OrderCode               sql.NullString `db:"order_code"`                // 订单编号
		AddTime                 sql.NullTime   `db:"add_time"`                  // 添加时间
		FirstInterviewTime      sql.NullTime   `db:"first_interview_time"`      // 首次面访完成时间
		CustomerInterviewStatus sql.NullInt64  `db:"customer_interview_status"` // 客户面访状态: -1/无面访 0/未完成面访 1/已完成面访
		UpdateTime              sql.NullTime   `db:"update_time"`               // 更新时间
		TmkStatus               sql.NullInt64  `db:"tmk_status"`                // 0 未操作  1无效  2新增有效 3未确定  4无效激活
		OldSchool               sql.NullString `db:"old_school"`                // 原在读学校
		OldClass                sql.NullString `db:"old_class"`                 // 原在读年纪
		OldMajor                sql.NullString `db:"old_major"`                 // 原在读专业
		BirthdayTime            sql.NullTime   `db:"birthday_time"`             // 生日日期
		DepositStatus           sql.NullInt64  `db:"deposit_status"`            // 1已付款 2已退款
		PhoneCarrier            sql.NullString `db:"phone_carrier"`             // 手机号运营商
		ActivityId              sql.NullInt64  `db:"activity_id"`               // 活动关联Id
		Nationality             sql.NullString `db:"nationality"`               // 国籍
		IsCustomerService       sql.NullInt64  `db:"is_customer_service"`       // 是否展示在客服列表 0 否 1 是
		ServiceTeacherId        sql.NullInt64  `db:"service_teacher_id"`        // 客服id
		ServiceName             sql.NullString `db:"service_name"`              // 客服名称
		FirstVisitTime          sql.NullTime   `db:"first_visit_time"`          // 首次回访时间
		LastVisitTime           sql.NullTime   `db:"last_visit_time"`           // 最近一次回访时间
		PlanVisitTime           sql.NullTime   `db:"plan_visit_time"`           // 计划下次回访时间
		FirstSignorderTime      sql.NullTime   `db:"first_signorder_time"`      // 首次下单时间
		CusSource               sql.NullString `db:"cus_source"`                // 客户来源 例3-45-77
		FirstConsultantTime     sql.NullTime   `db:"first_consultant_time"`     // 首次分配时间
		HighLavel               sql.NullInt64  `db:"high_lavel"`                // 客户最高评级
		QuoteStatus             sql.NullInt64  `db:"quote_status"`              // 报价状态 1已报价 0 未报价
		CustomerLabel           sql.NullInt64  `db:"customer_label"`            // 客户标签
		FirstEffective          sql.NullInt64  `db:"first_effective"`           // 首次是否有效 1是 0不是
		CustomerSignTab         sql.NullInt64  `db:"customer_sign_tab"`         // 客户来源原始常规标记
		FirstAffirmTime         sql.NullTime   `db:"first_affirm_time"`         // 第一次确认收款时间
		PolicyholderNumber      sql.NullString `db:"policyholder_number"`       // 阳光客户号
		AddCustomerServiceTime  sql.NullTime   `db:"add_customer_service_time"` // 进入到客服列表的时间
		IsShare                 sql.NullInt64  `db:"is_share"`                  // 是否共享 1共享 0不共享
		ShareId                 sql.NullInt64  `db:"share_id"`                  // 共享人id
		ShareName               sql.NullString `db:"share_name"`                // 共享人name
		ShareScale              sql.NullString `db:"share_scale"`               // 共享人比例
		ShareScaleOperate       sql.NullInt64  `db:"share_scale_operate"`       // 目前可操作人
		ShareScaleStatus        sql.NullInt64  `db:"share_scale_status"`        // 确认状态 0未确认 1已确认 2已更改
		HeadId                  sql.NullInt64  `db:"head_id"`                   // 负责人id
		HeadName                sql.NullString `db:"head_name"`                 // 负责人名称
		ChangeHeadReason        sql.NullString `db:"change_head_reason"`        // 转负责人的原因
		ResidePapers            sql.NullInt64  `db:"reside_papers"`             // 1 是  2否  是否有居住证件
		PapersCountry           sql.NullInt64  `db:"papers_country"`            // 居住证国家
		Residence               sql.NullString `db:"residence"`                 // 居住地
		SubmitTime              sql.NullTime   `db:"submit_time"`               // 转案修改时间
		TransferAddTime         sql.NullTime   `db:"transfer_add_time"`         // 转案补充完成时间
		CustomerStage           sql.NullInt64  `db:"customer_stage"`            // 客户阶段 1线索 2资源 3商机 4已签约客户 5商机待确认
		IsAgree                 sql.NullInt64  `db:"is_agree"`                  // 渠道和客户是否同意沟通 1是 2否
		ConfirmReceive          sql.NullString `db:"confirm_receive"`           // 公海池确认领取 领取商机用到的字段 存储确认人id
		ClueToResourceTime      sql.NullTime   `db:"clue_to_resource_time"`     // 线索转资源时间
		ResourceToBusinessTime  sql.NullTime   `db:"resource_to_business_time"` // 资源转商机时间
		SignStatusDetail        sql.NullInt64  `db:"sign_status_detail"`        // 签约状态详细信息：1未生成订单 2生成订单 3生成合同 4签署订单; 20财务确认收款
		DetermineStatus         sql.NullInt64  `db:"determine_status"`          // 判定状态:1已判定 2未判定
		DetermineTimeEnd        sql.NullTime   `db:"determine_time_end"`        // 需要判定的时间结束
		ImportLogId             sql.NullInt64  `db:"import_log_id"`             // 导入日志记录的id
		IsChangeConsultant      sql.NullInt64  `db:"is_change_consultant"`      // 责任顾问是否发生了变更：1是，0否
	}
)

func newMsCustomerModel(conn sqlx.SqlConn) *defaultMsCustomerModel {
	return &defaultMsCustomerModel{
		conn:  conn,
		table: "`ms_customer`",
	}
}

func (m *defaultMsCustomerModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMsCustomerModel) FindOne(ctx context.Context, id int64) (*MsCustomer, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", msCustomerRows, m.table)
	var resp MsCustomer
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

func (m *defaultMsCustomerModel) Insert(ctx context.Context, data *MsCustomer) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, msCustomerRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Type, data.SourceType, data.Sex, data.Name, data.IsCollision, data.ChannelId, data.ChannelName, data.CreatorId, data.CreatorName, data.AllotStatus, data.AllotTime, data.SignStatus, data.IsTmk, data.TmkId, data.TmkName, data.Phone, data.OtherPhone, data.FamilyPhoneOne, data.FamilyPhoneOneName, data.FamilyPhoneTwo, data.FamilyPhoneTwoName, data.Wechat, data.Qq, data.Email, data.ProvinceId, data.ProvinceName, data.CityId, data.CityName, data.CourseId, data.CourseName, data.SchoolId, data.ClassId, data.MajorId, data.Note, data.CompanyId, data.CompanyName, data.OrderCode, data.AddTime, data.FirstInterviewTime, data.CustomerInterviewStatus, data.TmkStatus, data.OldSchool, data.OldClass, data.OldMajor, data.BirthdayTime, data.DepositStatus, data.PhoneCarrier, data.ActivityId, data.Nationality, data.IsCustomerService, data.ServiceTeacherId, data.ServiceName, data.FirstVisitTime, data.LastVisitTime, data.PlanVisitTime, data.FirstSignorderTime, data.CusSource, data.FirstConsultantTime, data.HighLavel, data.QuoteStatus, data.CustomerLabel, data.FirstEffective, data.CustomerSignTab, data.FirstAffirmTime, data.PolicyholderNumber, data.AddCustomerServiceTime, data.IsShare, data.ShareId, data.ShareName, data.ShareScale, data.ShareScaleOperate, data.ShareScaleStatus, data.HeadId, data.HeadName, data.ChangeHeadReason, data.ResidePapers, data.PapersCountry, data.Residence, data.SubmitTime, data.TransferAddTime, data.CustomerStage, data.IsAgree, data.ConfirmReceive, data.ClueToResourceTime, data.ResourceToBusinessTime, data.SignStatusDetail, data.DetermineStatus, data.DetermineTimeEnd, data.ImportLogId, data.IsChangeConsultant)
	return ret, err
}

func (m *defaultMsCustomerModel) Update(ctx context.Context, data *MsCustomer) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, msCustomerRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Type, data.SourceType, data.Sex, data.Name, data.IsCollision, data.ChannelId, data.ChannelName, data.CreatorId, data.CreatorName, data.AllotStatus, data.AllotTime, data.SignStatus, data.IsTmk, data.TmkId, data.TmkName, data.Phone, data.OtherPhone, data.FamilyPhoneOne, data.FamilyPhoneOneName, data.FamilyPhoneTwo, data.FamilyPhoneTwoName, data.Wechat, data.Qq, data.Email, data.ProvinceId, data.ProvinceName, data.CityId, data.CityName, data.CourseId, data.CourseName, data.SchoolId, data.ClassId, data.MajorId, data.Note, data.CompanyId, data.CompanyName, data.OrderCode, data.AddTime, data.FirstInterviewTime, data.CustomerInterviewStatus, data.TmkStatus, data.OldSchool, data.OldClass, data.OldMajor, data.BirthdayTime, data.DepositStatus, data.PhoneCarrier, data.ActivityId, data.Nationality, data.IsCustomerService, data.ServiceTeacherId, data.ServiceName, data.FirstVisitTime, data.LastVisitTime, data.PlanVisitTime, data.FirstSignorderTime, data.CusSource, data.FirstConsultantTime, data.HighLavel, data.QuoteStatus, data.CustomerLabel, data.FirstEffective, data.CustomerSignTab, data.FirstAffirmTime, data.PolicyholderNumber, data.AddCustomerServiceTime, data.IsShare, data.ShareId, data.ShareName, data.ShareScale, data.ShareScaleOperate, data.ShareScaleStatus, data.HeadId, data.HeadName, data.ChangeHeadReason, data.ResidePapers, data.PapersCountry, data.Residence, data.SubmitTime, data.TransferAddTime, data.CustomerStage, data.IsAgree, data.ConfirmReceive, data.ClueToResourceTime, data.ResourceToBusinessTime, data.SignStatusDetail, data.DetermineStatus, data.DetermineTimeEnd, data.ImportLogId, data.IsChangeConsultant, data.Id)
	return err
}

func (m *defaultMsCustomerModel) tableName() string {
	return m.table
}
