package stock

type SingleStock struct {
	ZYZB []SingleZyzb //主要指标
	ZCFZ []SingleZcfz //资产负债
	LRB  []SingleLrb  //利润表
	THBJ *SingleTHBJ  //同行比较
	Gbyj *SingleGbyj  //股本研究
	JGTJ int          //机构推荐数
}

func ParseSingle(code string, icode string) (*SingleStock, error) {
	sigleStock, err := ParseReportData(code) //解析报表数据 -- 报表
	sigleStock.THBJ = ParseTHBJ(code, icode)
	sigleStock.Gbyj = ParseGbyj(code)
	sigleStock.JGTJ = ParseSingleJkdy(code)

	return sigleStock, err
}
