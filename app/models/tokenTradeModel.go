package models

import (
	"admin/utils"
	"fmt"
	"time"
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
	FeeCny    int64  `xorm:"comment('手续费折合CNY') BIGINT(20)"`
	TotalCny  int64  `xorm:"comment('总交易额折合CNY') BIGINT(20)"`
	EntrustId string `xorm:"VARCHAR(32)"`
	TokenAdmissionId  int `xorm:"comment('入账货币id') TINYINT(4)" json:"token_admission_id"`
}

type TradeReturn struct {
	Trade          `xorm:"extends"`
	TradeNum    	int64 `xorm:" comment('已成交') BIGINT(20)" json:"trade_num"`
	AllNum         int64   `xorm:"not null comment('总数量') BIGINT(20)"`  //总数
	SurplusNum     int64  `xorm:"not null comment('剩余数量') BIGINT(20)"` //余数
	ToAccountNum   string `xorm:"-" json:"to_account_num"`//实际到账数量
	FinishCount    string `xorm:"-" json:"finish_count"`               //已成
	FeeTrue        string `xorm:"-" json:"fee_true"`
	AllNumTrue     string `xorm:"-" json:"all_num_true"`
	SurplusNumTrue string `xorm:"-" json:"surplus_num_true"`
	PriceTrue      string `xorm:"-" json:"price_true"`
	TokenName      string `xorm:"-" json:"token_name"`
}

func (t *TradeReturn) TableName() string {
	return "trade"
}

type TradeEx struct {
	Trade          `xorm:"extends"`
	ConfigTokenCny `xorm:"extends"`
	TotalTrue      float64 //交易总额
	FeeTrue        float64 //交易手续费
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

type  DayCount struct {
	TokenId int64 `json:"token_id"`//货币id
	Total int64 `json:"total"`//数量
	TotalCny int64 `json:"total_cny"`//折合
	FeeTotal int64 `json:"fee_total"`
	FeeTotalCny int64 `json:"fee_total_cny"`
	Date int64 `xorm:"-" json:"date"` //日期
}

func (t *DayCount) TableName() string {
	return "trade"
}

func (t*Trade) Get(tid int , bt,et,opt int64)(*DayCount,error){
	engine :=utils.Engine_token
	dc:=new(DayCount)
	_,err:=engine.Select(" token_admission_id token_id ,FROM_UNIXTIME(deal_time,'%Y-%m-%d %H:%i:%s') date, SUM(num) total ,SUM(total_cny) total_cny,SUM(fee) fee_total,SUM(fee_cny) fee_total_cny ").Where("token_admission_id=? and opt=?  and deal_time between ? and ? ",tid,opt,bt,et).Get(dc)
	if err!=nil{
		utils.AdminLog.Errorln(err.Error())
		return nil,err
	}
	dc.Date =bt
	return dc,nil

}

//func (this *Trade) TotalTotalTradeList(page, rows int, date uint64) (*ModelList, error) {
//	fmt.Println("bibi 交易手续费一天汇总")
//	engine := utils.Engine_token
//	query := engine.Desc("deal_time")
//	query = query.Join("left", "config_token_cny p", "trade.token_id = p.token_id")
//	query = query.GroupBy("deal_time")
//	if date != 0 {
//		temp := date / 1000
//		query = query.Where("left(deal_time,7)=?", temp)
//	}
//	tempQuery := *query
//	buyQuery := *query
//	sellQuery := *query
//	count, err := tempQuery.Count(&Trade{})
//	if err != nil {
//		return nil, err
//	}
//	offset, mList := this.Paging(page, rows, int(count))
//	//买入总额
//	buyList := make([]TradeEx, 0)
//	err = buyQuery.Where("opt=1").Limit(mList.PageSize, offset).Find(&buyList)
//	if err != nil {
//		return nil, err
//	}
//	//卖出总额
//	sellList := make([]TradeEx, 0)
//	err = sellQuery.Where("opt=2").Limit(mList.PageSize, offset).Find(&sellList)
//	//var totalBuy uint64
//	//var totalSell uint64
//	//买卖总金额
//	totalDateList := make([]map[int64]*TotalTradeCNY, 0)
//	dateMap := make(map[int64]*TotalTradeCNY, 0)
//	for _, v := range buyList {
//		key := v.DealTime / 1000
//		for i, _ := range totalDateList {
//			if _, ok := totalDateList[i][key]; !ok {
//				dateMap[key] = &TotalTradeCNY{Date: v.DealTime}
//				totalDateList = append(totalDateList, dateMap)
//			}
//			strBuy := this.Int64MulInt64By8BitString(v.Num, v.ConfigTokenCny.Price)
//			buy, err := strconv.ParseUint(strBuy, 10, 64)
//			if err != nil {
//				continue
//			}
//			totalDateList[i][key].Buy += buy
//
//		}
//
//	}
//
//	for _, v := range sellList {
//		key := v.DealTime / 1000
//		for i, _ := range totalDateList {
//			if _, ok := totalDateList[i][key]; !ok {
//				dateMap[key] = &TotalTradeCNY{Date: v.DealTime}
//				totalDateList = append(totalDateList, dateMap)
//			}
//			strSell := this.Int64MulInt64By8BitString(v.Num, v.ConfigTokenCny.Price)
//			sell, err := strconv.ParseUint(strSell, 10, 64)
//			if err != nil {
//				continue
//			}
//			totalDateList[i][key].Sell += sell
//		}
//
//	}
//	mList.Items = totalDateList
//	return mList, nil
//}

func (this *Trade) GetTokenRecordList(page, rows, opt, uid int, bt, et uint64, name string) (*ModelList, error) {
	engine := utils.Engine_token
	 fmt.Println("这里到了没有啊 ")
	query := engine.Desc("t.entrust_id")
	query = query.Alias("t").Join("left", "entrust_detail e", "e.entrust_id= t.entrust_id")
	if name!=``{
		query = query.Where("(e.states=2 or e.states=1) and e.symbol=?", name) //交易对
	}
	tm:=time.Now().Unix()
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

	} else{
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
	fmt.Println("list",list)
	for i, v := range list {
		//allNum, surPlusNUm := this.SubductionZeroMethodInt64(v.AllNum, v.SurplusNum)
		list[i].AllNumTrue = convert.Int64ToStringBy8Bit(v.AllNum)
		list[i].SurplusNumTrue = convert.Int64ToStringBy8Bit(v.SurplusNum)
		list[i].FinishCount = convert.Int64ToStringBy8Bit(v.Num)
		list[i].ToAccountNum = convert.Int64ToStringBy8Bit(v.TradeNum)
		list[i].FeeTrue = convert.Int64ToStringBy8Bit(v.Fee)
		list[i].PriceTrue = convert.Int64ToStringBy8Bit(v.Price)
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
	query = query.Join("left", "config_token_cny p", "trade.token_id = p.token_id")
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
	for i, v := range list {
		list[i].TotalTrue = this.Int64ToFloat64By8Bit(v.TotalCny)
		list[i].FeeTrue = this.Int64ToFloat64By8Bit(v.FeeCny)
	}
	mlist.Items = list
	return mlist, nil
}
