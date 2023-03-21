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
)

func RealTimeStockWrite() {
	fmt.Println("會長股及時股價處理開始!!! Start.")
	ctx := context.Background()
	b, err := ioutil.ReadFile("./config/director-stock-selection-13a579b578f4.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly", "https://www.googleapis.com/auth/spreadsheets")
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
	// 會長股及時股票
	spreadsheetId := "1p5xzy_-c-qgitZ1LyPzdVAKfyUIlT-FoP7sSe7XbyQ4"
	readRange := "RealTime!A2:D"

	//先讀取會長股列表

	var results []map[string]interface{}
	if err := rootVar.DbConn.Table("director_selection").Select("director_selection.stock_id", "director_selection.stock_name", "stock_items.stock_type").
		Joins("left join stock_items on stock_items.stock_id = director_selection.stock_id").Where("on_list = ?", 1).
		Find(&results).Error; err != nil {
		log.Fatalln(err)
		panic("RealTimeStockWrite Error(取得會長選股錯誤)")
	}
	// The ranges to clear, in A1 notation.
	// 清除會長股及時股票
	ranges := []string{}
	ranges = append(ranges, readRange)
	rb := &sheets.BatchClearValuesRequest{
		Ranges: ranges,
	}
	// 如果清除失敗
	_, err = srv.Spreadsheets.Values.BatchClear(spreadsheetId, rb).Context(ctx).Do()
	if err != nil {
		log.Println(err)
		return
	}
	// 可以自訂 , 不會讓全部寫入都是字串
	valueInputOption := "USER_ENTERED"
	// 寫入資料
	writeRows := [][]interface{}{}
	for _, row := range results {
		rowWrite := make([]interface{}, 4)
		rowWrite[0] = row["stock_id"].(string)
		rowWrite[1] = row["stock_name"].(string)
		if row["stock_type"].(string) == "Tpex" {
			rowWrite[2] = "=importxml(\"https://histock.tw/stock/" + row["stock_id"].(string) + "\" , \"/html/body/form/div[4]/div[3]/div/div[1]/div[1]/div[2]/div[1]/div[2]/ul/li[1]/span/span/span\")"
		} else {
			rowWrite[2] = "=GOOGLEFINANCE(\"TPE:\"&" + row["stock_id"].(string) + ", \"price\")"
		}
		rowWrite[3] = row["stock_type"].(string)
		writeRows = append(writeRows, rowWrite)
	}

	// 設定請求
	batcRequest := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: valueInputOption,
	}

	batcRequest.Data = append(batcRequest.Data, &sheets.ValueRange{
		Range:  readRange,
		Values: writeRows,
	})

	// 執行
	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, batcRequest).Context(ctx).Do()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("會長股及時股價寫入完成!!! Done.")

}
