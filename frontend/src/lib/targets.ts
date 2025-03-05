import { callAPI } from "@/lib/api";

const getProjectTargets = async (projectId: string, access_token) => {
    return await callAPI(`/api/v1/projects/${projectId}/targets`, {
        method: 'GET',
        expected_status: 200,
        access_token
    })
}

const createProjectTarget = async (projectId: string, targetType: string, value: string, access_token) => {
    const body = {
        project_id: projectId,
        target_type: targetType,
        value
    }
    return await callAPI("/api/v1/targets", {
        method: 'POST',
        expected_status: 201,
        access_token,
        body: JSON.stringify(body)
    })
}

export {getProjectTargets, createProjectTarget}