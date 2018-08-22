package models

import (
	"admin/utils"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// 订单表
/*type Order struct {
	BaseModel      `xorm:"-"`
	SubductionZero `xorm:"-"`
	Id             uint64 `xorm:"not null pk autoincr comment('ID')  INT(10)"  json:"id"`
	OrderId        string `xorm:"not null pk comment('订单ID') INT(10)"   json:"order_id"` // hash( type_id, 6( user_id, + 时间秒）
	AdId           uint64 `xorm:"not null default 0 comment('广告ID') index INT(10)"  json:"ad_id"`
	AdType         uint32 `xorm:"not null default 0 comment('广告类型:1出售 2购买') TINYINT(1)"  json:"ad_type"`
	Price          int64  `xorm:"not null default 0 comment('价格') BIGINT(64)"   json:"price"`
	Num            int64  `xorm:"not null default 0 comment('数量') BIGINT(64)"   json:"num"`
	TokenId        uint64 `xorm:"not null default 0 comment('货币类型') INT(10)"       json:"token_id"`
	PayId          string `xorm:"not null default 0 comment('支付类型') VARCHAR(64)"       json:"pay_id"`
	SellId         uint64 `xorm:"not null default 0 comment('卖家id') INT(10)"         json:"sell_id"`
	SellName       string `xorm:"not null default '' comment('卖家昵称') VARCHAR(64)"  json:"sell_name"`
	BuyId          uint64 `xorm:"not null default 0 comment('买家id') INT(10)"    json:"buy_id"`
	BuyName        string `xorm:"not null default '' comment('买家昵称') VARCHAR(64)"   json:"buy_name"`
	Fee            int64  `xorm:"not null default 0 comment('手续费用') BIGINT(64)"  json:"fee"`
	States         uint32 `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(1)"   json:"states"`
	PayStatus      uint32 `xorm:"not null default 0 comment('支付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)') TINYINT(1)"  json:"pay_status"`
	CancelType     uint32 `xorm:"not null default 0 comment('取消类型: 1卖方 2 买方') TINYINT(1)"   json:"cancel_type"`
	CreatedTime    string `xorm:"not null comment('创建时间') DATETIME"  json:"created_time"`
	UpdatedTime    string `xorm:"comment('修改时间')     DATETIME"               json:"updated_time"`
	//ConfirmTime sql.NullString `xorm:"default null comment('确认支付时间')  DATETIME"     json:"confirm_time"`
	//ReleaseTime sql.NullString `xorm:"default null comment('放行时间')     DATETIME"     json:"release_time"`
	ConfirmTime string `xorm:"default null comment('确认支付时间')  DATETIME"     json:"confirm_time"`
	ReleaseTime string `xorm:"default null comment('放行时间')     DATETIME"     json:"release_time"`
}*/

type Order struct {
	BaseModel      `xorm:"-"`
	SubductionZero `xorm:"-"`
	Id             int    `xorm:"not null pk autoincr comment('ID') INT(10)" json:"id"`
	OrderId        string `xorm:"not null default '' comment('订单ID') unique VARCHAR(64)" json:"order_id"`
	AdId           int    `xorm:"not null default 0 comment('广告ID') index INT(10)" json:"ad_id"`
	AdType         int    `xorm:"not null default 0 comment('广告类型:1出售 2购买') TINYINT(1)" json:"ad_type"`
	Price          int64  `xorm:"not null default 0 comment('价格') BIGINT(64)" json:"price"`
	Num            int64  `xorm:"not null default 0 comment('数量') BIGINT(64)" json:"num"`
	NumTotalPrice  int64  `xorm:"default 0 comment('后台需要的数量总价格') BIGINT(64)" json:"num_total_price"` //折合人民币
	TokenId        int    `xorm:"not null default 0 comment('货币类型') INT(10)" json:"token_id"`
	PayId          string `xorm:"not null default '0' comment('支付类型') VARCHAR(64)" json:"pay_id"`
	SellId         int    `xorm:"not null default 0 comment('卖家id') INT(10)" json:"sell_id"`
	SellName       string `xorm:"not null default '' comment('卖家昵称') VARCHAR(64)" json:"sell_name"`
	BuyId          int    `xorm:"not null default 0 comment('买家id') INT(10)" json:"buy_id"`
	BuyName        string `xorm:"not null default '' comment('买家昵称') VARCHAR(64)" json:"buy_name"`
	Fee            int64  `xorm:"not null default 0 comment('手续费用') BIGINT(64)" json:"fee"`
	FeePrice       int64  `xorm:"default 0 comment('后台需要的计算出费用价格') BIGINT(64)" json:"fee_price"` //折合人民币
	States         int    `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(1)" json:"states"`
	PayStatus      int    `xorm:"not null default 0 comment('支付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)') TINYINT(1)" json:"pay_status"`
	CancelType     int    `xorm:"not null default 0 comment('取消类型: 1卖方 2 买方') TINYINT(1)" json:"cancel_type"`
	CreatedTime    string `xorm:"not null comment('创建时间') DATETIME" json:"created_time"`
	UpdatedTime    string `xorm:"comment('修改时间') DATETIME" json:"updated_time"`
	ExpiryTime     string `xorm:"comment('过期时间') DATETIME" json:"expiry_time"`
	ConfirmTime    string `xorm:"comment('确认支付时间') DATETIME" json:"confirm_time"`
	ReleaseTime    string `xorm:"comment('放行时间') DATETIME" json:"release_time"`
}

type OrderGroup struct {
	Order          `xorm:"extends"`
	Uid            uint64  `xorm:"INT(10)"     json:"uid"`
	TokenName      string  //货币名称
	BuyQuantity    float64 //buy数量
	BuyTotalPrice  int64   //总额
	SellQuantity   float64 //卖出数量
	SellTotalPrice int64   //总额
	Transfer       float64
}

type OrderAddName struct {
	Order `xorm:"extends"`
	Mark  string `xorm:"VARBINARY(20)" json:"Name"` // 货币标识
}

func (o *Order) TableName() string {
	return "order"
}

func (o *OrderAddName) TableName() string {
	return "order"
}

func (o *OrderGroup) TableName() string {
	return "order"
}

//查询个人的所有数据货币的交易记录
func (this *Order) GetOrderListOfUid(page, rows, uid, token_id int) (*ModelList, error) {

	engine := utils.Engine_currency
	//统计
	type statistics struct {
		TokenId      int     `json:"token_id"`
		TokenName    string  `json:"token_name"`
		BuyTotal     float64 `json:"buy_toal"`          //累计买入
		BuyTotalCny  float64 `json:"buy_toal_cny"`      //累计买入折合
		SellTotal    float64 `json:"sell_total"`        //累计卖出
		SellTotalCny float64 `json:"sell_total_cny"`    //累计卖出折合
		Transfer     float64 `json:"transfer" xorm:"-"` //累计划转
	}
	sql := fmt.Sprintf("SELECT o.token_id,SUM( IF( o.sell_id=%d, o.num_total_price, 0)) DIV 100000000 AS sell_total_cny , SUM( IF( o.sell_id=%d, o.num, 0)) DIV 100000000 AS sell_total, SUM( IF( o.sell_id=%d, o.num_total_price, 0))DIV 100000000 AS buy_total_cny, SUM( IF( o.sell_id=%d, o.num, 0))DIV 100000000 AS buy_total FROM `order` o  ", uid, uid, uid, uid)
	condition := ""
	if token_id != 0 {
		condition = fmt.Sprintf("WHERE o.token_id=%d ", token_id)
	} else {
		condition = fmt.Sprintf(" GROUP BY o.token_id ")
	}
	countSql := fmt.Sprintf("select count(t.token_id) num from(%s) t", sql+condition)
	count := &struct {
		Num int
	}{}
	_, err := engine.SQL(countSql).Get(count)
	if err != nil {
		return nil, err
	}
	offset, mList := this.Paging(page, rows, count.Num)
	limitStr := fmt.Sprintf("LIMIT %d, %d; ", offset, mList.PageSize)
	list := make([]statistics, 0)
	err = engine.SQL(sql + condition + limitStr).Find(&list)
	if err != nil {
		return nil, err
	}
	//查询所有币种名称及Id
	reslt, err := new(CommonTokens).GetTokenList()
	if err != nil {
		return nil, err
	}
	tidList := make([]int, 0)
	for index, tokenid := range list {
		//根据token_id 查找货币名称
		for _, value := range reslt {
			if value.Id == uint32(tokenid.TokenId) {
				list[index].TokenName = value.Mark
				tidList = append(tidList, tokenid.TokenId)
				break
			}
		}
	}

	type transfer struct {
		TokenId int64
		Num     int64
	}
	transferList := make([]transfer, 0)
	//查询划转的数量
	query := engine.Table("user_currency_history").Select("SUM(num)DIV 10000000 AS SUM ,token_id").In("token_id", transferList).GroupBy("token_id")
	queryCount := *query
	ct, err := queryCount.Count(&UserCurrencyHistory{})
	if err != nil {
		return nil, err
	}
	if ct == 0 {
		//根据用户uid 获取该用户的所有划转数量
		engine.Select("SUM(num)DIV 10000000 AS SUM ,token_id").Where("uid=?", uid)
	}
	toffset, tList := this.Paging(page, rows, int(ct))
	err = query.Limit(toffset, tList.PageSize).Find(&transferList)
	if err != nil {
		return nil, err
	}

	mList.Items = list
	return mList, nil
}

//
//根据 uid  及交易状态 多表查询拉取 所有相关订单的交易记录
func (this *Order) GetOrderId(uid []int, status int) ([]OrderGroup, error) {
	if len(uid) <= 0 {
		return nil, errors.New("uid [] is empty!!")
	}
	fmt.Println("GetOrderId", uid, status)
	list := make([]OrderGroup, 0)
	engine := utils.Engine_currency
	query := engine.Desc("order.id")
	query = query.Join("INNER", "ads", "order.ad_id=ads.id")
	query = query.In("uid", uid)
	query = query.Where("pay_status=?", status)
	err := query.Find(&list)

	//err := engine.In("uid", orderId).Where("status=?", status).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (this *Order) GetOrderList(Page, PageNum, AdType, States, TokenId int, search string) (*ModelList, error) {
	engine := utils.Engine_currency
	query := engine.Desc("order.id")
	query = query.Join("LEFT", "g_common.tokens t", "order.token_id=t.id")
	if AdType != 0 {
		query = query.Where("ad_type=?", AdType)
	}
	if States == 5 {
		query = query.AllCols().Where("states=?", 0)
	}
	if States != 0 {
		query = query.AllCols().Where("states=?", States)
	}
	if TokenId != 0 {
		query = query.Where("token_id=?", TokenId)
	}
	//if StartTime != `` {
	//	substr := StartTime[:11] + "23:59:59"
	//	temp := fmt.Sprintf("created_time BETWEEN '%s' AND '%s' ", StartTime, substr)
	//	query = query.Where(temp)
	//}
	if search != `` {
		temp := fmt.Sprintf(" concat(IFNULL(sell_name,''),IFNULL(buy_name,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}

	tmpQuery := *query
	count, err := tmpQuery.Count(&Order{})
	if err != nil {
		return nil, err
	}
	offset, modelList := this.Paging(Page, PageNum, int(count))
	//查询符合要求数据
	list := make([]OrderAddName, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	//所有符合要求的数据的函数

	if err != nil {
		return nil, err
	}
	//去掉零
	for i, v := range list {
		num, price := this.SubductionZeroMethodInt64(v.Num, v.Price)
		list[i].NumberTrue = num
		list[i].PriceTrue = price
	}
	modelList.Items = list
	return modelList, nil

}

//法币交易手续费 ---> 注：仪表盘 买卖都需要加起来。 获取当天的。
func (this *Order) GetOrderDayFee() (float64, error) {
	engine := utils.Engine_currency
	current := time.Now().Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("SELECT m.fee fee,c.price price FROM (SELECT t.days,t.fee,t.token_id FROM (SELECT SUBSTRING(confirm_time,1,10) days,fee,token_id FROM g_currency.`order` WHERE pay_status=3) t  WHERE t.days ='%s' GROUP BY t.token_id) m JOIN  g_token.`config_token_cny` c ON m.token_id= c.token_id", current[:10])
	type fee struct {
		Fee   int64
		Price int64
	}
	list := make([]fee, 0)
	err := engine.SQL(sql).Find(&list)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, v := range list {
		result := this.Int64MulInt64By8BitString(v.Fee, v.Price)
		float, err := strconv.ParseFloat(result, 64)
		if err != nil {
			utils.AdminLog.Println(err.Error())
			continue
		}
		total += float
	}
	return total, nil
}
