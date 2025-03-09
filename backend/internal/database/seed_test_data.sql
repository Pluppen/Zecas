-- seed_test_data.sql
-- Comprehensive script to populate the database with test data for the security scanner
-- Usage: psql -U scanuser -d scandb -f seed_test_data.sql

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Clean up existing data
DELETE FROM findings;
DELETE FROM services;
DELETE FROM target_relations;
DELETE FROM targets;
DELETE FROM scan_tasks;
DELETE FROM scans;
DELETE FROM scan_configs;
DELETE FROM projects;

-- Define static UUIDs for consistent referencing
DO $$
DECLARE
    -- Projects
    test_project_id UUID := '10000000-0000-0000-0000-000000000001';
    
    -- Scan Configs
    ping_config_id UUID := '20000000-0000-0000-0000-000000000001';
    nmap_config_id UUID := '20000000-0000-0000-0000-000000000002';
    dns_config_id UUID := '20000000-0000-0000-0000-000000000003';
    subdomain_config_id UUID := '20000000-0000-0000-0000-000000000004';
    
    -- Scans
    ping_scan_id UUID := '30000000-0000-0000-0000-000000000001';
    nmap_scan_id UUID := '30000000-0000-0000-0000-000000000002';
    dns_scan_id UUID := '30000000-0000-0000-0000-000000000003';
    
    -- Targets: Domains
    example_domain_id UUID := '40000000-0000-0000-0000-000000000001';
    github_domain_id UUID := '40000000-0000-0000-0000-000000000002';
    blog_subdomain_id UUID := '40000000-0000-0000-0000-000000000003';
    api_subdomain_id UUID := '40000000-0000-0000-0000-000000000004';
    
    -- Targets: IPs
    example_ip_id UUID := '40000000-0000-0000-0000-000000000005';
    github_ip1_id UUID := '40000000-0000-0000-0000-000000000006';
    github_ip2_id UUID := '40000000-0000-0000-0000-000000000007';
    google_dns_id UUID := '40000000-0000-0000-0000-000000000008';
    
    -- Targets: CIDRs
    internal_cidr_id UUID := '40000000-0000-0000-0000-000000000009';
    
    -- Services
    example_http_service_id UUID := '50000000-0000-0000-0000-000000000001';
    example_https_service_id UUID := '50000000-0000-0000-0000-000000000002';
    example_ssh_service_id UUID := '50000000-0000-0000-0000-000000000003';
    github_http_service_id UUID := '50000000-0000-0000-0000-000000000004';
    github_https_service_id UUID := '50000000-0000-0000-0000-000000000005';
    
    -- Current timestamp
    current_timestamp TIMESTAMP WITH TIME ZONE := NOW();
