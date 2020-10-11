package stock

const (
	file_stockGzfxformat = "stock_估值分析_%s.data"               //股票估值数据
	file_stockYjbbformat = "stock_业绩报表_%s.data"               //业绩报表
	file_stockZcfzformat = "stock_资产负债_%s.data"               //资产负债
	file_stockLrbformat  = "stock_利润表_%s.data"                //利润表
	file_stockXjllformat = "stock_现金流量_%s.data"               //现金流量
	report_singleformat  = "report_%s_%s.data"                //个股数据
	jgdy_singleformate   = "jgdy_%s_%s.data"                  //个股-机构调研
	kx_singleformate     = "kx_%s_%s.data"                    //个股-K线
	reportDate           = "2020-06-30"                       //报表时间
	token                = "894050c76af8597a853f5b408b759f5d" //访问用到的token
	TradeData            = "2020-10-09"                       //最后一个交易日
)

//http://data.eastmoney.com/gzfx/hylist.html		行业估值分析
//http://data.eastmoney.com/gzfx/list.html			估值分析 - 全部  股票列表
//http://data.eastmoney.com/bbsj/202006/yjbb.html	业绩报表 - 全部
//http://data.eastmoney.com/bbsj/202006/zcfz.html	资产负债列表 - 全部
//http://data.eastmoney.com/bbsj/202006/lrb.html	利润表 - 全部
//http://data.eastmoney.com/bbsj/202006/xjll.html	现金流量表

//http://f10.eastmoney.com/OperationsRequired/Index?type=web&code=SZ000069#	个股数据  操盘必读-财务分析
//http://data.eastmoney.com/stockdata/000069.html		个股数据
