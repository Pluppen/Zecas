import { API_URL } from "@/config";

const getProjectTargets = async (projectId: string) => {
    const response = await fetch(`${API_URL}/api/v1/projects/${projectId}/targets`)
    console.log(response);
    const targets = response.json();

    return targets
}

const createProjectTarget = async (projectId: string, targetType: string, value: string) => {
    const body = {
        project_id: projectId,
        target_type: targetType,
        value
    }
    const response = await fetch(`${API_URL}/api/v1/targets`, {
        method: 'POST',
        body: JSON.stringify(body)
    })
    if (response.status != 201) {
        console.error(response.statusText);
        return {error: "Something went wrong"}
    }
    return response.json()
}

export {getProjectTargets, createProjectTarget}