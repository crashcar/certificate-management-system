import request from '@/utils/request'

// 用户登录
export function userLogin(data) {
    return request({
        url: '/login',
        method: 'post',
        data
    })
}


//用户注册
export function userRegister(data) {
    return request({
        url: '/register',
        method: 'post',
        data
    })
}


//管理员登录
export function adminLogin(data) {
    return request({
        url: '/adminLogin',
        method: 'post',
        data
    })
}

//管理员注册
export function adminRegister(data) {
    return request({
        url: '/adminRegister',
        method: 'post',
        data
    })
}





