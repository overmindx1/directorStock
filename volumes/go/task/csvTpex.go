package task

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	ModelSocket "github.com/overmind/go_test/models/stockItemNoMem"
	"github.com/overmind/go_test/rootVar"
)

// 取得上櫃資料
func ParseTpexCSV() {
	//filename
	date := time.Now().Format("20060102")
	year, monthObj, dayInt := time.Now().Date()
	rocYear := year - 1911
	var month string
	var day string
	if int(monthObj) < 10 {
		month = "0" + strconv.Itoa(int(monthObj))
	} else {
		month = strconv.Itoa(int(monthObj))
	}

	if dayInt < 10 {
		day = "0" + strconv.Itoa(dayInt)
	} else {
		day = strconv.Itoa(dayInt)
	}
	rocDate := strconv.Itoa(rocYear) + "/" + month + "/" + day
	//寫入檔案的檔
	fileName := "./csv/" + date + "Tpex.csv"

	out, err := os.Create(fileName)
	if err != nil {
		fmt.Println("write csv file", err)
	}
	defer out.Close()
	log.Print("開始抓取Tpex 的今日 股票資料集 " + rocDate)
	resp, err := http.Get("http://www.tpex.org.tw/web/stock/aftertrading/daily_close_quotes/stk_quote_download.php?l=zh-tw&d=" + rocDate + "&s=0,asc,0")
	if err != nil {
		fmt.Println("http (I/O Timeout)", err)
		return
	}
	defer resp.Body.Close()

	if resp != nil {
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Println("write file error", err)
		}

		file, err := os.Open(fileName)

		if err != nil {
			fmt.Println("failed to open")

		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		//line by line
		scanner.Split(bufio.ScanLines)
		// var texts []string //slice
		isStart := false
		i := 0
		for scanner.Scan() {
			// text := append(text, scanner.Text())

			bytes := []byte(scanner.Text())
			//big5 => utf8
			bytes, err := decodeBig5(bytes)
			if err != nil {
				fmt.Println(err)
			}

			isFound := strings.Index(string(bytes), "代號")
			if isFound != -1 {
				isStart = true
			}

			// 跑到最後找關鍵
			isFound = strings.Index(string(bytes), "管理股票")
			if isFound != -1 {
				isStart = false
			}

			if isStart && i > 1 {
				var stockItems ModelSocket.StockItemsNoMem
				var count int64
				csvReader := csv.NewReader(strings.NewReader(strings.Replace(string(bytes), "=", "", -1)))
				//fmt.Println(strings.Replace(string(bytes), "=", "", -1))
				record, err := csvReader.Read()
				if err != nil {
					fmt.Println(err)
				}

				// 此行沒資料
				if len(record) < 10 {
					continue
				}

				// 找出有沒有股票代號  有的畫拉出修改 沒有的就新增
				stockResult := rootVar.DbConn.Where("stock_id = ?", record[0]).Limit(1).Find(&stockItems)
				// 數量
				stockResult.Count(&count)

				var check string
				if strings.Index(record[3], "+") != -1 {
					check = "+"
					strings.Replace(record[3], "+", "", -1)
				} else if strings.Index(record[3], "-") != -1 {
					check = "-"
					strings.Replace(record[3], "-", "", -1)
				} else {
					check = ""
				}
				if count > 0 {
					// 看看當天收盤價有沒有異常
					_, err = strconv.ParseFloat(record[2], 64)
					if err == nil {
						stockItems.ClosePrice = record[2]
					} else {
						//fmt.Println("csvTpexWritePriceError", err)
					}
					stockItems.OpenPrice = record[4]
					stockItems.UpDown = check
					stockItems.UpDownCount = record[3]
					stockItems.Date = time.Now().Format("2006-01-02")
					stockItems.StockType = "Tpex"
					rootVar.DbConn.Save(&stockItems)
				} else {
					// 看看當天收盤價有沒有異常
					_, err := strconv.ParseFloat(record[2], 64)
					if err == nil {
						record[2] = "0"
					} else {
						//fmt.Println("csvTpexWritePriceError", err)
					}
					stockInsert := ModelSocket.StockItemsNoMem{
						StockId:     record[0],
						StockName:   record[1],
						StockType:   "Tpex",
						OpenPrice:   record[4],
						ClosePrice:  record[2],
						UpDown:      check,
						UpDownCount: record[3],
						Date:        date,
					}
					rootVar.DbConn.Create(&stockInsert)
				}
			}

			// 證券代號的下一行開始
			if isStart {
				i++
			}
		}
	}
}
