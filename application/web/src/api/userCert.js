import request from '@/utils/request'

export function query_user_certificates(data) {
    return request({
        url: '/queryCertByUserOrg',
        method: 'post',
        data
    })
}

// 下载证书
export function user_download_certificate(userId, certificateId) {
    return request({
        url: '/userDownloadCertificate',
        method: 'post',
        params: {
            userId,
            certificateId
        },
        responseType: 'blob' // 指定响应数据的类型为 Blob
    })
}

//用户查询在当前机构所有的已审核通过的证书
export function user_apply_certificate(userID, realName, certType, file) {
    const formData = new FormData();
    formData.append('userID', userID);
    formData.append("realName", realName);
    formData.append('certType', certType);
    formData.append('file', file);

    return request({
        url: '/upload',
        method: 'post',
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
}


//获取用户待审核列表
export function user_query_checking_certificates(data) {
    return request({
        url: '/showApplicationList',
        method: 'post',
        data
    })
}

//单条申请详情
export function user_single_cert_certificates(data) {
    return request({
        url: '/applicationDetail',
        method: 'post',
        data
    })
}




//用户撤销审核或者删除已通过的审核记录
export function user_revoke_or_delete_apply_certificates(data) {
    return request({
        url: '/deleteRecord',
        method: 'post',
        data
    })
}