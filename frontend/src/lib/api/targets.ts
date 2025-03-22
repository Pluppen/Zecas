import { callAPI } from "@/lib/api";

export const getProjects = async (access_token: string) => {
    return await callAPI(`/api/v1/targets`, {
        method: "GET",
        access_token,
        expected_status: 200
    })
}

export const bulkCreateTargets = async () => {
    // TODO
}


export const createProjectTarget = async (projectId: string, targetType: string, value: string, access_token) => {
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

export const getTargetById = async (targetId: string, access_token: string) => {
    return await callAPI(`/api/v1/targets/${targetId}`, {
        access_token,
        expected_status: 200,
        method: "GET"
    })
}

export type Target = {
    id: string,
    target_type: string,
    value: string,
    metadata: JSON
}

export const updateTargetById = async (target: Target, access_token: string) => {
    const body = {
        target_type: target.target_type,
        value: target.value,
        metadata: target.metadata
    }

    return await callAPI(`/api/v1/targets/${target.id}`, {
        access_token,
        expected_status: 200,
        body: JSON.stringify(body),
        method: "PUT"
    })
}

export const deleteTargetById = async (targetId: string, access_token: string) => {
    return await callAPI(`/api/v1/targets/${targetId}`, {
        access_token,
        expected_status: 200,
        method: "DELETE"
    })
}

export const getTargetFindings = async (targetId: string, access_token: string) => {
    return await callAPI(`/api/v1/targets/${targetId}/findings`, {
        access_token,
        expected_status: 200,
        method: "GET"
    })
}

export const getTargetServices = async (targetId: string, access_token: string) => {
    return await callAPI(`/api/v1/targets/${targetId}/services`, {
        access_token,
        expected_status: 200,
        method: "GET"
    })
}

export const getTargetRelations = async (targetId: string, access_token: string) => {
    return await callAPI(`/api/v1/targets/${targetId}/relations`, {
        access_token,
        expected_status: 200,
        method: "GET"
    })
}