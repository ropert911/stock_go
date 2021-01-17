package main

import (
	"fmt"
	"math"
	"stock"
	"strings"
	"unicode"
	"util/file"
)

type StockInfo2 struct {
	zcfz stock.StockZCFZ
	yjbb stock.StockYJBB
	lrb  stock.StockLRB
	xjll stock.StockXJLL
	gzfx stock.StockGZFX
	xlyc stock.StockYLYC
}

func main() {
	mapStock := parserData()
	mapStock = filterData(mapStock)
	exportResult(mapStock)
	//finAnalyser2(mapStock)
}

//解析所有数据-估值分析、资产负债、业绩报表、利润表
func parserData() map[string]StockInfo2 {
	mapStock := make(map[string]StockInfo2)
	socksZcfz := stock.ReadStockZCFZ()
	socksYjbb := stock.ReadStockYJBB()
	stockLrb := stock.ReadStockLRB()
	stockXjll := stock.ReadStockXJLL()
	socksgzfx := stock.ReadStockGZFX()
	stockYlyc := stock.ReadStockYLYC()
	for i := 0; i < len(socksZcfz); i++ {
		var stockInfo StockInfo2
		stockInfo.zcfz = socksZcfz[i]
		for j := 0; j < len(socksYjbb); j++ {
			if socksZcfz[i].SECURITY_CODE == socksYjbb[j].SECURITY_CODE {
				stockInfo.yjbb = socksYjbb[j]
				break
			}
		}

		for z := 0; z < len(stockLrb); z++ {
			if socksZcfz[i].SECURITY_CODE == stockLrb[z].SECURITY_CODE {
				stockInfo.lrb = stockLrb[z]
				break
			}
		}
		for z := 0; z < len(stockXjll); z++ {
			if socksZcfz[i].SECURITY_CODE == stockXjll[z].SECURITY_CODE {
				stockInfo.xjll = stockXjll[z]
				break
			}
		}
		for z := 0; z < len(socksgzfx); z++ {
			if socksZcfz[i].SECURITY_CODE == socksgzfx[z].SECURITYCODE {
				stockInfo.gzfx = socksgzfx[z]
				break
			}
		}
		for z := 0; z < len(stockYlyc); z++ {
			if socksgzfx[i].SECURITYCODE == stockYlyc[z].STOCKCODE {
				stockInfo.xlyc = stockYlyc[z]
				break
			}
		}

		mapStock[socksZcfz[i].SECURITY_CODE] = stockInfo
	}

	return mapStock
}

