import request from '@/utils/request'

export function prove_certificate(data) {
    return request({
        url: '/queryCertByFullInfoSys',
        method: 'post',
        data
    })
}
