package models

import (
	models "admin/app/models/user"
	"admin/utils"
	"fmt"
)

// 买卖(广告)表
type Ads struct {
	Id          uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid         uint64 `xorm:"INT(10)" json:"uid"`              // 用户ID
	TypeId      uint32 `xorm:"TINYINT(1)" json:"type_id"`       // 类型:1出售 2购买
	TokenId     uint32 `xorm:"INT(10)" json:"token_id"`         // 货币类型
	TokenName   string `xorm:"VARBINARY(36)" json:"token_name"` // 货币名称
	Price       uint64 `xorm:"BIGINT(20)" json:"price"`         // 单价
	Num         uint64 `xorm:"BIGINT(20)" json:"num"`           // 数量
	Premium     int32  `xorm:"INT(10)" json:"premium"`          // 溢价
	AcceptPrice uint64 `xorm:"BIGINT(20)" json:"accept_price"`  // 可接受最低[高]单价
	MinLimit    uint32 `xorm:"INT(10)" json:"min_limit"`        // 最小限额
	MaxLimit    uint32 `xorm:"INT(10)" json:"max_limit"`        // 最大限额
	IsTwolevel  uint32 `xorm:"TINYINT(1)" json:"is_twolevel"`   // 是否要通过二级认证:0不通过 1通过
	Pays        string `xorm:"VARBINARY(50)" json:"pays"`       // 支付方式:以 , 分隔: 1,2,3
	Remarks     string `xorm:"VARBINARY(512)" json:"remarks"`   // 交易备注
	Reply       string `xorm:"VARBINARY(512)" json:"reply"`     // 自动回复问候语
	IsUsd       uint32 `xorm:"TINYINT(1)" json:"is_usd"`        // 是否美元支付:0否 1是
	States      uint32 `xorm:"TINYINT(1)" json:"states"`        // 状态:0下架 1上架
	CreatedTime string `xorm:"DATETIME" json:"created_time"`    // 创建时间
	UpdatedTime string `xorm:"DATETIME" json:"updated_time"`    // 修改时间
	IsDel       uint32 `xorm:"TINYINT(1)" json:"is_del"`        // 是否删除:0不删除 1删除
}

// 法币交易列表 - 用户虚拟币-订单统计 - 用户虚拟货币资产
type AdsUserCurrencyCount struct {
	Ads     `xorm:"extends"`
	Balance int64
	Freeze  int64
	Success uint32
	Uname   string
	Phone   string
	Email   string
	Ustatus uint32
}

type UserInfo struct {
	Uid              uint64 `xorm:"not null pk autoincr comment('用户ID') BIGINT(11)"`
	Account          string `xorm:"comment('账号') unique VARCHAR(64)"`
	Pwd              string `xorm:"comment('密码') VARCHAR(255)"`
	Country          string `xorm:"comment('地区号') VARCHAR(32)"`
	Phone            string `xorm:"comment('手机') unique VARCHAR(64)"`
	PhoneVerifyTime  int    `xorm:"comment('手机验证时间') INT(11)"`
	Email            string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	EmailVerifyTime  int    `xorm:"comment('邮箱验证时间') INT(11)"`
	GoogleVerifyId   string `xorm:"comment('谷歌私钥') VARCHAR(128)"`
	GoogleVerifyTime int    `xorm:"comment('谷歌验证时间') INT(255)"`
	SmsTip           int    `xorm:"default 0 comment('短信提醒') TINYINT(1)"`
	PayPwd           string `xorm:"comment('支付密码') VARCHAR(255)"`
	NeedPwd          int    `xorm:"comment('免密设置1开启0关闭') TINYINT(1)"`
	NeedPwdTime      int    `xorm:"comment('免密周期') INT(11)"`
	Status           uint32 `xorm:"default 0 comment('用户状态，0正常，1冻结') INT(11)"`
	SecurityAuth     int    `xorm:"comment('认证状态1110') TINYINT(8)"`
}

type UserOhterInfo struct {
	UserInfo `xorm:"extends"`
	NickName string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	//RegisterTime int64  `xorm:"comment('注册时间') BIGINT(20)"`
}

//g_common
func (w *UserOhterInfo) TableName() string {
	return "user"
}

func (UserInfo) TableName() string {
	return "user"
}

//g_currency
func (AdsUserCurrencyCount) TableName() string {
	return "ads"
}

//法币挂单信息参数信息
type Currency struct {
	Page    int `form:"page" json:"page" binding:"required"`
	PageNum int `form:"rows" json:"rows" `
	//g_backstage
	Uid     int    `form:"uid" json:"uid" `
	Uname   string `form:"uname" json:"uname" `
	Phone   string `form:"phone" json:"phone" `
	Email   string `form:"email" json:"email" `
	Ustatus int    `form:"ustatus" json:"ustatus" ` //用户登录状态
	/// g_currency
	Date    string `form:"date" json:"date" `         //挂单日期
	Verify  int    `form:"verify" json:"verify" `     //实名认证 二级认证 google 验证
	TokenId int    `form:"token_id" json:"token_id" ` //货币名称
	TradeId int    `form:"tid" json:"tid" `           //买卖ID
}

