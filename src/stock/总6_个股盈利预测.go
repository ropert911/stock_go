package stock

import (
	"encoding/json"
	"fmt"
	"strings"
	"util/file"
	"util/http"
)

type StockYLYC struct {
	STOCKCODE         string  //代码
	STOCKNAME         string  //股票名称
	LASTYEARACTUALEPS string  //去年每股收益
	THISYEAREPS       string  //今年预测每股收益
	NEXTYEAREPS       string  //明年预测每股收益
	EGR               float64 //预期收益增长
	TOTAL             float32 //研报数
	RATEBUY           float32 //研报-买入数
	RATEINCREASE      float32 //研报-增持数
	RATENEUTRAl       float32 //研报-中性数
	RATEREDUCE        float32 //研报-减持数
	RATESELLOUT       float32 //研报-卖出数
}

func DownloadStockYLYC() {
	var fileName = fmt.Sprintf(file_stockYlycformat, TradeData)
	if file.FileExist(fileName) {
		return
	}
	file.Delete(fileName)
	file.WriteFile(fileName, "[")

	var urlFormat = `http://reportapi.eastmoney.com/report/predic?pageNo=%d&pageSize=%d`
	var index = 1
	var first = true
	for {
		var url = fmt.Sprintf(urlFormat, index, 100)
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
		var end = strings.LastIndex(data, "]")
		if -1 == start || -1 == end {
			break
		}
		data = data[start+1 : end]

		//写内容
		if len(data) > 10 {
			if !first {
				file.AppendFile(fileName, ",")
			} else {
				first = false
			}
			file.AppendFile(fileName, strings.ToUpper(data))
		} else {
			break
		}
	}
	file.AppendFile(fileName, "]")
}

func ReadStockYLYC() []StockYLYC {
	var fileName = fmt.Sprintf(file_stockYlycformat, TradeData)
	if !file.FileExist(fileName) {
		DownloadStockYLYC()
	}

	//读取数据
	data, err2 := file.ReadFile_v1(fileName)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockYLYC
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	//for i := 0; i < len(stockArray2); i++ {
	//	stockArray2[i].EGR = 100 * (ToFloat(stockArray2[i].NEXTYEAREPS) - ToFloat(stockArray2[i].THISYEAREPS)) / ToFloat(stockArray2[i].THISYEAREPS)
	//	fmt.Println(
	//		" 股票代码=", stockArray2[i].STOCKCODE,
	//		" 简称=", stockArray2[i].STOCKNAME,
	//		" 去年每股收益=", stockArray2[i].LASTYEARACTUALEPS,
	//		" 预测今年每股收益=", stockArray2[i].THISYEAREPS,
	//		" 预测明年每股收益=", stockArray2[i].NEXTYEAREPS,
	//		" 预期收益增长=", stockArray2[i].EGR,
	//	)
	//}

	return stockArray2
}
