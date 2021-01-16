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
type StockZCFZ struct {
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
func DownloadStockZCFZ() {
	var fileName = fmt.Sprintf(file_stockZcfzformat, reportDate)
	if file.FileExist(fileName) {
		return
	}
	file.Delete(fileName)
	file.WriteFile(fileName, "[")

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

func ReadStockZCFZ() []StockZCFZ {
	var fileName = fmt.Sprintf(file_stockZcfzformat, reportDate)
	if !file.FileExist(fileName) {
		DownloadStockZCFZ()
	}

	//读取数据
	data, err2 := file.ReadFile_v1(fileName)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockZCFZ
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	//for i := 0; i < len(stockArray2); i++ {
	//	fmt.Println(
	//		" 股票代码=", stockArray2[i].SECURITY_CODE,
	//		" 简称=", stockArray2[i].SECURITY_NAME_ABBR,
	//		" 总资产=", stockArray2[i].TOTAL_ASSETS,
	//		" 货币资金（元）=", stockArray2[i].MONETARYFUNDS,
	//		" 应收账款（元）=", stockArray2[i].ACCOUNTS_RECE,
	//		" 存货（元）=", stockArray2[i].INVENTORY,
	//		" 总负债=", stockArray2[i].TOTAL_LIABILITIES,
	//		" 应付账款（元）=", stockArray2[i].ACCOUNTS_PAYABLE,
	//		" 股东权益合计(元)=", stockArray2[i].TOTAL_EQUITY,
	//		" 总负债同比(%)=", stockArray2[i].TOTAL_LIAB_RATIO,
	//		" 资产负债率(%)=", stockArray2[i].DEBT_ASSET_RATIO,
	//	)
	//}

	return stockArray2
}
