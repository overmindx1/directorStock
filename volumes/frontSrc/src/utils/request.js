import axios from 'axios'
// import { MessageBox, Message } from 'element-ui'
import { MessageBox, Message , Loading } from 'element-ui'

// create an axios instance
const service = axios.create({
    baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
    // withCredentials: true, // send cookies when cross-domain requests
    timeout: 5000 // request timeout
})

// request interceptor
service.interceptors.request.use(
    config => {
        // do something before request is sent
        Loading.service({
            text: 'Loading',
            spinner: 'el-icon-loading',
            background: 'rgba(0, 0, 0, 0.7)',
            fullscreen : true
        })
        return config
    },
    error => {
        // do something with request error
        console.log(error) // for debug
        return Promise.reject(error)
    }
)

// response interceptor
service.interceptors.response.use(
    /**
     * If you want to get http information such as headers or status
     * Please return  response => response
    */

    /**
     * Determine the request status by custom code
     * Here is just an example
     * You can also judge the status by HTTP Status Code
     */
    response => {
        Loading.service().close()
        const res = response.data

        // if the custom code is not 20000, it is judged as an error.
        if (res.code !== 0) {
            Message({
                message: res.msg || 'Error',
                type: 'error',
                duration: 5 * 1000
            }) 
            //message.error(res.msg || 'Error', 5)
            return Promise.reject(new Error(res.msg || 'Error'))
        } else {
            return res.data
        }
    },
    error => {
        Loading.service().close()
        console.log(error.response) // for debug
        Message({
            message: error.response.data.msg,
            type: 'error',
            duration: 5 * 1000
        }) 
        return Promise.reject(error)
    }
)

export default service