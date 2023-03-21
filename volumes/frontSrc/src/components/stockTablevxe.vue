<template>
    <vxe-table 
        ref="table"
        :data="tableData" 
        :height="tableHeight" 
        stripe border show-overflow highlight-hover-row highlight-hover-column 
        :sort-config="{sortMethod: sortMethods ,trigger: 'cell', orders: ['desc', 'asc', null]}">
        
        <vxe-column field="StockName" title="股名" min-width="96" fixed="left" class-name="rowClass">
            <template #default="{ row }">
                <a target="_blank" :href="`https://tw.stock.yahoo.com/q/bc?s=${row.StockId}`">
                    {{row.StockName}}
                </a>
            </template>
        </vxe-column>
        <vxe-column field="StockId" title="股號" min-width="96" sortable class-name="rowClass"></vxe-column>
        <vxe-column field="ClosePrice" title="現價" min-width="96" sortable class-name="rowClass"></vxe-column>
        <!-- 我的資料 -->
        <template v-for="(column , key) in tableColumeSetting">
            <vxe-column 
                :key="column.index"
                class-name="rowClass"
                v-if="column.show"
                :field="column.field"
                :title="column.label"
                :min-width="column.minWidth"
                :sortable="column.sortable"
            >
                <template #default="{ row }" >
                    <div v-if="column.field == 'MemberData.Profit'">
                        <span v-if="row.MemberData !== null">
                            <div class="directorPercent" :class="{amethystBg : row.DirectorSelection.Percent == '-', redBg : row.MemberData.CompareIsUp , greenBg : !row.MemberData.CompareIsUp}">
                                {{row.MemberData.Profit}}
                            </div>
                        </span>
                        <span v-else>  </span>
                    </div>
                    <div v-else-if="column.field == 'MemberData.Percent'">
                        <span v-if="row.MemberData !== null">
                            <div class="directorPercent" :class="{ amethystBg : row.DirectorSelection.Percent == '-' , redBg : row.MemberData.CompareIsUp , greenBg : !row.MemberData.CompareIsUp}">
                                {{row.MemberData.Percent == "-" ? "-" : `${row.MemberData.Percent} %`}}
                            </div>
                        </span>
                        <span v-else>  </span>
                    </div>
                    <!-- <div v-else-if="column.field == 'SharesRatio'">
                        {{ row | shareRatio }}
                    </div> -->
                    <div v-else-if="column.field == 'DirectorSelection.Percent'">
                        <span v-if="row.DirectorSelection !== null">
                            <div class="directorPercent" :class="{amethystBg : row.DirectorSelection.Percent == '-' ,redBg : row.DirectorSelection.CompareIsUp , greenBg : !row.DirectorSelection.CompareIsUp}">
                                {{row.DirectorSelection.Percent == "-" ? "-" : `${row.DirectorSelection.Percent} %`}}
                            </div>
                        </span>
                        <span v-else>  </span>
                    </div>
                    <div v-else>
                        {{ row | showValueByStringKey(column.field)  }}
                    </div>
                </template>

                
            </vxe-column>
        </template>
        
        <!-- <vxe-column field="MemberData.MemAvgCost" title="我的成本" min-width="112" sortable class-name="rowClass"
            v-if="tableColumes['MemberData.MemAvgCost'].show" ></vxe-column>
        <vxe-column field="MemberData.MemShares" title="我的股數" min-width="112" sortable  class-name="rowClass" v-if="tableColumes['MemberData.MemShares'].show" ></vxe-column>
        <vxe-column field="MemberData.CostValue" title="我的投資" min-width="112" sortable  class-name="rowClass" v-if="tableColumes['MemberData.CostValue'].show" ></vxe-column>
        <vxe-column field="MemberData.Profit" title="我的獲利" min-width="112" sortable  class-name="rowClass" v-if="tableColumes['MemberData.Profit'].show" >
            <template #default="{ row }">
                <span v-if="row.MemberData !== null">
                    <div class="directorPercent" :class="{redBg : row.MemberData.CompareIsUp , greenBg : !row.MemberData.CompareIsUp}">
                        {{row.MemberData.Profit}}
                    </div>
                </span>
                <span v-else>  </span>
            </template>
        </vxe-column>
        <vxe-column field="MemberData.Percent" title="我的收益" min-width="112" sortable class-name="rowClass" v-if="tableColumes['MemberData.Percent'].show" >
            <template #default="{ row }">
                <span v-if="row.MemberData !== null">
                    <div class="directorPercent" :class="{redBg : row.MemberData.CompareIsUp , greenBg : !row.MemberData.CompareIsUp}">
                        {{row.MemberData.Percent == "-" ? "-" : `${row.MemberData.Percent} %`}}
                    </div>
                </span>
                <span v-else>  </span>
           </template>
        </vxe-column>
        <vxe-column title="存股比例" min-width="112" class-name="rowClass">
            <template #default="{ row }">
                <div>
                    {{ row | shareRatio }}
                </div>
            </template>
        </vxe-column> -->
        <!-- 會長資料 -->
        <!-- <vxe-column field="DirectorSelection.Percent" title="會長收益" min-width="112" sortable class-name="rowClass" v-if="tableColumes['DirectorSelection.Percent'].show" >
            <template #default="{ row }">
                <span v-if="row.DirectorSelection !== null">
                    <div class="directorPercent" :class="{redBg : row.DirectorSelection.CompareIsUp , greenBg : !row.DirectorSelection.CompareIsUp}">
                        {{row.DirectorSelection.Percent == "-" ? "-" : `${row.DirectorSelection.Percent} %`}}
                    </div>
                </span>
                <span v-else>  </span>
           </template>
        </vxe-column>
        <vxe-column field="DirectorSelection.AvgCost" title="會長成本" min-width="112" sortable class-name="rowClass" v-if="tableColumes['DirectorSelection.AvgCost'].show" > </vxe-column>
        <vxe-column field="DirectorSelection.Shares" title="會長持股" min-width="112" v-if="tableColumes['DirectorSelection.Shares'].show" > </vxe-column>
         -->
    </vxe-table>
