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
	exportResult2(mapStock, []string{"300815"})
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
	var tip5 = fmt.Sprintf("%c[;;32m%s%c[0m", 0x1B, "补充", 0x1B)
	var tips = ""

	for i := range codes {
		var key = codes[i]
		var value = mapStock[codes[i]]
		var single, _ = stock.ParseSingle(codes[i], value.gzfx.ORIGINALCODE)
		var lrbP0 = single.LRB[0]
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

		fmt.Println("------------------------------------------------------------------------------------------")
		//	===============基本信息 编码、名称、"股价", "行业"
		fmt.Printf("\t%6s\t%-s\t股价: %6.2f元\t行业%-s\n", key, value.zcfz.SECURITY_NAME_ABBR, value.gzfx.NEW, value.gzfx.HYName)
		fmt.Println("=======1 收入-利润 增长快")
		var avg = avgGrow(stock.ToFloat(zyzbP2.YYZSRTBZZ), stock.ToFloat(zyzbP1.YYZSRTBZZ), stock.ToFloat(zyzbP0.YYZSRTBZZ))
		var yincomeY = incomePre(value.yjbb.TOTAL_OPERATE_INCOME, value.zcfz.REPORT_DATE)
		fmt.Printf("\t★近三年收入增长:\t%s(%s%%)\t%s(%s%%)\t今年%s(%s%%%s)\t预估今年%.2f亿\t平均增长:%.2f%%(%s)\t至少要求 %s\n",
			zyzbP2.YYZSR, zyzbP2.YYZSRTBZZ,
			zyzbP1.YYZSR, zyzbP1.YYZSRTBZZ,
			zyzbP0.YYZSR, zyzbP0.YYZSRTBZZ, curGrowGrade(stock.ToFloat(zyzbP0.YYZSRTBZZ)), yincomeY/(10000*10000),
			avg, avgGrowGrade(avg),
			tip2)

		avg = avgGrow(stock.ToFloat(zyzbP2.GSJLRTBZZ), stock.ToFloat(zyzbP1.GSJLRTBZZ), stock.ToFloat(zyzbP0.GSJLRTBZZ))
		fmt.Printf("\t★近三年净利润增长: \t%s(%s%%)\t%s(%s%%)\t今年%s(%s%%%s)\t\t\t\t\t平均增长:%.2f%%(%s)\t至少要求 %s\n",
			zyzbP2.GSJLR, zyzbP2.GSJLRTBZZ,
			zyzbP1.GSJLR, zyzbP1.GSJLRTBZZ,
			zyzbP0.GSJLR, zyzbP0.GSJLRTBZZ, avgProfGrowGrade(stock.ToFloat(zyzbP0.GSJLRTBZZ)),
			avg, avgProfGrowGrade(avg),
			tip2)
		fmt.Printf("\t成长性排名: %4s\n", strings.ReplaceAll(single.THBJ.CZXBJ.DATA[0].PM, "U003E", ">"))

		fmt.Println("=======2 同行毛利高-资金效率-现金多")
		avg = avgGrow(stock.ToFloat(zyzbP2.MLL), stock.ToFloat(zyzbP1.MLL), stock.ToFloat(zyzbP0.MLL))
		fmt.Printf("\t近三年毛利率: \t%s%%\t%s%%\t今年%s%%\t平均:%.2f%%(%s)\t\t\t\t至少要求 %s\n",
			zyzbP2.MLL, zyzbP1.MLL, zyzbP0.MLL,
			avg, mllGrade(avg),
			tip3)
		fmt.Printf("\t★ROE净资产收益率:\t%5.2f%%(%s)\t\t\t\t\t\t\t\t\t\t\t至少要求 %s\n",
			value.yjbb.WEIGHTAVG_ROE, roeGrade(value.yjbb.WEIGHTAVG_ROE), tip2)
		var yincome = stock.ToFloat(strings.ReplaceAll(zyzbP1.YYZSR, "亿", "")) * 10000 * 10000
		var xjbl = 100 * float64(value.zcfz.MONETARYFUNDS) / yincome
		fmt.Printf("\t★货币资金/去年营业收入:\t%6.1f亿/%-6s=%.2f%%(%s)\t\t\t\t\t\t至少要求 %s\n",
			value.zcfz.MONETARYFUNDS/(10000*10000), zyzbP1.YYZSR, xjbl, xjblGrade(xjbl),
			tip2)
		fmt.Printf("\t每股公积金/现价:%.2f/%.2f=%.1f%%\t\t每股未分配/现价:%.2f/%.2f=%.1f%%\t\t\t%s\n",
			stock.ToFloat(single.ZYZB[0].MGGJJ), float64(value.gzfx.NEW), 100*stock.ToFloat(single.ZYZB[0].MGGJJ)/float64(value.gzfx.NEW),
			stock.ToFloat(single.ZYZB[0].MGWFPLY), value.gzfx.NEW, 100*stock.ToFloat(single.ZYZB[0].MGWFPLY)/float64(value.gzfx.NEW),
			"%20以上值得重点关注")

		fmt.Println("=======3 产品分析")
		fmt.Printf("\t☆不同产品收入组成:(F10 看经营分析) \t\t\t\t\t\t\t\t★同业产品毛利率对比 (同花顺 公司概况-经营分析-按产品)\t\t\t\t%s\n", tip5)
		fmt.Printf("\t不同产品的收入增长:(同花顺 公司概况-经营分析-按产品) \t\t\t\t☆不同地区海内外收入比例: (F10 看经营分析)\t\t\t\t\t\t\t%s\n", tip5)

		fmt.Println("=======4 公司发展-行业分析")
		fmt.Printf("\t★研发投入/净利润: \t\t\t%.2f亿/%.2f亿=\t%-5.2f%%\t\t\t\t\t%s\n",
			stock.ToFloat(lrbP0.RDEXP)/(10000*10000), stock.ToFloat(lrbP0.PARENTNETPROFIT)/(10000*10000), 100*stock.ToFloat(lrbP0.RDEXP)/stock.ToFloat(lrbP0.PARENTNETPROFIT),
			"%20以上不错  30%以上很重视了")
		fmt.Printf("\t★研发投入/营业收入: \t\t\t%.2f亿/%.2f亿=\t%5.2f%%\t\t\t\t\t%s\n",
			stock.ToFloat(lrbP0.RDEXP)/(10000*10000), stock.ToFloat(lrbP0.TOTALOPERATEREVE)/(10000*10000), 100*stock.ToFloat(lrbP0.RDEXP)/stock.ToFloat(lrbP0.TOTALOPERATEREVE),
			"%3以上不错  10%以上就很有希望了")
		fmt.Printf("\t★财报分析 & 行业和题材分析: \t%s\t\t\t\t\t\t\t\t\t%s\n", tip5, "(同花顺 财务分析-分析)(同花顺研报-研报摘要)(同花顺 市场观点-机构调研中的调研报告)")

		fmt.Println("=======5 高管、机构都看好")
		fmt.Printf("\t☆高管增减持分析: \t\t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "(爱问财 搜如：三一重工高官持股变化点评)")
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
		fmt.Printf("\t☆社/流: %.2f%%\t\t\t\t\t\t\t\t\t\t\t\t%s\n", SBZB, tips)
		fmt.Printf("\t机/流: %.2f%%\n", JGZB)
		if JGTJ >= 30 {
			tips = tip1
		} else if JGTJ >= 10 {
			tips = tip2
		} else {
			tips = tip3
		}
		fmt.Printf("\t☆近6月推荐买入数: %d\t\t\t\t\t\t\t\t\t\t%s\n", JGTJ, tips)
		fmt.Printf("\t☆机构推荐数1-6月&调研数: \t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "(同花顺 市场观点-机构评级&机构调研)")
		fmt.Printf("\t☆机构目标价&业绩预测: \t\t\t\t\t\t\t\t\t\t%s\t%s\n", tip5, "(同花顺 研报-研报数据")

		fmt.Println("=======6 合理估值看时机")
		fmt.Printf("\t静市盈率/行: %5.1f/%-5.1f\t动市盈率/行: %5.1f/%-5.1f\n",
			value.gzfx.PE7, value.gzfx.HY_PE7,
			value.gzfx.PE9, value.gzfx.HY_PE9)
		fmt.Printf("\tPEG: %4s\t\t估值排名: %4s \n", single.THBJ.GZBJ.DATA[0].PEG, strings.ReplaceAll(single.THBJ.GZBJ.DATA[0].PM, "U003E", ">"))
		fmt.Printf("\t☆股东人数 %s=%s %s=%s %s=%s 当前和前2期比:%.2f%%\t%s\n",
			single.Gbyj.GDRS[0].RQ, single.Gbyj.GDRS[0].GDRS,
			single.Gbyj.GDRS[1].RQ, single.Gbyj.GDRS[1].GDRS,
			single.Gbyj.GDRS[2].RQ, single.Gbyj.GDRS[2].GDRS,
			100*stock.ToFloat(strings.ReplaceAll(single.Gbyj.GDRS[0].GDRS, "万", ""))/stock.ToFloat(strings.ReplaceAll(single.Gbyj.GDRS[2].GDRS, "万", ""))-100,
			"-20%就已经比较集中了")

		fmt.Println()
		//fmt.Println("总结：")
		//fmt.Println("过去的成长性：")
		//fmt.Println("同行比收入高不高赚钱多不多(是不是有龙头的特点)：")
		//fmt.Println("公司基本情况是否专一、有没有海外市场：")
		//fmt.Println("公司内生动力-行业前景是否支持公司大发展：")
		//fmt.Println("大家是不是都看好公司(高官、机构)")
		//fmt.Println("估值是否合理，当前是否适合进入")
		fmt.Println("")
	}
}

var tip1 = fmt.Sprintf("%c[;;36m%s%c[0m", 0x1B, "优秀", 0x1B)
var tip2 = fmt.Sprintf("%c[;;34m%s%c[0m", 0x1B, "良好", 0x1B)
var tip3 = "一般"
var tip4 = fmt.Sprintf("%c[;;31m%s%c[0m", 0x1B, "不好", 0x1B)
var tip5 = fmt.Sprintf("%c[;;32m%s%c[0m", 0x1B, "补充", 0x1B)

//今年收入增长评级
func curGrowGrade(tb0 float64) string {
	var tips = ""
	if tb0 >= 40 {
		tips = tip1
	} else if tb0 >= 20 {
		tips = tip2
	} else if tb0 >= 10 {
		tips = tip3
	} else {
		tips = tip4
	}
	return tips
}

//今年收入收季度预估
func incomePre(yincomeY float32, jbsj string) float32 {
	var _, date = stock.GetDate(jbsj)
	if date == 3 {
		yincomeY = yincomeY * 4
	} else if date == 6 {
		yincomeY = yincomeY * 2
	} else if date == 9 {
		yincomeY = yincomeY * 4 / 3
	}
	return yincomeY
}

//计算几看平均增长
func avgGrow(tb2 float64, tb1 float64, tb0 float64) float64 {
	return math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
}

//几年平均收入增长评级
func avgGrowGrade(avg float64) string {
	var tips = ""
	if avg >= 20 {
		tips = tip1
	} else if avg >= 15 {
		tips = tip2
	} else if avg >= 10 {
		tips = tip3
	} else {
		tips = tip4
	}
	return tips
}

//平均利润增长评级
func avgProfGrowGrade(avg float64) string {
	var tips = ""
	if avg >= 20 {
		tips = tip1
	} else if avg >= 15 {
		tips = tip2
	} else if avg >= 10 {
		tips = tip3
	} else {
		tips = tip4
	}
	return tips
}

//现金比例评级
func xjblGrade(xjbl float64) string {
	var tips = ""
	if xjbl >= 40 {
		tips = tip1
	} else if xjbl >= 20 {
		tips = tip2
	} else if xjbl >= 15 {
		tips = tip3
	} else {
		tips = tip4
	}
	return tips
}

//roe评级
func roeGrade(roe float32) string {
	var tips = ""
	if roe >= 20 {
		tips = tip1
	} else if roe >= 12 {
		tips = tip2
	} else if roe >= 8 {
		tips = tip3
	} else {
		tips = tip4
	}
	return tips
}

//毛利率评级
func mllGrade(avg float64) string {
	var tips = ""
	if avg >= 50 {
		tips = tip1
	} else if avg >= 30 {
		tips = tip2
	} else if avg >= 20 {
		tips = tip3
	} else {
		tips = tip4
	}
	return tips
}