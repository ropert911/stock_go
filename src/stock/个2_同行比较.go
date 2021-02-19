package stock

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"util/file"
	"util/http"
)

//行业资讯
type THBJHYZX struct {
	DATE  string //日期
	CODE  string //代号
	TITLE string //标题
}

//成长性比较数据
type THBJCZXBJDATA struct {
	PM           string //成长性排名(1,2,>137)	或  空
	DM           string //代码		或  行业
	JC           string //简称   		或  行业均值、行业中值
	JBMGSYZZLFH  string //基本每股收益增长率(%)	3年复合
	JBMGSYZZL    string //基本每股收益增长率(%)	基年增长
	JBMGSYZZLTTM string //基本每股收益增长率(%)	TTM滚动增长
	JBMGSYZZL1   string //基本每股收益增长率(%)	预期-下年
	JBMGSYZZL2   string //基本每股收益增长率(%)	预期-再下年
	JBMGSYZZL3   string //基本每股收益增长率(%)	预期-再三年
	YYSRZZLFH    string //营业收入增长率(%)	3年复合
	YYSRZZL      string //营业收入增长率(%)	基年增长
	YYSRZZLTTM   string //营业收入增长率(%)	TTM滚动增长
	YYSRZZL1     string //营业收入增长率(%)	预期-下年
	YYSRZZL2     string //营业收入增长率(%)	预期-再下年
	YYSRZZL3     string //营业收入增长率(%)	预期-再三年
	JLRZZLFH     string //净利润增长率(%)		3年复合
	JLRZZL       string //净利润增长率(%)		基年增长
	JLRZZLTTM    string //净利润增长率(%)		TTM滚动增长
	JLRZZL1      string //净利润增长率(%)		预期-下年
	JLRZZL2      string //净利润增长率(%)		预期-再下年
	JLRZZL3      string //净利润增长率(%)		预期-再三年
}

//成长性比较
type THBJCZXBJ struct {
	BASEYEAR uint //基年  上一年
	DATA     []THBJCZXBJDATA
}

type THBJGZBJDATA struct {
	PM      string //估值排名(1,2,>137)	或  空
	DM      string //代码
	JC      string //简称
	PEG     string //PEG
	SYL     string //市盈率  - 基年实际
	SYLTTM  string //市盈率  - 基年滚动增长TTM
	SYL1    string //市盈率 - 下年预期
	SYL2    string //市盈率 - 再一年预期
	SYL3    string //市盈率 - 再二年预期
	SSL     string //市销率 - 基年市销率
	SSLTTM  string //市销率 - 基年滚动增长TTM
	SSL1    string //市销率 - 下年预期
	SSL2    string //市销率 - 再一年预期
	SSL3    string //市销率 - 再二年预期
	SJL     string //市净率 - 基年
	SJLMRQ  string //市净率 - MRQ（上一交易日收盘价/最新每股净资产）
	SXL1    string //市现率① - 基年实际		市现率①=总市值/现金及现金等价物净增加额
	SXLTTM1 string //市现率① - 滚动市现率
	SXL2    string //市现率② - 基年实际		市现率②=总市值/经营活动产生的现金流量净额
	SXLTTM2 string //市现率② - 滚动市现率
	EV      string //EV/EBITDA  - 基年实际
	EVTTM   string //EV/EBITDA - TTM
}

//估值比较
type THBJGZBJ struct {
	BASEYEAR uint //基年  上一年
	DATA     []THBJGZBJDATA
}

//杜邦分析比较-数据
type THBJDBFXBJDATA struct {
	PM       string //排名(1,2,>137)	或  空
	DM       string //代码
	JC       string //简称   		或  行业均值、行业中值
	ROEPJ    string //ROE - 三年平均
	ROE      string //ROE(%) 前3年
	ROE1     string //ROE(%) 前2年
	ROE2     string //ROE(%) 前1年
	JLLPJ    string //净利率(%) 3年平均
	JLL      string //净利率(%) 前3年
	JLL1     string //净利率(%) 前2年
	JLL2     string //净利率(%) 前1年
	ZZCZZLPJ string //总资产周转率(%)	三年平均
	ZZCZZL   string //总资产周转率(%)	前3年
	ZZCZZL1  string //总资产周转率(%)	前2年
	ZZCZZL2  string //总资产周转率(%)	前1年
	QYCSPJ   string //权益乘数(%)	三年平均
	QYCS     string //权益乘数(%)	前3年
	QYCS1    string //权益乘数(%)	前2年
	QYCS2    string //权益乘数(%)	前1年

}

