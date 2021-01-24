package main

import (
	"fmt"
	"math"
	"stock"
	"strings"
	"unicode"
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
		fmt.Println("收入与利润增长")
		fmt.Printf("\t近三年收入: \t%s\t%s\t%s\n", zyzbP2.YYZSR, zyzbP1.YYZSR, zyzbP0.YYZSR)
		var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
		var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
		var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
		var avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
		fmt.Printf("\t近三年收入增长: \t%s%%\t%s%%\t%s%%\t平均:%.2f%%\n", zyzbP2.YYZSRTBZZ, zyzbP1.YYZSRTBZZ, zyzbP0.YYZSRTBZZ, avg)
		fmt.Printf("\t近三年毛利率: \t%s%%\t%s%%\t%s%%\n", zyzbP2.MLL, zyzbP1.MLL, zyzbP0.MLL)
		fmt.Printf("\t近三年净利润: \t%s\t%s\t%s\n", zyzbP2.GSJLR, zyzbP1.GSJLR, zyzbP0.GSJLR)
		tb0 = stock.ToFloat(zyzbP0.GSJLRTBZZ)
		tb1 = stock.ToFloat(zyzbP1.GSJLRTBZZ)
		tb2 = stock.ToFloat(zyzbP2.GSJLRTBZZ)
		avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
		fmt.Printf("\t近三年净利润同比: \t%s%%\t%s%%\t%s%%\n", zyzbP2.GSJLRTBZZ, zyzbP1.GSJLRTBZZ, zyzbP0.GSJLRTBZZ)
		fmt.Printf("\t成长性排名: %4s\n", strings.ReplaceAll(single.THBJ.CZXBJ.DATA[0].PM, "U003E", ">"))

		fmt.Println("资金有效率有钱")
		fmt.Printf("\tROE净资产收益率:\t%5.2f\n", value.yjbb.WEIGHTAVG_ROE)
		var yincome = stock.ToFloat(strings.ReplaceAll(zyzbP1.YYZSR, "亿", "")) * 10000 * 10000
		fmt.Printf("\t货币资金/去年营业收入:\t%.1f亿/%s=%.2f%%\n", value.zcfz.MONETARYFUNDS/(10000*10000), zyzbP1.YYZSR, 100*float64(value.zcfz.MONETARYFUNDS)/yincome)

		fmt.Println("估值")
		fmt.Printf("\t静态市盈率-行: %5.1f/%-5.1f\n", value.gzfx.PE7, value.gzfx.HY_PE7)
		fmt.Printf("\t动态市盈率-行: %5.1f/%-5.1f\n", value.gzfx.PE9, value.gzfx.HY_PE9)
		fmt.Printf("\tPEG: %4s\n", single.THBJ.GZBJ.DATA[0].PEG)
		fmt.Printf("\t估值排名: %4s\n", strings.ReplaceAll(single.THBJ.GZBJ.DATA[0].PM, "U003E", ">"))

		fmt.Println("主力")
		var SBZB = ""
		var JGZB = ""
		var JGTJ = ""
		{

			for i := 0; i < len(single.Gbyj.ZLCC); i++ {
				if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "社保") && "--" != single.Gbyj.ZLCC[i].ZLTGBL {
					var sbbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
					if sbbl >= 2 {
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
		fmt.Printf("\t社/流: %4s\n", SBZB)
		fmt.Printf("\t机/流: %4s\n", JGZB)
		fmt.Printf("\t推荐数: %3s\n", JGTJ)
		fmt.Println("")
	}
}

func ChineseCount2(str1 string) (count int) {
	for _, char := range str1 {
		if unicode.Is(unicode.Han, char) {
			count++
		}
	}

	return
}
