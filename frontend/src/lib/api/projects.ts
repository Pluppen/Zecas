import { callAPI } from ".";

export const getProjectById = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}`, {
        method: "GET",
        access_token,
        expected_status: 200
    })    
}

export const createNewProject = async (body: string, access_token: string) => {
    return await callAPI("/api/v1/projects", {
        method: "POST",
        access_token,
        body,
        expected_status: 201
    })    
}

export const updateProject = async (projectId: string, name: string, description: string | undefined, access_token: string) => {
    const body = {
        name,
        description
    }    

    return await callAPI(`/api/v1/projects/${projectId}`, {
        method: "PUT",
        access_token,
        body: JSON.stringify(body),
        expected_status: 200
    })    
}

export const getProjects = async (access_token: string) => {
    return await callAPI(`/api/v1/projects/`, {
        method: "GET",
        expected_status: 200,
        access_token
    })
}

export const deleteProject = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}`, {
        method: "DELETE",
        expected_status: 200,
        access_token
    })
}

export const getProjectTargets = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}/targets`, {
        method: 'GET',
        expected_status: 200,
        access_token
    })
}

export const getProjectScans = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}/scans`, {
        method: 'GET',
        expected_status: 200,
        access_token
    })
}

export const getProjectFindings = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}/findings`, {
        method: 'GET',
        access_token,
        expected_status: 200
    })
}

export const getProjectServices = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}/services`, {
        method: 'GET',
        access_token,
        expected_status: 200
    })
}

export const getProjectApplications = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}/applications`, {
        method: 'GET',
        access_token,
        expected_status: 200
    })
}

export const getProjectDNSRecords = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}/dns-records`, {
        method: 'GET',
        access_token,
        expected_status: 200
    })
}

export const getProjectCertificates = async (projectId: string, access_token: string) => {
    return await callAPI(`/api/v1/projects/${projectId}/certificates`, {
        method: 'GET',
        access_token,
        expected_status: 200
    })
}
