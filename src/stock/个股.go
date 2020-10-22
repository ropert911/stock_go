package stock

type SingleStock struct {
	ZYZB []SingleZyzb //主要指标
	ZCFZ []SingleZcfz //资产负债
	JGTJ int          //机构推荐数
	LAST float32      //去年涨幅
	THIS float32      //今年涨幅
	TWO  float32      //两年合计涨幅
	Gbyj *SingleGbyj
}

func DownloadSingle(code string) {
	DowloadKx(code)          //下转年K线数据
	DownloadReportData(code) //下载报表相差的
	DowloadJgdy(code)        //下载机构调研相关数据
	DowloadGbyj(code)        //股本研究
}

func ParseSingle(code string) (*SingleStock, error) {
	sigleStock, err := ParseReportData(code)                         //解析报表数据 -- 报表
	sigleStock.JGTJ = ParseSingleJkdy(code)                          //解析机构推荐数据 -- 日转月用
	sigleStock.LAST, sigleStock.THIS, sigleStock.TWO = ParseKy(code) //K线数据 -- 日转月用
	sigleStock.Gbyj = ParseGbyj(code)

	return sigleStock, err
}
