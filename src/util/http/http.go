package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	//"strings"
	//"unsafe"
)

func HttpGet(url string) (*string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("request error 1", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll error 1", err)
		return nil, err
	}

	var content = string(body)
	return &content, nil
}

func HttpGetRequst(url string) (*string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//处理返回结果
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("request error 1", err)
		return nil, err
	}
	defer response.Body.Close()

	body2, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("request error 2", err)
		return nil, err
	}

	var content = string(body2)
	return &content, nil
}

func HttpPost(urls string, v url.Values) (*string, error) {
	resp, err := http.PostForm(urls, v)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var content = string(respBytes)
	return &content, nil
}

func HttpPostRequestBytesHeaderRest(url string, bytesData []byte) (*string, error) {
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//Header
	request.Header.Set("Content-Type", "application/json")

	//发送
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	//回应
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var content = string(respBytes)
	return &content, nil
}

func HttpPostRequestValuesHeaderRest(url string, v url.Values) (*string, error) {
	reader := ioutil.NopCloser(strings.NewReader(v.Encode()))
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//Header
	request.Header.Set("Content-Type", "application/json")

	//发送
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	//回应
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var content = string(respBytes)
	return &content, nil
}

type GpsInfo struct {
	DeviceId  string  `bson:"deviceid,omitempty" json:"deviceid"`
	Direction uint16  `bson:"direction,omitempty" json:"direction"`
	Lat       float32 `bson:"lat,omitempty" json:"lat"`
	Lng       float32 `bson:"lng,omitempty" json:"lng"`
	Speed     float32 `bson:"speed,omitempty" json:"speed"`
	Time      string  `bson:"time,omitempty" json:"time"`
	Vehicle   string  `bson:"vehicle,omitempty" json:"vehicle"`
}

func Test() {
	str, err := HttpGet("http://192.168.20.46:10099/device/getDeviceExtras")
	if nil == err {
		fmt.Println(*str)
	}

	fmt.Println("")
	fmt.Println("===========")
	str, err = HttpGetRequst("http://192.168.20.46:10099/device/getDeviceExtras")
	if nil == err {
		fmt.Println(*str)
	}

	fmt.Println("")
	fmt.Println("===========")
	str, err = HttpPost("http://192.168.20.56:18004/itc/api/schoolbus/gps", url.Values{"Userlist": {"1145,1150"}})
	if nil == err {
		fmt.Println(*str)
	}

	fmt.Println("")
	fmt.Println("===========")
	var gps = GpsInfo{DeviceId: "7070730", Direction: 120, Lat: 103.95077, Lng: 30.777979, Speed: 50.125, Time: "2019-06-19 17:48:40", Vehicle: "川A123"}
	bytesData, _ := json.Marshal(gps)
	str, err = HttpPostRequestBytesHeaderRest("http://192.168.20.56:18004/itc/api/schoolbus/gps", bytesData)
	if nil == err {
		fmt.Println(*str)
	}

	fmt.Println("")
	fmt.Println("===========")
	v := url.Values{}
	v.Set("DeviceId", "7070730")
	v.Set("Direction", "120")
	v.Set("Time", "2019-06-19 17:48:40")
	v.Set("Vehicle", "川A123")
	str, err = HttpPostRequestValuesHeaderRest("http://www.baidu.com", v)
	if nil == err {
		fmt.Println(*str)
	}
}