//过滤掉不符合条件的
func filterData(mapStock map[string]StockInfo2) map[string]StockInfo2 {
	for key, value := range mapStock {
		//////////////////////资金效率高
		//净资产收益率>10%
		{
			if value.yjbb.WEIGHTAVG_ROE < 10 {
				delete(mapStock, key)
				continue
			}
		}
		//毛利率>=10%
		//{
		//	var tmp = value.yjbb.XSMLL
		//	if tmp < 10 {
		//		delete(mapStock, key)
		//		continue
		//	}
		//}
		//净利润率>=5%
		//{
		//	var tmp = 100 * value.yjbb.PARENT_NETPROFIT / value.yjbb.TOTAL_OPERATE_INCOME
		//	if tmp < 5 {
		//		delete(mapStock, key)
		//		continue
		//	}
		//}

		///////////////////收入多
		//净利润>0.2亿
		{
			if value.yjbb.PARENT_NETPROFIT < 20000000 {
				delete(mapStock, key)
				continue
			}
		}
		//扣非净利润>0.2亿
		{
			if value.lrb.DEDUCT_PARENT_NETPROFIT < 20000000 {
				delete(mapStock, key)
				continue
			}
		}
		//年收入大于10亿
		{
			var yincome = value.yjbb.TOTAL_OPERATE_INCOME
			var _, date = stock.GetDate(value.zcfz.REPORT_DATE)
			if date == 3 {
				yincome = yincome * 4
			} else if date == 6 {
				yincome = yincome * 2
			} else if date == 9 {
				yincome = yincome * 4 / 3
			}

			if yincome < 1000000000 {
				delete(mapStock, key)
				continue
			}
		}

		///////////////////有余钱
		//(货币资金/营业收入)>20%
		{
			var tmp = 100 * value.zcfz.MONETARYFUNDS / value.yjbb.TOTAL_OPERATE_INCOME
			if tmp < 20 {
				delete(mapStock, key)
				continue
			}
		}
		//总现金流>0.2亿
		{
			if value.xjll.CCE_ADD < 20000000 {
				delete(mapStock, key)
				continue
			}
		}
		//每股净资产>1 && 不为空
		{
			if value.yjbb.BPS < 1 {
				delete(mapStock, key)
				continue
			}
		}

		/////////////估值
		//100>市盈率(动静)>0不为空
		{
			if value.gzfx.PE9 >= 100 || value.gzfx.PE9 <= 0 || value.gzfx.PE7 >= 100 || value.gzfx.PE7 <= 0 {
				delete(mapStock, key)
				continue
			}
		}
		//市销率<10
		{
			if value.gzfx.PS9 >= 10 {
				delete(mapStock, key)
				continue
			}
		}

		/////////////////////风险
		//资产负债率<70%
		{
			if value.zcfz.DEBT_ASSET_RATIO >= 70 {
				delete(mapStock, key)
				continue
			}
		}

		///////////成长
		//收入同比>=15%
		{
			if value.yjbb.YSTZ < 15 {
				delete(mapStock, key)
				continue
			}
		}

		///////////////其它
		//非ST
		{
			if strings.HasPrefix(value.zcfz.SECURITY_NAME_ABBR, "*ST") || strings.HasPrefix(value.zcfz.SECURITY_NAME_ABBR, "ST") {
				delete(mapStock, key)
				continue
			}
		}
		//>5元 && 不为空
		{
			if value.gzfx.NEW < 5 {
				delete(mapStock, key)
				continue
			}
		}

		//========================================个股=================================
		var single, err = stock.ParseSingle(key)
		if nil != err {
			fmt.Println("Error parse single stock data ", err)
			continue
		}

		//未分配利润 >1亿
		{
			if stock.ToFloat(single.ZCFZ[0].RETAINEDEARNING) < 100000000 {
				delete(mapStock, key)
				continue
			}
		}
		//每股公积金>0.5元
		if stock.ToFloat(single.ZYZB[0].MGGJJ) < 0.5 {
			delete(mapStock, key)
			continue
		}

		//每股未分配利润>0.5元
		if stock.ToFloat(single.ZYZB[0].MGWFPLY) < 0.5 {
			delete(mapStock, key)
			continue
		}

		//成长要求
		{
			var zyzbP0 = single.ZYZB[0] //当前财务指标
			var zyzbP1 stock.SingleZyzb //上一年的主要指标
			var zyzbP2 stock.SingleZyzb //前2年的财务指标
			//找到前2年的财务指标
			var curYear, _ = stock.GetDate(zyzbP0.DATE)
			for i := 0; i < len(single.ZYZB); i++ {
				var year, date = stock.GetDate(single.ZYZB[i].DATE)
				if date == 12 {
					if 1 == curYear-year {
						zyzbP1 = single.ZYZB[i]
					} else if 2 == curYear-year {
						zyzbP2 = single.ZYZB[i]
					}
				}
			}

			//近3年不能有负增长
			{
				var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
				var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
				var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
				if tb0 < 0 || tb1 < 0 || tb2 < 0 {
					delete(mapStock, key)
					continue
				}
			}
			//近2年平均>15%
			{
				var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
				var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
				var last2avg = math.Cbrt((1+tb0/100)*(1+tb1/100))*100 - 100
				if last2avg < 15 {
					delete(mapStock, key)
					continue
				}
			}
		}

		////3成长 -- 营业收入>0不为空
		//if value.yjbb.TOTAL_OPERATE_INCOME <= 0 {
		//	delete(mapStock, key)
		//	continue
		//}

		//

		//
		////3成长 -- 近3年利润不为负
		//{
		//	var jlr0 = stock.ToFloat(single.ZYZB[0].GSJLR)
		//	var jlr1 = stock.ToFloat(zyzbP1.GSJLR)
		//	var jlr2 = stock.ToFloat(zyzbP2.GSJLR)
		//	if jlr0 < 0 || jlr1 < 0 || jlr2 < 0 {
		//		delete(mapStock, key)
		//		continue
		//	}
		//}

		////5其它 -- 近2年<150 每年小于60%
		//if single.THIS >= 60 || single.TWO >= 150 {
		//	delete(mapStock, key)
		//	continue
		//}
	}

	return mapStock
}

