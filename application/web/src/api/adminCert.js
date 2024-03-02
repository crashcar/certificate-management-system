import request from '@/utils/request'

// 机构所有已认证证书
export function query_institute_certificates(data) {
    return request({
        url: '/queryCertOrg',
        method: 'post',
        data
    })
}

// 机构所有带审查证书
export function query_institute_showCertList(data) {
    return request({
        url: '/showCertList',
        method: 'post',
        data
    })
}