/*
 * @Author: LPY
 * @Date: 2025-05-29 18:43:46
 * @LastEditors: LPY
 * @LastEditTime: 2025-08-26 09:24:09
 * @FilePath: \glkvm-cloud\web-ui\src\stores\modules\user.ts
 * @Description: 用户相关状态存储
 */

import { ref } from 'vue'
import { defineStore } from 'pinia'
import { LocalStorageKeys, useLocalStorage } from '@/hooks/useLocalStorage'
import { reqCheckLoginStatus, reqLogin, reqLogout} from '@/api/user'
import router from '@/router'
import { type LoginParams } from '@/models/user'
import { useAppStore } from './app'
import { getCookieToken, removeCookieToken, setCookieToken } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
    // token
    const token = ref(getCookieToken())

    // 登录状态
    const loginStatus = ref(true)

    /** 设置token */
    const setToken = (newToken: string) => {
        token.value = newToken
        setCookieToken(newToken)
    }

    /** 清除token */
    const clearToken = () => {
        token.value = ''
        removeCookieToken()
    }

    /** 登录 */
    const login = async (credentials: LoginParams) => {
        // 加密密码
        const params = {
            password: credentials.password,
        }
        const data = await reqLogin(params)
        console.log(data.info)
        loginStatus.value = true
        
        // setToken(newToken)
    }

    /** 合并登出方法，不可导出使用 */
    const logout = () => {
        // 清除token
        clearToken()
        
        // 清除侧边栏状态
        useAppStore().resetManualSetting()
        useLocalStorage(LocalStorageKeys.SIDEBAR_MANUAL_CONTROL_KEY).removeValue()

        loginStatus.value = false
        // 登出
        router.push({ path: '/login' })
    }
  
    /** 自动登出 */
    const autoLogout = () => {
        logout()
    }

    /** 手动登出 */
    const manualLogout = async () => {
        try {
            await reqLogout()
            logout()
        } catch (error) {
            console.log(error)
        }
    }

    /** 判断登录状态 */
    const checkLoginStatus = async () => {
        try {
            await reqCheckLoginStatus()
            loginStatus.value = true
        } catch (error) {
            autoLogout()
        }
    }

    checkLoginStatus()

    return { 
        token,
        loginStatus,
        setToken,
        clearToken,
        login,
        autoLogout,
        manualLogout,
    }
})