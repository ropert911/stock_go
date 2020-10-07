package https

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HTTPsGet(url string) (string, error) {
	//设置tls配置信息
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Get error:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

//func main() {
//	var content, err = HTTPsGet("https://192.168.20.46:10099/device/getDeviceExtras")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(content)
//	fmt.Println("")
//
//	var response map[string]interface{}
//	json.Unmarshal([]byte(content), &response)
//	var s = response["success"].(bool)
//	if s {
//		var data = response["data"].([]interface{})
//		for i := range data {
//			var deviceData = data[i].(map[string]interface{})
//			var id string
//			if deviceData["id"] != nil {
//				id = deviceData["id"].(string)
//			}
//			var model string
//			if deviceData["deviceType"] != nil {
//				model = deviceData["deviceType"].(string)
//			}
//
//			if len(id) > 0 && len(model) > 0 {
//				fmt.Println(id, model)
//			}
//		}
//	}
//
//}
