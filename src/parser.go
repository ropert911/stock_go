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
			if value.xjll.CCE_ADD < (0.2 * 10000 * 10000) {
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

		///////////成长
		//收入同比>=10%
		{
			if value.yjbb.YSTZ < 10 {
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

		///////////////其它
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

		//========================================个股=================================
		var single, err = stock.ParseSingle(key, value.gzfx.ORIGINALCODE)
		if nil != err {
			fmt.Println("Error parse single stock data ", err)
			continue
		}

		//三年平均净资产收益率>行业平均和行业中值
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
		//前一年的净资产收益>=行业中值*1.2
		{
			if stock.ToFloat(single.THBJ.DBFXBJ.DATA[0].ROE2) < (stock.ToFloat(single.THBJ.DBFXBJ.DATA[2].ROE2) * 1.2) {
				delete(mapStock, key)
				continue
			}
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
			//近2年平均>10%
			{
				var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
				var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
				var last2avg = math.Cbrt((1+tb0/100)*(1+tb1/100))*100 - 100
				if last2avg < 10 {
					delete(mapStock, key)
					continue
				}
			}
			//近3年平均>8%
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
		}
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

type SockInfoShow struct {
	//基本信息
	Code   string //编码
	Name   string //名称
	Price  string //股价
	HYName string //行业

	//重点关注
	MGGJJ         string //每股公积金
	MGWFPLY       string //每股未分配
	WEIGHTAVG_ROE string //净益率

	//估值
	RPB8   string //市净率估值
	RPE7   string //PE(静)估值
	RPE9   string //PE(TTM)估值
	RPS9   string //市销率估值
	MGLZGZ string //每股内在股价估值
	PEG    string //PEG
	//估值-值
	SJL  string //市净率
	PEJT string //静态市盈率
	PEDT string //动态市盈率
	PS9  string //市销率

	//成长
	YYZSRZZ  string //营业总收增长
	YYZSRAVG string //营业总收3年平均

	//主力研究
	SBZB string //社保占流通比
	JGZB string //机构合计占流通比
	JGTJ string //机构推荐数
}
type SockInfoShowArray []SockInfoShow

func (s SockInfoShowArray) Len() int           { return len(s) }
func (s SockInfoShowArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SockInfoShowArray) Less(i, j int) bool { return s[i].Price >= s[j].Price }

//显示结果
func exportResult(mapStock map[string]StockInfo2) {
	CreateExportEBK2(mapStock)

	var sockInfoShows = SockInfoShowArray{}

	//把符合条件股票中的亮点数据显示出来
	for key, value := range mapStock {
		var single, _ = stock.ParseSingle(key, value.gzfx.ORIGINALCODE)

		var sockInfoShow SockInfoShow
		sockInfoShow.Code = key
		sockInfoShow.Name = value.zcfz.SECURITY_NAME_ABBR
		sockInfoShow.Price = fmt.Sprint(value.gzfx.NEW)
		sockInfoShow.HYName = value.gzfx.HYName

		//重点关注
		//公积金*4>现价
		var mggjj = stock.ToFloat(single.ZYZB[0].MGGJJ)
		if mggjj*4 > float64(value.gzfx.NEW) {
			sockInfoShow.MGGJJ = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, mggjj, 0x1B)
		} else {
			sockInfoShow.MGGJJ = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		}
		//每股未分配利润 *4>现价
		var mgwfply = stock.ToFloat(single.ZYZB[0].MGWFPLY)
		if mgwfply*4 > float64(value.gzfx.NEW) {
			sockInfoShow.MGWFPLY = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, mgwfply, 0x1B)
		} else {
			sockInfoShow.MGWFPLY = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		}
		//ROE > 20
		if value.yjbb.WEIGHTAVG_ROE >= 15 {
			sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, value.yjbb.WEIGHTAVG_ROE, 0x1B)
		} else {
			sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		}

		////////////////////////////////估值
		//市净率/行业市净率<0.8
		var rPB8 = value.gzfx.PB8 / value.gzfx.HY_PB8
		if rPB8 < 0.8 {
			sockInfoShow.RPB8 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPB8, 0x1B)
		} else {
			sockInfoShow.RPB8 = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		}
		//市盈率（静）/行业市盈率（静）<0.8
		var rPE7 = value.gzfx.PE7 / value.gzfx.HY_PE7
		if rPE7 < 0.8 {
			sockInfoShow.RPE7 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPE7, 0x1B)
		} else {
			sockInfoShow.RPE7 = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		}
		//市盈率（动）/行业市盈率（动）<0.8
		var rPE9 = value.gzfx.PE9 / value.gzfx.HY_PE9
		if rPE9 < 0.8 {
			sockInfoShow.RPE9 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPE9, 0x1B)
		} else {
			sockInfoShow.RPE9 = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		}
		//市销率/行业市销率<0.8
		var rPS9 = value.gzfx.PS9 / value.gzfx.HY_PS9
		if rPS9 < 0.8 {
			sockInfoShow.RPS9 = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, rPS9, 0x1B)
		} else {
			sockInfoShow.RPS9 = fmt.Sprintf("%c[;;30m   %c[0m", 0x1B, 0x1B)
		}
		//公司内在估值/现价>1.3
		//E(2R+8.5)*4.4/Y
		// 	E:每股收益
		//	R:预期收益增长率
		//	8.5：平均市盈率，中国应该是22.5，按20来算
		//	4.4：平均利息
		//	Y:公司债/国债收益率（5年期）  约等于3.2%
		//	-->每股收益*(2*预期收益增长率+22.5)*4.4/3.2
		{
			mglgz := float64(value.yjbb.BASIC_EPS) * (2*value.xlyc.EGR + 20.5) * 4.4 / 3.2
			gzbl := mglgz / float64(value.gzfx.NEW)
			if gzbl > 1.30 {
				sockInfoShow.MGLZGZ = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, gzbl, 0x1B)
			} else {
				sockInfoShow.MGLZGZ = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
			}
		}
		//PEG<0.8
		{
			if !strings.HasPrefix(single.THBJ.GZBJ.DATA[0].PEG, "--") && stock.ToFloat(single.THBJ.GZBJ.DATA[0].PEG) < 0.8 {
				sockInfoShow.PEG = fmt.Sprintf("%c[;;36m%2s%c[0m", 0x1B, single.THBJ.GZBJ.DATA[0].PEG, 0x1B)
			} else {
				sockInfoShow.PEG = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
			}
		}

		////////////////估值-值
		//市净率
		sockInfoShow.SJL = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PB8, 0x1B)
		//动态市盈率
		sockInfoShow.PEDT = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PE9, 0x1B)
		//静态市盈率
		sockInfoShow.PEJT = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PE7, 0x1B)
		sockInfoShow.PS9 = fmt.Sprintf("%c[;;30m%.2f%c[0m", 0x1B, value.gzfx.PS9, 0x1B)

		////////////////////////////////成长
		//营业收入增长>50%
		{
			if stock.ToFloat(single.ZYZB[0].YYZSRTBZZ) > 50 {
				sockInfoShow.YYZSRZZ = fmt.Sprintf("%c[;;36m%2s%c[0m", 0x1B, single.ZYZB[0].YYZSRTBZZ, 0x1B)
			} else {
				sockInfoShow.YYZSRZZ = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
			}
		}
		//3年平均>=20%
		{
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
			var tb0 = stock.ToFloat(zyzbP0.YYZSRTBZZ)
			var tb1 = stock.ToFloat(zyzbP1.YYZSRTBZZ)
			var tb2 = stock.ToFloat(zyzbP2.YYZSRTBZZ)
			var avg = math.Cbrt((1+tb0/100)*(1+tb1/100)*(1+tb2/100))*100 - 100
			if avg >= 20 {
				sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, avg, 0x1B)
			} else {
				sockInfoShow.YYZSRAVG = fmt.Sprintf("%c[;;30m%c[0m", 0x1B, 0x1B)
			}
		}

		////////////////主力研究
		//社保/流通股 >= 3%   机构占流通股比例>40%
		sockInfoShow.SBZB = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		sockInfoShow.JGZB = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		for i := 0; i < len(single.Gbyj.ZLCC); i++ {
			if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "社保") && "--" != single.Gbyj.ZLCC[i].ZLTGBL {
				var sbbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
				if sbbl >= 3 {
					sockInfoShow.SBZB = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, sbbl, 0x1B)
				}
			} else if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "合计") {
				var hjbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
				if hjbl >= 40 {
					sockInfoShow.JGZB = fmt.Sprintf("%c[;;36m%.2f%c[0m", 0x1B, hjbl, 0x1B)
				}
			}
		}
		//机构推荐数>=20
		if single.JGTJ >= 20 {
			sockInfoShow.JGTJ = fmt.Sprintf("%c[;;36m%d%c[0m", 0x1B, single.JGTJ, 0x1B)
		} else {
			sockInfoShow.JGTJ = fmt.Sprintf("%c[;;30m  %c[0m", 0x1B, 0x1B)
		}

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

		sockInfoShows = append(sockInfoShows, sockInfoShow)
	}
	sort.Stable(sockInfoShows)

	fmt.Println("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
	fmt.Printf("┃%7s┃%8s┃%5s┃%8s", "Code", "名称", "现价", "行业")
	fmt.Printf("┃┃%4s┃%3s┃%3s", "公积金", "未分配", "净益率")
	fmt.Printf("┃┃%5s┃%4s┃%5s┃%5s┃%5s┃%6s", "市净率估", "PE估", "PET估", "市销率估", "内在股价估", "PEG")
	fmt.Printf("┃┃%5s┃%4s┃%6s┃%5s", "市净率", "PE静", "PE TTM", "市销率")
	fmt.Printf("┃┃%4s┃%3s", "收入增长", "三年avg增长")
	fmt.Printf("┃┃%4s┃%3s┃%3s", "社/流", "机/流", "推荐数")
	fmt.Printf("\n")
	fmt.Println("┃━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┃")
	for i := 0; i < len(sockInfoShows); i++ {
		var stock = sockInfoShows[i]
		fmt.Printf("┃")
		fmt.Printf("%7s", stock.Code)

		fmt.Printf("┃")
		var rLen = len(stock.Name) - ChineseCount2(stock.Name)
		var bLen = 10 - rLen
		for bLen > 0 {
			fmt.Printf(" ")
			bLen--
		}
		fmt.Printf("%s", stock.Name)

		fmt.Printf("┃%7s", stock.Price)

		fmt.Printf("┃")
		rLen = len(stock.HYName) - ChineseCount2(stock.HYName)
		bLen = 10 - rLen
		for bLen > 0 {
			fmt.Printf(" ")
			bLen--
		}
		fmt.Printf("%s", stock.HYName)

		//重点-------
		fmt.Printf("┃┃%18s", stock.MGGJJ)
		fmt.Printf("┃%16s", stock.MGWFPLY)
		fmt.Printf("┃%17s", stock.WEIGHTAVG_ROE)

		//估值
		fmt.Printf("┃┃%19s", stock.RPB8)
		fmt.Printf("┃%16s", stock.RPE7)
		fmt.Printf("┃%17s", stock.RPE9)
		fmt.Printf("┃%19s", stock.RPS9)
		fmt.Printf("┃%21s", stock.MGLZGZ)
		fmt.Printf("┃%17s", stock.PEG)

		//估值-值
		fmt.Printf("┃┃%18s", stock.SJL)
		fmt.Printf("┃%17s", stock.PEJT)
		fmt.Printf("┃%17s", stock.PEDT)
		fmt.Printf("┃%17s", stock.PS9)

		//高成长
		fmt.Printf("┃┃%19s", stock.YYZSRZZ)
		fmt.Printf("┃%21s", stock.YYZSRAVG)

		//主力研究
		fmt.Printf("┃┃%17s", stock.SBZB)
		fmt.Printf("┃%16s", stock.JGZB)
		fmt.Printf("┃%15s", stock.JGTJ)

		fmt.Println("┃")
	}
	fmt.Println("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛")
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
