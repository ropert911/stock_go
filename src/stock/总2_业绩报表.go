package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"util/file"
	"util/http"
)

//业绩报表 http://data.eastmoney.com/bbsj/202006/yjbb.html
//数据格式
//{
//"SECURITY_CODE":"300344",				股票代码
//"SECURITY_NAME_ABBR":"太空智造",		简称
//"TRADE_MARKET_CODE":"069001002002",
//"TRADE_MARKET":"深交所创业板",
//"SECURITY_TYPE_CODE":"058001001",
//"SECURITY_TYPE":"A股",
//"UPDATE_DATE":"2020-10-01 00:00:00",
//"REPORTDATE":"2020-06-30 00:00:00",
//"BASIC_EPS":-0.0874,					每股收益（元）
//"DEDUCT_BASIC_EPS":-0.0893,
//"TOTAL_OPERATE_INCOME":76755948.04,		营业收入
//"PARENT_NETPROFIT":-42325660,			净利润
//"WEIGHTAVG_ROE":-6.55,					净资产收益率%
//"YSTZ":-65.0955787682,					同比增长%
//"SJLTZ":-444.03,						净利润环比增长
//"BPS":1.230604809575,					每股净资产
//"MGJYXJJE":-0.027835091109,
//"XSMLL":23.9988262283,					销售毛利润%
//"YSHZ":47.8119,							季度环比增长
//"SJLHZ":-11.9622,						净利润季度环比增长
//"ASSIGNDSCRPT":"不分配不转增",			利润分配
//"PAYYEAR":"2020",
//"PUBLISHNAME":"软件服务",				所处行业
//"ZXGXL":null,
//"NOTICE_DATE":"2020-08-28 00:00:00",
//"ORG_CODE":"10119359",
//"TRADE_MARKET_ZJG":"0202",
//"ISNEW":"1",
//"QDATE":"2020Q2",
//"DATATYPE":"2020年 半年报",
//"DATAYEAR":"2020",
//"DATEMMDD":"半年报",
//"EITIME":"2020-08-27 16:08:42"
//}

type StockYjbb struct {
	SECURITY_CODE        string  //股票编码
	SECURITY_NAME_ABBR   string  //股票名
	TOTAL_OPERATE_INCOME float32 //营业收入
	YSTZ                 float32 //营业收入同比增长%
	YSHZ                 float32 //营业收入季度环比增长
	PARENT_NETPROFIT     float32 //净利润
	SJLTZ                float32 //净利润同比增长
	SJLHZ                float32 //净利润季度环比增长
	XSMLL                float32 //销售毛利润%
	BASIC_EPS            float32 //每股收益（元）
	BPS                  float32 //每股净资产
	WEIGHTAVG_ROE        float32 //净资产收益率%
}

//业绩报表
func DownloadStockYjbb() {
	var fileName = fmt.Sprintf(file_stockYjbbformat, reportDate)
	if file.FileExist(fileName) {
		return
	}
	file.WriteFile(fileName, "[")

	var urlFormat = `http://datacenter.eastmoney.com/api/data/get?type=RPT_LICO_FN_CPD&sty=ALL&p=%d&ps=%d&st=UPDATE_DATE,SECURITY_CODE&sr=-1,-1&var=sefpfESI&filter=(SECURITY_TYPE_CODE=058001001)(REPORTDATE='%s')&rt=%d`
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

func ReadStockYjbb() []StockYjbb {
	var fileName = fmt.Sprintf(file_stockYjbbformat, reportDate)
	if !file.FileExist(fileName) {
		DownloadStockYjbb()
	}

	//读取数据
	data, err2 := file.ReadFile_v1(fileName)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockYjbb
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	for i := 0; i < len(stockArray2); i++ {
		//fmt.Printf("%6s %6s 收入=%f 同比=%f 环比=%f 净利润=%f 同比=%f 环比=%f 销售毛利润=%f%% 每股收益（元）=%f 每股净资产=%f元 净资产收益率%%=%f\n",
		//	stockArray2[i].SECURITY_CODE,
		//	stockArray2[i].SECURITY_NAME_ABBR,
		//	stockArray2[i].TOTAL_OPERATE_INCOME,
		//	stockArray2[i].YSTZ,
		//	stockArray2[i].YSHZ,
		//	stockArray2[i].PARENT_NETPROFIT,
		//	stockArray2[i].SJLTZ,
		//	stockArray2[i].SJLHZ,
		//	stockArray2[i].XSMLL,
		//	stockArray2[i].BASIC_EPS,
		//	stockArray2[i].BPS,
		//	stockArray2[i].WEIGHTAVG_ROE,
		//)
	}

	return stockArray2
}
