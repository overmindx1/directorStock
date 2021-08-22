<template>
    <el-dialog
        title="選擇使用者"
        :visible="userConfigTaigger"
        :width="widthSet"
        :show-close="false"
    >
        <el-row>
            <el-col style="margin-bottom:1rem;text-align:left">
                <el-button size="mini" type="primary" @click="addUserTrigger = true">新增使用者</el-button>
            </el-col>
            <el-col>
                <vxe-table :data="userList" stripe border show-overflow >
                    <vxe-column title="項次" width="50">
                        <template #default="{   rowIndex }">
                            {{rowIndex + 1}}
                        </template>
                    </vxe-column>
                    <vxe-column field="nickname" title="使用者" width="100"></vxe-column>
                    <vxe-column title="操作" >
                        <template #default="{  row }">
                            <el-button-group>
                                <el-button size="mini" type="primary" @click="changeCurrentUser(row)" :disabled="row.nickname == currentUserData.nickname">更換</el-button>
                                <el-button size="mini" type="danger" @click="deleteUser(row)" v-if="row.nickname != currentUserData.nickname">刪除</el-button>                        
                            </el-button-group>
                        </template>
                    </vxe-column>
                </vxe-table>
            </el-col>
            
        </el-row>
        <el-row slot="footer" class="dialog-footer">
            <el-col >
                <el-button @click="closeDialog">取 消</el-button>
            </el-col>            
        </el-row>

        <el-dialog        
            title="新增使用者"
            :visible.sync="addUserTrigger"
            :width="widthSet"
            append-to-body>
            <el-row>
                <el-col>
                    <el-form :inline="true"  class="demo-form-inline">
                        <el-form-item label="使用者名稱">
                            <el-input v-model="addUserNickname" placeholder="使用者名稱"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button size="mini" type="primary" @click="insertNewUser" >新增</el-button>
                            <el-button size="mini"  @click="addUserTrigger = false" >取消</el-button>
                        </el-form-item>
                    </el-form>
                </el-col>
            </el-row>
        </el-dialog>
    </el-dialog>
</template>

<script>

export default {
    name : "tableEnableColumes",
    props : {
        dbInstance : {
            type : Object,
            required : true
        },
        currentUserData : {
            type : Object,
            required : true
        },
        userConfigTaigger : {
            type : Boolean,
            required : true
        },
        
    },
    watch : {
        triggerColumeSetting(newVal , oldVal) {
            if(newVal) {
                this.tableColumesTemp = JSON.parse(JSON.stringify(this.tableColumes))
            }
        }
    },
    data : () => ({
        userList : [],
        widthSet : "30%",
        addUserTrigger : false,
        addUserNickname : "",
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
        // 新增使用者
        async insertNewUser(){
            if(this.addUserNickname != "") {
                let execute = await this.dbInstance.addNewUser(this.addUserNickname);
                if(execute) {
                    this.$message(`使用者 資料新增!`);
                    this.userList = await this.dbInstance.findAllUsers();
                    this.addUserTrigger = false;
                    
                } else {
                    this.$message({
                        message: `使用者 資料新增失敗!`,
                        type: 'error'
                    });
                }
            } else {
                this.$message({
                    message: `使用者 名稱 必須輸入!`,
                    type: 'error'
                });
            }
        },
        // 刪除使用者
        async deleteUser(row){
            this.$confirm('此操作將永久删除該使用者, 是否繼續?', '提示', {
                confirmButtonText: '確定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(async () => {
                let execute = await this.dbInstance.deleteUser(row.nickname);
                if(execute) {
                    this.$message(`使用者 ${row.nickname} 資料刪除成功!`);
                    this.userList = await this.dbInstance.findAllUsers();
                } else {
                    this.$message({
                        message: `使用者 ${row.nickname} 資料刪除失敗!`,
                        type: 'error'
                    });
                }
            })
        },
        // 更換當前使用者
        changeCurrentUser(row){
            this.$emit("changeCurrentUser" , row) 
        },
        //關閉Dialog
        closeDialog() {
            this.$emit('closeDialog')
        },
        
    },
    async mounted(){
        window.onresize = this.getCurrentWidth;
        this.getCurrentWidth();

        this.userList = await this.dbInstance.findAllUsers();
    }
}
</script>

