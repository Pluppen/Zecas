import { API_URL } from "@/config";

const getProjectFindings = async (projectId: string) => {
    const response = await fetch(`${API_URL}/api/v1/projects/${projectId}/findings`)
    if (response.status != 200) {
        return {error: "Something went wrong fetching scans for project"};
    }
    return response.json();
}

interface CreateFindingParam {
    target_id: string,
    title: string,
    description?: string,
    severity: string,
    finding_type: string,
    details?: string,
    manual?: boolean,
    scan_id?: string
}

const createFinding = async (finding: CreateFindingParam) => {
    const body: CreateFindingParam = {
        target_id: finding.target_id,
        title: finding.title,
        severity: finding.severity,
        finding_type: finding.finding_type,
    }

    if (finding.description) {
        body.description = finding.description
    }

    if (finding.scan_id) {
        body.scan_id= finding.scan_id
    }

    if (finding.details) {
        body.details = finding.details
    }

    if (finding.manual) {
        body.manual = finding.manual
    }

    const response = await fetch(`${API_URL}/api/v1/findings`, {
        method: 'POST',
        body: JSON.stringify(body)
    })
    if (response.status != 201) {
        console.error(response.statusText);
        return {error: "Something went wrong"}
    }
    return response.json()
}

const removeFinding = async (findingId: string) => {
    const response = await fetch(`${API_URL}/api/v1/findings/${findingId}`, {
        method: 'DELETE'
    })
    if (response.status != 200) {
        console.error(response.statusText);
        return {error: "Something went wrong"}
    }
    return response.json()
}

export {getProjectFindings, createFinding, removeFinding}