<template>
    <div id="app">
        <el-row type="flex" justify="center" align="top" style="margin-top:1rem">
            <el-col :xs="22" :sm="22" :md="18" :lg="12" :xl="10">
                <el-form  :model="form" size="small" class="sheetForm">                    
                    <el-form-item label="表單">
                        <el-col :xs="11" :sm="12" :md="15" :lg="16" :xl="12">
                            <el-input v-model="currentUserData.userData.googleSpreadsheet" placeholder="" />
                        </el-col>
                        <el-col :xs="10" :sm="9" :md="7" :lg="5" :xl="5">
                            <el-button-group>
                                <el-button size="mini" type="primary" @click="sendMss">送出</el-button>
                                <el-button size="mini" type="primary" @click="drawerTrigger = true">選單</el-button>
                            </el-button-group>
                            <!-- <el-dropdown trigger="click" split-button type="primary"  @click="sendMss" @command="btnDropListCommand" size="small">
                                送出
                                <el-dropdown-menu slot="dropdown">
                                    <el-dropdown-item command="saveSheetUrl">記住表單資料</el-dropdown-item>
                                    <el-dropdown-item command="clearSheetUrl" divided>清除已儲存資料</el-dropdown-item>
                                    <el-dropdown-item command="showColumeSetting" divided>
                                        表格欄位設定
                                    </el-dropdown-item>
                                    <el-dropdown-item command="showStatic" divided>
                                        {{ staticTrigger == true ? '關閉' : '顯示' }}統計資料
                                    </el-dropdown-item>
                                    <el-dropdown-item command="wishList" divided>
                                        功能願望
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </el-dropdown> -->
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
                        <li>尚未實現損益 : {{ parseFloat(statics.director.totalProfit).toFixed(2) }}</li>
                    </ul>
                </div>
                <div class="dataBlock">
                    您的統計資料
                    <ul>
                        <li>您的股票筆數 : {{ statics.member.stockCount }}</li>
                        <li>您的投資成本 : {{ parseFloat(statics.member.totalInvestment).toFixed(2) }}</li>
                        <li>目前投資價值 : {{ parseFloat(statics.member.totalNowValue).toFixed(2) }}</li>
                        <li>尚未實現損益 : {{ parseFloat(statics.member.totalProfit).toFixed(2) }} (已扣除手續費及證交稅)</li>
                    </ul>
                </div>
            </el-col>
        </el-row>

        <el-row type="flex" justify="center" align="top" >
            <el-col :xs="22" :sm="22" :md="22" :lg="22" :xl="22">
                <el-form  :model="form" size="small"  >
                    <el-form-item label="">
                        <el-col :xs="9" :sm="4" :md="4" :lg="3" :xl="2" style="font-size:1rem;text-align:left">
                            <el-dropdown trigger="click"  @command="changeSearchType">
                                <el-button split-button type="primary">
                                    {{tableSearchLable}} : <i class="el-icon-arrow-down el-icon--right"></i>
                                </el-button>
                                <el-dropdown-menu slot="dropdown">
                                    <el-dropdown-item command="stockId" >股號搜尋</el-dropdown-item>
                                    <el-dropdown-item command="stockName" divided>股名搜尋</el-dropdown-item>
                                </el-dropdown-menu>
                            </el-dropdown>
                        </el-col>
                        <el-col :xs="14" :sm="8" :md="6" :lg="3" :xl="3" style="font-size:1rem">
                            <el-input v-model="form.stockSearch"></el-input>
                        </el-col>
                    </el-form-item>
                </el-form>
                <stock-tablevxe ref="vxeTable" :tableData="tableData" :tableColumes="currentUserData.userData.tableColumes"></stock-tablevxe>
                <!-- <stock-table :tableData="tableData" :tableColumes="tableColumes"></stock-table> -->
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

        <!-- 設定要顯示的欄位 -->
        <table-enable-columes :tableColumes="currentUserData.userData.tableColumes" :triggerColumeSetting="triggerColumeSetting" 
         @saveColumesShow="saveColumesShow" @resetColumeSetting="resetColumeSetting" @closeDialog="triggerColumeSetting = false"></table-enable-columes>

        <!-- 使用者相關 -->
        <change-user :dbInstance="dbInstance" :currentUserData="currentUserData" :userConfigTaigger="userConfigTaigger"
            @closeDialog="userConfigTaigger = false" @changeCurrentUser="changeCurrentUser"
        ></change-user>

        <!-- 抽屜 -->
        <drawer-menu :drawerTrigger="drawerTrigger" :currentUserData="currentUserData"
         @saveUserData="saveUserData" @closedDrawer="closedDrawer" @openFromDrawer="openFromDrawer"></drawer-menu>
    </div>
</template>

