package stock

import (
	"encoding/json"
	"fmt"
	"strings"
	"util/file"
	"util/http"
)

type SingleZyzb struct {
	DATE        string //财报日期
	JBMGSY      string //基本每股收益(元)
	KFMGSY      string //扣非每股收益(元)
	XSMGSY      string //稀释每股收益(元)
	MGJZC       string //每股净资产(元)
	MGGJJ       string //每股公积金(元)
	MGWFPLY     string //每股未分配利润(元)
	MGJYXJL     string //每股经营现金流(元)
	YYZSR       string //营业总收入(元)
	MLR         string //毛利润(元)
	GSJLR       string //归属净利润(元)
	KFJLR       string //扣非净利润(元)
	YYZSRTBZZ   string //营业总收入同比增长(%)
	GSJLRTBZZ   string //归属净利润同比增长(%)
	KFJLRTBZZ   string //扣非净利润同比增长(%)
	YYZSRGDHBZZ string //营业总收入滚动环比增长(%)
	GSJLRGDHBZZ string //归属净利润滚动环比增长(%)
	KFJLRGDHBZZ string //扣非净利润滚动环比增长(%)
	JQJZCSYL    string //加权净资产收益率(%)
	TBJZCSYL    string //摊薄净资产收益率(%)
	TBZZCSYL    string //摊薄总资产收益率(%)
	MLL         string //毛利率(%)
	JLL         string //净利率(%)
	SJSL        string //实际税率(%)
	YSKYYSR     string //预收款/营业收入
	XSXJLYYSR   string //销售现金流/营业收入
	JYXJLYYSR   string //经营现金流/营业收入
	ZZCZZY      string //总资产周转率(次)
	YSZKZZTS    string //应收账款周转天数(天)
	CHZZTS      string //存货周转天数(天)
	ZCFZL       string //资产负债率(%)
	LDZCZFZ     string //流动负债/总负债(%)
	LDBL        string //流动比率
	SDBL        string //速动比率
}

//资产负债
type SingleZcfz struct {
	REPORTDATE           string //报表日期
	MONETARYFUND         string //货币资金
	ACCOUNTBILLREC       string //应收票据及应收账款
	BILLREC              string //其中：	应收票据
	ACCOUNTREC           string //			应收账款
	ADVANCEPAY           string //预付款项
	TOTAL_OTHER_RECE     string //其他应收款合计
	INTERESTREC          string //其中：	应收利息
	DIVIDENDREC          string //			应收股利
	OTHERREC             string //			其他应收款
	INVENTORY            string //存货
	NONLASSETONEYEAR     string //一年内到期的非流动资产
	OTHERLASSET          string //其他流动资产
	SUMLASSET            string //流动资产合计
	LTREC                string //长期应收款
	LTEQUITYINV          string //长期股权投资
	ESTATEINVEST         string //投资性房地产
	FIXEDASSET           string //固定资产
	CONSTRUCTIONPROGRESS string //在建工程
	INTANGIBLEASSET      string //无形资产
	DEVELOPEXP           string //开发支出
	GOODWILL             string //商誉
	LTDEFERASSET         string //长期待摊费用
	DEFERINCOMETAXASSET  string //递延所得税资产
	OTHERNONLASSET       string //其他非流动资产
	SUMNONLASSET         string //非流动资产合计
	SUMASSET             string //资产总计
	STBORROW             string //短期借款
	ACCOUNTBILLPAY       string //应付票据及应付账款
	BILLPAY              string //其中：	应付票据
	ACCOUNTPAY           string //			应付账款
	ADVANCERECEIVE       string //预收款项
	SALARYPAY            string //应付职工薪酬
	TAXPAY               string //应交税费
	TOTAL_OTHER_PAYABLE  string //其他应付款合计
	INTERESTPAY          string //其中：	应付利息
	DIVIDENDPAY          string //			应付股利
	OTHERPAY             string //			其他应付款
	NONLLIABONEYEAR      string //一年内到期的非流动负债
	OTHERLLIAB           string //其他流动负债
	SUMLLIAB             string //流动负债合计
	LTBORROW             string //长期借款
	BONDPAY              string //应付债券
	LTACCOUNTPAY         string //长期应付款
	DEFERINCOME          string //递延收益
	DEFERINCOMETAXLIAB   string //递延所得税负债
	OTHERNONLLIAB        string //	其他非流动负债
	SUMNONLLIAB          string //非流动负债合计
	SUMLIAB              string //负债合计
	//所有者权益(或股东权益)
	SHARECAPITAL    string //实收资本（或股本）
	CAPITALRESERVE  string //资本公积
	INVENTORYSHARE  string //库存股
	SURPLUSRESERVE  string //盈余公积
	RETAINEDEARNING string //未分配利润
	SUMPARENTEQUITY string //归属于母公司股东权益合计
	MINORITYEQUITY  string //少数股东权益
	SUMSHEQUITY     string //股东权益合计
	SUMLIABSHEQUITY string //负债和股东权益合计
}

