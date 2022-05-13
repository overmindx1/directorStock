package task

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	ModelSocket "github.com/overmind/go_test/models/stockItemNoMem"
	"github.com/overmind/go_test/rootVar"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func RealTimeStockPrice() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("./config/director-stock-selection-13a579b578f4.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(oauth2.NoContext)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println("Unable to retrieve Sheets client: ", err)
	}

	spreadsheetId := "1p5xzy_-c-qgitZ1LyPzdVAKfyUIlT-FoP7sSe7XbyQ4"
	readRange := "RealTime!A1:D"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from sheet: ", err)
	}
	log.Print("開始抓取及時股價 Start")

	if resp != nil {
		if len(resp.Values) == 0 {
			fmt.Println("No data found.")
		} else {
			for _, row := range resp.Values {
				//將 interface 的資料都轉成 string
				var stockItems ModelSocket.StockItemsNoMem
				var count int64
				rowSlice := make([]string, 0)
				for _, str := range row {
					rowSlice = append(rowSlice, str.(string))
				}
				stockResult := rootVar.DbConn.Where("stock_id = ?", rowSlice[0]).Limit(1).Find(&stockItems)
				stockResult.Count(&count)

				// 有找到股票 - 更新
				if count > 0 {
					// 看一下是不是數字
					// fmt.Println(rowSlice)
					_, err = strconv.ParseFloat(rowSlice[2], 64)
					if err == nil {
						stockItems.ClosePrice = rowSlice[2]
						rootVar.DbConn.Save(&stockItems)
						// fmt.Println("StockData", rowSlice)
					} else {
						//fmt.Println("updateStockPriceError", err)
					}
				}
			}
		}
	}
	log.Print("開始抓取及時股價 End")
}
