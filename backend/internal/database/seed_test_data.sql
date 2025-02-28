-- seed_test_data.sql
-- Populates the security scanner database with test data
-- Usage: psql -U scanuser -d scandb -f seed_test_data.sql

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Clean up existing data (optional)
DELETE FROM findings;
DELETE FROM scan_tasks;
DELETE FROM scans;
DELETE FROM scan_configs;
DELETE FROM targets;
DELETE FROM projects;

-- Generate static UUIDs for consistent referencing
-- Project ID
DO $$
DECLARE
    project_id UUID := '11111111-1111-1111-1111-111111111111';
    scan_id UUID := '22222222-2222-2222-2222-222222222222';
    ping_config_id UUID := '33333333-3333-3333-3333-333333333333';
    nmap_config_id UUID := '44444444-4444-4444-4444-444444444444';
    dns_config_id UUID := '55555555-5555-5555-5555-555555555555';
    target_id_1 UUID := 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
    target_id_2 UUID := 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb';
    target_id_3 UUID := 'cccccccc-cccc-cccc-cccc-cccccccccccc';
    target_id_4 UUID := 'dddddddd-dddd-dddd-dddd-dddddddddddd';
    target_id_5 UUID := 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee';
    current_timestamp TIMESTAMP WITH TIME ZONE := NOW();