//利润表
type SingleLrb struct {
	SECURITYCODE      string //编码
	SECURITYSHORTNAME string //名称
	REPORTDATE        string //报告期
	TOTALOPERATEREVE  string //营业总收入
	TOTALOPERATEEXP   string //营业总成本
	OPERATEEXP        string //    营业成本
	RDEXP             string //    研发费用
	OPERATETAX        string //    营业税金及附加
	SALEEXP           string //    销售费用
	MANAGEEXP         string //    管理费用
	FINANCEEXP        string //    财务费用
	INVESTINCOME      string //其它经营收益：投资收益
	OPERATEPROFIT     string //营业利润
	NONOPERATEREVE    string //    加:营业外收入
	NONOPERATEEXP     string //    减:营业外支出
	SUMPROFIT         string //利润总额
	INCOMETAX         string //    减:所得税费用
	NETPROFIT         string //净利润
	PARENTNETPROFIT   string //    其中:归属于母公司股东的净利润
	BASICEPS          string //每股收益-基本EPS
	DILUTEDEPS        string //每股收益-稀释EPS
	SUMCINCOME        string //综合收益总额
	PARENTCINCOME     string //    归属于母公司所有者的综合收益总额
}

//个股主要-指标数据
func DownloadSingleZyzb(code string) (*string, error) {
	var urlFormat = `http://f10.eastmoney.com/NewFinanceAnalysis/MainTargetAjax?type=0&code=%s`

	var url = fmt.Sprintf(urlFormat,
		getSCByCode(code))
	fmt.Println(url)

	//获取数据
	content, err := http.HttpGet(url)
	if nil != err {
		fmt.Println("http get failed url=", url, " error=", err)
		return nil, err
	}

	return content, nil
}

//个股--资产负债表
func DownloadSingleZcfz(code string) (*string, error) {
	var urlFormat = `http://f10.eastmoney.com/NewFinanceAnalysis/zcfzbAjax?companyType=4&reportDateType=0&reportType=1&endDate=&code=%s`

	var url = fmt.Sprintf(urlFormat,
		getSCByCode(code))
	fmt.Println(url)

	//获取数据
	content, err := http.HttpGet(url)
	if nil != err {
		fmt.Println("http get failed url=", url, " error=", err)
		return nil, err
	}

	return content, nil
}

//个股--利润
func DownloadSingleLrb(code string) (*string, error) {
	var urlFormat = `http://f10.eastmoney.com/NewFinanceAnalysis/lrbAjax?companyType=4&reportDateType=0&reportType=1&endDate=&code=%s`

	var url = fmt.Sprintf(urlFormat,
		getSCByCode(code))
	fmt.Println(url)

	//获取数据
	content, err := http.HttpGet(url)
	if nil != err {
		fmt.Println("http get failed url=", url, " error=", err)
		return nil, err
	}

	return content, nil
}

