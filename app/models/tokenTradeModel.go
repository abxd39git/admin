package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

//bibi委托表
type EntrustDetail struct {
	BaseModel   `xorm:"-"`
	EntrustId   string `xorm:"not null pk comment('委托记录表（委托管理）') VARCHAR(64)"`
	Uid         int64  `xorm:"not null comment('用户id') BIGINT(32)"`
	Symbol      string `xorm:"comment('队列') VARCHAR(64)"`
	TokenId     int    `xorm:"not null comment('货币id') INT(32)"`
	AllNum      int64  `xorm:"not null comment('总数量') BIGINT(20)"`
	SurplusNum  int64  `xorm:"not null comment('剩余数量') BIGINT(20)"`
	Price       int64  `xorm:"not null comment('实际平均价格(卖出价格）') BIGINT(20)"`
	Opt         int    `xorm:"not null comment('类型 卖出单1 还是买入单0') TINYINT(4)"`
	Type        int    `xorm:"comment('交易类型') TINYINT(4)"`
	OnPrice     int64  `xorm:"not null comment('委托价格(挂单价格全价格 卖出价格是扣除手续费的）') BIGINT(20)"`
	Fee         int64  `xorm:"not null comment('手续费比例') BIGINT(20)"`
	States      int    `xorm:"not null comment('0是挂单，1是部分成交,2成交， 3撤销') TINYINT(4)"`
	CreatedTime int64  `xorm:"not null comment('添加时间') BIGINT(10)"`
	Mount       int64  `xorm:"comment('总金额') BIGINT(20)"`
}

//bibi 交易表
type Trade struct {
	BaseModel    `xorm:"-"`
	TradeId      int    `xorm:"not null pk autoincr comment('交易表的id') INT(11)"`
	TradeNo      string `xorm:"comment('订单号') unique(uni_reade_no) VARCHAR(32)"`
	Uid          int64  `xorm:"comment('买家uid') index BIGINT(11)"`
	TokenId      int    `xorm:"comment('主货币id') index INT(11)"`
	TokenTradeId int    `xorm:"comment('交易币种') INT(11)"`
	TokenName    string `xorm:"not null comment('交易对 名称 例如USDT/BTC') VARCHAR(10)"`
	Price        int64  `xorm:"comment('价格') BIGINT(20)"`
	Num          int64  `xorm:"comment('数量') BIGINT(20)"`
	Money        int64  `xorm:"BIGINT(20)"`
	Fee          int64  `xorm:"comment('手续费') BIGINT(20)"`
	Opt          int    `xorm:"comment(' buy  1或sell 2') index unique(uni_reade_no) TINYINT(4)"`
	DealTime     int64  `xorm:"comment('成交时间') BIGINT(11)"`
	States       int    `xorm:"comment('0是挂单，1是部分成交,2成交， -1撤销') INT(11)"`
}

func (this *EntrustDetail) IsExist(symbol string) (bool, error) {
	engine := utils.Engine_token
	query := engine.Desc("entrust_id")
	return query.Where("symbol=?", symbol).Exist(&EntrustDetail{})
}

func (this *EntrustDetail) EvacuateOder(uid int, odid string) error {
	engine := utils.Engine_token
	//query := engine.Desc("")
	temp := fmt.Sprintf("uid=%d AND entrust_id =%s", uid, odid)
	query := engine.Where(temp)
	TempQuery := *query
	has, err := TempQuery.Exist(&EntrustDetail{})
	if err != nil {
		return err
	}
	if !has {
		return errors.New("订单不存在！！")
	}
	_, err = query.Update(&EntrustDetail{
		States: -1,
	})
	if err != nil {
		return err
	}
	return nil

}

func (this *EntrustDetail) GetTokenOrderList(page, rows, ad_id, status, start_t, uid int, symbo, trade_id string) (*ModelList, error) {
	engine := utils.Engine_token
	//
	query := engine.Desc("entrust_id")
	if trade_id != `` {
		query = query.Where("entrust_id=?", trade_id)
	}
	if symbo != `` {
		query = query.Where("symbol=?", symbo)
	}
	if ad_id != 0 {
		query = query.Where("opt=?", ad_id)
	}
	if status != 0 {
		query = query.Where("states=?", status)
	}
	if start_t != 0 {
		query = query.Where("created_time  BETWEEN ? AND ? ", start_t, start_t+86400)
	}
	if uid != 0 {
		query = query.Where("uid=?", uid)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&EntrustDetail{})
	if err != nil {
		return nil, err
	}
	offset, modelList := this.Paging(page, rows, int(count))
	list := make([]EntrustDetail, offset)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

func (this *Trade) GetTokenRecordList(page, rows, trade_id, trade_duad, ad_id, uid int, start_t string) (*ModelList, error) {
	engine := utils.Engine_token

	query := engine.Desc("uid")
	if trade_id != 0 {
		query = query.Where("trade_id=?", trade_id)
	}
	if trade_duad != 0 {
		query = query.Where("token_trade_id=?", trade_duad) //交易对
	}
	if ad_id != 0 {
		query = query.Where("opt=?", ad_id) //交易方向
	}
	if uid != 0 {
		query = query.Where("uid=?", uid)
	}
	if start_t != `` {

	}
	tempQuery := *query

	count, err := tempQuery.Count(&Trade{})
	if err != nil {
		return nil, err
	}
	offset, modelList := this.Paging(page, rows, int(count))
	list := make([]Trade, offset)
	//fmt.Printf("$$$$$$$$$$$$$$$%#v\n", rows)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

//p5-1-0-1币币交易手续费明细
/********************************
* id 兑币id
* trade_type 交易方向 1 卖 2买
* search 筛选
 */
func (this *Trade) GetFeeInfoList(page, rows, uid, opt int, date uint64, name string) (*ModelList, error) {
	engine := utils.Engine_token
	query :=engine.Desc("deal_time")
	query = query.Join("left","config_token_cny","token_id = id")

	if uid != 0 {
		query = query.Where("uid=?", uid)
	}
	if date != 0 {
		query = query.Where("deal_time BETWEEN ? AND ?", date, date+864000)
	}
	if opt != 0 {
		query = query.Where("opt=?", opt)
	}
	if name != `` {
		query = query.Where("token_name=?", name)
	}
	ValuQuery := *query
	count, err := query.Distinct("deal_time").Count(&Trade{})
	if err != nil {
		return nil, err
	}
	offset, mlist := this.Paging(page, rows, int(count))
	list := make([]Trade, offset)
	err = ValuQuery.Limit(mlist.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	//未完待续 折合成人民币
	//fmt.Println("len=",len(list))
	mlist.Items = list
	return mlist, nil
}
