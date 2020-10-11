package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"util/file"
	"util/http"
)

//估值分析
//{
//"SECURITYCODE":"000001",				股票编码
//"SName":"平安银行",						公司简称
//"CompanyCode":"10000590",				公司编号
//"MKT":"szzb",							市场
//"HYName":"银行",						行业
//"HYCode":"016029",						行业编码
//"TRADEDATE":"2020-09-30T00:00:00",		交易日期
//"PE9":11.1216,							PE(TTM)
//"PE7":10.4411,							PE(静)
//"PB8":1.046,							市静率
//"PB7":1.0782,
//"PCFJYXJL7":-7.3551,
//"PCFJYXJL9":-3.1956,
//"PS7":2.1339,
//"PS9":1.983,							市销率	和通达信差别很大
//"PEG1":0.811217082898694,				PEG值
//"ZSZ":294387779063.66,					总市值
//"AGSZBHXS":294385268155.6,
//"ZGB":19405918198.0,
//"LTAG":19405752680.0,
//"NEW":15.17,							最新价
//"CHG":2.5,								涨幅
//"HY_PE9":6.51325182207927,				行业--PE(TTM)
//"HY_PE7":6.23559040885457,				行业--PE(静)
//"HY_PB8":0.772921522694968,				行业--市静率
//"HY_PB7":0.803297501254901,
//"HY_PCFJYXJL7":30.1625509569441,
//"HY_PCFJYXJL9":6.80684983826845,
//"HY_PS7":2.09050593250638,
//"HY_PS9":2.01950019161583,
//"HY_PEG1":0.984226600308373,
//"HY_ZSZ":269666778962.719,
//"HY_AGSZBHXS":170231480129.876,
//"HY_ZGB":50485617839.1667,
//"HY_LTAG":31308292166.6944,"ORIGINALCODE":"475"
//},
type StockGzfx struct {
	SECURITYCODE string  //股票编码
	SName        string  //股票名
	NEW          float32 //最新价
	HYName       string  //行业名
	HYCode       string  //行业编码

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
func DownloadStockGzfx() {
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

func ReadStockGzfx() []StockGzfx {
	var fileName = fmt.Sprintf(file_stockGzfxformat, TradeData)
	if !file.FileExist(fileName) {
		DownloadStockGzfx()
	}

	//读取数据
	data, err2 := file.ReadFile_v1(fileName)
	if nil != err2 {
		fmt.Println("error read file ", err2)
		return nil
	}

	var stockArray2 []StockGzfx
	err := json.Unmarshal([]byte(data), &stockArray2)
	if nil != err {
		fmt.Println("22222 json unmarshal failed!!!!", err)
		return nil
	}

	for i := 0; i < len(stockArray2); i++ {
		//fmt.Printf("%6s %6s %6s %06s 动=%f 静=%f 市净=%f 市销=%f\n",
		//	stockArray2[i].SECURITYCODE,
		//	stockArray2[i].SName,
		//	stockArray2[i].HYName,
		//	stockArray2[i].HYCode,
		//	stockArray2[i].PE9,
		//	stockArray2[i].PE7,
		//	stockArray2[i].PB8,
		//	stockArray2[i].PS9,
		//)
	}

	return stockArray2
}
