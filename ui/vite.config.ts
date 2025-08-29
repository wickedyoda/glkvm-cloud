import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import removeConsole from 'vite-plugin-remove-console'

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, `${process.cwd()}/env`)

    return {
        plugins: [
            vue(),
            env.VITE_APP_USE_SCRIPT === 'prod' ? removeConsole({
                external: ['src/assets/iconfont/iconfont.js'],  // 排除 iconfont.js 文件
            }) : null,
        ],
        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url)),
            },
        },
        server: {
            host: '0.0.0.0',
            port: 3011,
            proxy: {
                '/devs': {
                    target: 'https://49.7.174.146:1443',
                    secure: false,
                },
                '/signin': {
                    target: 'https://49.7.174.146:1443/',
                    secure: false,
                },
                '/signout': {
                    target: 'https://49.7.174.146:1443',
                    secure: false,
                },
                '/alive': {
                    target: 'https://49.7.174.146:1443',
                    secure: false,
                },
                '/get': {
                    target: 'https://49.7.174.146:1443',
                    secure: false,
                },
                '^/cmd/.*': {
                    target: 'https://49.7.174.146:1443',
                    secure: false,
                },
                '^/connect/.*': {
                    ws: true,
                    target: 'https://49.7.174.146:1443',
                },
                '^/web/*': {
                    target: 'https://49.7.174.146:1443',
                },
            },
        },
        // 关闭 css 警告
        css: {
            preprocessorOptions: {
                scss: {
                    api: 'modern-compiler',
                },
            },
        },
    }
})
