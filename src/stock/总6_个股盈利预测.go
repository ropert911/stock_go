package stock

//盈利预测

//数据格式
//{"stockName":"长城汽车",
//"stockCode":"601633",
//"total":99,			研报数
//"rateBuy":78,		买入
//"rateIncrease":21,	增持
//"rateNeutral":0,	中性
//"rateReduce":0,		减持
//"rateSellout":0,	卖出
//"ratekanduo":99,
//"lastYearEps":"",
//"lastYearPe":"",
//"lastYearProfit":"",
//"thisYearEps":"0.6291",				今年预测每股收益
//"thisYearPe":"",
//"thisYearProfit":"5.7713761E9",
//"nextYearEps":"0.765",				明年预测每股收益
//"nextYearPe":"",
//"nextYearProfit":"7.0211103E9",
//"afterYearEps":"",
//"afterYearPe":"",
//"afterYearProfit":"",
//"lastYearActualProfit":"4.496875E9",
//"lastYearActualEps":"0.4901",			去年每股收益
//"beforeYearActualProfit":"",
//"beforeYearActualEps":"",
//"aimPriceT":"37.18",
//"aimPriceL":"10.0",
//"updateTime":"2021-01-06 19:30:14.000","hyBK":"481",
//"gnBK":["682","707","718","802","900","813","815","815","815","815","815","815","816","817","817","817","821","845","867","879","499","500","567","570","574","596","596","596","596","596","596","596","596","612"],
//"dyBK":"199003",
//"market":["HU"],
//"total_1":10,
//"rateBuy_1":8,
//"rateIncrease_1":2,
//"rateNeutral_1":0,
//"rateReduce_1":0,
//"rateSellout_1":0,
//"total_3":39,
//"rateBuy_3":32,
//"rateIncrease_3":7,
//"rateNeutral_3":0,
//"rateReduce_3":0,
//"rateSellout_3":0,
//"total_12":173,
//"rateBuy_12":128,
//"rateIncrease_12":45,
//"rateNeutral_12":0,
//"rateReduce_12":0,
//"rateSellout_12":0},

//type StockYlyc struct {
//	STOCKCODE         string  //代码
//	STOCKNAME         string  //股票名称
//	LASTYEARACTUALEPS string  //去年每股收益
//	THISYEAREPS       string  //今年预测每股收益
//	NEXTYEAREPS       string  //明年预测每股收益
//	EGR               float64 //预期收益增长
//	TOTAL             float32 //研报数
//	RATEBUY           float32 //研报-买入数
//	RATEINCREASE      float32 //研报-增持数
//	RATENEUTRAl       float32 //研报-中性数
//	RATEREDUCE        float32 //研报-减持数
//	RATESELLOUT       float32 //研报-卖出数
//}

//盈利预测
//func DownloadStockYlyc() {
//	var fileName = fmt.Sprintf(file_stockYlycformat, TradeData)
//	if file.FileExist(fileName) {
//		return
//	}
//	file.Delete(fileName)
//	file.WriteFile(fileName, "[")
//
//	var urlFormat = `http://reportapi.eastmoney.com/report/predic?pageNo=%d&pageSize=%d`
//	var index = 1
//	var first = true
//	for {
//		var url = fmt.Sprintf(urlFormat, index, 100)
//		fmt.Println(url)
//		index++
//
//		//获取数据
//		content, err := http.HttpGet(url)
//		if nil != err {
//			fmt.Println("http get failed url=", url, " error=", err)
//			break
//		}
//
//		//得到真实内容
//		data := *content
//		var start = strings.IndexAny(data, "[")
//		var end = strings.LastIndex(data, "]")
//		if -1 == start || -1 == end {
//			break
//		}
//		data = data[start+1 : end]
//
//		//写内容
//		if len(data) > 10 {
//			if !first {
//				file.AppendFile(fileName, ",")
//			} else {
//				first = false
//			}
//			file.AppendFile(fileName, strings.ToUpper(data))
//		} else {
//			break
//		}
//	}
//	file.AppendFile(fileName, "]")
//}

//func ReadStockYlyc() []StockYlyc {
//	var fileName = fmt.Sprintf(file_stockYlycformat, TradeData)
//	if !file.FileExist(fileName) {
//		DownloadStockYlyc()
//	}
//
//	//读取数据
//	data, err2 := file.ReadFile_v1(fileName)
//	if nil != err2 {
//		fmt.Println("error read file ", err2)
//		return nil
//	}
//
//	var stockArray2 []StockYlyc
//	err := json.Unmarshal([]byte(data), &stockArray2)
//	if nil != err {
//		fmt.Println("22222 json unmarshal failed!!!!", err)
//		return nil
//	}
//
//	for i := 0; i < len(stockArray2); i++ {
//		stockArray2[i].EGR = 100 * (ToFloat(stockArray2[i].NEXTYEAREPS) - ToFloat(stockArray2[i].THISYEAREPS)) / ToFloat(stockArray2[i].THISYEAREPS)
//		//fmt.Println(
//		//	" 股票代码=", stockArray2[i].STOCKCODE,
//		//	" 简称=", stockArray2[i].STOCKNAME,
//		//	" 去年每股收益=", stockArray2[i].LASTYEARACTUALEPS,
//		//	" 预测今年每股收益=", stockArray2[i].THISYEAREPS,
//		//	" 预测明年每股收益=", stockArray2[i].NEXTYEAREPS,
//		//	" 预期收益增长=", stockArray2[i].EGR,
//		//)
//	}
//
//	return stockArray2
//}
