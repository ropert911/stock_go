package stock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"util/file"
	"util/http"
)

//个股主要指标
//{
//"date":"2020-06-30",		财报日期
//"jbmgsy":"0.2606",			基本每股收益(元)
//"kfmgsy":"0.2481",			扣非每股收益(元)
//"xsmgsy":"0.2606",			稀释每股收益(元)
//"mgjzc":"8.1662",			每股净资产(元)
//"mggjj":"0.6331",			每股公积金(元)
//"mgwfply":"6.3103",			每股未分配利润(元)
//"mgjyxjl":"-1.1698",		每股经营现金流(元)
//"yyzsr":"171亿",			营业总收入(元)
//"mlr":"49.1亿",				毛利润(元)
//"gsjlr":"21.4亿",			归属净利润(元)
//"kfjlr":"20.4亿",			扣非净利润(元)
//"yyzsrtbzz":"-3.09",		营业总收入同比增长(%)
//"gsjlrtbzz":"-23.91",		归属净利润同比增长(%)
//"kfjlrtbzz":"-23.41",		扣非净利润同比增长(%)
//"yyzsrgdhbzz":"0.94",		营业总收入滚动环比增长(%)
//"gsjlrgdhbzz":"-2.61",		归属净利润滚动环比增长(%)
//"kfjlrgdhbzz":"-3.07",		扣非净利润滚动环比增长(%)
//"jqjzcsyl":"3.19",			加权净资产收益率(%)
//"tbjzcsyl":"3.19",			摊薄净资产收益率(%)
//"tbzzcsyl":"0.52",			摊薄总资产收益率(%)
//"mll":"54.88",				毛利率(%)
//"jll":"12.02",				净利率(%)
//"sjsl":"29.58",				实际税率(%)
//"yskyysr":"0.03",			预收款/营业收入
//"xsxjlyysr":"2.18",			销售现金流/营业收入
//"jyxjlyysr":"-0.56",		经营现金流/营业收入
//"zzczzy":"0.04",			总资产周转率(次)
//"yszkzzts":"3.29",			应收账款周转天数(天)
//"chzzts":"4743.23",			存货周转天数(天)
//"zcfzl":"77.16",			资产负债率(%)
//"ldzczfz":"65.83",			流动负债/总负债(%)
//"ldbl":"1.58",				流动比率
//"sdbl":"0.53"				速动比率
//}

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

