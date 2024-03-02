import request from '@/utils/request'

// 获取当前机构证书类型
export function getReviewTypes() {
    return request({
        url: '/reviewTypes',
        method: 'get'
    })
}