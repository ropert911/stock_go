package main

import (
	"stock"
)

func main() {
	stock.DownloadStockGzfx() //估计分析表
	stock.DownloadStockYjbb() //业绩报表
	stock.DownloadStockZcfz() //资产负债表
	stock.DownloadStockLrb()  //利润表
}