//{
//\"SECURITYCODE\":\"000069.SZ\",
//\"REPORTTYPE\":\"1\",
//\"REPORTDATETYPE\":\"0\",
//\"TYPE\":\"1\",
//\"REPORTDATE\":\"2020/6/30 0:00:00\",	报表日期
//\"CURRENCY\":\"人民币\",				币种
//\"MONETARYFUND\":\"38468278102.22\",			货币资金
//\"ACCOUNTBILLREC\":\"389949086.36\",			应收票据及应收账款
//\"BILLREC\":\"50000000\",						其中：	应收票据
//\"ACCOUNTREC\":\"339949086.36\",						应收账款
//\"ADVANCEPAY\":\"19908186412.13\",				预付款项
//\"TOTAL_OTHER_RECE\":\"40370095586.77\",		其他应收款合计
//\"INTERESTREC\":\"1062151061.78\",				其中：	应收利息
//\"DIVIDENDREC\":\"138429820.9\",						应收股利
//\"OTHERREC\":\"39169514704.09\",						其他应收款
//\"INVENTORY\":\"221475859676.48\",				存货
//\"NONLASSETONEYEAR\":\"117224881.44\",			一年内到期的非流动资产
//\"OTHERLASSET\":\"12628510732.13\",				其他流动资产
//\"SUMLASSET\":\"333358104477.53\",			流动资产合计
//\"LTREC\":\"539683299.11\",						长期应收款
//\"LTEQUITYINV\":\"14950800543.6\",				长期股权投资
//\"ESTATEINVEST\":\"13119635581.81\",			投资性房地产
//\"FIXEDASSET\":\"15760105074.69\",				固定资产
//\"CONSTRUCTIONPROGRESS\":\"5866465162.53\",		在建工程
//\"INTANGIBLEASSET\":\"9463881256.9\",			无形资产
//\"DEVELOPEXP\":\"\",							开发支出
//\"GOODWILL\":\"79663843.85\",					商誉
//\"LTDEFERASSET\":\"868177185.17\",				长期待摊费用
//\"DEFERINCOMETAXASSET\":\"11922996425.31\",		递延所得税资产
//\"OTHERNONLASSET\":\"5506997755.65\",			其他非流动资产
//\"SUMNONLASSET\":\"80832768658.91\",		非流动资产合计
//\"SUMASSET\":\"414190873136.44\",			资产总计
//\"STBORROW\":\"28108132711.83\",				短期借款
//\"ACCOUNTBILLPAY\":\"14343426702.3\",			应付票据及应付账款
//\"BILLPAY\":\"309627804.14\",					其中：	应付票据
//\"ACCOUNTPAY\":\"14033798898.16\",						应付账款
//\"ADVANCERECEIVE\":\"451321588.44\",			预收款项
//\"SALARYPAY\":\"1188731816.69\",				应付职工薪酬
//\"TAXPAY\":\"3296489358.06\",					应交税费
//\"TOTAL_OTHER_PAYABLE\":\"67445095717.29\",		其他应付款合计
//\"INTERESTPAY\":\"1807138770.41\",				其中：	应付利息
//\"DIVIDENDPAY\":\"1180692466.11\",						应付股利
//\"OTHERPAY\":\"64457264480.77\",						其他应付款
//\"NONLLIABONEYEAR\":\"7909419858.29\",			一年内到期的非流动负债
//\"OTHERLLIAB\":\"6075718976.32\",				其他流动负债
//\"SUMLLIAB\":\"210371302024.27\",			流动负债合计
//\"LTBORROW\":\"94732569136.1\",					长期借款
//\"BONDPAY\":\"12977695919.9\",					应付债券
//\"LTACCOUNTPAY\":\"117506212.05\",				长期应付款
//\"DEFERINCOME\":\"1230533597.07\",				递延收益
//\"DEFERINCOMETAXLIAB\":\"94129943.11\",			递延所得税负债
//\"OTHERNONLLIAB\":\"\",							其他非流动负债
//\"SUMNONLLIAB\":\"109198628072.56\",		非流动负债合计
//\"SUMLIAB\":\"319569930096.83\",			负债合计
//											所有者权益(或股东权益)
//\"SHARECAPITAL\":\"8201793915\",				实收资本（或股本）
//\"CAPITALRESERVE\":\"5192295032.44\",			资本公积
//\"INVENTORYSHARE\":\"900127746.7\",				库存股
//\"SURPLUSRESERVE\":\"4292981786.04\",			盈余公积
//\"RETAINEDEARNING\":\"51755589182.65\",			未分配利润
//\"SUMPARENTEQUITY\":\"66977101992.52\",		归属于母公司股东权益合计
//\"MINORITYEQUITY\":\"27643841047.09\",			少数股东权益
//\"SUMSHEQUITY\":\"94620943039.61\",			股东权益合计
//\"SUMLIABSHEQUITY\":\"414190873136.44\",	负债和股东权益合计
//\"SETTLEMENTPROVISION\":\"\",
//\"LENDFUND\":\"\",
//\"FVALUEFASSET\":\"\",
//\"TRADEFASSET\":\"\",
//\"DEFINEFVALUEFASSET\":\"\",
//\"PREMIUMREC\":\"\",
//\"RIREC\":\"\",
//\"RICONTACTRESERVEREC\":\"\",
//\"EXPORTREBATEREC\":\"\",
//\"SUBSIDYREC\":\"\",
//\"INTERNALREC\":\"\",
//\"BUYSELLBACKFASSET\":\"\",
//\"CLHELDSALEASS\":\"\",
//\"LOANADVANCES\":\"\",
//\"SALEABLEFASSET\":\"\",
//\"HELDMATURITYINV\":\"\",
//\"CONSTRUCTIONMATERIAL\":\"\",
//\"LIQUIDATEFIXEDASSET\":\"\",
//\"PRODUCTBIOLOGYASSET\":\"\",
//\"OILGASASSET\":\"\",
//\"BORROWFROMCBANK\":\"\",
//\"DEPOSIT\":\"\",
//\"BORROWFUND\":\"\",
//\"FVALUEFLIAB\":\"\",
//\"TRADEFLIAB\":\"\",
//\"DEFINEFVALUEFLIAB\":\"\",
//\"SELLBUYBACKFASSET\":\"\",
//\"COMMPAY\":\"\",
//\"RIPAY\":\"\",
//\"INTERNALPAY\":\"\",
//\"ANTICIPATELLIAB\":\"\",
//\"CONTACTRESERVE\":\"\",
//\"AGENTTRADESECURITY\":\"\",
//\"AGENTUWSECURITY\":\"\",
//\"DEFERINCOMEONEYEAR\":\"\",
//\"STBONDREC\":\"\",
//\"CLHELDSALELIAB\":\"\",
//\"PREFERSTOCBOND\":\"\",
//\"SUSTAINBOND\":\"\",
//\"LTSALARYPAY\":\"\",
//\"SPECIALPAY\":\"\",
//\"ANTICIPATELIAB\":\"\",
//\"OTHEREQUITY\":\"\",
//\"PREFERREDSTOCK\":\"\",
//\"SUSTAINABLEDEBT\":\"\",
//\"OTHEREQUITYOTHER\":\"\",
//\"SPECIALRESERVE\":\"\",
//\"GENERALRISKPREPARE\":\"\",
//\"UNCONFIRMINVLOSS\":\"\",
//\"PLANCASHDIVI\":\"\",
//\"DIFFCONVERSIONFC\":\"\",
//\"MARGINOUTFUND\":\"\",
//\"DERIVEFASSET\":\"\",
//\"AMORCOSTFASSET\":\"\",
//\"FVALUECOMPFASSET\":\"\",
//\"CONTRACTASSET\":\"\",
//\"HELDSALEASS\":\"\",
//\"LASSETOTHER\":\"\",
//\"LASSETBALANCE\":\"\",
//\"CREDINV\":\"\",
//\"AMORCOSTFASSETFLD\":\"\",
//\"OTHCREDINV\":\"\",
//\"FVALUECOMPFASSETFLD\":\"\",
//\"OTHEREQUITYINV\":\"\",
//\"OTHERNONFASSET\":\"\",
//\"NONLASSETOTHER\":\"\",
//\"NONLASSETBALANCE\":\"\",
//\"ASSETOTHER\":\"\",
//\"ASSETBALANCE\":\"\",
//\"DERIVEFLIAB\":\"\",
//\"CONTRACTLIAB\":\"\",
//\"AMORCOSTFLIAB\":\"\",
//\"HELDSALELIAB\":\"\",
//\"LLIABOTHER\":\"\",
//\"LLIABBALANCE\":\"\",
//\"AMORCOSTFLIABFLD\":\"\",
//\"NONLLIABOTHER\":\"\",
//\"NONLLIABBALANCE\":\"\",
//\"LIABOTHER\":\"\",
//\"LIABBALANCE\":\"\",
//\"OTHERCINCOME\":\"\",
//\"PARENTEQUITYOTHER\":\"\",
//\"PARENTEQUITYBALANCE\":\"\",
//\"SHEQUITYOTHER\":\"\",
//\"SHEQUITYBALANCE\":\"\",
//\"LIABSHEQUITYOTHER\":\"\",
//\"LIABSHEQUITYBALANCE\":\"\",
//\"TRADE_FINASSET_NOTFVTPL\":\"\",
//\"TRADE_FINLIAB_NOTFVTPL\":\"\",
//\"MONETARYFUND_YOY\":\"9.5060818670156\",
//\"SETTLEMENTPROVISION_YOY\":\"\",
//\"LENDFUND_YOY\":\"\",
//\"FVALUEFASSET_YOY\":\"\",
//\"TRADEFASSET_YOY\":\"\",
//\"DEFINEFVALUEFASSET_YOY\":\"\",
//\"BILLREC_YOY\":\"\",
//\"ACCOUNTREC_YOY\":\"-7.00573612261548\",
//\"ADVANCEPAY_YOY\":\"83.4927955016016\",
//\"PREMIUMREC_YOY\":\"\",
//\"RIREC_YOY\":\"\",
//\"RICONTACTRESERVEREC_YOY\":\"\"
//\"INTERESTREC_YOY\":\"74.7069464674184\",
//\"DIVIDENDREC_YOY\":\"33.1916082281107\",
//\"OTHERREC_YOY\":\"14.0280376477332\",
//\"EXPORTREBATEREC_YOY\":\"\",
//\"SUBSIDYREC_YOY\":\"\",
//\"INTERNALREC_YOY\":\"\",
//\"BUYSELLBACKFASSET_YOY\":\"\",
//\"INVENTORY_YOY\":\"23.374374284057\",
//\"CLHELDSALEASS_YOY\":\"\",
//\"NONLASSETONEYEAR_YOY\":\"\",
//\"OTHERLASSET_YOY\":\"150.453119721499\",
//\"SUMLASSET_YOY\":\"25.3394414257684\",
//\"LOANADVANCES_YOY\":\"\",
//\"SALEABLEFASSET_YOY\":\"\",
//\"HELDMATURITYINV_YOY\":\"\",
//\"LTREC_YOY\":\"-60.8843474845114\",
//\"LTEQUITYINV_YOY\":\"-4.46544628147395\",
//\"ESTATEINVEST_YOY\":\"100.887761763689\",
//\"FIXEDASSET_YOY\":\"9.95853519444658\",
//\"CONSTRUCTIONPROGRESS_YOY\":\"45.8036869591123\",
//\"CONSTRUCTIONMATERIAL_YOY\":\"\",
//\"LIQUIDATEFIXEDASSET_YOY\":\"\",
//\"PRODUCTBIOLOGYASSET_YOY\":\"\",
//\"OILGASASSET_YOY\":\"\",
//\"INTANGIBLEASSET_YOY\":\"26.4508700874987\",
//\"DEVELOPEXP_YOY\":\"\",
//\"GOODWILL_YOY\":\"-10.8261437636469\",
//\"LTDEFERASSET_YOY\":\"19.8337675381008\",
//\"DEFERINCOMETAXASSET_YOY\":\"25.3107636201523\",
//\"OTHERNONLASSET_YOY\":\"-4.69381119521709\",
//\"SUMNONLASSET_YOY\":\"17.7662525535727\",
//\"SUMASSET_YOY\":\"23.7859240148403\",
//\"STBORROW_YOY\":\"0.330331733172637\",
//\"BORROWFROMCBANK_YOY\":\"\",
//\"DEPOSIT_YOY\":\"\",
//\"BORROWFUND_YOY\":\"\",
//\"FVALUEFLIAB_YOY\":\"\",
//\"TRADEFLIAB_YOY\":\"\",
//\"DEFINEFVALUEFLIAB_YOY\":\"\",
//\"BILLPAY_YOY\":\"172.471862849395\",
//\"ACCOUNTPAY_YOY\":\"26.2636484495597\",
//\"ADVANCERECEIVE_YOY\":\"-99.2255405887603\",
//\"SELLBUYBACKFASSET_YOY\":\"\",
//\"COMMPAY_YOY\":\"\",
//\"SALARYPAY_YOY\":\"14.4645177615642\",
//\"TAXPAY_YOY\":\"30.7869203618484\",
//\"INTERESTPAY_YOY\":\"28.8594222073834\",
//\"DIVIDENDPAY_YOY\":\"3847.54565130802\",
//\"RIPAY_YOY\":\"\",
//\"INTERNALPAY_YOY\":\"\",
//\"OTHERPAY_YOY\":\"11.0492608415428\",
//\"ANTICIPATELLIAB_YOY\":\"\",
//\"CONTACTRESERVE_YOY\":\"\",
//\"AGENTTRADESECURITY_YOY\":\"\",
//\"AGENTUWSECURITY_YOY\":\"\",
//\"DEFERINCOMEONEYEAR_YOY\":\"\",
//\"STBONDREC_YOY\":\"\",
//\"CLHELDSALELIAB_YOY\":\"\",
//\"NONLLIABONEYEAR_YOY\":\"921.624884821751\",
//\"OTHERLLIAB_YOY\":\"\",
//\"SUMLLIAB_YOY\":\"26.6762204632065\",
//\"LTBORROW_YOY\":\"25.3757827659394\",
//\"BONDPAY_YOY\":\"-4.43431609937568\",
//\"PREFERSTOCBOND_YOY\":\"\",
//\"SUSTAINBOND_YOY\":\"\",
//\"LTACCOUNTPAY_YOY\":\"48.7659361640255\",
//\"LTSALARYPAY_YOY\":\"\",
//\"SPECIALPAY_YOY\":\"\",
//\"ANTICIPATELIAB_YOY\":\"\",
//\"DEFERINCOME_YOY\":\"-7.07911579578517\",
//\"DEFERINCOMETAXLIAB_YOY\":\"153.221520633371\",
//\"OTHERNONLLIAB_YOY\":\"\",
//\"SUMNONLLIAB_YOY\":\"20.5559417600167\",
//\"SUMLIAB_YOY\":\"24.5161909721855\",
//\"SHARECAPITAL_YOY\":\"-0.00868636931142223\",
//\"OTHEREQUITY_YOY\":\"\",
//\"PREFERREDSTOCK_YOY\":\"\",
//\"SUSTAINABLEDEBT_YOY\":\"\",
//\"OTHEREQUITYOTHER_YOY\":\"\",
//\"CAPITALRESERVE_YOY\":\"-5.78712476389578\",
//\"INVENTORYSHARE_YOY\":\"\",\"SPECIALRESERVE_YOY\":\"\",
//\"SURPLUSRESERVE_YOY\":\"9.01565480234912\",
//\"GENERALRISKPREPARE_YOY\":\"\",
//\"UNCONFIRMINVLOSS_YOY\":\"\",
//\"RETAINEDEARNING_YOY\":\"20.2744612142823\",
//\"PLANCASHDIVI_YOY\":\"\",
//\"DIFFCONVERSIONFC_YOY\":\"\",
//\"SUMPARENTEQUITY_YOY\":\"11.881971248539\",
//\"MINORITYEQUITY_YOY\":\"52.8196430687319\",
//\"SUMSHEQUITY_YOY\":\"21.3816308199223\",
//\"SUMLIABSHEQUITY_YOY\":\"23.7859240148403\",
//\"MARGINOUTFUND_YOY\":\"\",
//\"DERIVEFASSET_YOY\":\"\",
//\"ACCOUNTBILLREC_YOY\":\"6.67193909532955\",
//\"AMORCOSTFASSET_YOY\":\"\",
//\"FVALUECOMPFASSET_YOY\":\"\",
//\"CONTRACTASSET_YOY\":\"\",
//\"HELDSALEASS_YOY\":\"\",
//\"LASSETOTHER_YOY\":\"\",
//\"LASSETBALANCE_YOY\":\"\",
//\"CREDINV_YOY\":\"\",
//\"AMORCOSTFASSETFLD_YOY\":\"\",
//\"OTHCREDINV_YOY\":\"\",
//\"FVALUECOMPFASSETFLD_YOY\":\"\",
//\"OTHEREQUITYINV_YOY\":\"\",
//\"OTHERNONFASSET_YOY\":\"\",
//\"NONLASSETOTHER_YOY\":\"\",
//\"NONLASSETBALANCE_YOY\":\"\",
//\"ASSETOTHER_YOY\":\"\",
//\"ASSETBALANCE_YOY\":\"\",
//\"DERIVEFLIAB_YOY\":\"\",
//\"ACCOUNTBILLPAY_YOY\":\"27.743354546273\",
//\"CONTRACTLIAB_YOY\":\"\",
//\"AMORCOSTFLIAB_YOY\":\"\",
//\"HELDSALELIAB_YOY\":\"\",
//\"LLIABOTHER_YOY\":\"\",
//\"LLIABBALANCE_YOY\":\"\",
//\"AMORCOSTFLIABFLD_YOY\":\"\",
//\"NONLLIABOTHER_YOY\":\"\",
//\"NONLLIABBALANCE_YOY\":\"\",
//\"LIABOTHER_YOY\":\"\",
//\"LIABBALANCE_YOY\":\"\",
//\"OTHERCINCOME_YOY\":\"\",
//\"PARENTEQUITYOTHER_YOY\":\"\",
//\"PARENTEQUITYBALANCE_YOY\":\"\",
//\"SHEQUITYOTHER_YOY\":\"\",
//\"SHEQUITYBALANCE_YOY\":\"\",
//\"LIABSHEQUITYOTHER_YOY\":\"\",
//\"LIABSHEQUITYBALANCE_YOY\":\"\",
//\"TOTAL_OTHER_RECE_YOY\":\"15.1369711221371\",
//\"TOTAL_OTHER_PAYABLE_YOY\":\"13.3985208021313\",
//\"TRADE_FINASSET_NOTFVTPL_YOY\":\"\",
//\"AUDITOPINIONSDOMESTIC\":\"\",
//\"AUDITOPINIONSDOMESTICJW\":\"\",
//\"TRADE_FINLIAB_NOTFVTPL_YOY\":\"\"
//}

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

