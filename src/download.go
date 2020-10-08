package main

import (
	"stock"
)

func main() {
	stock.DownloadStockGzfx() //估计分析表 -- 和现价有关
	stock.DownloadStockYjbb() //业绩报表	-- 报表
	stock.DownloadStockZcfz() //资产负债表	-- 报表
	stock.DownloadStockLrb()  //利润表	-- 报表
	stock.DownloadStockXjll() //现金流量表	-- 报表
}
