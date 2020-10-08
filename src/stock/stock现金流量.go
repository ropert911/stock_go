package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"util/file"
	"util/http"
)

//现金流量表

//数据格式
//{
//"SECUCODE":"688013.SH",
//"SECURITY_CODE":"688013",					代码
//"INDUSTRY_CODE":"016047",
//"ORG_CODE":"10000071414",
//"SECURITY_NAME_ABBR":"天臣医疗",			简称
//"INDUSTRY_NAME":"医疗行业",					行业
//"MARKET":"kcb",
//"SECURITY_TYPE_CODE":"058001001",
//"TRADE_MARKET_CODE":"069001001006",
//"DATE_TYPE_CODE":"002",
//"REPORT_TYPE_CODE":"001",
//"DATA_STATE":"2",
//"NOTICE_DATE":"2020-09-09 00:00:00",		公告日期
//"REPORT_DATE":"2020-06-30 00:00:00",		报表日期
//"NETCASH_OPERATE":15470566.93,				经营现金流净额
//"NETCASH_OPERATE_RATIO":115.6940896681,		经营现金流-净现金流占比
//"SALES_SERVICES":76531523.21,
//"SALES_SERVICES_RATIO":572.3284058531,
//"PAY_STAFF_CASH":18022913.13,
//"PSC_RATIO":134.7813908292,
//"NETCASH_INVEST":-3349126.05,				投资现金流净额
//"NETCASH_INVEST_RATIO":-25.0458881883,		投资现金流-净现金流占比
//"RECEIVE_INVEST_INCOME":29006.85,
//"RII_RATIO":0.2169229557,
//"CONSTRUCT_LONG_ASSET":3466385.3,
//"CLA_RATIO":25.9227921987,
//"NETCASH_FINANCE":null,						融资现金流净额
//"NETCASH_FINANCE_RATIO":null,				融资现金流-净现金流占比
//"CCE_ADD":13371959.6,						净现金流（总现金流）
//"CCE_ADD_RATIO":197.7202906035,				净现金流同比增长
//"CUSTOMER_DEPOSIT_ADD":null,
//"CDA_RATIO":null,
//"DEPOSIT_IOFI_OTHER":null,
//"DIO_RATIO":null,"LOAN_ADVANCE_ADD":null,
//"LAA_RATIO":null,
//"RECEIVE_INTEREST_COMMISSION":null,
//"RIC_RATIO":null,"INVEST_PAY_CASH":null,
//"IPC_RATIO":null,
//"BEGIN_CCE":null,
//"BEGIN_CCE_RATIO":null,
//"END_CCE":null,
//"END_CCE_RATIO":null,
//"RECEIVE_ORIGIC_PREMIUM":null,
//"ROP_RATIO":null,
//"PAY_ORIGIC_COMPENSATE":null,
//"POC_RATIO":null
//}

type StockXjll struct {
	SECURITY_CODE         string  //代码
	SECURITY_NAME_ABBR    string  //简称
	INDUSTRY_NAME         string  //行业
	REPORT_DATE           string  //报表日期
	NETCASH_OPERATE       float32 //经营现金流净额
	NETCASH_OPERATE_RATIO float32 //经营现金流-净现金流占比
	NETCASH_INVEST        float32 //投资现金流净额
	NETCASH_INVEST_RATIO  float32 //投资现金流-净现金流占比
	NETCASH_FINANCE       float32 //融资现金流净额
	NETCASH_FINANCE_RATIO float32 //融资现金流-净现金流占比
	CCE_ADD               float32 //净现金流（总现金流）
	CCE_ADD_RATIO         float32 //净现金流同比增长
}

//现金流量
func DownloadStockXjll() {
	file.Delete(file_stockXjll)
	file.WriteFile(file_stockXjll, "[")

	var urlFormat = `http://datacenter.eastmoney.com/api/data/get?type=RPT_DMSK_FN_CASHFLOW&sty=ALL&p=%d&ps=%d&st=NOTICE_DATE,SECURITY_CODE&sr=-1,-1&var=FPETiANd&filter=(REPORT_DATE='%s')&rt=%d`
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
				file.AppendFile(file_stockXjll, ",")
			} else {
				first = false
			}
			file.AppendFile(file_stockXjll, data)
		} else {
			break
		}
	}
	file.AppendFile(file_stockXjll, "]")
}

func ReadStockXjll() []StockXjll {
	//读取数据
	data, err2 := file.ReadFile_v1(file_stockXjll)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockXjll
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	for i := 0; i < len(stockArray2); i++ {
		//fmt.Println(
		//	" 股票代码=", stockArray2[i].SECURITY_CODE,
		//	" 简称=", stockArray2[i].SECURITY_NAME_ABBR,
		//	" 行业=", stockArray2[i].INDUSTRY_NAME,
		//	" 报表日期=", stockArray2[i].REPORT_DATE,
		//	" 经营现金流净额=", stockArray2[i].NETCASH_OPERATE,
		//	" 经营现金流-净现金流占比=", stockArray2[i].NETCASH_OPERATE_RATIO,
		//	" 投资现金流净额=", stockArray2[i].NETCASH_INVEST,
		//	" 投资现金流-净现金流占比=", stockArray2[i].NETCASH_INVEST_RATIO,
		//	" 融资现金流净额=", stockArray2[i].NETCASH_FINANCE,
		//	" 融资现金流-净现金流占比=", stockArray2[i].NETCASH_FINANCE_RATIO,
		//	" 净现金流（总现金流）=", stockArray2[i].CCE_ADD,
		//	" 净现金流同比增长=", stockArray2[i].CCE_ADD_RATIO,
		//)
	}

	return stockArray2
}
