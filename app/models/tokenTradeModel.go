package models

import (
	"fmt"
	"time"

	"admin/errors"
	"admin/utils"
	"admin/utils/convert"
)

//bibi 交易表
//type Trade struct {
//	BaseModel    `xorm:"-"`
//	TradeId      int    `xorm:"not null pk autoincr comment('交易表的id') INT(11)"`
//	TradeNo      string `xorm:"comment('订单号') unique(uni_reade_no) VARCHAR(32)"`
//	Uid          int64  `xorm:"comment('买家uid') index BIGINT(11)"`
//	TokenId      int    `xorm:"comment('主货币id') index INT(11)"`
//	TokenTradeId int    `xorm:"comment('交易币种') INT(11)"`
//	TokenName    string `xorm:"not null comment('交易对 名称 例如USDT/BTC') VARCHAR(10)"`
//	Price        int64  `xorm:"comment('价格') BIGINT(20)"`
//	Num          int64  `xorm:"comment('数量') BIGINT(20)"`
//	Fee          int64  `xorm:"comment('手续费') BIGINT(20)"`
//	Opt          int    `xorm:"comment(' buy  1或sell 2') index unique(uni_reade_no) TINYINT(4)"`
//	DealTime     int64  `xorm:"comment('成交时间') BIGINT(11)"`
//	States       int    `xorm:"comment('0是挂单，1是部分成交,2成交， -1撤销') INT(11)"`
//	FeeCny       int64  `xorm:"comment( '手续费折合CNY') BIGINT(20)"`
//	TotalCny     int64  `xorm:"comment( '总交易额折合CNY') BIGINT(20)"`
//}

type Trade struct {
	BaseModel    `xorm:"-"`
	TradeId      int    `xorm:"not null pk autoincr comment('交易表的id') INT(11)"`
	TradeNo      string `xorm:"comment('订单号') unique(uni_reade_no) VARCHAR(32)"`
	Uid          int64  `xorm:"comment('买家uid') index BIGINT(11)"`
	TokenId      int    `xorm:"comment('主货币id') index INT(11)"`
	TokenTradeId int    `xorm:"comment('交易币种') INT(11)"`
	Symbol       string `xorm:"not null default 'BTC' comment('交易对 名称 例如USDT/BTC') VARCHAR(16)" `
	Price        int64  `xorm:"comment('价格') BIGINT(20)"`
	Num          int64  `xorm:"comment('入账数量') BIGINT(20)"`
	Fee          int64  `xorm:"comment('手续费') BIGINT(20)"`
	Opt          int    `xorm:"comment(' buy  1或sell 2') index unique(uni_reade_no) TINYINT(4)"`
	DealTime     int64  `xorm:"comment('成交时间') index BIGINT(11)"`
	//States       int    `xorm:"comment('0是挂单，1是部分成交,2成交， -1撤销') INT(11)"`
	FeeCny           int64  `xorm:"comment('手续费折合CNY') BIGINT(20)"`
	TotalCny         int64  `xorm:"comment('总交易额折合CNY') BIGINT(20)"`
	EntrustId        string `xorm:"VARCHAR(32)"`
	TokenAdmissionId int    `xorm:"comment('入账货币id') TINYINT(4)" json:"token_admission_id"`
}

type TradeReturn struct {
	Trade          `xorm:"extends"`
	TradeNum       int64  `xorm:" comment('已成交') BIGINT(20)" json:"trade_num"`
	AllNum         int64  `xorm:"not null comment('总数量') BIGINT(20)"`  //总数
	SurplusNum     int64  `xorm:"not null comment('剩余数量') BIGINT(20)"` //余数
	ToAccountNum   string `xorm:"-" json:"to_account_num"`             //实际到账数量
	FinishCount    string `xorm:"-" json:"finish_count"`               //已成
	FeeTrue        string `xorm:"-" json:"fee_true"`
	AllNumTrue     string `xorm:"-" json:"all_num_true"`
	SurplusNumTrue string `xorm:"-" json:"surplus_num_true"`
	PriceTrue      string `xorm:"-" json:"price_true"`
	TokenName      string `xorm:"-" json:"token_name"`
	DealTimeStr      string `xorm:"-" json:"deal_time_str"`
}

func (t *TradeReturn) TableName() string {
	return "trade"
}

type TradeEx struct {
	Trade          `xorm:"extends"`
	//ConfigTokenCny `xorm:"extends"`
	TotalTrue     float64 `xorm:"-"` //交易总额
	FeeTrue       float64 `xorm:"-"` //交易手续费
	Mark 		string    `xorm:"-" json:"mark"`
	DealTimeStr      string `xorm:"-" json:"deal_time_str"`
}

type TotalTradeCNY struct {
	Date  int64  //日期
	Buy   uint64 //买入总额
	Sell  uint64 //卖出总额
	Total uint64 // 买卖总金额
}

func (t *TradeEx) TableName() string {
	return "trade"
}

