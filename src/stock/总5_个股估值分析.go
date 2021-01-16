package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"util/file"
	"util/http"
)

type StockGZFX struct {
	SECURITYCODE string  //股票编码
	SName        string  //股票名
	NEW          float32 //最新价
	HYName       string  //行业名
	HYCode       string  //行业编码
	ORIGINALCODE string  //原始行业编码

	PB8 float32 //市净率
	PE7 float32 //静态市盈率
	PE9 float32 //动态市盈率
	PS9 float32 //市销率

	HY_PB8 float32 //行业--市净率
	HY_PE7 float32 //行业--PE(静)
	HY_PE9 float32 //行业--PE(TTM)
	HY_PS9 float32 //行业-市销率
}

//估值分析表
func DownloadStockGZFX() {
	var fileName = fmt.Sprintf(file_stockGzfxformat, TradeData)
	if file.FileExist(fileName) {
		return
	}

	file.Delete(fileName)
	file.WriteFile(fileName, "[")

	var urlFormat = "http://dcfm.eastmoney.com/EM_MutiSvcExpandInterface/api/js/get?type=GZFX_GGZB&token=%s&st=SECURITYCODE&sr=1&p=%d&ps=%d&filter=(TRADEDATE=^%s^)&rt=%d"
	var index = 1
	var first = true
	for {
		var url = fmt.Sprintf(urlFormat,
			token, //
			index, //第几页
			100,   //每页条数
			TradeData,
			rand.Intn(899999)+100000)
		fmt.Println(url)
		index++

		//获取数据
		content, err := http.HttpGet(url)
		if nil != err {
			fmt.Println("http get failed url=", url, " error=", err)
			break
		}

		//写内容
		data := *content
		if len(data) > 10 {
			if !first {
				file.AppendFile(fileName, ",")
			} else {
				first = false
			}
			data = data[1 : len(data)-1]
			file.AppendFile(fileName, data)
		} else {
			break
		}
	}
	file.AppendFile(fileName, "]")
}

func ReadStockGZFX() []StockGZFX {
	var fileName = fmt.Sprintf(file_stockGzfxformat, TradeData)
	if !file.FileExist(fileName) {
		DownloadStockGZFX()
	}

	//读取数据
	data, err2 := file.ReadFile_v1(fileName)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockGZFX
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	//for i := 0; i < len(stockArray2); i++ {
	//	fmt.Printf("%6s %6s %6s %06s 动=%f 静=%f 市净=%f 市销=%f\n",
	//		stockArray2[i].SECURITYCODE,
	//		stockArray2[i].SName,
	//		stockArray2[i].HYName,
	//		stockArray2[i].HYCode,
	//		stockArray2[i].PE9,
	//		stockArray2[i].PE7,
	//		stockArray2[i].PB8,
	//		stockArray2[i].PS9,
	//	)
	//}

	return stockArray2
}
