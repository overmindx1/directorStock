import Dexie from 'dexie';

class indexedDbWarp {

    // 建立 IndexedDb
    constructor(){
        //Def IndexedDB
        this.db = new Dexie("directorStock");
        this.db.version(1).stores({
            user: "++id , nickname",
            userData: "++id , nickname , googleSpreadsheet , tableColumes",
            defaultUser: "id , nickname",
            mssData : "++id , nickname , responseData"
        });
        this.db.open();
        (async () => {
            try {
                let count = await this.db.defaultUser.count(); 
                if(count == 0) {
                    let id = await this.db.user.add({nickname : "預設使用者"});
                    await this.db.defaultUser.add({id :1 ,nickname : "預設使用者"})
                    if(id) {
                        this.db.userData.add({nickname : "預設使用者" , googleSpreadsheet : "" ,tableColumes : {} });
                        this.db.mssData.add({nickname : "預設使用者" , responseData : {} })
                    }
                }            
            } catch (error) {
                console.error([ "indexedDB Class constructor error ", error])
            }
        })();
    }

    // 取得所有使用者資料
    async findAllUsers(){
        try {
            return await this.db.user.toArray();
        } catch (error) {
            console.error([ "indexedDB Class findAllUsers error ", error]);
            return [];
        }
    }

    // 設定當前使用者
    async setDefaultUser(name = "預設使用者"){
        try {
            let user = await this.db.user.get({nickname : name});
            if(user === undefined) {
               return false;
            }
            await this.db.defaultUser.put({id : 1 , nickname : name});
            return this.getDefaultUserData(name);
        } catch (error) {
            console.error([ "indexedDB Class setDefaultUser error ", error]);
            return false;
        }
    }

    // 取得當前使用者資料
    async getDefaultUserData(name = "預設使用者") {
        try {
            let defaultUser = await this.db.defaultUser.get(1);
            if(defaultUser === undefined) {
                return false;
            }
            // console.log(defaultUser)
            let userData = await this.db.userData.get({nickname : defaultUser.nickname});
            return userData;
        } catch (error) {
            console.error([ "indexedDB Class getDefaultUserData error ", error]);
            return false;
        }
    }

    // 更新使用者資料
    async updateUserData(name , userData) {
        try {
            // 找看看有沒有使用者資料
            let find = await this.db.userData.get({
                nickname : name
            });
            // 找到資料
            if(find) {
                await this.db.userData.put(userData)
                return true;
            }
            return false
        } catch (error) {
            console.error([ "indexedDB Class updateUserData error ", error]);
            return false;
        }
    }

    // 新增使用者
    async addNewUser(name) {
        try {
            let id = await this.db.user.add({nickname : name});
            if(id) {
                let execute = await this.db.userData.add({nickname : name , googleSpreadsheet : "" ,tableColumes : {}})
                this.db.mssData.add({nickname : name , responseData : {} })
                if(execute) {
                    return true
                }
            }
            return false;            
        } catch (error) {
            console.error([ "indexedDB Class addNewUser error ", error]);
            return false;
        }
    }

    // 刪除使用者
    async deleteUser(name){
        try {
            let user = await this.db.user.get({nickname : name});
            if(user === undefined) {
                return false;
            }
            await this.db.userData.delete(user.id);
            await this.db.user.delete(user.id);
            await this.db.mssData.delete(user.id);
            return true;
        } catch (error) {
            console.error([ "indexedDB Class deleteUser error ", error]);
            return false;
        }
        
    }
}

let indexedDbInstance = new indexedDbWarp();

export default indexedDbInstance ;