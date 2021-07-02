<template>
    <el-table 
        :data="tableData" 
        height="calc(100vh - 9.5rem)" 
        size="mini"
        stripe
        border
        :default-sort = "{prop: 'DirectorSelection.Percent', order: 'descending'}"
    >
        <!-- 一般資料 -->
        <el-table-column
            prop="StockName"
            label="股名"
            fixed
            class-name="rowClass"
            min-width="96"
            >
            <template slot-scope="scope">
                <a target="_blank" :href="`https://tw.stock.yahoo.com/q/bc?s=${scope.row.StockId}`">
                    {{scope.row.StockName}}
                </a>
            </template>
        </el-table-column>
        <el-table-column
            prop="StockId"
            label="股號"
            sortable
            class-name="rowClass"
            >
        </el-table-column>
        <el-table-column
            prop="ClosePrice"
            label="收盤價"
            sortable
            class-name="rowClass"
            min-width="96"
            >
        </el-table-column>
        <!-- 會員資料 -->
        <el-table-column
            prop="MemberData.MemAvgCost"
            label="我的成本"
            sortable
            class-name="rowClass"
            min-width="112"
            >
        </el-table-column>
        <el-table-column
            prop="MemberData.MemShares"
            label="我的股數"
            class-name="rowClass"
            min-width="112"
            >
        </el-table-column>
        <el-table-column
            prop="MemberData.Percent"
            label="我的收益"
            sortable
            :sort-method="sortPercentMemberData"
            class-name="rowClass"
            min-width="112"
            >
            <template slot-scope="scope">
                <span v-if="scope.row.MemberData !== null">
                    <div class="directorPercent" :class="{redBg : scope.row.MemberData.CompareIsUp , greenBg : !scope.row.MemberData.CompareIsUp}">
                        {{scope.row.MemberData.Percent}}
                    </div>
                </span>
                <span v-else>  </span>
           </template>
        </el-table-column>
        <el-table-column
            prop="MemberData.Profit"
            label="我的獲利"
            class-name="rowClass"
            min-width="112"
            >
            <template slot-scope="scope">
                <span v-if="scope.row.MemberData !== null">
                    <div class="directorPercent" :class="{redBg : scope.row.MemberData.CompareIsUp , greenBg : !scope.row.MemberData.CompareIsUp}">
                        {{scope.row.MemberData.Profit}}
                    </div>
                </span>
                <span v-else>  </span>
           </template>
        </el-table-column>
        <!-- 會長資料 -->
        <el-table-column
            prop="DirectorSelection.Percent"
            label="會長收益"
            sortable
            :sort-method="sortPercentDirector"
            class-name="rowClass"
            min-width="112"
           >
           <template slot-scope="scope">
               <span v-if="scope.row.DirectorSelection !== null">
                    <div class="directorPercent" :class="{redBg : scope.row.DirectorSelection.CompareIsUp , greenBg : !scope.row.DirectorSelection.CompareIsUp}">
                        {{scope.row.DirectorSelection.Percent}}
                    </div>
                </span>
                <span v-else>  </span>                           
           </template>
        </el-table-column>
        <el-table-column
            prop="DirectorSelection.AvgCost"
            label="會長成本"
            sortable
            class-name="rowClass"
            min-width="112"
           >
        </el-table-column>
        <el-table-column
            prop="DirectorSelection.Shares"
            label="會長持股"
            class-name="rowClass"
            min-width="112"
            >
        </el-table-column>
        
    </el-table>
</template>

<script>
export default {
    name : "StockTable",
    props : {
        tableData : {
            type : Array,
            required : true
        }
    },
    methods : {
        sortPercentDirector(a , b ) {
            
            if(b.DirectorSelection.Percent == "-" ) {
                return -1
            }
            if(b.DirectorSelection != null && a.DirectorSelection != null) {
                return b.DirectorSelection.Percent - a.DirectorSelection.Percent
            }
        },
        sortPercentMemberData(a , b ) {
            if(b.MemberData.Percent == "-" ) {
                return -1
            }
            if(b.MemberData != null && a.MemberData != null) {
                return b.MemberData.Percent - a.MemberData.Percent
            }
        }
    }
}
</script>

<style>

</style>