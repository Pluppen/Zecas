import { callAPI } from "@/lib/api";
import { z } from "zod";

export type ScannerType = "nmap" | "dns" | "subdomain" | "nuclei" | "httpx";

const NmapParametersSchema = z.object({
    scan_type: z.enum([
        "quick",
        "comprehensive",
        "service",
        "all_ports",
        "basic",
    ]),
    port_range: z.string(),
    timing: z.enum(["0", "1", "2", "3", "4", "5"]),
});

const allowedRecordTypes = ["A", "AAAA", "CNAME", "MX", "TXT", "NS"] as const;
const DNSParametersSchema = z.object({
    record_types: z.array(z.enum(allowedRecordTypes)),
});


const SubdomainParametersSchema = z.object({
    recursive: z.boolean(),
    resolve_ip: z.boolean(),
    wordlist: z.string(),
    timeout: z.number(),
});

const NucleiParametersSchema = z.object({
    template_tags: z.array(z.string()),
    template_paths: z.array(z.string()),
    template_exclude: z.array(z.string()),
    severity: z.string(),
    timeout: z.number(),
    rate_limit: z.number(),
    bulk_size: z.number(),
    templates_dir: z.string(),
    headless: z.boolean(),
    include_all: z.boolean(),
});

const HttpxParametersSchema = z.object({
    timeout: z.number(),
    threads: z.number(),
    follow_redirects: z.boolean(),
    tech_detect: z.boolean(),
    status_code: z.boolean(),
    title: z.boolean(),
    web_server: z.boolean(),
    content_type: z.boolean(),
    tls: z.boolean(),
    favicon: z.boolean(),
    jarm: z.boolean(),
    probe: z.boolean(),
    ports: z.string(),
    http2: z.boolean(),
    security_headers: z.boolean(),
    extract_cname: z.boolean(),
});

const ScanConfigBaseSchema = z.object({
    id: z.string().optional(),
    name: z.string(),
    active: z.boolean(),
});

export const ScanConfigSchema = z.discriminatedUnion("scanner_type", [
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("nmap"),
        parameters: NmapParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("dns"),
        parameters: DNSParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("subdomain"),
        parameters: SubdomainParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("nuclei"),
        parameters: NucleiParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("httpx"),
        parameters: HttpxParametersSchema,
    }),
]);

export const getScanConfigs = async (access_token: string) => {
    return await callAPI("/api/v1/scan-configs", {
        method: "GET",
        expected_status: 200,
        access_token,
    });
};

export const createScanConfig = async (
    scanConfig: ScanConfig,
    access_token: string,
) => {
    const body = {
        name: scanConfig.name,
        scanner_type: scanConfig.scanner_type,
        parameters: scanConfig.parameters,
        active: scanConfig.active,
    };

    return await callAPI(`/api/v1/scan-configs`, {
        method: "POST",
        access_token,
        body: JSON.stringify(body),
        expected_status: 201,
    });
};

export const getScanConfigById = async (
    scanConfigId: string,
    access_token: string,
) => {
    return await callAPI(`/api/v1/scan-configs/${scanConfigId}`, {
        method: "GET",
        access_token,
        expected_status: 200,
    });
};

export const updateScanConfigById = async (
    scanConfig: ScanConfig,
    access_token: string,
) => {
    const body = {
        name: scanConfig.name,
        scanner_type: scanConfig.scanner_type,
        parameters: scanConfig.parameters,
        active: scanConfig.active,
    };

    return await callAPI(`/api/v1/scan-configs/${scanConfig.id}`, {
        method: "PUT",
        access_token,
        body: JSON.stringify(body),
        expected_status: 200,
    });
};

export const deleteScanConfigById = async (
    scanId: string,
    access_token: string,
) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}`, {
        method: "DELETE",
        access_token,
        expected_status: 200,
    });
};

// SCANS

export const startNewScan = async (
    projectId: string,
    scanConfigId: string,
    targetIds: string[],
    access_token: string,
) => {
    const body = {
        project_id: projectId,
        scan_config_id: scanConfigId,
        target_ids: targetIds,
    };
    return await callAPI("/api/v1/scans", {
        method: "POST",
        expected_status: 202,
        body: JSON.stringify(body),
        access_token,
    });
};

export const getScanById = async (scanId: string, access_token: string) => {
    return await callAPI(`/api/v1/scans/${scanId}`, {
        method: "GET",
        access_token,
        expected_status: 200,
    });
};

export const getScans = async (access_token: string) => {
    return await callAPI(`/api/v1/scan/`, {
        method: "GET",
        access_token,
        expected_status: 200,
    });
};

export const cancelScanById = async (scanId: string, access_token: string) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}/cancel`, {
        method: "POST",
        access_token,
        expected_status: 200,
    });
};

export const getFindingsByScanId = async (
    scanId: string,
    access_token: string,
) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}/findings`, {
        method: "POST",
        access_token,
        expected_status: 200,
    });
};
