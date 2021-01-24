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
	CreateExportEBK(mapStock)
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
			if socksZcfz[i].SECURITY_CODE == stockYlyc[z].STOCKCODE {
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

		//////////////////////公司概况===================================
		//非ST
		{
			if strings.HasPrefix(value.zcfz.SECURITY_NAME_ABBR, "*ST") || strings.HasPrefix(value.zcfz.SECURITY_NAME_ABBR, "ST") {
				delete(mapStock, key)
				continue
			}
		}
		//现价>5元 && 不为空
		{
			if value.gzfx.NEW < 5 {
				delete(mapStock, key)
				continue
			}
		}
		//营业收入>0不为空
		if value.yjbb.TOTAL_OPERATE_INCOME <= 0 {
			delete(mapStock, key)
			continue
		}
		//资产负债率<70%
		{
			if value.zcfz.DEBT_ASSET_RATIO >= 70 {
				delete(mapStock, key)
				continue
			}
		}
		//////////////////////龙头业绩和增长===================================
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

			if yincome < (10 * 10000 * 10000) {
				delete(mapStock, key)
				continue
			}
		}
		//收入同比>=10%
		{
			if value.yjbb.YSTZ < 10 {
				delete(mapStock, key)
				continue
			}
		}

		//近3年不能有负增长 --- 单个里
		//近2年平均>10% --- 单个里
		//近3年平均>8% --- 单个里

		//净利润>0.2亿
		{
			if value.yjbb.PARENT_NETPROFIT < (0.2 * 10000 * 10000) {
				delete(mapStock, key)
				continue
			}
		}
		//扣非净利润>0.2亿
		{
			if value.lrb.DEDUCT_PARENT_NETPROFIT < (0.2 * 10000 * 10000) {
				delete(mapStock, key)
				continue
			}
		}
		//净资产收益率>10%
		{
			if value.yjbb.WEIGHTAVG_ROE < 10 {
				delete(mapStock, key)
				continue
			}
		}
		//三年平均ROE>行业平均和行业中值 --- 单个里
		//前一年的净资产收益>=行业中值*1.2--- 单个里
		//(货币资金/营业收入)>20%
		{
			var tmp = 100 * value.zcfz.MONETARYFUNDS / value.yjbb.TOTAL_OPERATE_INCOME
			if tmp < 20 {
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
		//未分配利润 >1亿 --- 单个里
		//每股公积金>0.5元 --- 单个里
		//每股未分配利润>0.5元 --- 单个里

		//////////////////////估值===================================
		//100>市盈率(动静)>0不为空
		{
			if value.gzfx.PE9 >= 100 || value.gzfx.PE9 <= 0 || value.gzfx.PE7 >= 100 || value.gzfx.PE7 <= 0 {
				delete(mapStock, key)
				continue
			}
		}

		//单个的==============
		{
			var single, err = stock.ParseSingle(key, value.gzfx.ORIGINALCODE)
			if nil != err {
				fmt.Println("Error parse single stock data ", err)
				continue
			}
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

			//近3年不能有负增长 --- 单个里
			{
				var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
				var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
				var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
				if tb0 < 0 || tb1 < 0 || tb2 < 0 {
					delete(mapStock, key)
					continue
				}
			}
			//近2年平均>10% --- 单个里
			{
				var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
				var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
				var last2avg = math.Cbrt((1+tb0/100)*(1+tb1/100))*100 - 100
				if last2avg < 10 {
					delete(mapStock, key)
					continue
				}
			}
			//近3年平均>8% --- 单个里
			{
				var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
				var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
				var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
				var last2avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
				if last2avg < 8 {
					delete(mapStock, key)
					continue
				}
			}
			//三年平均ROE>行业平均和行业中值 --- 单个里
			{
				if stock.ToFloat(single.THBJ.DBFXBJ.DATA[0].ROEPJ) < stock.ToFloat(single.THBJ.DBFXBJ.DATA[1].ROEPJ) {
					delete(mapStock, key)
					continue
				}
				if stock.ToFloat(single.THBJ.DBFXBJ.DATA[0].ROEPJ) < stock.ToFloat(single.THBJ.DBFXBJ.DATA[2].ROEPJ) {
					delete(mapStock, key)
					continue
				}
			}
			//前一年的净资产收益>=行业中值*1.2--- 单个里
			{
				if stock.ToFloat(single.THBJ.DBFXBJ.DATA[0].ROE2) < (stock.ToFloat(single.THBJ.DBFXBJ.DATA[2].ROE2) * 1.2) {
					delete(mapStock, key)
					continue
				}
			}
			//未分配利润 >1亿 --- 单个里
			{
				if stock.ToFloat(single.ZCFZ[0].RETAINEDEARNING) < (10000 * 10000) {
					delete(mapStock, key)
					continue
				}
			}
			//每股公积金>0.5元 --- 单个里
			if stock.ToFloat(single.ZYZB[0].MGGJJ) < 0.5 {
				delete(mapStock, key)
				continue
			}
			//每股未分配利润>0.5元 --- 单个里
			if stock.ToFloat(single.ZYZB[0].MGWFPLY) < 0.5 {
				delete(mapStock, key)
				continue
			}
		}
	}

	return mapStock
}

