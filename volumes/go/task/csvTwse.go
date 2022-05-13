package task

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"

	ModelSocket "github.com/overmind/go_test/models/stockItemNoMem"
	"github.com/overmind/go_test/rootVar"
)

func ParseCSV() {
	//filename
	date := time.Now().Format("20060102")
	//寫入檔案的檔
	fileName := "./csv/" + date + "Twse.csv"
	//fileName := "./csv/20210609.csv"
	out, err := os.Create(fileName)
	if err != nil {
		fmt.Println("write csv file", err)
	}
	defer out.Close()
	// download twse csv file
	//resp, err := http.Get("https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=" + date + "&type=ALL")
	log.Print("開始抓取TWSE 的今日 股票資料集")
	resp, err := http.Get("https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=" + date + "&type=ALL")
	if err != nil {
		fmt.Println("http (sqlConnect)", err)
	}
	defer resp.Body.Close()

	// Write the body to file
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

		isFound := strings.Index(string(bytes), "證券代號")
		if isFound != -1 {
			isStart = true
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

			// 找出有沒有股票代號  有的畫拉出修改 沒有的就新增
			stockResult := rootVar.DbConn.Where("stock_id = ?", record[0]).Limit(1).Find(&stockItems)
			// 數量
			stockResult.Count(&count)
			if count > 0 {
				// 看看當天收盤價有沒有異常
				_, err = strconv.ParseFloat(record[8], 64)
				if err == nil {
					stockItems.ClosePrice = record[8]
				} else {
					//fmt.Println("csvTwseWritePriceError", err)
				}
				stockItems.OpenPrice = record[5]
				stockItems.UpDown = record[9]
				stockItems.UpDownCount = record[10]
				stockItems.Date = time.Now().Format("2006-01-02")
				stockItems.StockType = "Twse"
				rootVar.DbConn.Save(&stockItems)
			} else {
				// 看看當天收盤價有沒有異常
				_, err := strconv.ParseFloat(record[8], 64)
				if err == nil {
					record[8] = "0"
				} else {
					//fmt.Println("csvTwseWritePriceError", err)
				}
				stockInsert := ModelSocket.StockItemsNoMem{
					StockId:     record[0],
					StockName:   record[1],
					StockType:   "Twse",
					OpenPrice:   record[5],
					ClosePrice:  record[8],
					UpDown:      record[9],
					UpDownCount: record[10],
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

func ParseDirectorSelection() {

}

//convert BIG5 to UTF-8
func decodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
