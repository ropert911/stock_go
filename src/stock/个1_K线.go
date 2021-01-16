package stock

//func getUCodeByCode(code string) string {
//	if strings.HasPrefix(code, "6") {
//		return "1." + code
//	} else {
//		return "0." + code
//	}
//}
//func DowloadKx(code string) {
//	t1 := time.Now()
//	endDate := t1.Format("2006-01-02")
//	var month = endDate[0:7]
//	var fileName = fmt.Sprintf(kx_singleformate, code, month)
//	if file.FileExist(fileName) {
//		return
//	}
//
//	var urlFormat = `http://push2his.eastmoney.com/api/qt/stock/kline/get?fields1=f1&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61&beg=0&end=20500101&ut=fa5fd1943c7b386f172d6893dbfba10b&rtntype=6&secid=%s&klt=106&fqt=1&cb=jsonp1602132750107`
//	var url = fmt.Sprintf(urlFormat,
//		//	//rand.Intn(8999999)+1000000,
//		getUCodeByCode(code),
//		//	//startDate,
//		//	//endDate,
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
//	var start = strings.IndexAny(data, "[")
//	var end = strings.IndexAny(data, "]")
//	if -1 == start || -1 == end {
//		return
//	}
//	data = data[start : end+1]
//	//fmt.Println(data)
//
//	file.WriteFile(fileName, data)
//}

//2020-09-30,
//5.96,		开盘价
//5.55,		收盘价
//6.55,		最高
//3.92,		最低
//18010351,	成交量
//9147389925.00,	成交额
//43.69,			振幅
//-7.81,			涨跌幅
//-0.47,			涨跌额
//583.91"			换手率
//返回 去年，今年，近2年涨幅
//func ParseKy(code string) (float32, float32, float32) {
//	t1 := time.Now()
//	endDate := t1.Format("2006-01-02")
//
//	var thisYear = ToInt(endDate[0:4])
//	var lastYear = thisYear - 1
//
//	var month = endDate[0:7]
//	var fileName = fmt.Sprintf(kx_singleformate, code, month)
//	data, err := file.ReadFile_v1(fileName)
//	if nil != err {
//		fmt.Println("error read file ", err)
//		return 0, 0, 0
//	}
//
//	var klint []string
//	err = json.Unmarshal([]byte(data), &klint)
//	if nil != err {
//		fmt.Println("json unmarshal failed!!!!", err)
//		return 0, 0, 0
//	}
//
//	var this float32 = 0
//	var last float32 = 0
//	var two float32 = 0
//	for i := 0; i < len(klint); i++ {
//		var attr = strings.Split(klint[i], ",")
//		var year = ToInt(attr[0][0:4])
//		if lastYear == year {
//			last = float32(ToFloat(attr[8]))
//		} else if thisYear == year {
//			this = float32(ToFloat(attr[8]))
//		}
//		//fmt.Println(
//		//	"时间=", attr[0],
//		//	" 开盘价=", attr[1],
//		//	" 收盘价=", attr[2],
//		//	" 涨跌幅=", attr[8],
//		//)
//	}
//
//	if this != 0 && last != 0 {
//		two = (100+last)*(1+this/100) - 100
//	}
//	if last == 0 { //可能是今年上市的，不受此限制
//		this = 0
//		two = 0
//	}
//
//	return last, this, two
//}
