package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"util/file"
	"util/http"
)

//资产负债表

//数据格式
//{
//"SECUCODE":"300344.SZ",
//"SECURITY_CODE":"300344",			股票代码
//"INDUSTRY_CODE":"016057",
//"ORG_CODE":"10119359",
//"SECURITY_NAME_ABBR":"太空智造",	简称
//"INDUSTRY_NAME":"软件服务",
//"MARKET":"cyb",
//"SECURITY_TYPE_CODE":"058001001",
//"TRADE_MARKET_CODE":"069001002002",
//"DATE_TYPE_CODE":"002",
//"REPORT_TYPE_CODE":"001",
//"DATA_STATE":"2",
//"NOTICE_DATE":"2020-10-01 00:00:00",
//"REPORT_DATE":"2020-06-30 00:00:00",
//"TOTAL_ASSETS":983875933.36,			总资产
//"FIXED_ASSET":25345642.2,				固定资产
//"MONETARYFUNDS":9155387.31,				货币资金（元）
//"MONETARYFUNDS_RATIO":-54.0070712048,
//"ACCOUNTS_RECE":211213607.22,			应收账款（元）
//"ACCOUNTS_RECE_RATIO":-4.8736053788,
//"INVENTORY":95395487.02,				存货（元）
//"INVENTORY_RATIO":-21.535827604,
//"TOTAL_LIABILITIES":371193102.01,		总负债
//"ACCOUNTS_PAYABLE":49758808.91,				应付账款（元）
//"ACCOUNTS_PAYABLE_RATIO":-33.4888620932,
//"ADVANCE_RECEIVABLES":null,					预收账款（元）
//"ADVANCE_RECEIVABLES_RATIO":null,
//"TOTAL_EQUITY":612682831.35,				股东权益合计(元)
//"TOTAL_EQUITY_RATIO":3.0084567925,
//"TOTAL_ASSETS_RATIO":-4.4085057144,
//"TOTAL_LIAB_RATIO":-14.5625132416,			总负债同比(%)
//"CURRENT_RATIO":101.1115771215,
//"DEBT_ASSET_RATIO":37.7276330708,			资产负债率(%)
//"CASH_DEPOSIT_PBC":null,
//"CDP_RATIO":null,
//"LOAN_ADVANCE":null,
//"LOAN_ADVANCE_RATIO":null,
//"AVAILABLE_SALE_FINASSET":null,
//"ASF_RATIO":null,
//"LOAN_PBC":null,
//"LOAN_PBC_RATIO":null,
//"ACCEPT_DEPOSIT":null,
//"ACCEPT_DEPOSIT_RATIO":null,
//"SELL_REPO_FINASSET":null,
//"SRF_RATIO":null,
//"SETTLE_EXCESS_RESERVE":null,
//"SER_RATIO":null,
//"BORROW_FUND":null,
//"BORROW_FUND_RATIO":null,
//"AGENT_TRADE_SECURITY":null,
//"ATS_RATIO":null,
//"PREMIUM_RECE":null,
//"PREMIUM_RECE_RATIO":null,
//"SHORT_LOAN":null,
//"SHORT_LOAN_RATIO":null,
//"ADVANCE_PREMIUM":null,
//"ADVANCE_PREMIUM_RATIO":null
//},

type StockZcfz struct {
	SECURITY_CODE      string  //股票代码
	SECURITY_NAME_ABBR string  //简称
	TOTAL_ASSETS       float32 //总资产
	MONETARYFUNDS      float32 //货币资金（元）
	ACCOUNTS_RECE      float32 //应收账款（元）
	INVENTORY          float32 //存货（元）
	TOTAL_LIABILITIES  float32 //总负债
	ACCOUNTS_PAYABLE   float32 //应付账款（元）
	TOTAL_EQUITY       float32 //股东权益合计(元)
	TOTAL_LIAB_RATIO   float32 //总负债同比(%)
	DEBT_ASSET_RATIO   float32 //资产负债率(%)

}

//资产负债
func DownloadStockZcfz() {
	file.Delete(file_stockZcfz)
	file.WriteFile(file_stockZcfz, "[")

	var urlFormat = `http://datacenter.eastmoney.com/api/data/get?type=RPT_DMSK_FN_BALANCE&sty=ALL&p=%d&ps=%d&st=NOTICE_DATE,SECURITY_CODE&sr=-1,-1&var=TjVwCDHJ&filter=(REPORT_DATE='%s')&rt=%d`
	var index = 1
	var first = true
	for {
		var url = fmt.Sprintf(urlFormat,
			index, //第几页
			100,   //每页条数
			reportDate,
			rand.Intn(899999)+100000)
		fmt.Println(url)
		index++

		//获取数据
		content, err := http.HttpGet(url)
		if nil != err {
			fmt.Println("http get failed url=", url, " error=", err)
			break
		}

		//得到真实内容
		data := *content
		var start = strings.IndexAny(data, "[")
		var end = strings.IndexAny(data, "]")
		if -1 == start || -1 == end {
			break
		}
		data = data[start+1 : end]
		//fmt.Println(data)

		//写内容
		if len(data) > 10 {
			if !first {
				file.AppendFile(file_stockZcfz, ",")
			} else {
				first = false
			}
			file.AppendFile(file_stockZcfz, data)
		} else {
			break
		}
	}
	file.AppendFile(file_stockZcfz, "]")
}

func ReadStockZcfz() []StockZcfz {
	//读取数据
	data, err2 := file.ReadFile_v1(file_stockZcfz)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockZcfz
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	for i := 0; i < len(stockArray2); i++ {
		//fmt.Println(
		//	" 股票代码=", stockArray2[i].SECURITY_CODE,
		//	" 简称=", stockArray2[i].SECURITY_NAME_ABBR,
		//	" 总资产=", stockArray2[i].TOTAL_ASSETS,
		//	" 货币资金（元）=", stockArray2[i].MONETARYFUNDS,
		//	" 应收账款（元）=", stockArray2[i].ACCOUNTS_RECE,
		//	" 存货（元）=", stockArray2[i].INVENTORY,
		//	" 总负债=", stockArray2[i].TOTAL_LIABILITIES,
		//	" 应付账款（元）=", stockArray2[i].ACCOUNTS_PAYABLE,
		//	" 股东权益合计(元)=", stockArray2[i].TOTAL_EQUITY,
		//	" 总负债同比(%)=", stockArray2[i].TOTAL_LIAB_RATIO,
		//	" 资产负债率(%)=", stockArray2[i].DEBT_ASSET_RATIO,
		//)
	}

	return stockArray2
}
