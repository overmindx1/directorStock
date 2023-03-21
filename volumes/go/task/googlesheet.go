package task

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/overmind/go_test/rootVar"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"

	Models "github.com/overmind/go_test/models"
)

func GetDirectorStockSelection() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("./config/director-stock-selection-13a579b578f4.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return
	}
	client := config.Client(oauth2.NoContext)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println("Unable to retrieve Sheets client: ", err)
		return
	}

	spreadsheetId := "1uddRr--H95Texz9sEwY0z7GH7BBAZXMwlMaHUu7CIPs"
	readRange := "會長存股彙整表!A2:L"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from sheet: ", err)
		return
	}
	log.Print("開始抓取會長存股彙整表資料Start")
	if resp != nil {
		if len(resp.Values) == 0 {
			fmt.Println("No data found.")
		} else {
			var directorSelection Models.DirectorSelection
			// 先全部當成沒有在會長的選股狀態
			rootVar.DbConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&directorSelection).Updates(Models.DirectorSelection{OnList: 0, Shares: "0"})

			for _, row := range resp.Values {
				// Print columns A and E, which correspond to indices 0 and 4.
				//fmt.Printf("%s, %s\n", row[0], row[4])
				//將 interface 的資料都轉成 string
				rowSlice := make([]string, 0)
				for _, str := range row {
					rowSlice = append(rowSlice, str.(string))
				}

				var selection Models.DirectorSelection
				var count int64
				// row9  = 股票id
				stockResult := rootVar.DbConn.Where("stock_id = ?", row[9]).Limit(1).Find(&selection)
				stockResult.Count(&count)
				// 有找到會長存股
				if count > 0 {
					selection.Shares = rowSlice[3]
					selection.AvgCost = rowSlice[7]
					selection.OnList = 1
					rootVar.DbConn.Save(&selection)
				} else {
					selectionInsert := Models.DirectorSelection{
						StockId:   rowSlice[9],
						StockName: rowSlice[1],
						AvgCost:   rowSlice[7],
						Shares:    rowSlice[3],
						OnList:    1,
					}
					rootVar.DbConn.Create(&selectionInsert)
				}
			}
		}
	}
	log.Print("開始抓取會長存股彙整表資料 End!")
}
