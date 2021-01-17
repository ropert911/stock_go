package stock

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"util/file"
	"util/http"
)

//股东人数
type SingleGDRS struct {
	RQ          string //季报
	GDRS        string //股东人数
	GDRS_JSQBH  string //较上期变化(%)
	RJLTG       string //人均流通股(股)
	RJLTG_JSQBH string //较上期变化(%)
	CMJZD       string //筹码集中度
	GJ          string //股价(元)
	RJCGJE      string //人均持股金额(元)
	QSDGDCGHJ   string //前十大股东持股合计(%)
	QSDLTGDCGHJ string //前十大流通股东持股合计(%)
}

type SingleJJCG2 struct {
	ID    string
	JJDM  string //基金代码
	JJMC  string //基金名称
	CGS   string //持股数(股)
	CGSZ  string //持仓市值(元)
	ZZGBB string //占总股本比
	ZLTB  string //占流通比
	ZJZB  string //占净值比
	ORDER string
}

//机构持仓
type SingleJJCC struct {
	RQ   string //季报
	JJCG []SingleJJCG2
}

//基构持股
type SingleJGCG struct {
	RQ      string
	JGLX    string //机构类型
	CCJS    string //持仓家数
	CCGS    string //持仓股数(股)
	ZLTGBL  string //占流通股比例
	ZLTGBBL string //占总股本比例
}

//股本研究
type SingleGbyj struct {
	GDRS []SingleGDRS
	JJCG []SingleJJCC
	ZLCC []SingleJGCG
}

func DowloadGbyj(code string) {
	var month = time.Now().Format("2006-01-02")[0:7]
	var fileName = fmt.Sprintf(gbyj_singleformate, code, month)
	if file.FileExist(fileName) {
		return
	}

	var urlFormat = `http://f10.eastmoney.com/ShareholderResearch/ShareholderResearchAjax?code=%s`
	var url = fmt.Sprintf(urlFormat,
		getSCByCode(code),
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

func ParseGbyj(code string) *SingleGbyj {
	var month = time.Now().Format("2006-01-02")[0:7]
	var fileName = fmt.Sprintf(gbyj_singleformate, code, month)
	if !file.FileExist(fileName) {
		DowloadGbyj(code)
	}

	data, err := file.ReadFile_v1(fileName)
	if nil != err {
		fmt.Println("error read file ", err)
		return nil
	}
	//fmt.Println(data)
	data = strings.ReplaceAll(data, `\U`, "U")
	data = strings.ReplaceAll(data, `NULL`, "null")

	var singleGbyj SingleGbyj
	err = json.Unmarshal([]byte(data), &singleGbyj)
	if nil != err {
		fmt.Println(" json unmarshal failed!!!! ", err, " data=", data)
		return nil
	}

	return &singleGbyj
}