//杜邦分析比较
type THBJDBFXBJ struct {
	STARTYEAR uint //开始年限
	DATA      []THBJDBFXBJDATA
}

//公司规模 - 总市值
type THBJGSGMZSZ struct {
	PM   string //排名 或 空
	DM   string //代码 或 空
	JC   string //名称 或 行业平均 行业中值
	ZSZ  string //总市值
	LTSZ string //流通市
	YYSR string //营业收入
	JLR  string //营业收入
	BGQ  string //净利润
}

//公司规模-流动市值
type THBJGSGMLTSZ struct {
	PM   string
	DM   string
	JC   string
	ZSZ  string
	LTSZ string
	YYSR string
	JLR  string
	BGQ  string
}

//公司规模-净利润
type THBJGSGMJLR struct {
	PM   string
	DM   string
	JC   string
	ZSZ  string
	LTSZ string
	YYSR string
	JLR  string
	BGQ  string
}

//同行比较
type SingleTHBJ struct {
	HYZX     []THBJHYZX     //资讯
	CZXBJ    THBJCZXBJ      //成长性比较
	GZBJ     THBJGZBJ       //估值比较
	DBFXBJ   THBJDBFXBJ     //杜邦分析比较
	GSGMZSZ  []THBJGSGMZSZ  //公司规模 - 总市值
	GSGMLTSZ []THBJGSGMLTSZ //公司规模-流动市值
	GSGMJLR  []THBJGSGMJLR  //公司规模-净利润
}

func DowloadTHBJ(code string, icode string) {
	var month = time.Now().Format("2006-01-02")[0:7]
	var fileName = fmt.Sprintf(thbj_singleformate, code, month)
	if file.FileExist(fileName) {
		return
	}

	var urlFormat = `http://f10.eastmoney.com/IndustryAnalysis/IndustryAnalysisAjax?code=%s&icode=%s`
	var url = fmt.Sprintf(urlFormat,
		getSCByCode(code), icode,
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
	data = strings.ToUpper(data)

	file.WriteFile(fileName, data)
}

func ParseTHBJ(code string, icode string) *SingleTHBJ {
	var month = time.Now().Format("2006-01-02")[0:7]
	var fileName = fmt.Sprintf(thbj_singleformate, code, month)
	if !file.FileExist(fileName) {
		DowloadTHBJ(code, icode)
	}

	data, err := file.ReadFile_v1(fileName)
	if nil != err {
		fmt.Println("error read file ", err)
		return nil
	}
	data = strings.ReplaceAll(data, `\U`, "U")
	data = strings.ReplaceAll(data, `NULL`, "null")

	var singleTHBJ SingleTHBJ
	err = json.Unmarshal([]byte(data), &singleTHBJ)
	if nil != err {
		//fmt.Println(" json unmarshal failed!!!! ", err, " data=", data)
		return nil
	}

	//fmt.Printf("资讯：标题 %s\n", singleTHBJ.HYZX[0].TITLE)
	//fmt.Printf("成长性比较：基年增长 %s\n", singleTHBJ.CZXBJ.DATA[0].JBMGSYZZL)
	//fmt.Printf("估值比较：简称 %s 市盈率 %s \n", singleTHBJ.GZBJ.DATA[0].JC, singleTHBJ.GZBJ.DATA[0].SYL)
	//fmt.Printf("杜邦分析：简称-%s 三年平均ROE  %s \n", singleTHBJ.DBFXBJ.DATA[0].JC, singleTHBJ.DBFXBJ.DATA[0].ROEPJ)

	return &singleTHBJ
}
