<template>
    <div id="app">
        <el-row type="flex" justify="center" align="top" style="margin-top:1rem">
            <el-col :xs="22" :sm="22" :md="18" :lg="12" :xl="10">
                <el-form  :model="form" size="small" class="sheetForm">                    
                    <el-form-item label="表單">
                        <el-col :xs="13" :sm="13" :md="17" :lg="16" :xl="12">
                            <el-input v-model="form.spreadsheed" placeholder="" />
                        </el-col>
                        <el-col :xs="7" :sm="7" :md="5" :lg="4" :xl="3">
                            <el-dropdown split-button type="primary"  @click="sendMss" @command="btnDropListCommand" size="mini">
                                送出
                                <el-dropdown-menu slot="dropdown">
                                    <el-dropdown-item command="saveSheetUrl">記住表單資料</el-dropdown-item>
                                    <el-dropdown-item command="clearSheetUrl" divided>清除已儲存資料</el-dropdown-item>
                                    <el-dropdown-item command="showStatic" divided>
                                        {{ staticTrigger == true ? '關閉' : '顯示' }}統計資料
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </el-dropdown>
                            <!-- <el-button>送出</el-button> -->
                        </el-col>
                    </el-form-item> 
                </el-form>
            </el-col>
        </el-row>

        <el-row type="flex" justify="center" align="top" >
            <el-col :xs="22" :sm="22" :md="22" :lg="22" :xl="22">
                <el-alert title="會長存股列表" type="success" :closable="false" />
            </el-col>
        </el-row>
        
        <el-row type="flex" justify="center" align="top">
            <el-col :xs="22" :sm="22" :md="22" :lg="22" :xl="22" class="staticClass" :class="{staticClassOpen : staticTrigger}">
                <div class="dataBlock">
                    會長統計資料
                    <ul>
                        <li>會長股票筆數 : {{ statics.director.stockCount }}</li>
                        <li>會長投資成本 : {{ parseFloat(statics.director.totalInvestment).toFixed(2) }}</li>
                        <li>目前投資價值 : {{ parseFloat(statics.director.totalNowValue).toFixed(2) }}</li>
                        <li>目前投資營利 : {{ parseFloat(statics.director.totalProfit).toFixed(2) }}</li>
                    </ul>
                </div>
                <div class="dataBlock">
                    您的統計資料
                    <ul>
                        <li>您的股票筆數 : {{ statics.member.stockCount }}</li>
                        <li>您的投資成本 : {{ parseFloat(statics.member.totalInvestment).toFixed(2) }}</li>
                        <li>目前投資價值 : {{ parseFloat(statics.member.totalNowValue).toFixed(2) }}</li>
                        <li>目前投資營利 : {{ parseFloat(statics.member.totalProfit).toFixed(2) }} (已扣除手續費及證交稅)</li>
                    </ul>
                </div>
            </el-col>
        </el-row>

        <el-row type="flex" justify="center" align="top" >
            <el-col :xs="22" :sm="22" :md="22" :lg="22" :xl="22">
                <el-form  :model="form" size="small" label-width="5rem">
                    <el-form-item label="股號搜尋">
                        <el-col :xs="22" :sm="8" :md="6" :lg="3" :xl="2" style="font-size:1rem">
                            <el-input v-model="form.stockId"></el-input>
                        </el-col>
                    </el-form-item>
                </el-form>
                <stock-table :tableData="tableData"></stock-table>
            </el-col>
        </el-row>

        <el-row type="flex" justify="center" align="top">
            <el-col :xs="22" :sm="22" :md="22" :lg="22" :xl="22">
                <el-alert
                    center
                    title="免責聲明"
                    description="計算的百分比跟獲利僅供參考 , Google Sheet 資料會經過服務主機 , 資料風險須自行判斷!"
                    type="info"
                    show-icon
                    :closable="false">
                </el-alert>
            </el-col>
        </el-row>
    </div>
</template>

