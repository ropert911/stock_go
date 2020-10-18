package main

import (
	"fmt"
	"math"
	"sort"
	"stock"
	"strings"
	"unicode"
	"util/file"
)

type StockInfo struct {
	gzfx stock.StockGzfx
	yjbb stock.StockYjbb
	zcfz stock.StockZcfz
	lrb  stock.StockLrb
	xjll stock.StockXjll
}

func main() {
	mapStock := parserData()
	mapStock = filterData(mapStock)
	exportResult(mapStock)
}

//解析所有数据-估值分析、资产负债、业绩报表、利润表
func parserData() map[string]StockInfo {
	mapStock := make(map[string]StockInfo)
	socksgzfx := stock.ReadStockGzfx()
	socksYjbb := stock.ReadStockYjbb()
	socksZcfz := stock.ReadStockZcfz()
	stockLrb := stock.ReadStockLrb()
	stockXjll := stock.ReadStockXjll()
	for i := 0; i < len(socksgzfx); i++ {
		var stockInfo StockInfo
		stockInfo.gzfx = socksgzfx[i]
		for j := 0; j < len(socksYjbb); j++ {
			if socksgzfx[i].SECURITYCODE == socksYjbb[j].SECURITY_CODE {
				stockInfo.yjbb = socksYjbb[j]
				break
			}
		}
		for z := 0; z < len(socksZcfz); z++ {
			if socksgzfx[i].SECURITYCODE == socksZcfz[z].SECURITY_CODE {
				stockInfo.zcfz = socksZcfz[z]
				break
			}
		}
		for z := 0; z < len(stockLrb); z++ {
			if socksgzfx[i].SECURITYCODE == stockLrb[z].SECURITY_CODE {
				stockInfo.lrb = stockLrb[z]
				break
			}
		}
		for z := 0; z < len(stockXjll); z++ {
			if socksgzfx[i].SECURITYCODE == stockXjll[z].SECURITY_CODE {
				stockInfo.xjll = stockXjll[z]
				break
			}
		}

		mapStock[socksgzfx[i].SECURITYCODE] = stockInfo
	}

	return mapStock
}

