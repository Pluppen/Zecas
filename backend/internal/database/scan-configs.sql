-- Create scan configurations
INSERT INTO scan_configs (name, scanner_type, parameters, active, created_at)
VALUES 
    (
        'Heavy Port Scan',
        'nmap',
        '{"scan_type": "service", "port_range": "1-65535", "timing": "4"}'::jsonb,
        true,
        current_timestamp
    );