package stock

//股本研究
//{"gdrs":[	//股东人数
//{
//"rq":"2020-06-30",		//季报
//"gdrs":"35.78万",		//股东人数
//"gdrs_jsqbh":"-4.44",	//较上期变化(%)
//"rjltg":"1.8万",		//人均流通股(股)
//"rjltg_jsqbh":"25.76",	//较上期变化(%)
//"cmjzd":"非常集中",		//筹码集中度
//"gj":"20.20",			//股价(元)
//"rjcgje":"37万",		//人均持股金额(元)
//"qsdgdcghj":"36.32",	//前十大股东持股合计(%)			****
//"qsdltgdcghj":"18.95"	//前十大流通股东持股合计(%)		****
//},
//......
//],
//"sdltgd":[	//十大流动股东
//{
//"rq":"2020-06-30",
//"sdltgd":[
//{"rq":"2020-06-30","mc":"1","gdmc":"其实","gdxz":"个人","gflx":"A股","cgs":"443,023,243","zltgbcgbl":"6.77%","zj":"73837207","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"2","gdmc":"香港中央结算有限公司","gdxz":"其它","gflx":"A股","cgs":"230,682,712","zltgbcgbl":"3.52%","zj":"116156230","bdbl":"101.42%"},
//{"rq":"2020-06-30","mc":"3","gdmc":"沈友根","gdxz":"个人","gflx":"A股","cgs":"216,710,776","zltgbcgbl":"3.31%","zj":"36118463","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"4","gdmc":"陆丽丽","gdxz":"个人","gflx":"A股","cgs":"212,308,838","zltgbcgbl":"3.24%","zj":"35384806","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"5","gdmc":"中央汇金资产管理有限责任公司","gdxz":"其它","gflx":"A股","cgs":"116,707,184","zltgbcgbl":"1.78%","zj":"19451197","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"6","gdmc":"中国建设银行股份有限公司-国泰中证全指证券公司交易型开放式指数证券投资基金","gdxz":"证券投资基金","gflx":"A股","cgs":"91,073,980","zltgbcgbl":"1.39%","zj":"10742663","bdbl":"13.37%"},
//{"rq":"2020-06-30","mc":"7","gdmc":"全国社保基金一一三组合","gdxz":"全国社保基金","gflx":"A股","cgs":"63,699,044","zltgbcgbl":"0.97%","zj":"新进","bdbl":"--"},
//{"rq":"2020-06-30","mc":"8","gdmc":"中国建设银行股份有限公司-华宝中证全指证券公司交易型开放式指数证券投资基金","gdxz":"证券投资基金","gflx":"A股","cgs":"53,925,217","zltgbcgbl":"0.82%","zj":"6978100","bdbl":"14.86%"},
//{"rq":"2020-06-30","mc":"9","gdmc":"全国社保基金一一二组合","gdxz":"全国社保基金","gflx":"A股","cgs":"50,883,903","zltgbcgbl":"0.78%","zj":"新进","bdbl":"--"},
//{"rq":"2020-06-30","mc":"10","gdmc":"基本养老保险基金八零二组合","gdxz":"基本养老基金","gflx":"A股","cgs":"47,825,829","zltgbcgbl":"0.73%","zj":"新进","bdbl":"--"}
//]
//},
//......
//],
//"sdgd":[	//十大股东
//{
//"rq":"2020-06-30",
//"sdgd":[
//{"rq":"2020-06-30","mc":"1","gdmc":"其实","gflx":"流通A股,限售流通A股","cgs":"1,772,092,973","zltgbcgbl":"21.99%","zj":"295348829","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"2","gdmc":"香港中央结算有限公司","gflx":"流通A股","cgs":"230,682,712","zltgbcgbl":"2.86%","zj":"116156230","bdbl":"101.42%"},
//{"rq":"2020-06-30","mc":"3","gdmc":"沈友根","gflx":"流通A股","cgs":"216,710,776","zltgbcgbl":"2.69%","zj":"36118463","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"4","gdmc":"陆丽丽","gflx":"流通A股","cgs":"212,308,838","zltgbcgbl":"2.63%","zj":"35384806","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"5","gdmc":"中央汇金资产管理有限责任公司","gflx":"流通A股","cgs":"116,707,184","zltgbcgbl":"1.45%","zj":"19451197","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"6","gdmc":"中国建设银行股份有限公司-国泰中证全指证券公司交易型开放式指数证券投资基金","gflx":"流通A股","cgs":"91,073,980","zltgbcgbl":"1.13%","zj":"10742663","bdbl":"13.37%"},
//{"rq":"2020-06-30","mc":"7","gdmc":"鲍一青","gflx":"流通A股,限售流通A股","cgs":"90,414,193","zltgbcgbl":"1.12%","zj":"15069032","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"8","gdmc":"史佳","gflx":"流通A股,限售流通A股","cgs":"78,968,724","zltgbcgbl":"0.98%","zj":"13161454","bdbl":"20.00%"},
//{"rq":"2020-06-30","mc":"9","gdmc":"全国社保基金一一三组合","gflx":"流通A股","cgs":"63,699,044","zltgbcgbl":"0.79%","zj":"新进","bdbl":"--"},
//{"rq":"2020-06-30","mc":"10","gdmc":"中国建设银行股份有限公司-华宝中证全指证券公司交易型开放式指数证券投资基金","gflx":"流通A股","cgs":"53,925,217","zltgbcgbl":"0.67%","zj":"6978100","bdbl":"14.86%"}
//]
//},
//......
//],
//"jjcg":[		//基金持股
//{
//"rq":"2020-06-30",
//"jjcg":[
//{
//"id":"",
//"jjdm":"512880",					基金代码
//"jjmc":"国泰中证全指证券公司ETF",	基金名称
//"cgs":"91,073,980.00",				持股数(股)
//"cgsz":"1,839,694,396.00",			持仓市值(元)
//"zzgbb":"1.13%",					占总股本比
//"zltb":"1.13%",						占流通比
//"zjzb":"9.31%",						占净值比
//"order":"91073980"
//},
//.......
//]
//},
//.......
//],
//"zlcc":[		//机构持仓
//{
//"rq":"2020-06-30",
//"jglx":"基金",			机构类型
//"ccjs":"924",			持仓家数
//"ccgs":"1025995259",	持仓股数(股)
//"zltgbl":"12.73%",		占流通股比例
//"zltgbbl":"12.73%"		占总股本比例
//},
//{"rq":"2020-06-30","jglx":"保险","ccjs":"--","ccgs":"--","zltgbl":"--","zltgbbl":"--"},
//{"rq":"2020-06-30","jglx":"券商","ccjs":"--","ccgs":"--","zltgbl":"--","zltgbbl":"--"},
//{"rq":"2020-06-30","jglx":"QFII","ccjs":"--","ccgs":"--","zltgbl":"--","zltgbbl":"--"},
//{"rq":"2020-06-30","jglx":"社保基金","ccjs":"2","ccgs":"114582947","zltgbl":"1.42%","zltgbbl":"1.42%"},
//{"rq":"2020-06-30","jglx":"信托","ccjs":"--","ccgs":"--","zltgbl":"--","zltgbbl":"--"},
//{"rq":"2020-06-30","jglx":"其他机构","ccjs":"3","ccgs":"395215725","zltgbl":"4.90%","zltgbbl":"4.90%"},
//{"rq":"2020-06-30","jglx":"合计","ccjs":"929","ccgs":"1535793931","zltgbl":"19.06%","zltgbbl":"19.06%"}
//],
//"zlcc_rz":["2020-06-30","2020-03-31","2019-12-31","2019-09-30","2019-06-30"],
//"kggx":{
//"sjkzr":"其实",		//实控人
//"cgbl":"21.99%"}	//持股比例
//}

