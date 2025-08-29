/*
 * @Author: LPY
 * @Date: 2025-05-29 18:41:20
 * @LastEditors: LPY
 * @LastEditTime: 2025-07-01 16:17:55
 * @FilePath: /kvm-cloud-frontend/src/hooks/useLanguage.ts
 * @Description: 语言hook
 */

import i18n from '@/lang'
import { languageLabelMap } from '@/models/setting'
import { computed } from 'vue'
import { LocalStorageKeys, useLocalStorage } from './useLocalStorage'
import { Languages } from '@gl/main'
import { isForeignEnv } from '@/utils'

/** 语言hook */
export default function useLanguage () {
    const {getValue, setValue} = useLocalStorage(LocalStorageKeys.STORAGE_LANGUAGE_KEY, isForeignEnv.value ? Languages.EN : Languages.ZH)

    const currentLang = computed(() => i18n.global.locale.value as Languages)

    const setLanguage = (lang = getValue()) => {
        i18n.global.locale.value = lang
        setValue(lang)
    }

    const isZh = computed(() => currentLang.value === Languages.ZH)

    const currentLangLabel = computed(() => languageLabelMap.get(currentLang.value))

    return {currentLang, setLanguage, isZh, t: i18n.global.t, currentLangLabel}
}

export const t = i18n.global.t