func (t *TotalTradeCNY) TableName() string {
	return "trade"
}

type DayCount struct {
	TokenId     int64  `json:"token_id"`  //货币id
	Total       string `json:"total"`     //数量
	TotalCny    string `json:"total_cny"` //折合
	FeeTotal    string `json:"fee_total"`
	FeeTotalCny string `json:"fee_total_cny"`
	Date        int64  `xorm:"-" json:"date"` //日期
}

func (t *DayCount) TableName() string {
	return "trade"
}

func (t *Trade) Get(tid int, bt, et, opt int64) (*DayCount, error) {
	engine := utils.Engine_token
	dc := new(DayCount)
	_, err := engine.Select(" token_admission_id token_id ,FROM_UNIXTIME(deal_time,'%Y-%m-%d %H:%i:%s') date, SUM(num) total ,SUM(total_cny) total_cny,SUM(fee) fee_total,SUM(fee_cny) fee_total_cny ").Where("token_admission_id=? and opt=?  and deal_time between ? and ? ", tid, opt, bt, et).Get(dc)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		fmt.Println(err.Error())
		return nil, err
	}
	dc.Date = bt
	return dc, nil

}


func (this *Trade) GetTokenRecordList(page, rows, opt, uid int, bt, et uint64, name string) (*ModelList, error) {
	engine := utils.Engine_token
	fmt.Println("这里到了没有啊 ")
	query := engine.Desc("t.entrust_id")
	query = query.Alias("t").Join("left", "entrust_detail e", "e.entrust_id= t.entrust_id")
	if name != `` {
		query = query.Where("(e.states=2 or e.states=1) and e.symbol=?", name) //交易对
	}
	tm := time.Now().Unix()
	if opt != 0 {
		query = query.Where("t.opt=?", opt) //交易方向
	}
	if uid != 0 {
		query = query.Where("t.uid=?", uid)
	}
	if bt != 0 {
		if et != 0 {
			query = query.Where("t.deal_time BETWEEN ? AND ? ", bt, et+86400)
		} else {
			query = query.Where("t.deal_time BETWEEN ? AND ? ", bt, bt+86400)
		}

	} else {
		query = query.Where("t.deal_time BETWEEN ? AND ? ", tm-86400, tm)
	}
	tempQuery := *query

	count, err := tempQuery.Count(&TradeReturn{})
	if err != nil {
		return nil, err
	}
	offset, modelList := this.Paging(page, rows, int(count))
	list := make([]TradeReturn, 0)
	fmt.Printf("$$$$$$$$$$$$$$$%#v\n", rows)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	fmt.Println("list", list)
	for i, v := range list {
		//allNum, surPlusNUm := this.SubductionZeroMethodInt64(v.AllNum, v.SurplusNum)
		list[i].AllNumTrue = convert.Int64ToStringBy8Bit(v.AllNum)
		list[i].SurplusNumTrue = convert.Int64ToStringBy8Bit(v.SurplusNum)
		list[i].FinishCount = convert.Int64ToStringBy8Bit(v.Num)
		list[i].ToAccountNum = convert.Int64ToStringBy8Bit(v.TradeNum)
		list[i].FeeTrue = convert.Int64ToStringBy8Bit(v.Fee)
		list[i].PriceTrue = convert.Int64ToStringBy8Bit(v.Price)
		list[i].DealTimeStr = time.Unix(v.DealTime,0).Format("2006-01-02 15:04:05")
	}
	modelList.Items = list
	return modelList, nil
}

//p5-1-0-1币币交易手续费明细
/********************************
* id 兑币id
* trade_type 交易方向 2 卖 1买
* search 筛选
 */
func (this *Trade) GetFeeInfoList(page, rows, uid, opt int, date uint64, name string) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Desc("trade.token_id")
	//query = query.Join("left", "config_token_cny p", "trade.token_id = p.token_id")
	query =query.Where("fee !=0")
	tm := time.Now().Unix()
	toBeCharge := time.Now().Format("2006-01-02 ") + "00:00:00"
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	unix := theTime.Unix()
	if uid != 0 {
		query = query.Where("uid=?", uid)
	}
	if date != 0 {
		query = query.Where("deal_time BETWEEN ? AND ?", date, date+86400)
	} else {
		query = query.Where("deal_time BETWEEN ? AND ?", unix, tm)
	}
	if opt != 0 {
		query = query.Where("opt=?", opt)
	}
	if name != `` {
		query = query.Where("symbol=?", name)
	}
	ValuQuery := *query
	count, err := query.Distinct("deal_time").Count(&Trade{})
	if err != nil {
		return nil, err
	}
	offset, mlist := this.Paging(page, rows, int(count))
	list := make([]TradeEx, 0)
	err = ValuQuery.Limit(mlist.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	//fmt.Println("len=",len(list))

	tidList,err:=new(Tokens).GetTokensList()
	if err!=nil{
		return nil,err
	}
	for i, v := range list {
		list[i].TotalTrue = convert.Int64ToFloat64By8Bit(v.Num)
		list[i].FeeTrue = convert.Int64ToFloat64By8Bit(v.Fee)
		for _,vt:=range tidList{
			if vt.Id == v.TokenAdmissionId{
				list[i].Mark = vt.Mark
				break
			}
		}
		list[i].DealTimeStr  =time.Unix(v.DealTime,0).Format("2006-01-02 15:04:05")
	}
	mlist.Items = list
	return mlist, nil
}

