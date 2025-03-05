import { callAPI } from "@/lib/api";

const getScanConfigs = async (access_token) => {
    return await callAPI("/api/v1/scan-configs", {
        method: 'GET',
        expected_status: 200,
        access_token
    })
}

const getProjectScans = async (projectId: string, access_token) => {
    return await callAPI(`/api/v1/projects/${projectId}/scans`, {
        method: 'GET',
        expected_status: 200,
        access_token
    })
}

const startNewScan = async (projectId: string, scanConfigId: string, targetIds: string[], access_token) => {
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

export {getScanConfigs, startNewScan, getProjectScans}