import request from '@/utils/request'

// 用户登录
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
