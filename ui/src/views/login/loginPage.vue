<!--
 * @Author: LPY
 * @Date: 2025-05-30 10:48:43
 * @LastEditors: LPY
 * @LastEditTime: 2025-08-26 17:11:12
 * @FilePath: \glkvm-cloud\web-ui\src\views\login\loginPage.vue
 * @Description: 登录页面
-->
<template>
    <BaseWhitePage>
        <LoginBox>
            <BaseText type="large-title-m" style="margin: 24px 0 16px;">
                {{ $t('login.authorizationRequired') }}
            </BaseText>
            <AForm 
                class="dense-form" 
                ref="formRef"
                :model="state.formModel"
                :rules="formRules"
                :validateTrigger="['blur', 'change']"
                style="width: 100%;"
                @validate="handleValidate"
            >
                <AFormItem name="password">
                    <GlPassword
                        name="password"
                        :useDefaultValidateRule="false"
                        @pressEnter="handleLogin"
                        v-model:value="state.formModel.password"
                        :placeholder="$t('login.password')"
                    />
                </AFormItem>
            </AForm>
            <BaseButton medium type="primary" style="width: 100%;margin-top: 16px;" :loading="state.loading" @click="handleLogin">
                {{ $t('login.signIn') }}
            </BaseButton>
        </LoginBox>
    </BaseWhitePage>
</template>

<script setup lang="ts">
import BaseWhitePage from '@/components/base/baseWhitePage.vue'
import LoginBox from './components/loginBox.vue'
import { computed, reactive, ref } from 'vue'
import { t } from '@/hooks/useLanguage'
import { useUserStore } from '@/stores/modules/user'
import { useValidateInfo, type FormRules } from 'gl-web-main'
import { GlPassword } from 'gl-web-main/components'
import { useRouter } from 'vue-router'
import { LoginParams } from '@/models/user'
import { message } from 'ant-design-vue'

const router = useRouter()

const { handleValidate } = useValidateInfo<LoginParams>()

const formRef = ref(null)

const state = reactive<{formModel: LoginParams, loading: boolean}>({
    formModel: {
        password: '',
    },
    loading: false,
})

const formRules = computed<FormRules<LoginParams>>(() => {
    return {
        password: [{ required: true, message: 'login.enterPwdTip'}],
    }
})

// 登录按钮
const handleLogin = () => {
    formRef.value.validate().then(async () => {
        state.loading = true
        try {
            await useUserStore().login(state.formModel)
            // 登录成功后跳转到首页或之前尝试访问的页面
            const redirect = router.currentRoute.value.query.redirect as string || '/' 
            console.log(redirect)
            
            state.loading = false
            router.push(redirect)
        } catch (error) {
            console.log(error)
            state.loading = false
            message.error(t('login.incorrectPwd'))
        }
    })
}
</script>

<style scoped lang="scss">
</style>