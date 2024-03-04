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
        url: '/showApplicationList',
        method: 'post',
        data
    })
}


//审查通过单个个列表
export function check_pass_single_Cert(data) {
    return request({
        url: '/approveCert',
        method: 'post',
        data
    })
}

//审查拒绝单个个列表
export function check_reject_single_Cert(data) {
    return request({
        url: '/denialCert',
        method: 'post',
        data,
    })
}