//生成导出信息
func CreateExportEBK(mapStock map[string]StockInfo2) {
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

//显示结果
func exportResult(mapStock map[string]StockInfo2) {
	//把符合条件股票中的亮点数据显示出来
	fmt.Printf("%4s\t%-8s\t%4s\t%-7s", "编码", "名称", "股价", "行业")
	fmt.Printf("┃┃%s %s %s %s %s", "收入亿/净利润", "增长50%/三年20%", "毛利率", "净利率-排", "ROE")
	fmt.Printf("┃┃%4s", "成长排行")
	fmt.Printf("┃┃%3s %5s %5s  %5s  %5s %4s", "市净率-行", "静态市盈率-行", "动态市盈率-行", "市销率-行", "PEG", "估值排名")
	fmt.Printf("┃┃%4s%4s\t%4s", "社/流", "机/流", "推荐数")

	fmt.Println("")
	for key, value := range mapStock {
		var single, _ = stock.ParseSingle(key, value.gzfx.ORIGINALCODE)
		var zyzbP0 = single.ZYZB[0]
		var zyzbP1 stock.SingleZyzb
		var zyzbP2 stock.SingleZyzb
		var curYear, _ = stock.GetDate(single.ZYZB[0].DATE)
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

		//	var stock = sockInfoShows[i]
		//	===============基本信息
		fmt.Printf("%6s\t%-s", key, value.zcfz.SECURITY_NAME_ABBR) //编码、名称
		var rLen = len(value.zcfz.SECURITY_NAME_ABBR) - ChineseCount2(value.zcfz.SECURITY_NAME_ABBR)
		var bLen = 10 - rLen
		for bLen > 0 {
			fmt.Printf(" ")
			bLen--
		}
		fmt.Printf("\t%6.2f\t%-s", value.gzfx.NEW, value.gzfx.HYName) //"股价", "行业"
		rLen = len(value.gzfx.HYName) - ChineseCount2(value.gzfx.HYName)
		bLen = 9 - rLen
		for bLen > 0 {
			fmt.Printf(" ")
			bLen--
		}

		//收入与增长
		fmt.Printf("\t%5.0f/%-4.0f", value.yjbb.TOTAL_OPERATE_INCOME/(10000*10000), value.yjbb.PARENT_NETPROFIT/(10000*10000)) //收入/净利润
		{
			var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
			var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
			var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
			var avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
			if tb0 >= 50 && avg >= 20 {
				fmt.Printf("\t%3.0f/%-3.0f", tb0, avg) //增长/三年
			} else if tb0 >= 50 {
				fmt.Printf("\t%3.0f/   ", tb0) //增长/三年
			} else if avg >= 20 {
				fmt.Printf("\t   /%-3.0f", avg) //增长/三年
			} else {
				fmt.Printf("\t   /   ")
			}
		}
		fmt.Printf("\t%5.2f", value.yjbb.XSMLL)                         //毛利润
		fmt.Printf("\t%5s/%-3s", zyzbP0.JLL, single.THBJ.GSGMJLR[0].PM) //净利率-排
		fmt.Printf("\t%5.2f", value.yjbb.WEIGHTAVG_ROE)                 //ROE

		//	===============发展
		fmt.Printf("\t%4s", strings.ReplaceAll(single.THBJ.CZXBJ.DATA[0].PM, "U003E", ">")) //成长性排名

		//	===============估值
		fmt.Printf("\t%4.1f/%-4.1f", value.gzfx.PB8, value.gzfx.HY_PB8)                    //市净率-行
		fmt.Printf("\t%5.1f/%-5.1f", value.gzfx.PE7, value.gzfx.HY_PE7)                    //静态市盈率-行
		fmt.Printf("\t%5.1f/%-5.1f", value.gzfx.PE9, value.gzfx.HY_PE9)                    //动态市盈率-行
		fmt.Printf("\t%4.1f/%-4.1f", value.gzfx.PS9, value.gzfx.HY_PS9)                    //市销率-行
		fmt.Printf("\t%4s", single.THBJ.GZBJ.DATA[0].PEG)                                  //PEG
		fmt.Printf("\t%4s", strings.ReplaceAll(single.THBJ.GZBJ.DATA[0].PM, "U003E", ">")) //估值排名

		//	===============主力
		//社保/流通股 >= 3%   机构占流通股比例>40%
		var SBZB = ""
		var JGZB = ""
		var JGTJ = ""
		{

			for i := 0; i < len(single.Gbyj.ZLCC); i++ {
				if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "社保") && "--" != single.Gbyj.ZLCC[i].ZLTGBL {
					var sbbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
					if sbbl >= 3 {
						SBZB = fmt.Sprintf("%.1f", sbbl)
					}
				} else if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "合计") {
					var hjbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
					if hjbl >= 50 {
						JGZB = fmt.Sprintf("%.1f", hjbl)
					}
				}
			}
			if single.JGTJ >= 20 {
				JGTJ = fmt.Sprintf("%d", single.JGTJ)
			} else {
				JGTJ = fmt.Sprintf("")
			}
		}
		fmt.Printf("\t%4s\t%4s\t%3s", SBZB, JGZB, JGTJ) // "社/流", "机/流", "推荐数"
		fmt.Println("")
	}
	fmt.Println("符合条件的股票有=", len(mapStock), "个")
}

func ChineseCount2(str1 string) (count int) {
	for _, char := range str1 {
		if unicode.Is(unicode.Han, char) {
			count++
		}
	}

	return
}