//生成导出信息
func CreateExportEBK2(mapStock map[string]StockInfo2) {
	if len(mapStock) > 0 {
		var exportName = fmt.Sprintf("%s.EBK", stock.TradeData)
		file.WriteFile(exportName, `
`)

		for key, _ := range mapStock {
			var code = stock.GetExportCodeByCode(key)
			file.AppendFile(exportName, fmt.Sprintf(`%s
`, code))
		}
	}
}

type SockInfoShow2 struct {
	Code          string
	Name          string
	Price         string
	HYName        string
	YYZSRAVG      string //营业总收3年平均
	MLL           string //毛利率(%)
	JLL           string //净利率(%)
	WEIGHTAVG_ROE string //净益率
	MGGJJ         string //每股公积金
	MGWFPLY       string //每股未分配

	//估值
	SJL    string //市净率
	PEJT   string //静态市盈率
	PEDT   string //动态市盈率
	PS9    string //市销率
	RPB8   string //市净率估值
	RPE7   string //PE(静)估值
	RPE9   string //PE(TTM)估值
	RPS9   string //市销率估值
	MGLZGZ string //每股内在股价
	//主力研究
	QSDGDCGHJ   string  //前十大股东持股合计
	QSDLTGDCGHJ string  //前十大流通股东持股合计
	SBZB        string  //社保占流通比
	JGZB        float64 //机构合计占流通比
	JGTJ        string  //机构推荐数
}
type SockInfoShowArray []SockInfoShow2

func (s SockInfoShowArray) Len() int           { return len(s) }
func (s SockInfoShowArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SockInfoShowArray) Less(i, j int) bool { return s[i].JGZB >= s[j].JGZB }

