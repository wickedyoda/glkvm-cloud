<!--
 * @Author: shufei.han
 * @Date: 2025-06-11 12:04:48
 * @LastEditors: LPY
 * @LastEditTime: 2025-08-26 15:40:52
 * @FilePath: \glkvm-cloud\web-ui\src\views\device\components\deviceListView.vue
 * @Description: 
-->
<template>
    <BaseLoadingContainer :spinning="deviceStore.state.getDeviceLoading">
        <div class="device-list-view full-height">
            <div class="header">
                <a-input-search
                    :placeholder="$t('device.searchTip')"
                    style="width: 316px"
                    @search="deviceStore.handleSearch"
                />

                <div>
                    <BaseButton size="middle" style="margin-right: 12px;" @click="refresh">{{ $t('common.refresh') }}</BaseButton>
                    <BaseButton size="middle" @click="executeCommand">{{ $t('device.executeCommand') }}</BaseButton>
                </div>
            </div>
            <div class="content">
                <BaseTable 
                    :data-source="deviceStore.deviceList"
                    :columns="deviceColumns"
                    rowKey="id"
                    :rowSelection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }"
                >
                    <template #connected="{ record }">
                        {{ 
                            record.connected ? 
                                (dayjs(record.connected).format('D') + $t('common.d') + ' ' +
                                    dayjs(record.connected).format('H') + $t('common.h') + ' ' +
                                    dayjs(record.connected).format('m') + $t('common.m') + ' ' +
                                    dayjs(record.connected).format('s') + $t('common.s')) : ''
                        }}
                    </template>
                    <template #uptime="{ record }">
                        {{ 
                            record.uptime ?
                                (dayjs(record.uptime).format('D') + $t('common.d') + ' ' +
                                    dayjs(record.uptime).format('H') + $t('common.h') + ' ' +
                                    dayjs(record.uptime).format('m') + $t('common.m') + ' ' +
                                    dayjs(record.uptime).format('s') + $t('common.s')) : ''
                        }}
                    </template>
                    <template #action="{ record }">
                        <div class="flex-start">
                            <a target="_blank" rel="noopener noreferrer" @click="handleRemoteSSH(record.id)">{{ $t('device.remoteSSH') }}</a>
                            <a target="_blank" rel="noopener noreferrer" style="margin-left: 16px;" @click="handleRemoteControl(record.id)">
                                {{ $t('device.remoteControl') }}
                            </a>
                        </div>
                    </template>
                </BaseTable>
            </div>
            <div class="pagination flex-end items-end">
                <BasePagination 
                    :total="deviceStore.pageLink.total"
                    :pageSize="deviceStore.pageLink.size"
                    v-model:current="page" />
            </div>
        </div>

        <!-- 批量编辑弹窗 -->
        <ExecuteCommandDialog 
            v-model:open="executeCommandOpen"
            :selection="state.selectedRows"
            @handleApply="executeCommandApply"
        />

        <!-- 执行命令结果弹窗 -->
        <CommandResponseDialog 
            v-if="commandResponseOpen"
            v-model:open="commandResponseOpen"
            :selection="state.selectedRows"
            :formData="executeCommandFormData"
        />
    </BaseLoadingContainer>
</template> 

<script setup lang="ts">
import BaseLoadingContainer from '@/components/base/baseLoadingContainer.vue'
import BasePagination from '@/components/base/basePagination.vue'
import BaseTable from '@/components/base/baseTable.vue'
import { t } from '@/hooks/useLanguage'
import { useDeviceStore } from '@/stores/modules/device'
import { message, type TableColumnType } from 'ant-design-vue'
import { computed, reactive, ref } from 'vue'
import ExecuteCommandDialog from './executeCommandDialog.vue'
import dayjs from 'dayjs'
import { DeviceInfo, ExecuteCommandFormData } from '@/models/device'
import CommandResponseDialog from './commandResponseDialog.vue'

const deviceStore  = useDeviceStore()

const deviceColumns = computed<TableColumnType[]>(() => { 
    return [
        {title: t('device.deviceID'), dataIndex: 'id', ellipsis: true},
        {title: t('device.connectedTime'), dataIndex: 'connected', ellipsis: true},
        {title: t('device.uptime'), dataIndex: 'uptime', ellipsis: true},
        {title: t('device.IPAddress'), dataIndex: 'ipaddr', ellipsis: true},
        {title: t('device.description'), dataIndex: 'description', ellipsis: true},
        {title: t('common.action'), dataIndex: 'action', width: 220},
    ]
})

const page = computed({
    get: () => deviceStore.pageLink.page,
    set: (val) => {
        deviceStore.pageLink.changePage(val)
    },
})

type Key = string | number;
const state = reactive<{
  selectedRowKeys: Key[];
  selectedRows: DeviceInfo[];
}>({
    selectedRowKeys: [], // Check here to configure the default column
    selectedRows: [],
})

const onSelectChange = (selectedRowKeys: Key[], selectedRows: DeviceInfo[]) => {
    console.log('selectedRowKeys changed: ', selectedRowKeys)
    state.selectedRowKeys = selectedRowKeys
    state.selectedRows = selectedRows
}

/** 刷新列表 */
const refresh = () => {
    deviceStore.getDeviceList()
    message.success(t('device.refreshSuccess'))
}

/** 批量配置 */
const executeCommandOpen = ref(false)
const executeCommand = () => {
    // 判断是否选择了设备
    if (!state.selectedRowKeys.length) {
        message.error(t('device.selectDeviceTips'))
        return
    }
    executeCommandOpen.value = true
}

/** 批量配置结果弹窗 */
const commandResponseOpen = ref(false)

const executeCommandFormData = ref<ExecuteCommandFormData>()
/** 执行命令弹窗提交 */
const executeCommandApply = (formData: ExecuteCommandFormData) => {
    executeCommandFormData.value = formData
    executeCommandOpen.value = false
    commandResponseOpen.value = true
}

/** 远程SSH */
const handleRemoteSSH = async (id: string) => {
    try {
        let url = `/#/rtty/${id}`
        window.open(url)
    } catch (error) {
        console.log(error)
    }
}

/** 远程控制 */
const handleRemoteControl = async (id: string) => {
    try {
        let proto = 'https'
        let ipaddr = '127.0.0.1'
        let port = 443
        let path = '/'
        const addr = encodeURIComponent(`${ipaddr}:${port}${path}`)
        window.open(`/web/${id}/${proto}/${addr}`)
    } catch (error) {
        console.log(error)
    }
}
</script> 

<style lang="scss" scoped>
.device-list-view {
    height: 100%;
    .header {
        height: 36px;
        margin-bottom: 16px;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    .content {
        overflow-x: hidden;
        overflow-y: auto;
        height: calc(100% - 92px);
    }
    .pagination {
        height: 40px;
    }
}
</style>