</template>

<script>

export default {
    name : "StockTable",
    data : () => ({
        tableHeight : 400 // default
    }),
    props : {
        tableData : {
            type : Array,
            required : true
        },
        tableColumes : {
            type : Object,
            required : true
        }
    },
    computed : {
        tableColumeSetting(){
            let arr = [];
            for (const key in this.tableColumes) {
                let obj = this.tableColumes[key]
                obj.field = key
                arr.push(obj)
            };
            let array = arr.sort( (a, b) => {
                return a.index - b.index;
            })
            return array;
        }
    },
    filters : {
        showValueByStringKey(row , key) {
            let [key1 , key2] = key.split('.');
            if(row[key1] == undefined) {
                return '';
            } else if (row[key1] != undefined && row[key1][key2] == undefined) {
                if(key1 == 'MemberData' || key1 == 'DirectorSelection') {
                    return '-';
                }
                return row[key1];
            } else {
                return row[key1][key2];
            }
        },
        // shareRatio(row){
        //     if(row.MemberData == null || row.MemberData.MemStockName == undefined || row.DirectorSelection.StockName == undefined) {
        //         return '-';
        //     } else {
        //         let memberShare = row.MemberData.MemShares;
        //         let directorShare = row.DirectorSelection.Shares;
        //         directorShare = row.DirectorSelection.Shares.replace(',','');
        //         return ((memberShare/directorShare) * 100).toFixed(2) + ' %'
        //     }
        // }
    },
    methods : {
        calcTableHeight(vh){
            let h = Math.max(document.documentElement.clientHeight, window.innerHeight || 0);
            this.tableHeight = ( (vh * h) / 100 ) - 152;
        },
        sortMethods({ data, sortList }){
            const sortItem = sortList[0]
            // 取出第一个排序的列
            const { property, order } = sortItem
            let list = []
            if (order === 'asc' || order === 'desc') {
                let [mainKey , subKey] = property.split('.');
                if (property === 'MemberData.MemAvgCost' || property === 'MemberData.Percent' || property === 'MemberData.MemShares' || property === 'MemberData.Profit' || property === 'MemberData.CostValue') {
                    list = data.sort((a, b) => {
                        if(b[mainKey].Percent == "-" ) {
                            return -1
                        }
                        if(b[mainKey] != null && a[mainKey] != null ) {
                            return b[mainKey][subKey] - a[mainKey][subKey]
                        }
                    })
                } else if (property === 'DirectorSelection.AvgCost' || property === 'DirectorSelection.Percent' ) {
                    list = data.sort((a, b) => {
                        if(b[mainKey].Percent == "-" ) {
                            return -1
                        }
                        if(b[mainKey] != null && a[mainKey] != null ) {
                            return b[mainKey][subKey] - a[mainKey][subKey]
                        }
                    });
                } else {
                    list = data.sort((a, b) => {
                        if(b[property] != null && a[property] != null ) {
                            return b[property] - a[property]
                        }
                    });
                }
            }
            if (order === 'desc') {
              list.reverse()
            }
            
            return list
        },
        
    },
    mounted(){
        setTimeout(() => {
            this.calcTableHeight(100);
            window.addEventListener('resize', () => {
                this.calcTableHeight(100);
            });
        }, 100)
    }
}
</script>