// 币s币交易合计
type TokenTradeTotal struct {
	// 交易次数
	TotalTime            int64 `xorm:"total_time"`               // 交易总次数
	TodayTotalTime       int64 `xorm:"today_total_time"`         // 今日交易次数
	YesterdayTotalTime   int64 `xorm:"yesterday_total_time"`     // 上日交易次数
	LastWeekDayTotalTime int64 `xorm:"last_week_day_total_time"` // 上周同日交易次数

	// 交易量
	TotalNum            string `xorm:"total_num"`               // 总计交易量
	TodayTotalNum       string `xorm:"today_total_num"`         // 今日交易量
	YesterdayTotalNum   string `xorm:"yesterday_total_num"`     // 上日交易量
	LastWeekDayTotalNum string `xorm:"last_week_day_total_num"` // 上周同日交易量

	// 交易手续费
	TotalFee            string `xorm:"total_fee"`               // 手续费总计
	TodayTotalFee       string `xorm:"today_total_fee"`         // 今日合计手续费
	YesterdayTotalFee   string `xorm:"yesterday_total_fee"`     // 上日合计手续费
	LastWeekDayTotalFee string `xorm:"last_week_day_total_fee"` // 上周同日合计手续费
}

// 交易次数、数量、手续费合计
// 今日、上日、上周同日
func (this *Trade) TradeTotal() (*TokenTradeTotal, error) {
	// 计算日期
	todayDate := time.Now().Format(utils.LAYOUT_DATE)

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, errors.NewSys(err)
	}
	datetime, err := time.ParseInLocation(utils.LAYOUT_DATE, todayDate, loc)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	// 零点
	todayZeroUnix := datetime.Unix()
	yesterdayZeroUnix := todayZeroUnix - 24*60*60
	lastWeekDayZeroUnix := todayZeroUnix - 7*24*60*60

	// 开始合计
	//1. 总计
	feeTotal := &TokenTradeTotal{}
	_, err = utils.Engine_token.
		Table(this).
		Select("COUNT(trade_id) total_time, IFNULL(SUM(num+fee), 0) total_num, IFNULL(SUM(fee), 0) total_fee").
		Get(feeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//2. 今日
	todayFeeTotal := &TokenTradeTotal{}
	_, err = utils.Engine_token.
		Table(this).
		Select("COUNT(trade_id) today_total_time, IFNULL(SUM(num+fee), 0) today_total_num, IFNULL(SUM(fee), 0) today_total_fee").
		And("deal_time>=?", todayZeroUnix).
		Get(todayFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//3. 上日
	yesFeeTotal := &TokenTradeTotal{}
	_, err = utils.Engine_token.
		Table(this).
		Select("COUNT(trade_id) yesterday_total_time, IFNULL(SUM(num+fee), 0) yesterday_total_num, IFNULL(SUM(fee), 0) yesterday_total_fee").
		And("deal_time>=?", yesterdayZeroUnix).
		And("deal_time<?", yesterdayZeroUnix+24*60*60). // 小于！！！
		Get(yesFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//4. 上周同日
	lastWeekFeeTotal := &TokenTradeTotal{}
	_, err = utils.Engine_token.
		Table(this).
		Select("COUNT(trade_id) last_week_day_total_time, IFNULL(SUM(num+fee), 0) last_week_day_total_num, IFNULL(sum(fee), 0) last_week_day_total_fee").
		And("deal_time>=?", lastWeekDayZeroUnix).
		And("deal_time<?", lastWeekDayZeroUnix+24*60*60). // 小于！！！
		Get(lastWeekFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	// 合并
	feeTotal.TodayTotalTime = todayFeeTotal.TodayTotalTime
	feeTotal.TodayTotalNum = todayFeeTotal.TodayTotalNum
	feeTotal.TodayTotalFee = todayFeeTotal.TodayTotalFee

	feeTotal.YesterdayTotalTime = yesFeeTotal.YesterdayTotalTime
	feeTotal.YesterdayTotalNum = yesFeeTotal.YesterdayTotalNum
	feeTotal.YesterdayTotalFee = yesFeeTotal.YesterdayTotalFee

	feeTotal.LastWeekDayTotalTime = lastWeekFeeTotal.LastWeekDayTotalTime
	feeTotal.LastWeekDayTotalNum = lastWeekFeeTotal.LastWeekDayTotalNum
	feeTotal.LastWeekDayTotalFee = lastWeekFeeTotal.LastWeekDayTotalFee

	return feeTotal, nil
}