//显示结果
func exportResult(mapStock map[string]StockInfo2) {
	CreateExportEBK2(mapStock)

	//var sockInfoShows = SockInfoShowArray{}

	//把符合条件股票中的亮点数据显示出来
	//for key, value := range mapStock {
	//stock.DownloadSingle(key)
	//var single, _ = stock.ParseSingle(key)
	//
	//var sockInfoShow SockInfoShow2
	//sockInfoShow.Code = key
	//sockInfoShow.Name = value.zcfz.SECURITY_NAME_ABBR
	//sockInfoShow.Price = fmt.Sprint(value.gzfx.NEW)
	//sockInfoShow.HYName = value.gzfx.HYName

	//2积累--每股公积金
	//var mggjj = stock.ToFloat(single.ZYZB[0].MGGJJ)
	//if mggjj*3 > float64(value.gzfx.NEW) || mggjj > 10 {
	//	sockInfoShow.MGGJJ = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, mggjj, 0x1B)
	//} else if mggjj*4 > float64(value.gzfx.NEW) || mggjj > 8 {
	//	sockInfoShow.MGGJJ = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, mggjj, 0x1B)
	//} else {
	//	sockInfoShow.MGGJJ = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, mggjj, 0x1B)
	//}
	//1积累--每股未分配利润
	//var mgwfply = stock.ToFloat(single.ZYZB[0].MGWFPLY)
	//if mgwfply*3 > float64(value.gzfx.NEW) || mgwfply > 10 {
	//	sockInfoShow.MGWFPLY = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, mgwfply, 0x1B)
	//} else if mgwfply*4 > float64(value.gzfx.NEW) || mgwfply > 8 {
	//	sockInfoShow.MGWFPLY = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, mgwfply, 0x1B)
	//} else {
	//	sockInfoShow.MGWFPLY = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, mgwfply, 0x1B)
	//}

	//2估值-市净率
	//if value.gzfx.PB8 <= 5 {
	//	sockInfoShow.SJL = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.gzfx.PB8, 0x1B)
	//} else {
	//	sockInfoShow.SJL = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PB8, 0x1B)
	//}
	//2估值-动态市盈率
	//if value.gzfx.PE9 <= 15 {
	//	sockInfoShow.PEDT = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, value.gzfx.PE9, 0x1B)
	//} else if value.gzfx.PE9 <= 25 {
	//	sockInfoShow.PEDT = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.gzfx.PE9, 0x1B)
	//} else {
	//	sockInfoShow.PEDT = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PE9, 0x1B)
	//}
	//2估值-静态市盈率
	//if value.gzfx.PE7 <= 15 {
	//	sockInfoShow.PEJT = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, value.gzfx.PE7, 0x1B)
	//} else if value.gzfx.PE7 <= 25 {
	//	sockInfoShow.PEJT = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.gzfx.PE7, 0x1B)
	//} else {
	//	sockInfoShow.PEJT = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PE7, 0x1B)
	//}
	//2估值-市销率
	//if value.gzfx.PS9 <= 2 {
	//	sockInfoShow.PS9 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, value.gzfx.PS9, 0x1B)
	//} else if value.gzfx.PS9 <= 3 {
	//	sockInfoShow.PS9 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.gzfx.PS9, 0x1B)
	//} else {
	//	sockInfoShow.PS9 = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PS9, 0x1B)
	//}
	//2估值-每股内在股价
	//E(2R+8.5)*4.4/Y
	// 	E:每股收益
	//	R:预期收益增长率
	//	8.5：平均市盈率，中国应该是22.5，按20来算
	//	4.4：平均利息
	//	Y:公司债/国债收益率（5年期）  约等于3.2%
	//	-->每股收益*(2*预期收益增长率+22.5)*4.4/3.2
	//mglgz := float64(value.yjbb.BASIC_EPS) * (2*value.xlyc.EGR + 20.5) * 4.4 / 3.2
	//gzbl := mglgz / float64(value.gzfx.NEW)
	//if gzbl > 1.30 {
	//	sockInfoShow.MGLZGZ = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, mglgz, 0x1B)
	//} else {
	//	sockInfoShow.MGLZGZ = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
	//}

	//3成长-净益率 净资产收益率
	//if value.yjbb.WEIGHTAVG_ROE >= 20 {
	//	sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, value.yjbb.WEIGHTAVG_ROE, 0x1B)
	//} else if value.yjbb.WEIGHTAVG_ROE >= 10 {
	//	sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.yjbb.WEIGHTAVG_ROE, 0x1B)
	//} else {
	//	sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.yjbb.WEIGHTAVG_ROE, 0x1B)
	//}
	//3成长-毛利率
	//var mll = stock.ToFloat(single.ZYZB[0].MLL)
	//if mll >= 30 {
	//	sockInfoShow.MLL = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, mll, 0x1B)
	//} else {
	//	sockInfoShow.MLL = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, mll, 0x1B)
	//}
	//3成长-净利润率
	//var jll = stock.ToFloat(single.ZYZB[0].JLL)
	//if jll >= 20 {
	//	sockInfoShow.JLL = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, jll, 0x1B)
	//} else {
	//	sockInfoShow.JLL = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, jll, 0x1B)
	//}
	//3成长 -- 3年平均>=20%
	//var zyzbP1 stock.SingleZyzb
	//var zyzbP2 stock.SingleZyzb
	//var curYear, _ = stock.GetDate(single.ZYZB[0].DATE)
	//for i := 0; i < len(single.ZYZB); i++ {
	//	var year, date = stock.GetDate(single.ZYZB[i].DATE)
	//	if date == 12 {
	//		if 1 == curYear-year {
	//			zyzbP1 = single.ZYZB[i]
	//		} else if 2 == curYear-year {
	//			zyzbP2 = single.ZYZB[i]
	//		}
	//	}
	//}
	//var tb0 = stock.ToFloat(single.ZYZB[0].YYZSRTBZZ)
	//var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
	//var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
	//var avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
	//if avg >= 40 {
	//	sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;31m%.2f%c[0m", 0x1B, avg, 0x1B)
	//} else if avg >= 30 {
	//	sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;33m%.2f%c[0m", 0x1B, avg, 0x1B)
	//} else if avg >= 20 {
	//	sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, avg, 0x1B)
	//} else {
	//	sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, avg, 0x1B)
	//}
	//
	//var rPB8 = value.gzfx.PB8 / value.gzfx.HY_PB8
	//if rPB8 < 0.7 {
	//	sockInfoShow.RPB8 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, rPB8, 0x1B)
	//} else if rPB8 < 1 {
	//	sockInfoShow.RPB8 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPB8, 0x1B)
	//} else {
	//	sockInfoShow.RPB8 = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
	//}
	//
	//var rPE7 = value.gzfx.PE7 / value.gzfx.HY_PE7
	//if rPE7 < 0.8 {
	//	sockInfoShow.RPE7 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, rPE7, 0x1B)
	//} else if rPE7 < 1 {
	//	sockInfoShow.RPE7 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPE7, 0x1B)
	//} else {
	//	sockInfoShow.RPE7 = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
	//}
	//
	//var rPE9 = value.gzfx.PE9 / value.gzfx.HY_PE9
	//if rPE9 < 0.8 {
	//	sockInfoShow.RPE9 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, rPE9, 0x1B)
	//} else if rPE9 < 1 {
	//	sockInfoShow.RPE9 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPE9, 0x1B)
	//} else {
	//	sockInfoShow.RPE9 = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
	//}
	//
	//var rPS9 = value.gzfx.PS9 / value.gzfx.HY_PS9
	//if rPS9 < 0.8 {
	//	sockInfoShow.RPS9 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, rPS9, 0x1B)
	//} else if rPS9 < 1 {
	//	sockInfoShow.RPS9 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPS9, 0x1B)
	//} else {
	//	sockInfoShow.RPS9 = fmt.Sprintf("%c[;;30m   %c[0m", 0x1B, 0x1B)
	//}
	//
	//前十大股东持股合计
	//if len(single.Gbyj.GDRS) > 0 && "--" != single.Gbyj.GDRS[0].QSDGDCGHJ {
	//	var sdgdzb = stock.ToFloat(single.Gbyj.GDRS[0].QSDGDCGHJ)
	//	if sdgdzb > 70 {
	//		sockInfoShow.QSDGDCGHJ = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, sdgdzb, 0x1B)
	//	} else if sdgdzb > 40 {
	//		sockInfoShow.QSDGDCGHJ = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, sdgdzb, 0x1B)
	//	} else {
	//		sockInfoShow.QSDGDCGHJ = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, sdgdzb, 0x1B)
	//	}
	//} else {
	//	sockInfoShow.QSDGDCGHJ = fmt.Sprintf("%c[;;30m%c[0m", 0x1B, 0x1B)
	//}
	//
	//前十大流通股东持股合计
	//if len(single.Gbyj.GDRS) > 0 && "--" != single.Gbyj.GDRS[0].QSDLTGDCGHJ {
	//	var sdltgdzb = stock.ToFloat(single.Gbyj.GDRS[0].QSDLTGDCGHJ)
	//	if sdltgdzb > 45 {
	//		sockInfoShow.QSDLTGDCGHJ = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, sdltgdzb, 0x1B)
	//	} else if sdltgdzb > 30 {
	//		sockInfoShow.QSDLTGDCGHJ = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, sdltgdzb, 0x1B)
	//	} else {
	//		sockInfoShow.QSDLTGDCGHJ = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, sdltgdzb, 0x1B)
	//	}
	//} else {
	//	sockInfoShow.QSDLTGDCGHJ = fmt.Sprintf("%c[;;30m%c[0m", 0x1B, 0x1B)
	//}
	//
	//sockInfoShow.SBZB = ""
	//sockInfoShow.JGZB = 0
	//for i := 0; i < len(single.Gbyj.ZLCC); i++ {
	//	if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "社保") && "--" != single.Gbyj.ZLCC[i].ZLTGBL {
	//		var sbbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
	//		sockInfoShow.SBZB = fmt.Sprintf("%.2f", sbbl)
	//	} else if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "合计") {
	//		var hjbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
	//		sockInfoShow.JGZB = hjbl
	//	}
	//}

	//机构推荐数
	//if single.JGTJ >= 20 {
	//	sockInfoShow.JGTJ = fmt.Sprintf("%c[;;35m%d%c[0m", 0x1B, single.JGTJ, 0x1B)
	//} else if single.JGTJ >= 6 {
	//	sockInfoShow.JGTJ = fmt.Sprintf("%c[;;36m%d%c[0m", 0x1B, single.JGTJ, 0x1B)
	//} else {
	//	sockInfoShow.JGTJ = fmt.Sprintf("%c[;;30m%c[0m", 0x1B, 0x1B)
	//}
	//
	//	sockInfoShows = append(sockInfoShows, sockInfoShow)
	//}
	//sort.Stable(sockInfoShows)
	//
	//fmt.Println("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
	//fmt.Printf("┃%7s┃%8s┃%5s┃%8s", "Code", "行业", "现价", "行业")
	//fmt.Printf("┃┃%4s┃%3s┃%3s┃%4s┃%4s┃%3s", "收入增长", "毛利润", "净利率", "净益率", "公积金", "未分配")
	//fmt.Printf("┃┃%4s┃%5s┃%5s┃%3s┃%3s┃%4s┃%4s┃%3s┃%4s", "市净", "PE静", "PE动", "市销率", "市净比", "PE比", "PET比", "市销比", "内在估值")
	//fmt.Printf("┃┃%3s┃%3s┃%5s┃%4s┃%2s┃\n", "十股占", "十流占", "社/流", "机/流", "推荐")
	//fmt.Println("┃━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┃")
	//for i := 0; i < len(sockInfoShows); i++ {
	//	var stock = sockInfoShows[i]
	//	fmt.Printf("┃")
	//	fmt.Printf("%7s", stock.Code)
	//
	//	fmt.Printf("┃")
	//	var rLen = len(stock.Name) - ChineseCount2(stock.Name)
	//	var bLen = 10 - rLen
	//	for bLen > 0 {
	//		fmt.Printf(" ")
	//		bLen--
	//	}
	//	fmt.Printf("%s", stock.Name)
	//
	//	fmt.Printf("┃%7s", stock.Price)

	//fmt.Printf("┃")
	//rLen = len(stock.HYName) - ChineseCount2(stock.HYName)
	//bLen = 10 - rLen
	//for bLen > 0 {
	//	fmt.Printf(" ")
	//	bLen--
	//}
	//fmt.Printf("%s", stock.HYName)

	//成长-有利润
	//fmt.Printf("┃┃%19s", stock.YYZSRAVG)
	//fmt.Printf("┃%17s", stock.MLL)
	//fmt.Printf("┃%17s", stock.JLL)
	//fmt.Printf("┃%17s", stock.WEIGHTAVG_ROE)
	//fmt.Printf("┃%17s", stock.MGGJJ)
	//fmt.Printf("┃%17s", stock.MGWFPLY)

	//估值类
	//fmt.Printf("┃┃%17s", stock.SJL)
	//fmt.Printf("┃%17s", stock.PEJT)
	//fmt.Printf("┃%17s", stock.PEDT)
	//fmt.Printf("┃%17s", stock.PS9)
	//fmt.Printf("┃%16s", stock.RPB8)
	//fmt.Printf("┃%16s", stock.RPE7)
	//fmt.Printf("┃%16s", stock.RPE9)
	//fmt.Printf("┃%16s", stock.RPS9)
	//fmt.Printf("┃%18s", stock.MGLZGZ)

	////主力研究
	//fmt.Printf("┃┃%17s", stock.QSDGDCGHJ)
	//fmt.Printf("┃%17s", stock.QSDLTGDCGHJ)
	//fmt.Printf("┃%6s", stock.SBZB)
	//fmt.Printf("┃%5.2f", stock.JGZB)
	//fmt.Printf("┃%15s", stock.JGTJ)
	//
	//fmt.Println("┃")
	//}
	//fmt.Println("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛")
	fmt.Println("符合条件的股票有=", len(mapStock), "个")
}

