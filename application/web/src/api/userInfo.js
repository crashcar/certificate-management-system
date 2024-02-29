import request from '@/utils/request'

// 用户登录
export function userLogin(data) {
    return request({
        url: '/userLogin',
        method: 'post',
        data
    })
}


//用户注册
export function userRegister(data) {
    return request({
        url: '/userRegister',
        method: 'post',
        data
    })
}





