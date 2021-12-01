<template>
    <el-dialog
        title="欄位選擇"
        :visible="triggerColumeSetting"
        :width="widthSet"
        :show-close="false"
    >
        <el-row v-for="(item , key) in tableColumesTemp" :key="key" style="margin : .5rem 0;font-size:.8rem">
            <el-col :span="5" style="line-height:2rem">
                {{item.label}}:  
            </el-col>
            <el-col :span="7">
                <el-select v-model="item.show"  placeholder="请选择" size="small">
                    <el-option label="顯示" :value="true">  </el-option>
                    <el-option label="不顯示" :value="false">  </el-option>
                </el-select>
            </el-col>
            <el-col :span="6" style="line-height:2rem">顯示順序: </el-col>
            <el-col :span="6">
                <el-select v-model="item.index" filterable placeholder="请选择" size="small">
                    <el-option v-for="n in 9" :key="n" :label="n" :value="n">  </el-option>
                </el-select>
            </el-col>
        </el-row>
        <el-row slot="footer" class="dialog-footer">
            <el-col :xs="8" :sm="8" :md="4" :lg="4" :xl="4" style="text-align:left">
                <el-button @click="resetColumeSetting" type="warning">清 除</el-button>
            </el-col>
            <el-col :xs="16" :sm="16" :md="20" :lg="20" :xl="20" style="text-align:right">
                <el-button @click="closeDialog">取 消</el-button>
                <el-button type="primary" @click="saveColumesShow">確 定</el-button>
            </el-col>            
        </el-row>
    </el-dialog>
</template>

<script>

export default {
    name : "tableEnableColumes",
    props : {
        triggerColumeSetting : {
            type : Boolean,
            required : true
        },
        tableColumes : {
            type : Object,
            required : true
        }
    },
    watch : {
        triggerColumeSetting(newVal , oldVal) {
            if(newVal) {
                this.tableColumesTemp = JSON.parse(JSON.stringify(this.tableColumes))
            }
        }
    },
    data : () => ({
        tableColumesTemp : {},
        widthSet : "30%"
    }),
    methods : {
        getCurrentWidth(e = null) {
            if(window.innerWidth < 460) {
                this.widthSet = "90%";
            } else if (window.innerWidth < 769) {
                this.widthSet = "70%";
            } else if (window.innerWidth < 1300) {
                this.widthSet = "45%";
            } else {
                this.widthSet = "30%";
            }
        },
        //關閉Dialog
        closeDialog() {
            this.$emit('closeDialog')
        },
        // 計入要顯示的表格欄位
        saveColumesShow(){
            // localStorage.setItem('tableColumes' , JSON.stringify(this.tableColumesTemp));
            this.$emit('saveColumesShow' , this.tableColumesTemp);
            this.closeDialog();
        },
        // 刪除儲存資料
        resetColumeSetting(){
            // localStorage.removeItem('tableColumes')
            this.$emit('resetColumeSetting' , this.tableColumesTemp);
        }
    },
    mounted(){
        window.onresize = this.getCurrentWidth;
        this.getCurrentWidth();
    }
}
</script>

