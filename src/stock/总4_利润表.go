package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"util/file"
	"util/http"
)

//数据格式
//{"SECUCODE":"688013.SH",
//"SECURITY_CODE":"688013",				代码
//"INDUSTRY_CODE":"016047",
//"ORG_CODE":"10000071414",
//"SECURITY_NAME_ABBR":"天臣医疗",		简称
//"INDUSTRY_NAME":"医疗行业",
//"MARKET":"kcb",
//"SECURITY_TYPE_CODE":"058001001",
//"TRADE_MARKET_CODE":"069001001006",
//"DATE_TYPE_CODE":"002",
//"REPORT_TYPE_CODE":"001",
//"DATA_STATE":"2",
//"NOTICE_DATE":"2020-09-09 00:00:00",
//"REPORT_DATE":"2020-06-30 00:00:00",
//"PARENT_NETPROFIT":16291105.44,			净利润
//"TOTAL_OPERATE_INCOME":70142426.22,		营业总收
//"TOTAL_OPERATE_COST":54966269.71,		营业收支出
//"TOE_RATIO":-4.4515461177,
//"OPERATE_COST":28928635.76,				营业支持
//"OPERATE_EXPENSE":28928635.76,
//"OPERATE_EXPENSE_RATIO":-12.1617147768,
//"SALE_EXPENSE":11471822.03,				销售费用
//"MANAGE_EXPENSE":8683771.1,				管理费用
//"FINANCE_EXPENSE":-1236201.38,			财务费用
//"OPERATE_PROFIT":15198280.27,			营业利润
//"TOTAL_PROFIT":18588710.1,				利润总额
//"INCOME_TAX":2297604.66,
//"OPERATE_INCOME":null,
//"INTEREST_NI":null,
//"INTEREST_NI_RATIO":null,
//"FEE_COMMISSION_NI":null,
//"FCN_RATIO":null,
//"OPERATE_TAX_ADD":930243.23,
//"MANAGE_EXPENSE_BANK":null,
//"FCN_CALCULATE":null,
//"INTEREST_NI_CALCULATE":null,
//"EARNED_PREMIUM":null,
//"EARNED_PREMIUM_RATIO":null,
//"INVEST_INCOME":null,
//"SURRENDER_VALUE":null,
//"COMPENSATE_EXPENSE":null,
//"TOI_RATIO":-13.8440888338,				营业总收同比
//"OPERATE_PROFIT_RATIO":-35.052892917803,
//"PARENT_NETPROFIT_RATIO":-24.12,		净利润同比
//"DEDUCT_PARENT_NETPROFIT":13309634.13,	扣非净利润
//"DPN_RATIO":-34.709659973599
//},

type StockLrb struct {
	SECURITY_CODE           string  //代码
	SECURITY_NAME_ABBR      string  //简称
	PARENT_NETPROFIT        float32 //净利润
	TOTAL_OPERATE_INCOME    float32 //营业总收入
	TOI_RATIO               float32 //营业总收入同比
	TOTAL_OPERATE_COST      float32 //营业收支出
	OPERATE_COST            float32 //营业支持
	SALE_EXPENSE            float32 //销售费用
	MANAGE_EXPENSE          float32 //管理费用
	FINANCE_EXPENSE         float32 //财务费用
	OPERATE_PROFIT          float32 //营业利润
	TOTAL_PROFIT            float32 //利润总额
	PARENT_NETPROFIT_RATIO  float32 //净利润同比
	DEDUCT_PARENT_NETPROFIT float32 //扣非净利润
}

func DownloadStockLrb() {
	var fileName = fmt.Sprintf(file_stockLrbformat, reportDate)
	if file.FileExist(fileName) {
		return
	}

	file.Delete(fileName)
	file.WriteFile(fileName, "[")

	var urlFormat = `http://datacenter.eastmoney.com/api/data/get?type=RPT_DMSK_FN_INCOME&sty=ALL&p=%d&ps=%d&st=NOTICE_DATE,SECURITY_CODE&sr=-1,-1&var=lkQKYzdA&filter=(REPORT_DATE='%s')&rt=%d`
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
				file.AppendFile(fileName, ",")
			} else {
				first = false
			}
			file.AppendFile(fileName, data)
		} else {
			break
		}
	}
	file.AppendFile(fileName, "]")
}

func ReadStockLrb() []StockLrb {
	var fileName = fmt.Sprintf(file_stockLrbformat, reportDate)
	if !file.FileExist(fileName) {
		DownloadStockLrb()
	}

	//读取数据
	data, err2 := file.ReadFile_v1(fileName)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockLrb
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	for i := 0; i < len(stockArray2); i++ {
		//fmt.Println(
		//	" 代码=", stockArray2[i].SECURITY_CODE,
		//	" 简称=", stockArray2[i].SECURITY_NAME_ABBR,
		//	" 净利润=", stockArray2[i].PARENT_NETPROFIT,
		//	" 营业总收入=", stockArray2[i].TOTAL_OPERATE_INCOME,
		//	" 营业总收入同比=", stockArray2[i].TOI_RATIO,
		//	" 营业收支出=", stockArray2[i].TOTAL_OPERATE_COST,
		//	" 营业支出=", stockArray2[i].OPERATE_COST,
		//	" 销售费用=", stockArray2[i].SALE_EXPENSE,
		//	" 管理费用=", stockArray2[i].MANAGE_EXPENSE,
		//	" 财务费用=", stockArray2[i].FINANCE_EXPENSE,
		//	" 营业利润=", stockArray2[i].OPERATE_PROFIT,
		//	" 利润总额=", stockArray2[i].TOTAL_PROFIT,
		//	" 净利润同比=", stockArray2[i].PARENT_NETPROFIT_RATIO,
		//	" 扣非净利润=", stockArray2[i].DEDUCT_PARENT_NETPROFIT,
		//)
	}

	return stockArray2
}
