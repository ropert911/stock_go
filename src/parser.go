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

type SockInfoShow struct {
	//基本信息
	Code   string  //编码
	Name   string  //名称
	Price  string  //股价
	pri    float32 //股价
	HYName string  //行业

	//龙头
	HPMJLR        string //净利润排名
	XSMLL         string //毛利率
	WEIGHTAVG_ROE string //ROE

	//业绩增长
	YYZSRAVG string //营业总收3年平均
	YYZSRZZ  string //营业总收增长
	HPMCZX   string //成长性排名

	//估值
	SJL   string //市净率
	PEJT  string //静态市盈率
	PEDT  string //动态市盈率
	PS9   string //市销率
	RPB8  string //市净率估值
	RPE7  string //PE(静)估值
	RPE9  string //PE(TTM)估值
	RPS9  string //市销率估值
	PEG   string //PEG 估值
	HPMGZ string //估值排名

	//主力
	SBZB string //社保占流通比
	JGZB string //机构合计占流通比
	JGTJ string //机构推荐数
}
type SockInfoShowArray []SockInfoShow

func (s SockInfoShowArray) Len() int           { return len(s) }
func (s SockInfoShowArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SockInfoShowArray) Less(i, j int) bool { return s[i].pri >= s[j].pri }

//显示结果
func exportResult(mapStock map[string]StockInfo2) {
	var sockInfoShows = SockInfoShowArray{}

	//把符合条件股票中的亮点数据显示出来
	for key, value := range mapStock {
		var single, _ = stock.ParseSingle(key, value.gzfx.ORIGINALCODE)

		var sockInfoShow SockInfoShow

		////基本信息
		sockInfoShow.Code = key
		sockInfoShow.Name = value.zcfz.SECURITY_NAME_ABBR
		sockInfoShow.Price = fmt.Sprint(value.gzfx.NEW)
		sockInfoShow.pri = value.gzfx.NEW
		sockInfoShow.HYName = value.gzfx.HYName

		//龙头
		sockInfoShow.HPMJLR = single.THBJ.GSGMJLR[0].PM //净利润排名
		sockInfoShow.HPMJLR = strings.ReplaceAll(sockInfoShow.HPMJLR, "U003E", ">")
		sockInfoShow.XSMLL = fmt.Sprintf("%.2f", value.yjbb.XSMLL) //毛利率
		//ROE > 20
		if value.yjbb.WEIGHTAVG_ROE >= 15 {
			sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("%.2f", value.yjbb.WEIGHTAVG_ROE)
		} else {
			sockInfoShow.WEIGHTAVG_ROE = fmt.Sprintf("")
		}

		//业绩增长
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
				sockInfoShow.YYZSRAVG = fmt.Sprintf("%.2f", avg)
			} else {
				sockInfoShow.YYZSRAVG = fmt.Sprintf("")
			}
		}
		//营业收入增长>50%
		{
			if stock.ToFloat(single.ZYZB[0].YYZSRTBZZ) > 50 {
				sockInfoShow.YYZSRZZ = fmt.Sprintf("%2s", single.ZYZB[0].YYZSRTBZZ)
			} else {
				sockInfoShow.YYZSRZZ = fmt.Sprintf("")
			}
		}
		sockInfoShow.HPMCZX = single.THBJ.CZXBJ.DATA[0].PM //成长性排名
		sockInfoShow.HPMCZX = strings.ReplaceAll(sockInfoShow.HPMCZX, "U003E", ">")

		//估值
		//市净率
		sockInfoShow.SJL = fmt.Sprintf("%.2f", value.gzfx.PB8)
		//静态市盈率
		sockInfoShow.PEJT = fmt.Sprintf("%.2f", value.gzfx.PE7)
		//动态市盈率
		sockInfoShow.PEDT = fmt.Sprintf("%.2f", value.gzfx.PE9)
		sockInfoShow.PS9 = fmt.Sprintf("%.2f", value.gzfx.PS9)
		//市净率/行业市净率<0.8
		var rPB8 = value.gzfx.PB8 / value.gzfx.HY_PB8
		if rPB8 < 0.8 {
			sockInfoShow.RPB8 = fmt.Sprintf("%.2f", rPB8)
		} else {
			sockInfoShow.RPB8 = fmt.Sprintf("")
		}
		//市盈率（静）/行业市盈率（静）<0.8
		var rPE7 = value.gzfx.PE7 / value.gzfx.HY_PE7
		if rPE7 < 0.8 {
			sockInfoShow.RPE7 = fmt.Sprintf("%.2f", rPE7)
		} else {
			sockInfoShow.RPE7 = fmt.Sprintf("")
		}
		//市盈率（动）/行业市盈率（动）<0.8
		var rPE9 = value.gzfx.PE9 / value.gzfx.HY_PE9
		if rPE9 < 0.8 {
			sockInfoShow.RPE9 = fmt.Sprintf("%.2f", rPE9)
		} else {
			sockInfoShow.RPE9 = fmt.Sprintf("")
		}
		//市销率/行业市销率<0.8
		var rPS9 = value.gzfx.PS9 / value.gzfx.HY_PS9
		if rPS9 < 0.8 {
			sockInfoShow.RPS9 = fmt.Sprintf("%.2f", rPS9)
		} else {
			sockInfoShow.RPS9 = fmt.Sprintf("")
		}
		//PEG<0.8
		{
			if !strings.HasPrefix(single.THBJ.GZBJ.DATA[0].PEG, "--") && stock.ToFloat(single.THBJ.GZBJ.DATA[0].PEG) < 0.8 {
				sockInfoShow.PEG = fmt.Sprintf("%2s", single.THBJ.GZBJ.DATA[0].PEG)
			} else {
				sockInfoShow.PEG = fmt.Sprintf("")
			}
		}
		sockInfoShow.HPMGZ = single.THBJ.GZBJ.DATA[0].PM //估值排名
		sockInfoShow.HPMGZ = strings.ReplaceAll(sockInfoShow.HPMGZ, "U003E", ">")

		//主力
		//社保/流通股 >= 3%   机构占流通股比例>40%
		sockInfoShow.SBZB = fmt.Sprintf("")
		//机构合计占流通比
		sockInfoShow.JGZB = fmt.Sprintf("")
		for i := 0; i < len(single.Gbyj.ZLCC); i++ {
			if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "社保") && "--" != single.Gbyj.ZLCC[i].ZLTGBL {
				var sbbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
				if sbbl >= 3 {
					sockInfoShow.SBZB = fmt.Sprintf("%.2f", sbbl)
				}
			} else if strings.Contains(single.Gbyj.ZLCC[i].JGLX, "合计") {
				var hjbl = stock.ToFloat(strings.ReplaceAll(single.Gbyj.ZLCC[i].ZLTGBL, "%", ""))
				if hjbl >= 50 {
					sockInfoShow.JGZB = fmt.Sprintf("%.2f", hjbl)
				}
			}
		}
		//机构推荐数
		if single.JGTJ >= 20 {
			sockInfoShow.JGTJ = fmt.Sprintf("%d", single.JGTJ)
		} else {
			sockInfoShow.JGTJ = fmt.Sprintf("")
		}

		sockInfoShows = append(sockInfoShows, sockInfoShow)
	}
	sort.Stable(sockInfoShows)

	fmt.Printf(" 基本                                                     龙头分析                        成长分析                        估值分析                                                                               主力分析\n")
	fmt.Printf("%4s\t%-8s\t%4s\t%-8s", "编码", "名称", "股价", "行业")
	fmt.Printf("\t%5s\t%3s\t%6s", "净利润排", "毛利率", "ROE")
	fmt.Printf("\t%7s\t%4s\t%5s\t%5s", "avg3增长", "增长", "研发投入", "成长排名")
	fmt.Printf("\t%6s\t%7s\t%7s\t%4s\t%6s\t%5s", "市净率-估", "市盈率静-估", "市盈率动-估", "PEG", "市销率-估", "估值排名")
	fmt.Printf("\t%3s\t%4s\t%4s\t%4s\n", "增持", "社/流", "机/流", "推荐数")
	for i := 0; i < len(sockInfoShows); i++ {
		var stock = sockInfoShows[i]
		//===============基本信息
		fmt.Printf("%6s\t%-s", stock.Code, stock.Name)
		var rLen = len(stock.Name) - ChineseCount2(stock.Name)
		var bLen = 10 - rLen
		for bLen > 0 {
			fmt.Printf(" ")
			bLen--
		}
		fmt.Printf("\t%6s", stock.Price)
		fmt.Printf("\t%-s", stock.HYName)
		rLen = len(stock.HYName) - ChineseCount2(stock.HYName)
		bLen = 10 - rLen
		for bLen > 0 {
			fmt.Printf(" ")
			bLen--
		}

		//==============龙头业绩和增长
		fmt.Printf("\t%8s\t%6s\t%6s",
			stock.HPMJLR,
			stock.XSMLL,
			stock.WEIGHTAVG_ROE)

		//=================业绩增长
		fmt.Printf("\t%8s\t%6s\t%9s\t%8s", stock.YYZSRAVG, stock.YYZSRZZ, "", stock.HPMCZX)

		//=================估值
		fmt.Printf("\t%5s/%4s", stock.SJL, stock.RPB8)
		fmt.Printf("\t%6s/%4s", stock.PEJT, stock.RPE7)
		fmt.Printf("\t%6s/%4s", stock.PEDT, stock.RPE9)
		fmt.Printf("\t%4s", stock.PEG)
		fmt.Printf("\t%5s/%4s", stock.PS9, stock.RPS9)
		fmt.Printf("\t%8s", stock.HPMGZ)

		//=================主力
		//fmt.Printf("┃\t%3s┃\t%4s┃\t%4s┃\t%4s\n", "增持", "社/流", "机/流", "推荐数")
		fmt.Printf("\t%5s", "")
		fmt.Printf("\t%6s", stock.SBZB)
		fmt.Printf("\t%6s", stock.JGZB)
		fmt.Printf("\t%7s", stock.JGTJ)

		fmt.Println("")
	}
	fmt.Println("┃━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┃")
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