//过滤掉不符合条件的
func filterData(mapStock map[string]StockInfo) map[string]StockInfo {
	for key, value := range mapStock {
		//1有积累---非ST
		if strings.HasPrefix(value.gzfx.SName, "*ST") || strings.HasPrefix(value.gzfx.SName, "ST") {
			delete(mapStock, key)
			continue
		}
		//1有积累--- 现价不为空
		if value.gzfx.NEW < 5 {
			delete(mapStock, key)
			continue
		}
		//1有积累 --- 每股净资产>1 && 不为空
		if value.yjbb.BPS < 1 {
			delete(mapStock, key)
			continue
		}

		//2估值 --- 市盈率(动静)>0不为空
		if value.gzfx.PE9 <= 0 || value.gzfx.PE7 <= 0 {
			delete(mapStock, key)
			continue
		}
		//2估值 --- 市销率<10
		if value.gzfx.PS9 >= 10 {
			delete(mapStock, key)
			continue
		}

		//3成长 -- 营业收入>0不为空
		if value.yjbb.TOTAL_OPERATE_INCOME <= 0 {
			delete(mapStock, key)
			continue
		}
		//3成长 --净利润>0.2亿
		if value.yjbb.PARENT_NETPROFIT < 20000000 {
			delete(mapStock, key)
			continue
		}
		//3成长 --扣非净利润>0.2亿
		if value.lrb.DEDUCT_PARENT_NETPROFIT < 20000000 {
			delete(mapStock, key)
			continue
		}
		//3成长 -- 收入同比>=15%
		if value.yjbb.YSTZ < 15 {
			delete(mapStock, key)
			continue
		}

		//4财报 -- 资产负债率<70%
		if value.zcfz.DEBT_ASSET_RATIO >= 70 {
			delete(mapStock, key)
			continue
		}
		//4财报 -- 资产负债率>65%的，净资产收益率>10%
		if value.zcfz.DEBT_ASSET_RATIO >= 65 && value.yjbb.WEIGHTAVG_ROE < 10 {
			delete(mapStock, key)
			continue
		}
		//4财报 -- 总现金流>0.2亿
		if value.xjll.CCE_ADD < 20000000 {
			delete(mapStock, key)
			continue
		}

		stock.DownloadSingle(key)
		var single, err = stock.ParseSingle(key)
		if nil != err {
			fmt.Println("Error parse single stock data ", err)
			continue
		}
		//1有积累 -- 每股公积金>0.5元
		if stock.ToFloat(single.ZYZB[0].MGGJJ) < 0.5 {
			delete(mapStock, key)
			continue
		}
		//1有积累 -- 未分配利润 >1亿
		if stock.ToFloat(single.ZCFZ[0].RETAINEDEARNING) < 100000000 {
			delete(mapStock, key)
			continue
		}
		//1有积累 -- 每股未分配利润>0.5元
		if stock.ToFloat(single.ZYZB[0].MGWFPLY) < 0.5 {
			delete(mapStock, key)
			continue
		}

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

		//3成长 -- 近3年利润不为负
		{
			var jlr0 = stock.ToFloat(single.ZYZB[0].GSJLR)
			var jlr1 = stock.ToFloat(zyzbP1.GSJLR)
			var jlr2 = stock.ToFloat(zyzbP2.GSJLR)
			if jlr0 < 0 || jlr1 < 0 || jlr2 < 0 {
				delete(mapStock, key)
				continue
			}
		}
		//3成长 -- 前2年>10% 或近 3年平均20%
		{
			var tb0 = stock.ToFloat(single.ZYZB[0].YYZSRTBZZ)
			var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
			var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
			var avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
			if tb0 < 15 { //当年大于15%
				delete(mapStock, key)
				continue
			}
			if tb1 < 0 || tb2 < 0 { //近3年不能有负增长
				delete(mapStock, key)
				continue
			}
			if avg < 15 {
				delete(mapStock, key)
				continue
			}
		}
		//3成长 -- 毛利率>=10%
		var mll = stock.ToFloat(single.ZYZB[0].MLL)
		if mll < 10 {
			delete(mapStock, key)
			continue
		}
		//3成长 -- 净利润率>=5%
		var jll = stock.ToFloat(single.ZYZB[0].JLL)
		if jll < 5 {
			delete(mapStock, key)
			continue
		}
		//5其它 -- 近2年<150 每年小于60%
		if single.THIS >= 60 || single.TWO >= 150 {
			delete(mapStock, key)
			continue
		}
	}

	return mapStock
}

