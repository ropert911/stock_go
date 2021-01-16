package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"util/file"
	"util/http"
)

type StockLRB struct {
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

func DownloadStockLRB() {
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

func ReadStockLRB() []StockLRB {
	var fileName = fmt.Sprintf(file_stockLrbformat, reportDate)
	if !file.FileExist(fileName) {
		DownloadStockLRB()
	}

	//读取数据
	data, err2 := file.ReadFile_v1(fileName)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockLRB
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	//for i := 0; i < len(stockArray2); i++ {
	//	fmt.Println(
	//		" 代码=", stockArray2[i].SECURITY_CODE,
	//		" 简称=", stockArray2[i].SECURITY_NAME_ABBR,
	//		" 净利润=", stockArray2[i].PARENT_NETPROFIT,
	//		" 营业总收入=", stockArray2[i].TOTAL_OPERATE_INCOME,
	//		" 营业总收入同比=", stockArray2[i].TOI_RATIO,
	//		" 营业收支出=", stockArray2[i].TOTAL_OPERATE_COST,
	//		" 营业支出=", stockArray2[i].OPERATE_COST,
	//		" 销售费用=", stockArray2[i].SALE_EXPENSE,
	//		" 管理费用=", stockArray2[i].MANAGE_EXPENSE,
	//		" 财务费用=", stockArray2[i].FINANCE_EXPENSE,
	//		" 营业利润=", stockArray2[i].OPERATE_PROFIT,
	//		" 利润总额=", stockArray2[i].TOTAL_PROFIT,
	//		" 净利润同比=", stockArray2[i].PARENT_NETPROFIT_RATIO,
	//		" 扣非净利润=", stockArray2[i].DEDUCT_PARENT_NETPROFIT,
	//	)
	//}

	return stockArray2
}