//func finAnalyser2(mapStock map[string]StockInfo2) {
//	for key, _ := range mapStock {
//		var single, _ = stock.ParseSingle(key)
//		var zcfz = single.ZCFZ[0]
//		//短期有息债
//		var dqyxz = stock.ToFloat(zcfz.STBORROW) + stock.ToFloat(zcfz.NONLLIABONEYEAR)
//
//		var xjbl = stock.ToFloat(zcfz.MONETARYFUND) / dqyxz
//		if xjbl > 0.7 {
//			//fmt.Printf("---- 现金比率=货币资金/短期有息债>70%%:%s 值=%.2f%%  (货币资金=%.2f亿 短期有息债=%.2f亿)\n",
//			//	"\033[36;4m"+"正常"+"\033[0m",
//			//	100*xjbl, stock.ToFloat(zcfz.MONETARYFUND)/100000000, dqyxz/100000000)
//		} else if xjbl < 0.5 {
//			//fmt.Printf("---- 现金比率=货币资金/短期有息债>70%%:%s 值=%.2f%%  (货币资金=%.2f亿 短期有息债=%.2f亿)\n",
//			//	"\033[31;4m"+"警告"+"\033[0m",
//			//	100*xjbl, stock.ToFloat(zcfz.MONETARYFUND)/100000000, dqyxz/100000000)
//		} else {
//			//fmt.Printf("---- 现金比率=货币资金/短期有息债>70%%:%s 值=%.2f%%  (货币资金=%.2f亿 短期有息债=%.2f亿)\n",
//			//	"注意",
//			//	100*xjbl, stock.ToFloat(zcfz.MONETARYFUND)/100000000, dqyxz/100000000)
//		}
//	}
//}

func ChineseCount2(str1 string) (count int) {
	for _, char := range str1 {
		if unicode.Is(unicode.Han, char) {
			count++
		}
	}

	return
}
