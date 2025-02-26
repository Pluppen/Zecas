import { API_URL } from "@/config";

const getScanConfigs = async () => {
    const response = await fetch(`${API_URL}/api/v1/scan-configs`)
    if (response.status != 200) {
        console.error(response.statusText);
        return {error: "Something went wrong"}
    }
    return response.json()
}

const getProjectScans = async (projectId: string) => {
    const response = await fetch(`${API_URL}/api/v1/projects/${projectId}/scans`)
    if (response.status != 200) {
        return {error: "Something went wrong fetching scans for project"};
    }
    return response.json();
}

const startNewScan = async (projectId: string, scanConfigId: string, targetIds: string[]) => {
    const body = {
        project_id: projectId,
        scan_config_id: scanConfigId,
        target_ids: targetIds
    }
    const response = await fetch(`${API_URL}/api/v1/scans`, {
        method: 'POST',
        body: JSON.stringify(body)
    })
    if (response.status != 202) {
        console.error(response.statusText);
        return {error: "Something went wrong"}
    }
    return response.json()
}

export {getScanConfigs, startNewScan, getProjectScans}