// 法币交易列表 - (广告(买卖))
func (this *Ads) GetAdsList(cur Currency) ([]AdsUserCurrencyCount, int, int, error) {
	if cur.PageNum <= 0 {
		cur.PageNum = 100
	}

	engine := utils.Engine_currency
	result, err := engine.Count(new(Ads))
	var total int
	total = int(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, 0, 0, err
	}
	//总页数
	var page int
	if total > cur.PageNum {
		page = total / cur.PageNum
		v := total % cur.PageNum
		if v != 0 {
			page = page + 1
		}
	} else {
		page = 1
	}

	if cur.Page <= 1 {
		cur.Page = 1
	}
	limit := 0
	if cur.Page > 1 {
		limit = int((cur.Page - 1) * cur.PageNum)
	}

	data := make([]AdsUserCurrencyCount, 0)
	// err = engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid AND ads.token_id=user_currency.token_id").
	// 	Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
	// 	Where("ads.type_id=? AND ads.token_id=?", TypeId, TokenId).
	// 	Desc("updated_time").
	// 	Limit(int(cur.PageNum), limit).
	// 	Find(&data)
	//
	//挂单日期
	if cur.TradeId != 0 {
		this.trade_id(cur.Page, cur.PageNum, cur.TradeId, &data)
	} else if cur.Uid != 0 {
		this.uid(cur.Page, cur.PageNum, cur.Uid, &data)
	} else if cur.TokenId != 0 {
		this.token_id(cur.Page, cur.PageNum, cur.TokenId, &data)
	} else if len(cur.Date) != 0 {
		this.date(cur.Page, cur.PageNum, cur.Date, &data)
	} else if len(cur.Phone) != 0 {
		this.phone(cur.Page, cur.PageNum, cur.Phone, &data)
	} else if len(cur.Email) != 0 {
		this.email(cur.Page, cur.PageNum, cur.Phone, &data)
	} else if cur.Ustatus != 0 {
		this.ustatus(cur.Page, cur.PageNum, cur.Ustatus, &data)
	} else if len(cur.Uname) != 0 {
		this.uname(cur.Page, cur.PageNum, cur.Uname, &data)
	} else if cur.Verify != 0 {
		this.verify(cur.Page, cur.PageNum, cur.Verify, &data)
	} else {
		err = engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
			Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
			Desc("updated_time").
			Limit(int(cur.PageNum), limit).
			Find(&data)
		if err != nil {
			utils.AdminLog.Errorln(err.Error())
			return nil, 0, 0, err
		}
	}

	//无条件判断

	uid := make([]uint64, 0)
	for _, v := range data {
		uid = append(uid, v.Uid)
	}
	ulist, uerr := this.getUserList(uid)
	if uerr != nil {
		return nil, 0, 0, err
	}
	for index, value := range data {
		for _, v := range ulist {
			if value.Uid == v.Uid {
				data[index].Phone = v.Phone
				data[index].Uname = v.NickName
				data[index].Ustatus = uint32(v.Status)
				data[index].Email = v.Email
				break
			}
		}
	}
	return data, page, total, nil
}

func (a *Ads) getUserList(uid []uint64) (ulist []models.UserGroup, err error) {
	engine := utils.Engine_common
	err = engine.Sql("select a.*,b.nick_name,b.register_time from `user` a left join user_ex b on a.uid=b.uid ").In("uid", uid).Find(&ulist)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return
	}
	//fmt.Printf("getUserList%#v", ulist)
	return
}

func (a *Ads) date(page, limit int, date string, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency

	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("date<?", date).
		Desc("updated_time").
		Limit(page, limit).
		Find(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

func (a *Ads) phone(page, limit int, phone string, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency

	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("phone=?", phone).
		Find(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

func (a *Ads) email(page, limit int, email string, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency
	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("email=?", email).
		Find(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

func (a *Ads) uid(page, limit int, uid int, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency
	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("uid=?", uid).
		Find(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

func (a *Ads) ustatus(page, limit int, status int, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency
	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("status=?", status).
		Find(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

func (a *Ads) uname(page, limit int, uname string, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency
	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("name=?", uname).
		Find(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

func (a *Ads) trade_id(page, limit int, tradeid int, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency
	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("type_id=?", tradeid).
		Find(result)
	fmt.Printf("tid%#v\n", tradeid)
	fmt.Printf("trade_id%#v\n", result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

func (a *Ads) token_id(page, limit int, tokenId int, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency
	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("token_id=?", tokenId).
		Find(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

func (a *Ads) verify(page, limit int, verify int, result *[]AdsUserCurrencyCount) error {
	engine := utils.Engine_currency
	err := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("security_auth=?", verify).
		Find(result)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return err
	}
	return nil
}

// 个人法币交易列表 - (广告(买卖))
// func (this *Ads) AdsUserList(Uid uint64, TypeId, Page, PageNum uint32) ([]AdsUserCurrencyCount, int64) {
// 	engine := utils.Engine_currency
// 	total, err := engine.Where("uid=? AND type_id=?", Uid, TypeId).Count(new(Ads))
// 	if err != nil {
// 		Log.Errorln(err.Error())
// 		return nil, 0
// 	}
// 	if total <= 0 {
// 		return nil, 0
// 	}

// 	limit := 0
// 	if Page > 0 {
// 		limit = int((Page - 1) * PageNum)
// 	}

// 	data := make([]AdsUserCurrencyCount, 0)
// 	err = engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid AND ads.token_id=user_currency.token_id").
// 		Where("ads.uid=? AND ads.type_id=?", Uid, TypeId).
// 		Desc("updated_time").
// 		Limit(int(PageNum), limit).
// 		Find(&data)

// 	if err != nil {
// 		Log.Errorln(err.Error())
// 		return nil, 0
// 	}

// 	return data, total
// }
