import { callAPI } from "@/lib/api";

export type ScanConfig = {
    id: string,
    name: string,
    scanner_type: string,
    parameters: JSON,
    active: boolean
}

export const getScanConfigs = async (access_token) => {
    return await callAPI("/api/v1/scan-configs", {
        method: 'GET',
        expected_status: 200,
        access_token
    })
}

export const createScanConfig = async (scanConfig: ScanConfig, access_token: string) => {
    const body = {
        name: scanConfig.name,
        scanner_type: scanConfig.scanner_type,
        parameters: scanConfig.parameters,
        active: scanConfig.active
    }    

    return await callAPI(`/api/v1/scan-configs`, {
        method: "POST",
        access_token,
        body: JSON.stringify(body),
        expected_status: 201
    })
}

export const getScanConfigById = async (scanConfigId: string, access_token: string) => {
    return await callAPI(`/api/v1/scan-configs/${scanConfigId}`, {
        method: "GET",
        access_token,
        expected_status: 200
    })
}

export const updateScanConfigById = async (scanConfig: ScanConfig, access_token: string) => {
    const body = {
        name: scanConfig.name,
        scanner_type: scanConfig.scanner_type,
        parameters: scanConfig.parameters,
        active: scanConfig.active
    }    

    return await callAPI(`/api/v1/scan-configs/${scanConfig.id}`, {
        method: "PUT",
        access_token,
        body: JSON.stringify(body),
        expected_status: 200
    })
}

export const deleteScanConfigById = async (scanId: string, access_token: string) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}`, {
        method: "DELETE",
        access_token,
        expected_status: 200
    })
}

// SCANS

export const startNewScan = async (projectId: string, scanConfigId: string, targetIds: string[], access_token) => {
    const body = {
        project_id: projectId,
        scan_config_id: scanConfigId,
        target_ids: targetIds
    }
    return await callAPI("/api/v1/scans", {
        method: 'POST',
        expected_status: 202,
        body: JSON.stringify(body),
        access_token
    })
}

export const getScanById = async (scanId: string, access_token: string) => {
    return await callAPI(`/api/v1/scans/${scanId}`, {
        method: "GET",
        access_token,
        expected_status: 200
    })
}

export const getScans = async (access_token: string) => {
    return await callAPI(`/api/v1/scan/`, {
        method: "GET",
        access_token,
        expected_status: 200
    })
}

export const cancelScanById = async (scanId: string, access_token: string) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}/cancel`, {
        method: "POST",
        access_token,
        expected_status: 200
    })
}

export const getFindingsByScanId = async (scanId: string, access_token: string) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}/findings`, {
        method: "POST",
        access_token,
        expected_status: 200
    })
}