//func DowloadGbyj(code string) {
//	var month = time.Now().Format("2006-01-02")[0:7]
//	var fileName = fmt.Sprintf(gbyj_singleformate, code, month)
//	if file.FileExist(fileName) {
//		return
//	}
//
//	var urlFormat = `http://f10.eastmoney.com/ShareholderResearch/ShareholderResearchAjax?code=%s`
//	var url = fmt.Sprintf(urlFormat,
//		getSCByCode(code),
//	)
//	fmt.Println(url)
//
//	//获取数据
//	content, err := http.HttpGet(url)
//	if nil != err {
//		fmt.Println("http get failed url=", url, " error=", err)
//		return
//	}
//
//	//得到真实内容
//	data := *content
//	data = strings.ToUpper(data)
//
//	file.WriteFile(fileName, data)
//}

//股东人数
//type SingleGDRS struct {
//	RQ          string //季报
//	GDRS        string //股东人数
//	GDRS_JSQBH  string //较上期变化(%)
//	RJLTG       string //人均流通股(股)
//	RJLTG_JSQBH string //较上期变化(%)
//	CMJZD       string //筹码集中度
//	GJ          string //股价(元)
//	RJCGJE      string //人均持股金额(元)
//	QSDGDCGHJ   string //前十大股东持股合计(%)
//	QSDLTGDCGHJ string //前十大流通股东持股合计(%)
//}

//type SingleJJCG2 struct {
//	ID    string
//	JJDM  string //基金代码
//	JJMC  string //基金名称
//	CGS   string //持股数(股)
//	CGSZ  string //持仓市值(元)
//	ZZGBB string //占总股本比
//	ZLTB  string //占流通比
//	ZJZB  string //占净值比
//	ORDER string
//}

//机构持仓
//type SingleJGCC struct {
//	RQ   string //季报
//	JJCG []SingleJJCG2
//}

//基金持股
//type SingleJJCG struct {
//	RQ      string
//	JGLX    string //机构类型
//	CCJS    string //持仓家数
//	CCGS    string //持仓股数(股)
//	ZLTGBL  string //占流通股比例
//	ZLTGBBL string //占总股本比例
//}

//股本研究
//type SingleGbyj struct {
//	GDRS []SingleGDRS
//	JJCG []SingleJGCC
//	ZLCC []SingleJJCG
//}

//func ParseGbyj(code string) *SingleGbyj {
//	var month = time.Now().Format("2006-01-02")[0:7]
//	var fileName = fmt.Sprintf(gbyj_singleformate, code, month)
//	if !file.FileExist(fileName) {
//		DowloadGbyj(code)
//	}
//
//	data, err := file.ReadFile_v1(fileName)
//	if nil != err {
//		fmt.Println("error read file ", err)
//		return nil
//	}
//	//fmt.Println(data)
//	data = strings.ReplaceAll(data, `\U`, "U")
//	data = strings.ReplaceAll(data, `NULL`, "null")
//
//	var singleGbyj SingleGbyj
//	err = json.Unmarshal([]byte(data), &singleGbyj)
//	if nil != err {
//		fmt.Println(" json unmarshal failed!!!! ", err, " data=", data)
//		return nil
//	}
//
//	return &singleGbyj
//}
