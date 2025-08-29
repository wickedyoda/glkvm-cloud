/*
 * @Author: LPY
 * @Date: 2025-06-03 12:21:21
 * @LastEditors: LPY
 * @LastEditTime: 2025-08-26 09:06:31
 * @FilePath: \glkvm-cloud\web-ui\src\api\user.ts
 * @Description: 用户相关请求api
 */

import request from './request'
import type { LoginParams } from '@/models/user'

/** 登录 */
export function reqLogin (data: LoginParams) {
    return request<void>({
        url: '/signin',
        method: 'POST',
        data,
    })
}

/** 退出登录 */
export function reqLogout () {
    return request<void>({
        url: '/signout',
    })
}

/** 判断登录状态 */
export function reqCheckLoginStatus () {
    return request<void>({
        url: '/alive',
    })
}