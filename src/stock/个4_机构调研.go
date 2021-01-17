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

type SingleJkdy struct {
	HITS int
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

func ParseSingleJkdy(code string) int {
	endDate := time.Now().Format("2006-01-02")
	var month = endDate[0:7]
	var fileName = fmt.Sprintf(jgdy_singleformate, code, month)
	if !file.FileExist(fileName) {
		DowloadJgdy(code)
	}

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
