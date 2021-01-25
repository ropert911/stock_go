package main

import (
	"fmt"
	"math"
	"stock"
	"strings"
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
	mapStock := parserData2()
	exportResult2(mapStock, []string{"002918"})
}

//解析所有数据-估值分析、资产负债、业绩报表、利润表
func parserData2() map[string]StockInfo2 {
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

//显示结果
func exportResult2(mapStock map[string]StockInfo2, codes []string) {
	var tip1 = fmt.Sprintf("%c[;;36m%s%c[0m", 0x1B, "优秀", 0x1B)
	var tip2 = fmt.Sprintf("%c[;;34m%s%c[0m", 0x1B, "良好", 0x1B)
	var tip3 = "一般"
	var tip4 = fmt.Sprintf("%c[;;31m%s%c[0m", 0x1B, "不好", 0x1B)
	var tip5 = fmt.Sprintf("%c[;;31m%s%c[0m", 0x1B, "补充", 0x1B)
	var tips = ""

	for i := range codes {
		var key = codes[i]
		var value = mapStock[codes[i]]
		var single, _ = stock.ParseSingle(codes[i], value.gzfx.ORIGINALCODE)
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

		//	===============基本信息 编码、名称、"股价", "行业"
		fmt.Printf("\t%6s\t%-s\t股价: %6.2f\t%-s\n", key, value.zcfz.SECURITY_NAME_ABBR, value.gzfx.NEW, value.gzfx.HYName) //
		fmt.Println("有增长-有效率-有钱")
		fmt.Printf("\t近三年收入: \t%s\t%s\t%s\n", zyzbP2.YYZSR, zyzbP1.YYZSR, zyzbP0.YYZSR)
		var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
		var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
		var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
		var avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
		if avg >= 20 {
			tips = tip1
		} else if avg >= 15 {
			tips = tip2
		} else if avg >= 10 {
			tips = tip3
		} else {
			tips = tip4
		}
		fmt.Printf("\t★近三年收入增长: \t%s%%\t%s%%\t%s%%\t平均:%.2f%%\t\t%s\t至少要求 %s\n", zyzbP2.YYZSRTBZZ, zyzbP1.YYZSRTBZZ, zyzbP0.YYZSRTBZZ, avg, tips, tip2)
		if tb0 >= 40 {
			tips = tip1
		} else if tb0 >= 20 {
			tips = tip2
		} else if tb0 >= 10 {
			tips = tip3
		} else {
			tips = tip4
		}
		fmt.Printf("\t★近一年收入增长: \t%.2f%%\t\t\t\t\t\t\t\t\t%s\t至少要求 %s\n", tb0, tips, tip3)
		tb0 = stock.ToFloat(zyzbP0.MLL)
		tb1 = stock.ToFloat(zyzbP1.MLL)
		tb2 = stock.ToFloat(zyzbP2.MLL)
		avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
		if avg >= 50 {
			tips = tip1
		} else if avg >= 30 {
			tips = tip2
		} else if avg >= 20 {
			tips = tip3
		} else {
			tips = tip4
		}
		fmt.Printf("\t★近三年毛利率: \t%s%%\t%s%%\t%s%%\t平均:%.2f%%\t\t\t%s\t至少要求 %s\n", zyzbP2.MLL, zyzbP1.MLL, zyzbP0.MLL, avg, tips, tip3)

		fmt.Printf("\t近三年净利润: \t%s\t%s\t%s\n", zyzbP2.GSJLR, zyzbP1.GSJLR, zyzbP0.GSJLR)
		tb0 = stock.ToFloat(zyzbP0.GSJLRTBZZ)
		tb1 = stock.ToFloat(zyzbP1.GSJLRTBZZ)
		tb2 = stock.ToFloat(zyzbP2.GSJLRTBZZ)
		avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
		if avg >= 20 {
			tips = tip1
		} else if avg >= 15 {
			tips = tip2
		} else if avg >= 10 {
			tips = tip3
		} else {
			tips = tip4
		}
		fmt.Printf("\t★近三年净利润同比: \t%s%%\t%s%%\t%s%%\t平均:%.2f%%\t\t%s\t至少要求 %s\n", zyzbP2.GSJLRTBZZ, zyzbP1.GSJLRTBZZ, zyzbP0.GSJLRTBZZ, avg, tips, tip3)
		fmt.Printf("\t成长性排名: %4s\n", strings.ReplaceAll(single.THBJ.CZXBJ.DATA[0].PM, "U003E", ">"))

		if value.yjbb.WEIGHTAVG_ROE >= 20 {
			tips = tip1
		} else if value.yjbb.WEIGHTAVG_ROE >= 12 {
			tips = tip2
		} else if value.yjbb.WEIGHTAVG_ROE >= 8 {
			tips = tip3
		} else {
			tips = tip4
		}
		fmt.Printf("\t★ROE净资产收益率:\t%5.2f\t\t\t\t\t\t\t\t\t%s\t至少要求 %s\n", value.yjbb.WEIGHTAVG_ROE, tips, tip2)
		var yincome = stock.ToFloat(strings.ReplaceAll(zyzbP1.YYZSR, "亿", "")) * 10000 * 10000
		var xjbl = 100 * float64(value.zcfz.MONETARYFUNDS) / yincome
		if xjbl >= 40 {
			tips = tip1
		} else if xjbl >= 25 {
			tips = tip2
		} else if xjbl >= 15 {
			tips = tip3
		} else {
			tips = tip4
		}
		fmt.Printf("\t★货币资金/去年营业收入:\t%6.1f亿/%-6s=%.2f%%\t\t\t\t%s\t至少要求 %s\n", value.zcfz.MONETARYFUNDS/(10000*10000), zyzbP1.YYZSR, xjbl, tips, tip2)

		fmt.Println("估值")
		fmt.Printf("\t静态市盈率-行: %5.1f/%-5.1f\n", value.gzfx.PE7, value.gzfx.HY_PE7)
		fmt.Printf("\t动态市盈率-行: %5.1f/%-5.1f\n", value.gzfx.PE9, value.gzfx.HY_PE9)
		fmt.Printf("\tPEG: %4s\n", single.THBJ.GZBJ.DATA[0].PEG)
		fmt.Printf("\t估值排名: %4s\n", strings.ReplaceAll(single.THBJ.GZBJ.DATA[0].PM, "U003E", ">"))
		fmt.Printf("\t★人工估值分析 \t\t\t\t\t\t\t\t\t\t\t\t%4s\t>=50的:得看是不是行业是不是战略转换类，增长很明确,增长也很快\n", tip5)
		fmt.Printf("\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t30-50的:行业增长也明确,企业增长也快,只是行业战略性不高\n")
		fmt.Printf("\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t30以下的，企业增长快就好，如果行业也不错最好\n")
		fmt.Printf("\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t动态小于静态 表明今天比昨天低\n")

		fmt.Println("主力")
		var SBZB = float64(0)
		var JGZB = float64(0)
		var JGTJ = 0
		{

			for i := 0; i < len(single.Gbyj.ZLCC); i++ {
				if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "社保") && "--" != single.Gbyj.ZLCC[i].ZLTGBL {
					var sbbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
					SBZB = sbbl
				} else if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "合计") {
					var hjbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
					JGZB = hjbl
				}
			}

			JGTJ = single.JGTJ
		}
		if SBZB >= 4 {
			tips = tip1
		} else if SBZB >= 1 {
			tips = tip2
		} else {
			tips = tip3
		}
		fmt.Printf("\t☆社/流: %.2f\t\t\t\t\t\t\t\t\t\t\t\t%s\n", SBZB, tips)
		fmt.Printf("\t机/流: %.2f\n", JGZB)
		if JGTJ >= 30 {
			tips = tip1
		} else if JGTJ >= 10 {
			tips = tip2
		} else {
			tips = tip3
		}
		fmt.Printf("\t☆推荐数: %d\t\t\t\t\t\t\t\t\t\t\t\t\t%s\n", JGTJ, tips)

		fmt.Printf("\t☆不同行业收入的组成(是否专一): \t\t\t\t\t\t\t\t%s\t%s\n", tip5, "F10 看经营分析")
		fmt.Printf("\t☆不同地区海内外收入比例: \t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "F10 看经营分析")
		fmt.Printf("\t☆研发投入/营业收入: \t\t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "http://f10.eastmoney.com/NewFinanceAnalysis/Index?type=web&code=SZ002918#lrb-0")
		fmt.Printf("\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t%s\n", "%3以上不错  10%以上就很有希望了")
		fmt.Printf("\t★同行毛利率对比: \t\t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "要求前3 (同花顺 公司概况-经营分析-按产品)")
		fmt.Printf("\t不同产品的收入增长: \t\t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "(同花顺 公司概况-经营分析-按产品)")
		fmt.Printf("\t☆高管增减持分析: \t\t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "(爱问财 搜如：三一重工高官持股变化点评)")
		fmt.Printf("\t☆机构推荐数&调研数: \t\t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "(同花顺 市场观点-机构评级&机构调研)")
		fmt.Printf("\t★财报分析: \t\t\t\t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "(同花顺 财务分析-分析)")
		fmt.Printf("\t★行业和题材分析: \t\t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "(同花顺 市场观点-机构调研中的调研报告)")
		fmt.Printf("\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t%s\n", "更多的要长期积累")
		fmt.Println("")
	}
}