<script>
import{ getDirectorSelection , sendMmeberSpreadsheet} from '@/utils/api'
import StockTable from '@/components/stockTable'
import StockTablevxe from '@/components/stockTablevxe'
import TableEnableColumes from '@/components/tableEnableColumes'
import ChangeUser from '@/components/changeUser'
import DrawerMenu from '@/components/drawerMenu'
import DefaultTableColumes  from '@/utils/tableColumes'
import dbInstance from '@/utils/indexedDbWarp'
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
    components : {StockTable , StockTablevxe , TableEnableColumes , ChangeUser , DrawerMenu} ,
    data : () => ({
        form : {
            stockSearch : "",
            searchType : "stockId" // stockId 股號搜索 stockName 股名搜索
        },
        tableColumes : DefaultTableColumes, //欄位設定
        triggerColumeSetting : false, //欄位設定開關
        staticTrigger : false, //統計資料開關
        userConfigTaigger :false, // 使用者相關
        drawerTrigger : false, //測欄選單開關
        
        statics : JSON.parse(JSON.stringify(DEFAULT_STATIC)),
        rawData : [],
        // 使用者資料
        currentUserData : {
            nickname : '',
            userData : {
                id : 0,
                googleSpreadsheet : '',
                tableColumes : {}
            }
        },
        getAllUser : [],
        dbInstance : dbInstance
    }),
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
                    // console.log([memberProfit , item])
                }
            });

            this.statics.director.totalProfit = this.statics.director.totalNowValue - this.statics.director.totalInvestment
            this.statics.member.totalProfit = parseFloat( memberProfit.toFixed(2) )
        }
    },
    computed : {
        tableData(){
            let find = []
            if(this.form.searchType == "stockId") {
                if(this.form.stockSearch == "" || this.form.stockSearch.length < 2) {
                    return this.rawData
                } 
                find = this.rawData.filter( item => {
                    if(item.StockId.indexOf(this.form.stockSearch) == 0) {
                        return item
                    }
                }) 
            } else {
                if(this.form.stockSearch == "" || this.form.stockSearch.length < 2) {
                    return this.rawData
                } 
                find = this.rawData.filter( item => {
                    if(item.StockName.indexOf(this.form.stockSearch) != -1) {
                        return item
                    }
                }) 
            }

            return find;
        },
        tableSearchLable () {
            if(this.form.searchType == "stockId") {
                return "股號搜尋";
            } else {
                return "股名搜尋";
            }
        }
    },
    methods : {
        // 取得相關的股票資料
        async GDS(){
            let data = await getDirectorSelection();
            this.rawData = data;
            
        },
        async sendMss(){
            if(this.currentUserData.userData.googleSpreadsheet.indexOf("https://docs.google.com/spreadsheets/d/") !== -1) {
                let spreadsheetId = this.currentUserData.userData.googleSpreadsheet.split("/")[5]
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
            if(command == 'showColumeSetting') {
                this.triggerColumeSetting = true;
            }
            if(command == 'showStatic') {
                this.staticTrigger = !this.staticTrigger
            }
            if(command == 'wishList') {
                window.open("https://forms.gle/4zPrpWXsSxSohusHA")
            }
        },
        changeSearchType(command) {
            this.form.searchType = command
        },
        // 存下已修改的顯示欄位
        saveColumesShow(obj) {
            this.currentUserData.userData.tableColumes = JSON.parse(JSON.stringify(obj));
            
            // 存在indexedDb
            this.dbInstance.updateUserData(this.currentUserData.userData.nickname , this.currentUserData.userData);
            setTimeout(() => {
                location.reload();
            },100)
            
        },

        resetColumeSetting() {
            this.currentUserData.userData.tableColumes = {};
            // 存在indexedDb
            this.dbInstance.updateUserData(this.currentUserData.userData.nickname , this.currentUserData.userData);
            setTimeout(() => {
                location.reload();
            },100)
        },
        
        // 儲存使用者資料
        async saveUserData(){
            let execute = await this.dbInstance.updateUserData(this.currentUserData.nickname , this.currentUserData.userData);
            if(execute) {
                this.$message(`使用者 ${this.currentUserData.nickname} 資料更新!`);
                this.drawerTrigger = false;
            } else {
                this.$message({
                    message: `使用者 ${this.currentUserData.nickname} 資料更新失敗!`,
                    type: 'error'
                });
            }
        },
        // 單向資料流 更換使用者
        async changeCurrentUser(row) {
            let userData = await this.dbInstance.setDefaultUser(row.nickname);
            if(userData) {
                this.currentUserData.nickname = userData.nickname;
                this.currentUserData.userData = userData;
                if(Object.keys(this.currentUserData.userData.tableColumes).length == 0) {
                    this.currentUserData.userData.tableColumes = DefaultTableColumes
                }
                this.$message(`更換使用者 ${row.nickname} 成功!`);
                this.drawerTrigger = false;
            } else {
                this.$message({
                    message: `更換使用者 ${row.nickname} 失敗!`,
                    type: 'error'
                });
            }
        },
        // 
        closedDrawer(){
            this.drawerTrigger = false;
        },
        // 
        openFromDrawer(trigger) {
            if(trigger == "staticTrigger") {
                this.staticTrigger = !this.staticTrigger
            } else {
                this[trigger] = true
            }
        }

    },
    async mounted() {
        this.GDS();
        let userData = await this.dbInstance.getDefaultUserData();
        this.currentUserData.nickname = userData.nickname;
        this.currentUserData.userData = userData;
        if(Object.keys(this.currentUserData.userData.tableColumes).length == 0) {
            this.currentUserData.userData.tableColumes = DefaultTableColumes
        } else {
            // 更新欄位
            if(this.currentUserData.userData.tableColumes.SharesRatio == undefined) {
                this.currentUserData.userData.tableColumes = DefaultTableColumes
            }
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

.amethystBg {
    background: #9b59b6 !important;
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
    @media (max-width:400px) {
        font-size: 0.9rem;
    }
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
