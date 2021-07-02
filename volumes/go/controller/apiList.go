package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/overmind/go_test/models"
	"github.com/overmind/go_test/rootVar"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func GetTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"test": "test",
	})
}

func GetDirectorSelection(c *gin.Context) {
	var stockItems []models.StockItems
	var results []map[string]interface{}
	var stockIds []string

	if err := rootVar.DbConn.Table("director_selection").Select("stock_id").Where("on_list = ?", 1).Find(&results).Error; err != nil {
		log.Fatalln(err)
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Api GetDirectorSelection Error",
		})
		panic("Api GetDirectorSelection Error(取得會長選股錯誤)")
	}

	for _, row := range results {
		stockIds = append(stockIds, row["stock_id"].(string))
	}

	if err := rootVar.DbConn.Preload("DirectorSelection").Where("stock_id in ?", stockIds).Find(&stockItems).Error; err != nil {
		log.Fatalln(err)
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Api GetDirectorSelection Error",
		})
		panic("Api GetDirectorSelection Error(取得股票資訊錯誤)")
	}
	dataJson := make(map[string]interface{})
	dataJson["list"] = stockItems

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "",
		"data": stockItems,
	})
}

type MemberSpreadsheet struct {
	SpreadsheetId string `json:"spreadsheetId"`
}

func GetMemberSpreadsheetCompare(c *gin.Context) {
	var mss MemberSpreadsheet

	c.BindJSON(&mss)
	//mss.SpreadsheetId
	// init  google api
	ctx := context.Background()
	b, err := ioutil.ReadFile("./config/director-stock-selection-13a579b578f4.json")
	if err != nil {
		log.Println("Unable to read client secret file: , ", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Println("Unable to parse client secret file to config: ", err)
	}
	client := config.Client(oauth2.NoContext)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println("Unable to retrieve Sheets client: ", err)
	}

	spreadsheetId := mss.SpreadsheetId
	//代號,名稱,庫存數量,平均買進成本,手續費則扣,手續費低消,券商
	readRange := "A2:G"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	log.Print("抓取資料Google Sheet Id : ", spreadsheetId)
	if err != nil {
		errOrgMsg := err.Error()
		errMsg := ""

		if strings.Index(errOrgMsg, "503") != -1 {
			errMsg = "服務目前不可用"
		} else if strings.Index(errOrgMsg, "403") != -1 {
			errMsg = "表單權限沒有開放"
		}

		log.Println("Unable to retrieve data from sheet: ", err)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "Unable to retrieve data from sheet Id" + errMsg,
		})
		panic("Unable to retrieve data from sheet " + spreadsheetId)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		//先排好陣列 ...
		memberRowData := make(map[string][]models.MemberData, 0)
		sIds := make([]string, 0)
		for _, item := range resp.Values {
			if len(item) > 5 {
				sId := item[0].(string)
				if memberRowData[sId] != nil {
					var rowData models.MemberData
					rowData.MemStockId = item[0].(string)
					rowData.MemStockName = item[1].(string)
					rowData.MemShares = item[2].(string)
					rowData.MemAvgCost = item[3].(string)
					rowData.MemFeeDiscount = item[4].(string)
					rowData.MemLowestFee = item[5].(string)
					memberRowData[sId] = append(memberRowData[sId], rowData)
				} else {
					var rowData models.MemberData
					rowData.MemStockId = item[0].(string)
					rowData.MemStockName = item[1].(string)
					rowData.MemShares = item[2].(string)
					rowData.MemAvgCost = item[3].(string)
					rowData.MemFeeDiscount = item[4].(string)
					rowData.MemLowestFee = item[5].(string)
					memberRowData[sId] = make([]models.MemberData, 0)
					memberRowData[sId] = append(memberRowData[sId], rowData)
					sIds = append(sIds, sId)
				}
			}
		}
		//算平均值
		for i, mapArr := range memberRowData {
			// 如果同一個股票有一筆以上 則來算平均值
			count := len(mapArr)
			if count > 1 {
				shares := 0
				avgCosts := float64(0)
				cost := float64(0)
				for _, row := range mapArr {
					share, _ := strconv.Atoi(row.MemShares)
					shares = shares + share
					avgCost, _ := strconv.ParseFloat(row.MemAvgCost, 64)
					cost += float64(share) * avgCost
					// fmt.Println(avgCost, shares)
				}
				//平均成本價
				avgCosts = cost / float64(shares)
				var newRow models.MemberData
				newRow.MemStockId = memberRowData[i][0].MemStockId
				newRow.MemStockName = memberRowData[i][0].MemStockName
				newRow.MemShares = strconv.Itoa(shares)
				newRow.MemAvgCost = fmt.Sprintf("%.2f", avgCosts)
				newRow.MemFeeDiscount = memberRowData[i][0].MemFeeDiscount
				newRow.MemLowestFee = memberRowData[i][0].MemLowestFee
				memberRowData[i] = make([]models.MemberData, 0)
				memberRowData[i] = append(memberRowData[i], newRow)
			}
		}

		// 先找出所有的會長股列表
		var results []map[string]interface{}
		if err := rootVar.DbConn.Table("director_selection").Select("stock_id").Where("on_list = ?", 1).Find(&results).Error; err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "Api GetDirectorSelection Error",
			})
			panic("Api GetDirectorSelection Error(取得會長選股錯誤)")
		}
		//  找出所有要抓取的股票列表
		for _, row := range results {
			bools := containsString(sIds, row["stock_id"].(string))
			if !bools {
				sIds = append(sIds, row["stock_id"].(string))
			}
		}

		// 開始找出股票資料
		var stockItems []models.StockItems
		if err := rootVar.DbConn.Preload("DirectorSelection").Where("stock_id in ?", sIds).Find(&stockItems).Error; err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "Api GetDirectorSelection Error",
			})
			panic("Api GetDirectorSelection Error(取得股票資訊錯誤)")
		}

		// 會長股加上自己的股票所有整合在一起
		for i, stockItem := range stockItems {
			if memberRowData[stockItem.StockId] != nil {
				stockItem.MemberData = &memberRowData[stockItem.StockId][0]
				stockItems[i] = stockItem
			}
		}

		dataJson := make(map[string]interface{})
		dataJson["list"] = stockItems
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "",
			"data": stockItems,
		})
	}

}

func containsString(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
