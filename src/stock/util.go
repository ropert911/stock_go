package stock

import (
	"strconv"
	"strings"
)

func ToFloat(data string) float64 {
	f, _ := strconv.ParseFloat(data, 64)
	return f
}

func ToInt(data string) int {
	i, _ := strconv.Atoi(data)
	return i
}

//根据代号确认url查询里的市场+代号
func getSCByCode(code string) string {
	switch code[0:1] {
	case "6":
		return "SH" + code
	default:
		return "SZ" + code
	}
}

func GetExportCodeByCode(code string) string {
	switch code[0:1] {
	case "6":
		return "1" + code
	default:
		return "0" + code
	}
}

//得到年月日
func GetDate(data string) (int, int) {
	data = strings.ReplaceAll(data, `/`, `-`)
	data = data[0:7]
	arr := strings.Split(data, "-")
	return ToInt(arr[0]), ToInt(arr[1])
}