type SingleStock struct {
	ZYZB []SingleZyzb //主要指标
	ZCFZ []SingleZcfz //资产负债
	JGTJ int          //机构推荐数
	LAST float32      //去年涨幅
	THIS float32      //今年涨幅
	TWO  float32      //两年合计涨幅
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

//下载报表相关数据
func DownloadReportData(code string) {
	var (
		zyzb  *string
		zyzbs string
		zcfz  *string
		zcfzs string
		err   error
	)

	var fileName = fmt.Sprintf(report_singleformat, code, reportDate)
	if file.FileExist(fileName) {
		return
	}

	fmt.Println("download single for ", code, " ", reportDate)

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
	//利润表-略
	//现金流量表-略

	//保存到文件里
	file.WriteFile(fileName, `{
`)
	file.AppendFile(fileName, `"ZYZB":`)
	file.AppendFile(fileName, zyzbs)
	file.AppendFile(fileName, `,
"ZCFZ":`)
	file.AppendFile(fileName, zcfzs)
	file.AppendFile(fileName, `
}`)
}

//下载机构调研数据
func DowloadJgdy(code string) {
	t1 := time.Now()
	endDate := t1.Format("2006-01-02")
	t0 := t1.AddDate(0, -6, 0)
	startDate := t0.Format("2006-01-02")
	var month = endDate[0:7]
	var fileName = fmt.Sprintf(jgdy_singleformate, code, month)
	if file.FileExist(fileName) {
		return
	}

	var urlFormat = `http://reportapi.eastmoney.com/report/list?cb=datatable%d&pageNo=1&pageSize=500&code=%s&industryCode=*&industry=*&rating=*&ratingchange=*&beginTime=%s&endTime=%s&fields=&qType=0&_=1602117201213`
	var url = fmt.Sprintf(urlFormat,
		rand.Intn(8999999)+1000000,
		code,
		startDate,
		endDate,
	)
	fmt.Println(url)

	//获取数据
	content, err := http.HttpGet(url)
	if nil != err {
		fmt.Println("http get failed url=", url, " error=", err)
		return
	}

	//得到真实内容
	data := *content
	var start = strings.IndexAny(data, "(")
	if -1 == start {
		return
	}
	data = data[start+1 : len(data)-1]
	data = strings.ToUpper(data)
	//fmt.Println(data)

	file.WriteFile(fileName, data)
}

func DownloadSingle(code string) {
	DownloadReportData(code) //下载报表相差的
	DowloadJgdy(code)        //下载机构调研相关数据
	DowloadKx(code)          //下转年K线数据
}

func ParseReportData(code string) (*SingleStock, error) {
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

	for i := 0; i < len(sigleStock.ZYZB); i++ {
		//fmt.Println(
		//	" 财报日期=", sigleStock.ZYZB[i].DATE,
		//	" 基本每股收益(元)=", sigleStock.ZYZB[i].JBMGSY,
		//)
	}

	for i := 0; i < len(sigleStock.ZCFZ); i++ {
		//fmt.Println(
		//	" 报表日期=", sigleStock.ZCFZ[i].REPORTDATE,
		//	" 应付票据及应付账款=", sigleStock.ZCFZ[i].ACCOUNTBILLPAY,
		//)
	}

	return &sigleStock, nil
}

type SingleJkdy struct {
	HITS int
}

func ParseSingleJkdy(code string) int {
	endDate := time.Now().Format("2006-01-02")
	var month = endDate[0:7]
	var fileName = fmt.Sprintf(jgdy_singleformate, code, month)
	data, err := file.ReadFile_v1(fileName)
	if nil != err {
		fmt.Println("error read file ", err)
		return 0
	}

	var singleJkdy SingleJkdy
	data = strings.ReplaceAll(data, `(MISSING)`, "")
	err = json.Unmarshal([]byte(data), &singleJkdy)
	if nil != err {
		fmt.Println(" json unmarshal failed!!!! ", err, " data=", data)
		return 0
	}

	return singleJkdy.HITS
}

func ParseSingle(code string) (*SingleStock, error) {
	sigleStock, err := ParseReportData(code)                         //解析报表数据 -- 报表
	sigleStock.JGTJ = ParseSingleJkdy(code)                          //解析机构推荐数据 -- 日转月用
	sigleStock.LAST, sigleStock.THIS, sigleStock.TWO = ParseKy(code) //K线数据 -- 日转月用
	return sigleStock, err
}
