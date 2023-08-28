/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package dtos

type WsCode uint32

const (
    WsCodeDeviceLibraryUpgrade   WsCode = 10001 // 驱动下载/升级
    WsCodeDeviceServiceRunStatus WsCode = 10002 // 驱动重启
    WsCodeDeviceLibraryDelete    WsCode = 10003 // 驱动删除
    WsCodeDeviceServiceLog       WsCode = 10004 // 驱动日志
    
    WsCodeCloudServiceDownload  WsCode = 20001 // 云服务下载
    WsCodeCloudServiceRunStatus WsCode = 20002 // 云服务重启/停止
    WsCodeCloudServiceRunDelete WsCode = 20003 // 云服务删除
    WsCodeCloudServiceRunLog    WsCode = 20004 // 云服务日志
    
    WsCodeCheckLang WsCode = 30001 // 切换语言
    
    // 云端网络情况
    WsCodeCloudState WsCode = 10007 // 云端网络情况
    //OTA
    WsCodeOTAUpgradeProgress WsCode = 10100 // OTA升级进度
    WsCodeOTAFirmwareUpgrade WsCode = 10101 // OTA升级
    // 严重警告
    WsCodeSeriousAlert WsCode = 10200
)