BEGIN
    -- Create test project
    INSERT INTO projects (id, name, description, created_at, updated_at)
    VALUES (
        test_project_id,
        'Security Test Project',
        'A comprehensive test project with various targets, services, and findings',
        current_timestamp,
        current_timestamp
    );

    -- Create scan configurations
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
            '{"scan_type": "basic", "port_range": "1-1000", "timing": "4"}'::jsonb,
            true,
            current_timestamp
        ),
        (
            dns_config_id,
            'DNS Resolution',
            'dns',
            '{"record_types": ["A", "AAAA", "CNAME", "MX", "TXT", "NS"]}'::jsonb,
            true,
            current_timestamp
        ),
        (
            subdomain_config_id,
            'Subdomain Enumeration',
            'subdomain',
            '{"recursive": false, "resolve_ip": true}'::jsonb,
            true,
            current_timestamp
        );

    -- Create targets: Domains
    INSERT INTO targets (id, project_id, target_type, value, metadata, created_at, updated_at)
    VALUES 
        (
            example_domain_id,
            test_project_id,
            'domain',
            'example.com',
            '{"description": "Example Domain"}'::jsonb,
            current_timestamp,
            current_timestamp
        ),
        (
            github_domain_id,
            test_project_id,
            'domain',
            'github.com',
            '{"description": "GitHub"}'::jsonb,
            current_timestamp,
            current_timestamp
        ),
        (
            blog_subdomain_id,
            test_project_id,
            'domain',
            'blog.example.com',
            jsonb_build_object('description', 'Blog Subdomain', 'discovered_from', 'example.com', 'discovery_scan', 'subdomain'),
            current_timestamp,
            current_timestamp
        ),
        (
            api_subdomain_id,
            test_project_id,
            'domain',
            'api.example.com',
            jsonb_build_object('description', 'API Subdomain', 'discovered_from', 'example.com', 'discovery_scan', 'subdomain'),
            current_timestamp,
            current_timestamp
        );
    
    -- Create targets: IPs
    INSERT INTO targets (id, project_id, target_type, value, metadata, created_at, updated_at)
    VALUES 
        (
            example_ip_id,
            test_project_id,
            'ip',
            '93.184.216.34',  -- example.com IP
            jsonb_build_object('description', 'Example.com IP', 'discovered_from', 'example.com', 'discovery_scan', 'dns'),
            current_timestamp,
            current_timestamp
        ),
        (
            github_ip1_id,
            test_project_id,
            'ip',
            '140.82.121.3',  -- One of GitHub's IPs
            jsonb_build_object('description', 'GitHub IP 1', 'discovered_from', 'github.com', 'discovery_scan', 'dns'),
            current_timestamp,
            current_timestamp
        ),
        (
            github_ip2_id,
            test_project_id,
            'ip',
            '140.82.121.4',  -- Another of GitHub's IPs
            jsonb_build_object('description', 'GitHub IP 2', 'discovered_from', 'github.com', 'discovery_scan', 'dns'),
            current_timestamp,
            current_timestamp
        ),
        (
            google_dns_id,
            test_project_id,
            'ip',
            '8.8.8.8',  -- Google DNS
            '{"description": "Google DNS Server"}'::jsonb,
            current_timestamp,
            current_timestamp
        );
    
    -- Create targets: CIDRs
    INSERT INTO targets (id, project_id, target_type, value, metadata, created_at, updated_at)
    VALUES 
        (
            internal_cidr_id,
            test_project_id,
            'cidr',
            '192.168.1.0/24',
            '{"description": "Example Internal Network"}'::jsonb,
            current_timestamp,
            current_timestamp
        );

    -- Create target relationships
    -- Domain to IP relationships
    INSERT INTO target_relations (id, source_id, destination_id, relation_type, metadata, created_at, updated_at)
    VALUES 
        (
            uuid_generate_v4(),
            example_domain_id,
            example_ip_id,
            'resolves_to',
            jsonb_build_object('discovered_at', current_timestamp, 'record_type', 'A'),
            current_timestamp,
            current_timestamp
        ),
        (
            uuid_generate_v4(),
            github_domain_id,
            github_ip1_id,
            'resolves_to',
            jsonb_build_object('discovered_at', current_timestamp, 'record_type', 'A'),
            current_timestamp,
            current_timestamp
        ),
        (
            uuid_generate_v4(),
            github_domain_id,
            github_ip2_id,
            'resolves_to',
            jsonb_build_object('discovered_at', current_timestamp, 'record_type', 'A'),
            current_timestamp,
            current_timestamp
        );
    
    -- Domain to subdomain relationships
    INSERT INTO target_relations (id, source_id, destination_id, relation_type, metadata, created_at, updated_at)
    VALUES 
        (
            uuid_generate_v4(),
            example_domain_id,
            blog_subdomain_id,
            'parent_of',
            jsonb_build_object('discovered_at', current_timestamp),
            current_timestamp,
            current_timestamp
        ),
        (
            uuid_generate_v4(),
            example_domain_id,
            api_subdomain_id,
            'parent_of',
            jsonb_build_object('discovered_at', current_timestamp),
            current_timestamp,
            current_timestamp
        );
    
    -- Subdomain to IP relationships
    INSERT INTO target_relations (id, source_id, destination_id, relation_type, metadata, created_at, updated_at)
    VALUES 
        (
            uuid_generate_v4(),
            blog_subdomain_id,
            example_ip_id,
            'resolves_to',
            jsonb_build_object('discovered_at', current_timestamp, 'record_type', 'A'),
            current_timestamp,
            current_timestamp
        ),
        (
            uuid_generate_v4(),
            api_subdomain_id,
            example_ip_id,
            'resolves_to',
            jsonb_build_object('discovered_at', current_timestamp, 'record_type', 'A'),
            current_timestamp,
            current_timestamp
        );

    -- Create services
    INSERT INTO services (id, target_id, port, protocol, service_name, version, title, description, banner, raw_info, created_at, updated_at)
    VALUES 
        (
            example_http_service_id,
            example_ip_id,
            80,
            'tcp',
            'http',
            'Apache/2.4.41',
            'HTTP Service on example.com',
            'Apache web server running on example.com',
            'Apache/2.4.41 (Ubuntu)',
            jsonb_build_object('product', 'Apache httpd', 'version', '2.4.41', 'os', 'Ubuntu', 'discovered_at', current_timestamp),
            current_timestamp,
            current_timestamp
        ),
        (
            example_https_service_id,
            example_ip_id,
            443,
            'tcp',
            'https',
            'Apache/2.4.41',
            'HTTPS Service on example.com',
            'Apache web server with SSL running on example.com',
            'Apache/2.4.41 (Ubuntu)',
            jsonb_build_object('product', 'Apache httpd', 'version', '2.4.41', 'os', 'Ubuntu', 'discovered_at', current_timestamp),
            current_timestamp,
            current_timestamp
        ),
        (
            example_ssh_service_id,
            example_ip_id,
            22,
            'tcp',
            'ssh',
            'OpenSSH 8.2p1',
            'SSH Service on example.com',
            'OpenSSH server running on example.com',
            'SSH-2.0-OpenSSH_8.2p1 Ubuntu-4ubuntu0.5',
            jsonb_build_object('product', 'OpenSSH', 'version', '8.2p1', 'os', 'Ubuntu', 'discovered_at', current_timestamp),
            current_timestamp,
            current_timestamp
        ),
        (
            github_http_service_id,
            github_ip1_id,
            80,
            'tcp',
            'http',
            '',
            'HTTP Service on github.com',
            'Web server running on github.com',
            'Apache',
            jsonb_build_object('product', 'Apache', 'discovered_at', current_timestamp),
            current_timestamp,
            current_timestamp
        ),
        (
            github_https_service_id,
            github_ip1_id,
            443,
            'tcp',
            'https',
            '',
            'HTTPS Service on github.com',
            'Web server with SSL running on github.com',
            'Apache',
            jsonb_build_object('product', 'Apache', 'discovered_at', current_timestamp),
            current_timestamp,
            current_timestamp
        );

    -- Create service-to-target relationships
    INSERT INTO target_relations (id, source_id, destination_id, relation_type, metadata, created_at, updated_at)
    VALUES 
        (
            uuid_generate_v4(),
            example_ip_id,
            example_domain_id,
            'hosts_service',
            jsonb_build_object('service_id', example_http_service_id, 'port', 80, 'protocol', 'tcp'),
            current_timestamp,
            current_timestamp
        ),
        (
            uuid_generate_v4(),
            example_ip_id,
            example_domain_id,
            'hosts_service',
            jsonb_build_object('service_id', example_https_service_id, 'port', 443, 'protocol', 'tcp'),
            current_timestamp,
            current_timestamp
        ),
        (
            uuid_generate_v4(),
            github_ip1_id,
            github_domain_id,
            'hosts_service',
            jsonb_build_object('service_id', github_https_service_id, 'port', 443, 'protocol', 'tcp'),
            current_timestamp,
            current_timestamp
        );

    -- Create scans
    INSERT INTO scans (id, project_id, scan_config_id, status, started_at, completed_at, raw_results, created_at)
    VALUES 
        (
            ping_scan_id,
            test_project_id,
            ping_config_id,
            'completed',
            current_timestamp - INTERVAL '2 hour',
            current_timestamp - INTERVAL '1 hour 55 minutes',
            '{"result": "ping scan completed successfully"}'::jsonb,
            current_timestamp - INTERVAL '2 hours'
        ),
        (
            nmap_scan_id,
            test_project_id,
            nmap_config_id,
            'completed',
            current_timestamp - INTERVAL '1 hour',
            current_timestamp - INTERVAL '45 minutes',
            '{"result": "nmap scan completed successfully"}'::jsonb,
            current_timestamp - INTERVAL '1 hour'
        ),
        (
            dns_scan_id,
            test_project_id,
            dns_config_id,
            'completed',
            current_timestamp - INTERVAL '30 minutes',
            current_timestamp - INTERVAL '25 minutes',
            '{"result": "dns scan completed successfully"}'::jsonb,
            current_timestamp - INTERVAL '30 minutes'
        );

    -- Create findings
    -- Ping findings
    INSERT INTO findings (id, scan_id, target_id, service_id, title, description, severity, finding_type, details, discovered_at, verified, fixed, manual)
    VALUES
        (
            uuid_generate_v4(),
            ping_scan_id,
            example_ip_id,
            NULL,
            'Host example.com is reachable',
            'Host example.com (93.184.216.34) is reachable via ICMP ping with 0% packet loss.',
            'info',
            'ping',
            '{"target": "93.184.216.34", "reachable": true, "min_rtt": 20.5, "avg_rtt": 22.3, "max_rtt": 25.1, "packet_loss": 0}'::jsonb,
            current_timestamp - INTERVAL '1 hour 59 minutes',
            true,
            false,
            false
        ),
        (
            uuid_generate_v4(),
            ping_scan_id,
            github_ip1_id,
            NULL,
            'Host github.com is reachable',
            'Host github.com (140.82.121.3) is reachable via ICMP ping with 0% packet loss.',
            'info',
            'ping',
            '{"target": "140.82.121.3", "reachable": true, "min_rtt": 15.2, "avg_rtt": 17.8, "max_rtt": 19.5, "packet_loss": 0}'::jsonb,
            current_timestamp - INTERVAL '1 hour 58 minutes',
            true,
            false,
            false
        ),
        (
            uuid_generate_v4(),
            ping_scan_id,
            google_dns_id,
            NULL,
            'Host 8.8.8.8 is reachable',
            'Host 8.8.8.8 (Google DNS) is reachable via ICMP ping with 0% packet loss.',
            'info',
            'ping',
            '{"target": "8.8.8.8", "reachable": true, "min_rtt": 8.1, "avg_rtt": 10.2, "max_rtt": 12.5, "packet_loss": 0}'::jsonb,
            current_timestamp - INTERVAL '1 hour 57 minutes',
            true,
            false,
            false
        );

    -- Nmap findings
    INSERT INTO findings (id, scan_id, target_id, service_id, title, description, severity, finding_type, details, discovered_at, verified, fixed, manual)
    VALUES
        (
            uuid_generate_v4(),
            nmap_scan_id,
            example_ip_id,
            example_http_service_id,
            'Open port 80/tcp: http',
            'Port 80/tcp is open (syn-ack).\nService: http\nProduct: Apache httpd 2.4.41\nAdditional info: (Ubuntu)',
            'medium',
            'open_port',
            jsonb_build_object(
                'target', '93.184.216.34',
                'port', 80,
                'protocol', 'tcp',
                'service', 'http',
                'product', 'Apache httpd',
                'version', '2.4.41',
                'state', 'open',
                'reason', 'syn-ack',
                'service_id', example_http_service_id
            ),
            current_timestamp - INTERVAL '55 minutes',
            true,
            false,
            false
        ),
        (
            uuid_generate_v4(),
            nmap_scan_id,
            example_ip_id,
            example_https_service_id,
            'Open port 443/tcp: https',
            'Port 443/tcp is open (syn-ack).\nService: https\nProduct: Apache httpd 2.4.41\nAdditional info: (Ubuntu)',
            'medium',
            'open_port',
            jsonb_build_object(
                'target', '93.184.216.34',
                'port', 443,
                'protocol', 'tcp',
                'service', 'https',
                'product', 'Apache httpd',
                'version', '2.4.41',
                'state', 'open',
                'reason', 'syn-ack',
                'service_id', example_https_service_id
            ),
            current_timestamp - INTERVAL '54 minutes',
            true,
            false,
            false
        ),
        (
            uuid_generate_v4(),
            nmap_scan_id,
            example_ip_id,
            example_ssh_service_id,
            'Open port 22/tcp: ssh',
            'Port 22/tcp is open (syn-ack).\nService: ssh\nProduct: OpenSSH 8.2p1\nAdditional info: Ubuntu-4ubuntu0.5',
            'medium',
            'open_port',
            jsonb_build_object(
                'target', '93.184.216.34',
                'port', 22,
                'protocol', 'tcp',
                'service', 'ssh',
                'product', 'OpenSSH',
                'version', '8.2p1',
                'state', 'open',
                'reason', 'syn-ack',
                'service_id', example_ssh_service_id
            ),
            current_timestamp - INTERVAL '53 minutes',
            true,
            false,
            false
        ),
        (
            uuid_generate_v4(),
            nmap_scan_id,
            example_ip_id,
            NULL,
            'Host 93.184.216.34 has 3 open port(s)',
            'Nmap discovered 3 open port(s) on host 93.184.216.34. See individual findings for details.',
            'info',
            'port_summary',
            jsonb_build_object(
                'target', '93.184.216.34',
                'open_port_count', 3,
                'scan_type', 'basic',
                'ip_address', '93.184.216.34'
            ),
            current_timestamp - INTERVAL '52 minutes',
            true,
            false,
            false
        );

    -- DNS findings
    INSERT INTO findings (id, scan_id, target_id, service_id, title, description, severity, finding_type, details, discovered_at, verified, fixed, manual)
    VALUES
        (
            uuid_generate_v4(),
            dns_scan_id,
            example_domain_id,
            NULL,
            'A records for example.com',
            'The following A records were found for example.com:\n• 93.184.216.34\n\nThese IP addresses are the direct hosts for this domain.',
            'info',
            'dns_records',
            jsonb_build_object(
                'domain', 'example.com',
                'record_type', 'A',
                'records', ARRAY['93.184.216.34']
            ),
            current_timestamp - INTERVAL '28 minutes',
            true,
            false,
            false
        ),
        (
            uuid_generate_v4(),
            dns_scan_id,
            example_domain_id,
            NULL,
            'MX records for example.com',
            'The following MX records were found for example.com:\n• example-smtp.example.com (priority: 10)\n\nThese servers handle email for this domain (lower priority values are preferred).',
            'info',
            'dns_records',
            jsonb_build_object(
                'domain', 'example.com',
                'record_type', 'MX',
                'records', ARRAY['example-smtp.example.com (priority: 10)']
            ),
            current_timestamp - INTERVAL '27 minutes',
            true,
            false,
            false
        ),
        (
            uuid_generate_v4(),
            dns_scan_id,
            github_domain_id,
            NULL,
            'A records for github.com',
            'The following A records were found for github.com:\n• 140.82.121.3\n• 140.82.121.4\n\nThese IP addresses are the direct hosts for this domain.',
            'info',
            'dns_records',
            jsonb_build_object(
                'domain', 'github.com',
                'record_type', 'A',
                'records', ARRAY['140.82.121.3', '140.82.121.4']
            ),
            current_timestamp - INTERVAL '26 minutes',
            true,
            false,
            false
        ),
        (
            uuid_generate_v4(),
            dns_scan_id,
            google_dns_id,
            NULL,
            'Reverse DNS for 8.8.8.8',
            'IP 8.8.8.8 resolves to dns.google',
            'info',
            'dns_ptr',
            jsonb_build_object(
                'ip', '8.8.8.8',
                'hostnames', ARRAY['dns.google']
            ),
            current_timestamp - INTERVAL '25 minutes',
            true,
            false,
            false
        );

    -- Create scan tasks
    INSERT INTO scan_tasks (id, scan_id, task_type, parameters, status, result, created_at, updated_at)
    VALUES
        (
            uuid_generate_v4(),
            ping_scan_id,
            'ping',
            '{"target": "example.com", "count": 3, "timeout": 2}'::jsonb,
            'completed',
            '{"min_rtt": 20.5, "avg_rtt": 22.3, "max_rtt": 25.1, "packet_loss": 0, "reachable": true}'::jsonb,
            current_timestamp - INTERVAL '2 hours',
            current_timestamp - INTERVAL '1 hour 55 minutes'
        ),
        (
            uuid_generate_v4(),
            ping_scan_id,
            'ping',
            '{"target": "github.com", "count": 3, "timeout": 2}'::jsonb,
            'completed',
            '{"min_rtt": 15.2, "avg_rtt": 17.8, "max_rtt": 19.5, "packet_loss": 0, "reachable": true}'::jsonb,
            current_timestamp - INTERVAL '2 hours',
            current_timestamp - INTERVAL '1 hour 55 minutes'
        ),
        (
            uuid_generate_v4(),
            ping_scan_id,
            'ping',
            '{"target": "8.8.8.8", "count": 3, "timeout": 2}'::jsonb,
            'completed',
            '{"min_rtt": 8.1, "avg_rtt": 10.2, "max_rtt": 12.5, "packet_loss": 0, "reachable": true}'::jsonb,
            current_timestamp - INTERVAL '2 hours',
            current_timestamp - INTERVAL '1 hour 55 minutes'
        ),
        (
            uuid_generate_v4(),
            nmap_scan_id,
            'nmap',
            '{"target": "93.184.216.34", "scan_type": "basic", "port_range": "1-1000", "timing": "4"}'::jsonb,
            'completed',
            '{"open_ports": [22, 80, 443], "services": {"22": "ssh", "80": "http", "443": "https"}}'::jsonb,
            current_timestamp - INTERVAL '1 hour',
            current_timestamp - INTERVAL '45 minutes'
        ),
        (
            uuid_generate_v4(),
            dns_scan_id,
            'dns',
            '{"target": "example.com", "record_types": ["A", "AAAA", "CNAME", "MX", "TXT", "NS"]}'::jsonb,
            'completed',
            jsonb_build_object(
                'records', jsonb_build_object(
                    'A', ARRAY['93.184.216.34'],
                    'MX', ARRAY['example-smtp.example.com (priority: 10)'],
                    'NS', ARRAY['a.iana-servers.net', 'b.iana-servers.net']
                )
            ),
            current_timestamp - INTERVAL '30 minutes',
            current_timestamp - INTERVAL '25 minutes'
        );

    -- Create a pending scan for testing worker functionality
    INSERT INTO scans (id, project_id, scan_config_id, status, created_at)
    VALUES (
        uuid_generate_v4(),
        test_project_id,
        subdomain_config_id,
        'pending',
        current_timestamp
    );


    -- Add Nuclei scan configuration
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES (
        uuid_generate_v4(),
        'Basic Nuclei Vulnerability Scan',
        'nuclei',
        '{
            "template_tags": ["cve", "vulnerability", "config"],
            "template_exclude": ["dos", "fuzz"],
            "severity": ["medium", "high", "critical"],
            "timeout": 300,
            "rate_limit": 150,
            "bulk_size": 25,
            "headless": false
        }'::jsonb,
        true,
        NOW()
    );

    -- Add comprehensive Nuclei scan configuration
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES (
        uuid_generate_v4(),
        'Comprehensive Nuclei Vulnerability Scan',
        'nuclei',
        '{
            "template_tags": ["cve", "vulnerability", "config", "exposure", "misconfiguration", "takeover"],
            "template_exclude": ["dos", "fuzz"],
            "severity": ["low", "medium", "high", "critical"],
            "timeout": 600,
            "rate_limit": 100,
            "bulk_size": 25,
            "headless": true
        }'::jsonb,
        true,
        NOW()
    );

    -- Add tech detection Nuclei scan configuration
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES (
        uuid_generate_v4(),
        'Technology Detection Scan',
        'nuclei',
        '{
            "template_tags": ["tech", "technology"],
            "severity": ["info", "low", "medium", "high", "critical"],
            "timeout": 300,
            "rate_limit": 150,
            "bulk_size": 25
        }'::jsonb,
        true,
        NOW()
    );

    -- Add targeted CVE scan configuration
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES (
        uuid_generate_v4(),
        'CVE Targeted Scan',
        'nuclei',
        '{
            "template_tags": ["cve"],
            "severity": ["high", "critical"],
            "timeout": 300,
            "rate_limit": 150,
            "bulk_size": 25
        }'::jsonb,
        true,
        NOW()
    );

    -- Add HTTPX scan configuration
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES (
        uuid_generate_v4(),
        'Basic HTTPX Web Scan',
        'httpx',
        '{
            "timeout": 30,
            "threads": 50,
            "follow_redirects": true,
            "tech_detect": true,
            "status_code": true,
            "title": true,
            "web_server": true,
            "content_type": true,
            "tls": true,
            "favicon": true,
            "jarm": false,
            "probe": true,
            "http2": true,
            "security_headers": true,
            "extract_cname": true
        }'::jsonb,
        true,
        NOW()
    );

    -- Add HTTPX scan configuration with custom ports
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES (
        uuid_generate_v4(),
        'HTTPX Multi-Port Web Scan',
        'httpx',
        '{
            "timeout": 30,
            "threads": 50,
            "follow_redirects": true,
            "tech_detect": true,
            "status_code": true,
            "title": true,
            "web_server": true,
            "content_type": true,
            "tls": true,
            "favicon": true,
            "jarm": false,
            "probe": true,
            "ports": "80,81,443,3000,8000,8080,8443",
            "http2": true,
            "security_headers": true,
            "extract_cname": true
        }'::jsonb,
        true,
        NOW()
    );

    -- Add HTTPX technology detection scan configuration
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES (
        uuid_generate_v4(),
        'HTTPX Technology Detection Scan',
        'httpx',
        '{
            "timeout": 30,
            "threads": 50,
            "follow_redirects": true,
            "tech_detect": true,
            "status_code": true,
            "title": true,
            "web_server": true,
            "content_type": false,
            "tls": false,
            "favicon": true,
            "jarm": false,
            "probe": true,
            "http2": false,
            "security_headers": false,
            "extract_cname": false
        }'::jsonb,
        true,
        NOW()
    );

    -- Add HTTPX security header analysis scan configuration
    INSERT INTO scan_configs (id, name, scanner_type, parameters, active, created_at)
    VALUES (
        uuid_generate_v4(),
        'HTTPX Security Headers Analysis',
        'httpx',
        '{
            "timeout": 30,
            "threads": 50,
            "follow_redirects": true,
            "tech_detect": false,
            "status_code": true,
            "title": false,
            "web_server": true,
            "content_type": false,
            "tls": true,
            "favicon": false,
            "jarm": false,
            "probe": true,
            "http2": false,
            "security_headers": true,
            "extract_cname": false
        }'::jsonb,
        true,
        NOW()
    );

    RAISE NOTICE 'Test data creation completed';
    RAISE NOTICE 'Project ID: %', test_project_id;
    RAISE NOTICE 'Example domain target ID: %', example_domain_id;
END $$;