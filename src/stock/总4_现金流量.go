package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"util/file"
	"util/http"
)

type StockXJLL struct {
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

func DownloadStockXJLL() {
	var fileName = fmt.Sprintf(file_stockXjllformat, reportDate)
	if file.FileExist(fileName) {
		return
	}
	file.Delete(fileName)
	file.WriteFile(fileName, "[")

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

func ReadStockXJLL() []StockXJLL {
	var fileName = fmt.Sprintf(file_stockXjllformat, reportDate)
	if !file.FileExist(fileName) {
		DownloadStockXJLL()
	}

	//读取数据
	data, err2 := file.ReadFile_v1(fileName)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockXJLL
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	//for i := 0; i < len(stockArray2); i++ {
	//	fmt.Println(
	//		" 股票代码=", stockArray2[i].SECURITY_CODE,
	//		" 简称=", stockArray2[i].SECURITY_NAME_ABBR,
	//		" 行业=", stockArray2[i].INDUSTRY_NAME,
	//		" 报表日期=", stockArray2[i].REPORT_DATE,
	//		" 经营现金流净额=", stockArray2[i].NETCASH_OPERATE,
	//		" 经营现金流-净现金流占比=", stockArray2[i].NETCASH_OPERATE_RATIO,
	//		" 投资现金流净额=", stockArray2[i].NETCASH_INVEST,
	//		" 投资现金流-净现金流占比=", stockArray2[i].NETCASH_INVEST_RATIO,
	//		" 融资现金流净额=", stockArray2[i].NETCASH_FINANCE,
	//		" 融资现金流-净现金流占比=", stockArray2[i].NETCASH_FINANCE_RATIO,
	//		" 净现金流（总现金流）=", stockArray2[i].CCE_ADD,
	//		" 净现金流同比增长=", stockArray2[i].CCE_ADD_RATIO,
	//	)
	//}

	return stockArray2
}
