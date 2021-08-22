<template>
    <el-drawer        
        :visible.sync="drawerTrigger"
        direction="ltr"
        :size="drawerSize"
        :before-close="closedDrawer"
    >
        <el-row type="flex" justify="center">
            <el-col :span="20" class="currentUserName">
                當前使用者 :  {{currentUserData.nickname}}
            </el-col>
        </el-row>

        <el-row >
            <el-col :offset="2"  class="buttonGap">
                <el-button type="primary" icon="el-icon-edit" @click="saveUserData" 
                    style="width:80%">儲存設定</el-button>
            </el-col>
            <el-col :offset="2"  class="buttonGap">
                <el-button type="primary" icon="el-icon-user" @click="openFromDrawer('userConfigTaigger')" 
                    style="width:80%">更換使用者</el-button>
            </el-col>
            <el-col :offset="2" class="buttonGap">
                <el-button type="primary" icon="el-icon-set-up" @click="openFromDrawer('triggerColumeSetting')" 
                    style="width:80%">調整表格欄位</el-button>
            </el-col>
            <el-col :offset="2"  class="buttonGap">
                <el-button type="primary" icon="el-icon-tickets" @click="openFromDrawer('staticTrigger')" 
                    style="width:80%">顯示/關閉統計資料</el-button>
            </el-col>
            
        </el-row>
    </el-drawer>
</template>

<script>
export default {
    name : 'drawerMenu',
    props : {
        drawerTrigger : {
            type : Boolean,
            required : true
        },
        currentUserData : {
            type : Object,
            required : true
        },
    },
    data : () => ({
        drawerSize : "30%",
    }),
    methods : {
        closedDrawer(){
            this.$emit("closedDrawer");
        },
        saveUserData(){
            this.$emit("saveUserData");
        },
        openFromDrawer(trigger){
            this.$emit("openFromDrawer" , trigger)
            this.closedDrawer();
        },
        // ResizeEventHandle
        resizeHandle(){
            let width = document.body.getBoundingClientRect().width;
            console.log(width)
            if(width > 992 && width < 1279) {
                this.drawerSize = "40%"
            } else if (width > 450 && width <= 992) {
                this.drawerSize = "60%"
            } else if (width <= 450) {
                this.drawerSize = "80%"
            } else {
                this.drawerSize = "30%"
            }
        },
    },
    mounted(){
        window.addEventListener('resize', this.resizeHandle);
        this.resizeHandle();
    }
}
</script>

<style lang="scss" scoped>
.currentUserName {
    border-bottom: 1px solid #bdc3c7;
    padding: 1rem 0;
    text-align: left;
    font-size: 1.25rem;
}
.buttonGap {
    margin : .5rem 0 ;
}
</style>