BEGIN
    -- Create Project
    INSERT INTO projects (id, name, description, created_at, updated_at)
    VALUES (
        project_id,
        'Test Security Project',
        'A test project with various targets for demonstration purposes',
        current_timestamp,
        current_timestamp
    );

    -- Create Targets
    -- IP Targets
    INSERT INTO targets (id, project_id, target_type, value, metadata, created_at, updated_at)
    VALUES 
        (
            target_id_1,
            project_id,
            'ip',
            '8.8.8.8',
            '{"description": "Google DNS Server"}'::jsonb,
            current_timestamp,
            current_timestamp
        ),
        (
            target_id_2,
            project_id,
            'ip',
            '1.1.1.1',
            '{"description": "Cloudflare DNS Server"}'::jsonb,
            current_timestamp,
            current_timestamp
        ),
        (
            target_id_3,
            project_id,
            'cidr',
            '192.168.1.0/24',
            '{"description": "Example Internal Network"}'::jsonb,
            current_timestamp,
            current_timestamp
        ),
        (
            target_id_4,
            project_id,
            'domain',
            'example.com',
            '{"description": "Example Domain"}'::jsonb,
            current_timestamp,
            current_timestamp
        ),
        (
            target_id_5,
            project_id,
            'domain',
            'github.com',
            '{"description": "GitHub"}'::jsonb,
            current_timestamp,
            current_timestamp
        );

    -- Create Scan Configurations
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES 
        (
            ping_config_id,
            'Basic Ping Test',
            'ping',
            '{"count": 3, "timeout": 2}'::jsonb,
            true,
            current_timestamp
        ),
        (
            nmap_config_id,
            'Basic Port Scan',
            'nmap',
            '{"intensity": "quick"}'::jsonb,
            true,
            current_timestamp
        ),
        (
            dns_config_id,
            'DNS Enumeration',
            'dnsenum',
            '{"scan_type": "basic"}'::jsonb,
            true,
            current_timestamp
        );

    -- Create Scan
    INSERT INTO scans (id, project_id, scan_config_id, status, started_at, completed_at, raw_results, created_at)
    VALUES (
        scan_id,
        project_id,
        ping_config_id,
        'completed',
        current_timestamp - INTERVAL '1 hour',
        current_timestamp - INTERVAL '45 minutes',
        '{"example_result": "This is a simulated scan result"}'::jsonb,
        current_timestamp - INTERVAL '1 hour'
    );

    -- Create Findings
    -- Medium severity findings
    INSERT INTO findings (id, scan_id, target_id, title, description, severity, finding_type, details, discovered_at, verified, fixed)
    VALUES
        (
            uuid_generate_v4(),
            scan_id,
            target_id_3,
            'Simulated Medium Vulnerability',
            'This is a simulated medium severity vulnerability created for testing purposes.',
            'medium',
            'vulnerability',
            '{"simulated": true, "cve": "CVE-2023-12345", "attack_vector": "network"}'::jsonb,
            current_timestamp - INTERVAL '50 minutes',
            false,
            false
        ),
        (
            uuid_generate_v4(),
            scan_id,
            target_id_2,
            'Outdated Software Version',
            'System is running an outdated version of software that may contain known vulnerabilities.',
            'medium',
            'vulnerability',
            '{"simulated": true, "software": "OpenSSL", "current_version": "1.0.2", "recommended_version": "1.1.1"}'::jsonb,
            current_timestamp - INTERVAL '50 minutes',
            true,
            false
        );

    -- High severity finding
    INSERT INTO findings (id, scan_id, target_id, title, description, severity, finding_type, details, discovered_at, verified, fixed)
    VALUES
        (
            uuid_generate_v4(),
            scan_id,
            target_id_1,
            'Simulated Critical Port Exposure',
            'System has critical administrative ports exposed to the public internet.',
            'high',
            'port_exposure',
            '{"simulated": true, "port": 22, "service": "SSH", "version": "OpenSSH 7.4"}'::jsonb,
            current_timestamp - INTERVAL '50 minutes',
            true,
            false
        );

    -- Unknown severity finding
    INSERT INTO findings (id, scan_id, target_id, title, description, severity, finding_type, details, discovered_at, verified, fixed)
    VALUES
        (
            uuid_generate_v4(),
            scan_id,
            target_id_3,
            'Unidentified port 142',
            'System has an unidentified port',
            'unknown',
            'port_exposure',
            '{"simulated": true, "port": 22, "service": "SSH", "version": "OpenSSH 7.4"}'::jsonb,
            current_timestamp - INTERVAL '50 minutes',
            false,
            false
        );

    -- Critical severity finding
    INSERT INTO findings (id, scan_id, target_id, title, description, severity, finding_type, details, discovered_at, verified, fixed)
    VALUES
        (
            uuid_generate_v4(),
            scan_id,
            target_id_5,
            'Simulated Critical Port Exposure',
            'System has critical administrative ports exposed to the public internet.',
            'critical',
            'port_exposure',
            '{"simulated": true, "port": 22, "service": "SSH", "version": "OpenSSH 7.4"}'::jsonb,
            current_timestamp - INTERVAL '50 minutes',
            true,
            false
        );
    
    -- Info finding
    INSERT INTO findings (id, scan_id, target_id, title, description, severity, finding_type, details, discovered_at, verified, fixed)
    VALUES
        (
            uuid_generate_v4(),
            scan_id,
            target_id_1,
            'System Information Disclosure',
            'System information gathered during reconnaissance phase.',
            'info',
            'info_disclosure',
            '{"simulated": true, "os": "Linux 5.10", "hostname": "test-server-1", "uptime": "15 days"}'::jsonb,
            current_timestamp - INTERVAL '49 minutes',
            true,
            true
        );

    -- Create Scan Tasks
    INSERT INTO scan_tasks (id, scan_id, task_type, parameters, status, result, created_at, updated_at)
    VALUES
        (
            uuid_generate_v4(),
            scan_id,
            'ping',
            '{"target": "8.8.8.8", "count": 3}'::jsonb,
            'completed',
            '{"min_rtt": 10.2, "avg_rtt": 12.5, "max_rtt": 15.1, "packet_loss": 0}'::jsonb,
            current_timestamp - INTERVAL '1 hour',
            current_timestamp - INTERVAL '59 minutes'
        ),
        (
            uuid_generate_v4(),
            scan_id,
            'ping',
            '{"target": "1.1.1.1", "count": 3}'::jsonb,
            'completed',
            '{"min_rtt": 8.5, "avg_rtt": 9.7, "max_rtt": 11.2, "packet_loss": 0}'::jsonb,
            current_timestamp - INTERVAL '59 minutes',
            current_timestamp - INTERVAL '58 minutes'
        );

    -- Create a pending scan to test workers
    INSERT INTO scans (id, project_id, scan_config_id, status, created_at)
    VALUES (
        uuid_generate_v4(),
        project_id,
        ping_config_id,
        'pending',
        current_timestamp
    );

    RAISE NOTICE 'Test data creation completed!';
    RAISE NOTICE 'Project ID: %', project_id;
    RAISE NOTICE 'Scan ID: %', scan_id;
    RAISE NOTICE 'Ping Config ID: %', ping_config_id;
END $$;