//生成导出信息
func CreateExportEBK(mapStock map[string]StockInfo) {
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

type SockInfoShow struct {
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
	SJL  string //市净率
	PEJT string //静态市盈率
	PEDT string //动态市盈率
	PS9  string //市销率
	RPB8 string //市净率估值
	RPE7 string //PE(静)估值
	RPE9 string //PE(TTM)估值
	RPS9 string //市销率估值
	//主力研究
	QSDGDCGHJ   string  //前十大股东持股合计
	QSDLTGDCGHJ string  //前十大流通股东持股合计
	SBZB        string  //社保占流通比
	JGZB        float64 //机构合计占流通比
	JGTJ        string  //机构推荐数
}
type SockInfoShowArray []SockInfoShow

func (s SockInfoShowArray) Len() int           { return len(s) }
func (s SockInfoShowArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SockInfoShowArray) Less(i, j int) bool { return s[i].JGZB >= s[j].JGZB }

//显示结果
func exportResult(mapStock map[string]StockInfo) {
	CreateExportEBK(mapStock)

	var sockInfoShows = SockInfoShowArray{}

	//把符合条件股票中的亮点数据显示出来
	for key, value := range mapStock {
		stock.DownloadSingle(key)
		var single, _ = stock.ParseSingle(key)

		var sockInfoShow SockInfoShow
		sockInfoShow.Code = key
		sockInfoShow.Name = value.gzfx.SName
		sockInfoShow.Price = fmt.Sprint(value.gzfx.NEW)
		sockInfoShow.HYName = value.gzfx.HYName

		//2积累--每股公积金
		var mggjj = stock.ToFloat(single.ZYZB[0].MGGJJ)
		if mggjj*3 > float64(value.gzfx.NEW) || mggjj > 10 {
			sockInfoShow.MGGJJ = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, mggjj, 0x1B)
		} else if mggjj*4 > float64(value.gzfx.NEW) || mggjj > 8 {
			sockInfoShow.MGGJJ = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, mggjj, 0x1B)
		} else {
			sockInfoShow.MGGJJ = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, mggjj, 0x1B)
		}
		//1积累--每股未分配利润
		var mgwfply = stock.ToFloat(single.ZYZB[0].MGWFPLY)
		if mgwfply*3 > float64(value.gzfx.NEW) || mgwfply > 10 {
			sockInfoShow.MGWFPLY = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, mgwfply, 0x1B)
		} else if mgwfply*4 > float64(value.gzfx.NEW) || mgwfply > 8 {
			sockInfoShow.MGWFPLY = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, mgwfply, 0x1B)
		} else {
			sockInfoShow.MGWFPLY = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, mgwfply, 0x1B)
		}

		//2估值-市净率
		if value.gzfx.PB8 <= 5 {
			sockInfoShow.SJL = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.gzfx.PB8, 0x1B)
		} else {
			sockInfoShow.SJL = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PB8, 0x1B)
		}
		//2估值-动态市盈率
		if value.gzfx.PE9 <= 15 {
			sockInfoShow.PEDT = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, value.gzfx.PE9, 0x1B)
		} else if value.gzfx.PE9 <= 25 {
			sockInfoShow.PEDT = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.gzfx.PE9, 0x1B)
		} else {
			sockInfoShow.PEDT = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PE9, 0x1B)
		}
		//2估值-静态市盈率
		if value.gzfx.PE7 <= 15 {
			sockInfoShow.PEJT = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, value.gzfx.PE7, 0x1B)
		} else if value.gzfx.PE7 <= 25 {
			sockInfoShow.PEJT = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.gzfx.PE7, 0x1B)
		} else {
			sockInfoShow.PEJT = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PE7, 0x1B)
		}
		//2估值-市销率
		if value.gzfx.PS9 <= 2 {
			sockInfoShow.PS9 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, value.gzfx.PS9, 0x1B)
		} else if value.gzfx.PS9 <= 3 {
			sockInfoShow.PS9 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.gzfx.PS9, 0x1B)
		} else {
			sockInfoShow.PS9 = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PS9, 0x1B)
		}

		//3成长-净益率 净资产收益率
		if value.yjbb.WEIGHTAVG_ROE >= 20 {
			sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, value.yjbb.WEIGHTAVG_ROE, 0x1B)
		} else if value.yjbb.WEIGHTAVG_ROE >= 10 {
			sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.yjbb.WEIGHTAVG_ROE, 0x1B)
		} else {
			sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.yjbb.WEIGHTAVG_ROE, 0x1B)
		}
		//3成长-毛利率
		var mll = stock.ToFloat(single.ZYZB[0].MLL)
		if mll >= 30 {
			sockInfoShow.MLL = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, mll, 0x1B)
		} else {
			sockInfoShow.MLL = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, mll, 0x1B)
		}
		//3成长-净利润率
		var jll = stock.ToFloat(single.ZYZB[0].JLL)
		if jll >= 20 {
			sockInfoShow.JLL = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, jll, 0x1B)
		} else {
			sockInfoShow.JLL = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, jll, 0x1B)
		}
		//3成长 -- 3年平均>=20%
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
		var tb0 = stock.ToFloat(single.ZYZB[0].YYZSRTBZZ)
		var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
		var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
		var avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
		if avg >= 30 {
			sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, avg, 0x1B)
		} else if avg >= 20 {
			sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, avg, 0x1B)
		} else {
			sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, avg, 0x1B)
		}

		var rPB8 = value.gzfx.PB8 / value.gzfx.HY_PB8
		if rPB8 < 0.7 {
			sockInfoShow.RPB8 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, rPB8, 0x1B)
		} else if rPB8 < 1 {
			sockInfoShow.RPB8 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPB8, 0x1B)
		} else {
			sockInfoShow.RPB8 = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, rPB8, 0x1B)
		}

		var rPE7 = value.gzfx.PE7 / value.gzfx.HY_PE7
		if rPE7 < 0.8 {
			sockInfoShow.RPE7 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, rPE7, 0x1B)
		} else if rPE7 < 1 {
			sockInfoShow.RPE7 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPE7, 0x1B)
		} else {
			sockInfoShow.RPE7 = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, rPE7, 0x1B)
		}

		var rPE9 = value.gzfx.PE9 / value.gzfx.HY_PE9
		if rPE9 < 0.8 {
			sockInfoShow.RPE9 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, rPE9, 0x1B)
		} else if rPE9 < 1 {
			sockInfoShow.RPE9 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPE9, 0x1B)
		} else {
			sockInfoShow.RPE9 = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, rPE9, 0x1B)
		}

		var rPS9 = value.gzfx.PS9 / value.gzfx.HY_PS9
		if rPS9 < 0.8 {
			sockInfoShow.RPS9 = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, rPS9, 0x1B)
		} else if rPS9 < 1 {
			sockInfoShow.RPS9 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPS9, 0x1B)
		} else {
			sockInfoShow.RPS9 = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, rPS9, 0x1B)
		}

		//前十大股东持股合计
		if len(single.Gbyj.GDRS) > 0 && "--" != single.Gbyj.GDRS[0].QSDGDCGHJ {
			var sdgdzb = stock.ToFloat(single.Gbyj.GDRS[0].QSDGDCGHJ)
			if sdgdzb > 70 {
				sockInfoShow.QSDGDCGHJ = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, sdgdzb, 0x1B)
			} else if sdgdzb > 40 {
				sockInfoShow.QSDGDCGHJ = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, sdgdzb, 0x1B)
			} else {
				sockInfoShow.QSDGDCGHJ = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, sdgdzb, 0x1B)
			}
		} else {
			sockInfoShow.QSDGDCGHJ = fmt.Sprintf("%c[;;30m%c[0m", 0x1B, 0x1B)
		}

		//前十大流通股东持股合计
		if len(single.Gbyj.GDRS) > 0 && "--" != single.Gbyj.GDRS[0].QSDLTGDCGHJ {
			var sdltgdzb = stock.ToFloat(single.Gbyj.GDRS[0].QSDLTGDCGHJ)
			if sdltgdzb > 45 {
				sockInfoShow.QSDLTGDCGHJ = fmt.Sprintf("%c[;;35m%.2f%c[0m", 0x1B, sdltgdzb, 0x1B)
			} else if sdltgdzb > 30 {
				sockInfoShow.QSDLTGDCGHJ = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, sdltgdzb, 0x1B)
			} else {
				sockInfoShow.QSDLTGDCGHJ = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, sdltgdzb, 0x1B)
			}
		} else {
			sockInfoShow.QSDLTGDCGHJ = fmt.Sprintf("%c[;;30m%c[0m", 0x1B, 0x1B)
		}

		sockInfoShow.SBZB = ""
		sockInfoShow.JGZB = 0
		for i := 0; i < len(single.Gbyj.ZLCC); i++ {
			if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "社保") && "--" != single.Gbyj.ZLCC[i].ZLTGBL {
				var sbbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
				sockInfoShow.SBZB = fmt.Sprintf("%.2f", sbbl)
			} else if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "合计") {
				var hjbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
				sockInfoShow.JGZB = hjbl
			}
		}

		//机构推荐数
		if single.JGTJ >= 20 {
			sockInfoShow.JGTJ = fmt.Sprintf("%c[;;35m%d%c[0m", 0x1B, single.JGTJ, 0x1B)
		} else if single.JGTJ >= 6 {
			sockInfoShow.JGTJ = fmt.Sprintf("%c[;;36m%d%c[0m", 0x1B, single.JGTJ, 0x1B)
		} else {
			sockInfoShow.JGTJ = fmt.Sprintf("%c[;;30m%c[0m", 0x1B, 0x1B)
		}

		sockInfoShows = append(sockInfoShows, sockInfoShow)
	}
	sort.Stable(sockInfoShows)

	fmt.Println("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
	fmt.Printf("┃%7s┃%8s┃%5s┃%8s", "Code", "行业", "现价", "行业")
	fmt.Printf("┃┃%4s┃%3s┃%3s┃%4s┃%4s┃%3s", "收入增长", "毛利润", "净利率", "净益率", "公积金", "未分配")
	fmt.Printf("┃┃%4s┃%5s┃%5s┃%3s┃%3s┃%4s┃%4s┃%3s", "市净", "PE静", "PE动", "市销率", "市净比", "PE比", "PET比", "市销比")
	fmt.Printf("┃┃%3s┃%3s┃%5s┃%4s┃%2s┃\n", "十股占", "十流占", "社/流", "机/流", "推荐")
	fmt.Println("┃━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┃")
	for i := 0; i < len(sockInfoShows); i++ {
		var stock = sockInfoShows[i]
		fmt.Printf("┃")
		fmt.Printf("%7s", stock.Code)

		fmt.Printf("┃")
		var rLen = len(stock.Name) - ChineseCount(stock.Name)
		var bLen = 10 - rLen
		for bLen > 0 {
			fmt.Printf(" ")
			bLen--
		}
		fmt.Printf("%s", stock.Name)

		fmt.Printf("┃%7s", stock.Price)

		fmt.Printf("┃")
		rLen = len(stock.HYName) - ChineseCount(stock.HYName)
		bLen = 10 - rLen
		for bLen > 0 {
			fmt.Printf(" ")
			bLen--
		}
		fmt.Printf("%s", stock.HYName)

		//成长-有利润
		fmt.Printf("┃┃%19s", stock.YYZSRAVG)
		fmt.Printf("┃%17s", stock.MLL)
		fmt.Printf("┃%17s", stock.JLL)
		fmt.Printf("┃%17s", stock.WEIGHTAVG_ROE)
		fmt.Printf("┃%17s", stock.MGGJJ)
		fmt.Printf("┃%17s", stock.MGWFPLY)

		//估值类
		fmt.Printf("┃┃%17s", stock.SJL)
		fmt.Printf("┃%17s", stock.PEJT)
		fmt.Printf("┃%17s", stock.PEDT)
		fmt.Printf("┃%17s", stock.PS9)
		fmt.Printf("┃%16s", stock.RPB8)
		fmt.Printf("┃%16s", stock.RPE7)
		fmt.Printf("┃%16s", stock.RPE9)
		fmt.Printf("┃%16s", stock.RPS9)

		//主力研究
		fmt.Printf("┃┃%17s", stock.QSDGDCGHJ)
		fmt.Printf("┃%17s", stock.QSDLTGDCGHJ)
		fmt.Printf("┃%6s", stock.SBZB)
		fmt.Printf("┃%5.2f", stock.JGZB)
		fmt.Printf("┃%15s", stock.JGTJ)

		fmt.Println("┃")
	}
	fmt.Println("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛")
	var num = len(sockInfoShows)
	fmt.Println("符合条件的股票有=", num, "个")
}

func ChineseCount(str1 string) (count int) {
	for _, char := range str1 {
		if unicode.Is(unicode.Han, char) {
			count++
		}
	}

	return
}
