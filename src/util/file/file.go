package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

//判断文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//删除文件
func Delete(filename string) {
	os.Remove(filename)
}

//新建文件并写内容
func WriteFile(filename, data string) {
	var (
		err error
	)

	// 拿到一个文件对象
	// file对象肯定是实现了io.Reader,is.Writer
	fileObj, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(data)
	// 方式一
	_, _ = fmt.Fprintf(fileObj, "%s", data)
	fileObj.Close()

	// 方式二
	//writer := bufio.NewWriter(fileObj)
	//defer writer.Flush()
	//writer.WriteString(data)
}

//附加内容
func AppendFile(filename, data string) {
	fileObj, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("os open error:", err)
		return
	}

	defer fileObj.Close()
	//fmt.Println(data)
	_, err = fileObj.WriteString(data)
	if nil != err {
		fmt.Println("write error=", err)
	}
}

func ReadFile_v1(filename string) (string, error) {
	var (
		err     error
		content []byte
	)
	fileObj, err := os.Open(filename)
	if err != nil {
		fmt.Println("os open error:", err)
		return "", err
	}
	defer fileObj.Close()
	content, err = ioutil.ReadAll(fileObj)
	if err != nil {
		fmt.Println("ioutil.ReadAll error:", err)
		return "", err
	}

	return string(content), nil
}

// 还有种方法
func Readfile_v2(filename string) string {
	var (
		err     error
		content []byte
	)
	content, err = ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(content)
}