<script>
import{ getDirectorSelection , sendMmeberSpreadsheet} from '@/utils/api'
import StockTable from '@/components/stockTable'
const DEFAULT_STATIC = {
    director : {
        stockCount : 0 ,
        totalInvestment : 0,
        totalProfit : 0,
        totalNowValue : 0 
    },
    member : {
        stockCount : 0 ,
        totalInvestment : 0,
        totalProfit : 0,
        totalNowValue : 0 
    }
}
export default {
    name: "App",
    data : () => ({
        form : {
            spreadsheed : "",
            stockId : ""
        },
        staticTrigger : false,
        statics : JSON.parse(JSON.stringify(DEFAULT_STATIC)),
        rawData : []
    }),
    components : {StockTable},
    watch : {
        rawData(newVal , oldVal) {
            //計算統計資料
            this.statics = JSON.parse(JSON.stringify(DEFAULT_STATIC))
            let memberProfit = 0;
            this.rawData.map( item => {
                let directorCost = 0;
                let directorValue = 0;
                let memberCost = 0;
                let memberValue = 0;
                
                // 會長資料計算
                if(item.DirectorSelection != null  && Object.keys(item.DirectorSelection).length > 3) {
                    this.statics.director.stockCount++;
                    
                    let share = parseInt( item.DirectorSelection.Shares.replaceAll("," , ""))
                    directorCost = item.DirectorSelection.AvgCost * share
                    directorValue = item.ClosePrice * share
                    
                    this.statics.director.totalInvestment = this.statics.director.totalInvestment + directorCost
                    this.statics.director.totalNowValue = this.statics.director.totalNowValue + directorValue
                }
                
                // 會員統計計算
                if(item.MemberData !== null && Object.keys(item.MemberData).length > 3) {
                    this.statics.member.stockCount++;
                    let share = parseInt( item.MemberData.MemShares.replaceAll("," , ""))
                    memberCost = item.MemberData.MemAvgCost * share
                    memberValue = item.ClosePrice * share
                    this.statics.member.totalInvestment = this.statics.member.totalInvestment + memberCost
                    this.statics.member.totalNowValue = this.statics.member.totalNowValue + memberValue
                    memberProfit = memberProfit + item.MemberData.Profit;
                    
                }
            });

            this.statics.director.totalProfit = this.statics.director.totalNowValue - this.statics.director.totalInvestment
            this.statics.member.totalProfit = parseFloat( memberProfit.toFixed(2) )
        }
    },
    computed : {
        tableData(){
            if(this.form.stockId == "" || this.form.stockId.length < 4) {
                return this.rawData
            } 
            const find = this.rawData.find( item => {
                return item.StockId == this.form.stockId
            })
            return [find]
        }
    },
    methods : {
        // 取得相關的股票資料
        async GDS(){
            let data = await getDirectorSelection();
            this.rawData = data;
            
        },
        async sendMss(){
            //https://docs.google.com/spreadsheets/d/
            if(this.form.spreadsheed.indexOf("https://docs.google.com/spreadsheets/d/") !== -1) {
                let spreadsheetId = this.form.spreadsheed.split("/")[5]
                let data = await sendMmeberSpreadsheet({spreadsheetId})
                this.rawData = data;
            } else {
                this.$message('送出的試算表網址格式不對');
            }
        },
        btnDropListCommand(command){
            if(command == 'saveSheetUrl') {
                localStorage.setItem('googleSheetUrl' , this.form.spreadsheed)
                this.$message("已儲存您的資料, 下次來到本服務會幫你自動填上!");
            }
            if(command == 'clearSheetUrl') {
                localStorage.clear()
                this.$message("已清除您的資料");
            }
            if(command == 'showStatic') {
                this.staticTrigger = !this.staticTrigger
            }
        }
    },
    mounted() {
        this.GDS();
        if(localStorage.getItem('googleSheetUrl') !== null ) {
            this.form.spreadsheed = localStorage.getItem('googleSheetUrl');
        }
    }
};
</script>

<style lang="scss">
#app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
   
    
}

body , html {
    background-color: #f1f3f4 !important;
    margin : 0
}

.redBg {
    background: #e74c3c;
    color : white;
}

.greenBg {
    background: #2ecc71;
    color : white;
}
.directorPercent {
    text-align: center;
}
.rowClass {
    font-size: 1rem;
    text-align: center !important;
}
.staticClass {
    font-size: 1rem;
    text-align: left;
    background: #dfe6e9;
    margin : .25rem 0;
    
    border-radius: .5rem;
    ul {
        margin : 0
    }

    .dataBlock {
        padding : .5rem;
        margin : .5rem 0 ;
        &:last-of-type {
            margin-bottom: 1rem;
        }
    }
    height: 0px;
    transition: all ease-out .3s;
    overflow: hidden;
}
.staticClassOpen {
    height : 16.5rem
}
.sheetForm {
     .el-form-item {
        margin-bottom: 10px !important;
    }
}
</style>
