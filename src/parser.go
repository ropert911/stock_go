package main

import (
	"fmt"
	"stock"
	"strings"
	"util/file"
)

type StockInfo struct {
	gzfx stock.StockGzfx
	yjbb stock.StockYjbb
	zcfz stock.StockZcfz
	lrb  stock.StockLrb
}

func main() {
	mapStock := parserData()
	mapStock = fileData(mapStock)
	exportResult(mapStock)
}

//解析所有数据-估值分析、资产负债、业绩报表、利润表
func parserData() map[string]StockInfo {
	mapStock := make(map[string]StockInfo)
	socksgzfx := stock.ReadStockGzfx()
	socksYjbb := stock.ReadStockYjbb()
	socksZcfz := stock.ReadStockZcfz()
	stockLrb := stock.ReadStockLrb()
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

		mapStock[socksgzfx[i].SECURITYCODE] = stockInfo
	}

	return mapStock
}

//过滤掉不符合条件的
func fileData(mapStock map[string]StockInfo) map[string]StockInfo {
	for key, value := range mapStock {
		//1有积累---非ST
		if strings.HasPrefix(value.gzfx.SName, "*ST") || strings.HasPrefix(value.gzfx.SName, "ST") {
			delete(mapStock, key)
			continue
		}
		//1有积累--- 现价不为空
		if value.gzfx.NEW == 0 {
			delete(mapStock, key)
			continue
		}
		//1有积累 --- 每股净资产>0 && 不为空
		if value.yjbb.BPS <= 0 {
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

		stock.DownloadSingle(key)
		var single, err = stock.ParseSingle(key)
		if nil != err {
			fmt.Println("Error parse single stock data ", err)
			continue
		}
		//1有积累 -- 每股公积金>0.2元
		if stock.ToFloat(single.ZYZB[0].MGGJJ) < 0.5 {
			delete(mapStock, key)
			continue
		}
		//1有积累 -- 未分配利润 >1亿
		if stock.ToFloat(single.ZCFZ[0].RETAINEDEARNING) < 100000000 {
			delete(mapStock, key)
			continue
		}
		//1有积累 -- 每股未分配利润>0.2元
		if stock.ToFloat(single.ZYZB[0].MGWFPLY) < 0.5 {
			delete(mapStock, key)
			continue
		}

		//3成长 -- 前2年>10% 或近 3年平均20%
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
		if tb1 < -10 || tb2 < -10 { //近3年有下跌10%的忽略
			delete(mapStock, key)
			continue
		}
		if !((tb0 >= 15 && tb1 >= 10 && tb2 >= 10) || ((tb0+tb1+tb2)/3 >= 20)) {
			delete(mapStock, key)
			continue
		}
	}

	return mapStock
}

//显示结果
func exportResult(mapStock map[string]StockInfo) {
	var num = len(mapStock)
	//生成导出信息
	if num > 0 {
		var exportName = fmt.Sprintf("%s.EBK", stock.TradeData)
		file.WriteFile(exportName, `
`)

		for key, _ := range mapStock {
			var code = stock.GetExportCodeByCode(key)
			file.AppendFile(exportName, fmt.Sprintf(`%s
`, code))
		}
	}

	//把符合条件股票中的亮点数据显示出来
	for key, value := range mapStock {
		fmt.Print(
			"代码=", key,
			"  名称=", value.gzfx.SName,
			" 最新价=", value.gzfx.NEW,
			" 行业名=", value.gzfx.HYName,
		)

		var single, _ = stock.ParseSingle(key)

		//2积累--每股公积金
		var mggjj = stock.ToFloat(single.ZYZB[0].MGGJJ)
		if mggjj*4 > float64(value.gzfx.NEW) {
			fmt.Printf("%c[;;36m  每股公积金=%f%c[0m ", 0x1B, mggjj, 0x1B)
		} else if mggjj > 8 {
			fmt.Printf("%c[;;35m  每股公积金=%f%c[0m ", 0x1B, mggjj, 0x1B)
		} else {
			fmt.Printf(" 每股公积金=%f", mggjj)
		}
		//1积累--每股未分配利润
		var mgwfply = stock.ToFloat(single.ZYZB[0].MGWFPLY)
		if mgwfply*4 > float64(value.gzfx.NEW) {
			fmt.Printf("%c[;;36m  每股未分配利润=%f%c[0m ", 0x1B, mgwfply, 0x1B)
		} else if mgwfply > 8 {
			fmt.Printf("%c[;;35m  每股未分配利润=%f%c[0m ", 0x1B, mgwfply, 0x1B)
		} else {
			fmt.Printf(" 每股未分配利润=%f", mgwfply)
		}

		//2估值-市净率
		if value.gzfx.PB8 <= 5 {
			fmt.Printf("%c[;;36m  市净率=%f%c[0m ", 0x1B, value.gzfx.PB8, 0x1B)
		} else {
			fmt.Printf(" 市净率=%f", value.gzfx.PB8)
		}
		//2估值-动态市盈率
		if value.gzfx.PE9 <= 25 {
			fmt.Printf("%c[;;36m  动态市盈率=%f%c[0m ", 0x1B, value.gzfx.PE9, 0x1B)
		} else {
			fmt.Printf(" 动态市盈率=%f", value.gzfx.PE9)
		}
		//2估值-静态市盈率
		if value.gzfx.PE7 <= 25 {
			fmt.Printf("%c[;;36m  静态市盈率=%f%c[0m ", 0x1B, value.gzfx.PE7, 0x1B)
		} else {
			fmt.Printf(" 静态市盈率=%f", value.gzfx.PE7)
		}
		//2估值-市销率
		if value.gzfx.PS9 <= 3 {
			fmt.Printf("%c[;;36m  市销率=%f%c[0m ", 0x1B, value.gzfx.PS9, 0x1B)
		} else {
			fmt.Printf(" 市销率=%f", value.gzfx.PS9)
		}

		//3成长-净益率
		if value.yjbb.WEIGHTAVG_ROE >= 10 {
			fmt.Printf("%c[;;36m  净益率=%f%c[0m", 0x1B, value.yjbb.WEIGHTAVG_ROE, 0x1B)
		} else {
			fmt.Printf(" 净益率=%f", value.yjbb.WEIGHTAVG_ROE)
		}
		//3成长-毛利率
		var mll = stock.ToFloat(single.ZYZB[0].MLL)
		if mll >= 30 {
			fmt.Printf("%c[;;36m  毛利率=%f%c[0m", 0x1B, mll, 0x1B)
		} else {
			fmt.Printf(" 毛利率=%f", mll)
		}
		var jll = stock.ToFloat(single.ZYZB[0].JLL)
		if jll >= 20 {
			fmt.Printf("%c[;;36m  净利率=%f%c[0m", 0x1B, jll, 0x1B)
		} else {
			fmt.Printf(" 净利率=%f", jll)
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
		var avg = (tb0 + tb1 + tb2) / 3
		if avg >= 20 {
			fmt.Printf("%c[;;36m  3年收入同比=%f%c[0m", 0x1B, avg, 0x1B)
		} else {
			fmt.Printf(" 3年收入同比=%f", avg)
		}

		fmt.Println("")
	}
	fmt.Println("符合条件的股票有=", num, "个")

	//for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
	//	fmt.Printf("%c[;;%dm  f=%d  %c[0m ", 0x1B, f, f, 0x1B)
	//	fmt.Println("")
	//}
}