//下载报表相关数据
func DownloadReportData(code string) {
	var (
		zyzb  *string //主要指标
		zyzbs string
		zcfz  *string //资产负债
		zcfzs string
		lrb   *string //利润表
		lrbs  string
		err   error
	)

	var fileName = fmt.Sprintf(report_singleformat, code, reportDate)
	if file.FileExist(fileName) {
		return
	}

	fmt.Println("下载单个：主要指标和资产负债表 for ", code, " ", reportDate)

	//主要指标
	{
		zyzb, err = DownloadSingleZyzb(code)
		if nil != err {
			fmt.Println("Error get 主要指标", err)
			return
		}
		zyzbs = *zyzb
		zyzbs = strings.ToUpper(zyzbs)
		//fmt.Println(zyzbs)
	}

	//资产负债表
	{
		zcfz, err = DownloadSingleZcfz(code)
		if nil != err {
			fmt.Println("Error get 资产负债", err)
			return
		}
		zcfzs = *zcfz
		if strings.HasPrefix(zcfzs, `"`) {
			zcfzs = zcfzs[1 : len(zcfzs)-1]
		}
		//fmt.Println(zcfzs)
		zcfzs = strings.ReplaceAll(zcfzs, `\"`, `"`)
		zcfzs = strings.ToUpper(zcfzs)
		//fmt.Println(zcfzs)
	}
	//利润表
	{
		lrb, err = DownloadSingleLrb(code)
		if nil != err {
			fmt.Println("Error get 利润表", err)
			return
		}
		lrbs = *lrb
		if strings.HasPrefix(lrbs, `"`) {
			lrbs = lrbs[1 : len(lrbs)-1]
		}
		//fmt.Println(lrbs)
		lrbs = strings.ReplaceAll(lrbs, `\"`, `"`)
		lrbs = strings.ToUpper(lrbs)
		//fmt.Println(lrbs)
	}
	//现金流量表-略

	//保存到文件里
	file.WriteFile(fileName, `{
`)
	file.AppendFile(fileName, `"ZYZB":`)
	file.AppendFile(fileName, zyzbs)
	file.AppendFile(fileName, `,
"ZCFZ":`)
	file.AppendFile(fileName, zcfzs)
	file.AppendFile(fileName, `,
"LRB":`)
	file.AppendFile(fileName, lrbs)
	file.AppendFile(fileName, `
}`)
}

func ParseReportData(code string) (*SingleStock, error) {
	//下载报表相差的
	DownloadReportData(code)

	var fileName = fmt.Sprintf(report_singleformat, code, reportDate)
	//读取数据
	data, err := file.ReadFile_v1(fileName)
	if nil != err {
		fmt.Println("error read file ", err)
		return nil, err
	}

	//fmt.Println(data)
	var sigleStock SingleStock
	err = json.Unmarshal([]byte(data), &sigleStock)
	if nil != err {
		fmt.Println(" json unmarshal failed!!!!", err)
		return nil, err
	}

	//for i := 0; i < len(sigleStock.ZYZB); i++ {
	//	fmt.Println(
	//		"主要指标:: 财报日期=", sigleStock.ZYZB[i].DATE,
	//		" 基本每股收益(元)=", sigleStock.ZYZB[i].JBMGSY,
	//	)
	//}

	//for i := 0; i < len(sigleStock.ZCFZ); i++ {
	//	fmt.Println(
	//		"资财负责 报表日期=", sigleStock.ZCFZ[i].REPORTDATE,
	//		" 应付票据及应付账款=", sigleStock.ZCFZ[i].ACCOUNTBILLPAY,
	//	)
	//}

	//for i := 0; i < len(sigleStock.LRB); i++ {
	//	fmt.Println(
	//		"资财负责 报表日期=", sigleStock.LRB[i].REPORTDATE,
	//		" 营业总收入=", sigleStock.LRB[i].TOTALOPERATEREVE,
	//	)
	//}

	return &sigleStock, nil
}
