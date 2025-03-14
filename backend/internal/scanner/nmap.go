// internal/scanner/nmap.go
package scanner

// TODO: Write down a spec on what we want out of this scanner, what types of assets shoudl it create
// and also if this scanner should be able to generate findings and
// what is counted as a finding from this tool.

import (
	"backend/internal/models"
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

// NmapScanner implements the Scanner interface for performing Nmap scans
type NmapScanner struct {
	binPath string
}

// NmapXML represents the XML output structure from nmap
type NmapXML struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []Host   `xml:"host"`
}

// Host represents an nmap host
type Host struct {
	Status    Status  `xml:"status"`
	Address   Address `xml:"address"`
	Ports     Ports   `xml:"ports"`
	Hostnames struct {
		Hostnames []Hostname `xml:"hostname"`
	} `xml:"hostnames"`
}

// Status represents the status of a host
type Status struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
}

// Address represents an address (IP or MAC)
type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
}

// Ports represents a collection of ports
type Ports struct {
	Ports []Port `xml:"port"`
}

// Port represents a single port
type Port struct {
	Protocol string  `xml:"protocol,attr"`
	PortID   int     `xml:"portid,attr"`
	State    State   `xml:"state"`
	Service  Service `xml:"service"`
}

// State represents the state of a port
type State struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
}

// Service represents a service detected on a port
type Service struct {
	Name      string `xml:"name,attr"`
	Product   string `xml:"product,attr"`
	Version   string `xml:"version,attr"`
	ExtraInfo string `xml:"extrainfo,attr"`
}

// Hostname represents a hostname
type Hostname struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// NewNmapScanner creates a new Nmap scanner
func NewNmapScanner() *NmapScanner {
	return &NmapScanner{
		binPath: "nmap",
	}
}

// Initialize checks if nmap is available
func (s *NmapScanner) Initialize(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, s.binPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nmap command not available: %w", err)
	}

	// Check if it's really nmap
	if !strings.Contains(string(output), "Nmap") {
		return fmt.Errorf("command '%s' does not appear to be nmap", s.binPath)
	}

	return nil
}

// ConvertTarget converts a Target to a format suitable for nmap
func (s *NmapScanner) ConvertTarget(target models.Target) interface{} {
	return target.Value
}

// ConvertService converts a Service to a format suitable for nmap
func (s *NmapScanner) ConvertService(service models.Service) interface{} {
	// For a service, we need the target value and the port
	return nil // Not directly supported
}

// isCIDR checks if a target is a CIDR range
func (s *NmapScanner) isCIDR(targetValue string) bool {
	_, _, err := net.ParseCIDR(targetValue)
	return err == nil
}

// Scan performs an nmap scan against the target
func (s *NmapScanner) Scan(ctx context.Context, target interface{}, params models.JSONB) (*models.ScanResults, error) {
	targetValue := target.(string)
	scanResults := &models.ScanResults{
		Findings:        []models.Finding{},
		NewTargets:      []models.Target{},
		TargetRelations: []models.TargetRelation{},
		Services:        []models.Service{},
	}

	// Check if this is a CIDR range
	isCIDR := s.isCIDR(targetValue)

	// Default scan options
	scanType := "basic"
	portRange := "1-1000"
	timing := "4"

	// Override with provided parameters if available
	if val, ok := params["scan_type"].(string); ok {
		scanType = val
	}

	if val, ok := params["port_range"].(string); ok {
		portRange = val
	}

	if val, ok := params["timing"].(string); ok {
		timing = val
	}

	// Build nmap command based on scan type
	args := []string{"-oX", "-"} // Output XML to stdout

	// Add timing template
	args = append(args, "-T"+timing)

	// Configure scan type
	switch scanType {
	case "quick":
		args = append(args, "-F") // Fast mode - fewer ports
	case "comprehensive":
		args = append(args, "--top-ports", "2000")
		args = append(args, "-sV") // Service version detection
		args = append(args, "-O")  // OS detection
	case "service":
		args = append(args, "-sV") // Service version detection
		args = append(args, "-p", portRange)
	case "all_ports":
		args = append(args, "-p-") // All ports
	default: // "basic"
		args = append(args, "-p", portRange)
	}

	// Add target
	args = append(args, targetValue)

	// Execute nmap command with timeout
	cmd := exec.CommandContext(ctx, s.binPath, args...)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil || !strings.Contains(outputStr, "<nmaprun") {
		return nil, fmt.Errorf("nmap scan failed: %w", err)
	}

	// Parse XML output
	var result NmapXML
	log.Printf("%s", output)
	if err := xml.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse nmap output: %w", err)
	}

	// Process scan results for each host
	for _, host := range result.Hosts {
		// Skip hosts that are not up
		if host.Status.State != "up" {
			continue
		}

		// Get the IP address
		ipAddress := host.Address.Addr
		var ipTarget *models.Target
		var ipTargetID uuid.UUID

		// Create a new target for this IP if we're scanning a CIDR range
		if isCIDR {
			ipTarget = &models.Target{
				ID:         uuid.New(),
				TargetType: models.TargetTypeIP,
				Value:      ipAddress,
				Metadata: models.JSONB{
					"discovered_from": targetValue,
					"discovery_scan":  "nmap",
					"discovered_at":   time.Now().Format(time.RFC3339),
				},
			}
			scanResults.NewTargets = append(scanResults.NewTargets, *ipTarget)
			ipTargetID = ipTarget.ID

			// Create a relationship between the CIDR and this IP
			relation := models.TargetRelation{
				ID:            uuid.New(),
				SourceID:      uuid.Nil, // Will be set by worker to the original CIDR target ID
				DestinationID: ipTarget.ID,
				RelationType:  "contains",
				Metadata: models.JSONB{
					"discovered_at": time.Now().Format(time.RFC3339),
				},
			}
			scanResults.TargetRelations = append(scanResults.TargetRelations, relation)
		} else {
			// If not a CIDR, use the original target ID (will be set by worker)
			ipTargetID = uuid.Nil
		}

		// Process hostnames
		for _, hostname := range host.Hostnames.Hostnames {
			// Create a new target for the hostname
			if hostname.Name != "" && (hostname.Type == "user" || hostname.Type == "PTR") {
				hostnameTarget := models.Target{
					ID:         uuid.New(),
					TargetType: models.TargetTypeDomain,
					Value:      hostname.Name,
					Metadata: models.JSONB{
						"discovered_from": ipAddress,
						"discovery_scan":  "nmap",
						"hostname_type":   hostname.Type,
						"discovered_at":   time.Now().Format(time.RFC3339),
					},
				}
				scanResults.NewTargets = append(scanResults.NewTargets, hostnameTarget)

				destinationID := uuid.Nil
				if isCIDR {
					destinationID = ipTarget.ID
				}

				// Create relation between IP and hostname
				ipToHostnameRelation := models.TargetRelation{
					ID:            uuid.New(),
					SourceID:      hostnameTarget.ID,
					DestinationID: destinationID, // Use the new IP target if CIDR
					RelationType:  models.RelationResolvesTo,
					Metadata: models.JSONB{
						"discovered_at": time.Now().Format(time.RFC3339),
					},
				}
				scanResults.TargetRelations = append(scanResults.TargetRelations, ipToHostnameRelation)
			}
		}

		// Count open ports
		openPortCount := 0
		targetForFindings := uuid.Nil // If CIDR, use the new IP target ID, otherwise the original target
		if isCIDR {
			targetForFindings = ipTargetID
		}

		// Process open ports
		for _, port := range host.Ports.Ports {
			if port.State.State == "open" {
				openPortCount++

				// Create a service for each open port
				service := models.Service{
					ID:          uuid.New(),
					TargetID:    targetForFindings, // Will be set by worker
					Port:        port.PortID,
					Protocol:    port.Protocol,
					ServiceName: port.Service.Name,
					Version:     port.Service.Version,
					Title:       fmt.Sprintf("%s service on port %d", port.Service.Name, port.PortID),
					Description: s.generateServiceDescription(port),
					Banner:      port.Service.ExtraInfo,
					RawInfo: models.JSONB{
						"product":       port.Service.Product,
						"version":       port.Service.Version,
						"extra_info":    port.Service.ExtraInfo,
						"state":         port.State.State,
						"reason":        port.State.Reason,
						"target_value":  ipAddress,
						"discovered_at": time.Now().Format(time.RFC3339),
					},
				}
				scanResults.Services = append(scanResults.Services, service)

				// Create a finding for each service
				// TODO: Change this to actually be a finding.
				// finding := models.Finding{
				// 	Title:       fmt.Sprintf("Open port %d/%s: %s on %s", port.PortID, port.Protocol, port.Service.Name, ipAddress),
				// 	Description: s.generatePortDescription(port, ipAddress),
				// 	Severity:    s.determineSeverityForPort(port),
				// 	TargetID:    targetForFindings,
				// 	FindingType: "open_port",
				// 	Details: models.JSONB{
				// 		"target":        ipAddress,
				// 		"port":          port.PortID,
				// 		"protocol":      port.Protocol,
				// 		"service":       port.Service.Name,
				// 		"product":       port.Service.Product,
				// 		"version":       port.Service.Version,
				// 		"state":         port.State.State,
				// 		"reason":        port.State.Reason,
				// 		"service_id":    service.ID.String(),
				// 		"discovered_at": time.Now().Format(time.RFC3339),
				// 	},
				// }
				// scanResults.Findings = append(scanResults.Findings, finding)
			}
		}

		// Create summary finding for this host
		if openPortCount > 0 {
			finding := models.Finding{
				Title:       fmt.Sprintf("Host %s has %d open port(s)", ipAddress, openPortCount),
				Description: fmt.Sprintf("Nmap discovered %d open port(s) on host %s. See individual findings for details.", openPortCount, ipAddress),
				TargetID:    targetForFindings,
				Severity:    models.SeverityInfo,
				FindingType: "port_summary",
				Details: models.JSONB{
					"target":          ipAddress,
					"open_port_count": openPortCount,
					"scan_type":       scanType,
					"ip_address":      ipAddress,
				},
			}
			scanResults.Findings = append(scanResults.Findings, finding)
		} else {
			finding := models.Finding{
				Title:       fmt.Sprintf("No open ports found on %s", ipAddress),
				Description: fmt.Sprintf("Nmap did not discover any open ports on host %s within the specified parameters.", ipAddress),
				Severity:    models.SeverityLow,
				FindingType: "no_open_ports",
				Details: models.JSONB{
					"target":     ipAddress,
					"scan_type":  scanType,
					"port_range": portRange,
					"ip_address": ipAddress,
				},
			}
			scanResults.Findings = append(scanResults.Findings, finding)
		}
	}

	// If we didn't find any hosts in a CIDR range, add a finding about it
	if isCIDR && len(result.Hosts) == 0 {
		finding := models.Finding{
			Title:       fmt.Sprintf("No live hosts found in CIDR range %s", targetValue),
			Description: fmt.Sprintf("Nmap did not discover any live hosts in the CIDR range %s with the current scan parameters.", targetValue),
			Severity:    models.SeverityLow,
			FindingType: "no_live_hosts",
			Details: models.JSONB{
				"target":     targetValue,
				"scan_type":  scanType,
				"port_range": portRange,
			},
		}
		scanResults.Findings = append(scanResults.Findings, finding)
	}

	return scanResults, nil
}

// Type returns the scanner type identifier
func (s *NmapScanner) Type() string {
	return "nmap"
}

// SupportsTargetType indicates whether this scanner can handle the specified target type
func (s *NmapScanner) SupportsTargetType(targetType string) bool {
	switch targetType {
	case models.TargetTypeIP, models.TargetTypeDomain, models.TargetTypeCIDR:
		return true
	default:
		return false
	}
}

// SupportsServices indicates whether this scanner can scan services
func (s *NmapScanner) SupportsServices() bool {
	return false // Nmap scans hosts, not individual services
}

// generatePortDescription creates a human-readable description of a port
func (s *NmapScanner) generatePortDescription(port Port, ipAddress string) string {
	desc := fmt.Sprintf("Port %d/%s is open (%s) on host %s.", port.PortID, port.Protocol, port.State.Reason, ipAddress)

	if port.Service.Name != "" {
		desc += fmt.Sprintf("\nService: %s", port.Service.Name)
	}

	if port.Service.Product != "" {
		desc += fmt.Sprintf("\nProduct: %s", port.Service.Product)
		if port.Service.Version != "" {
			desc += fmt.Sprintf(" %s", port.Service.Version)
		}
	}

	if port.Service.ExtraInfo != "" {
		desc += fmt.Sprintf("\nAdditional info: %s", port.Service.ExtraInfo)
	}

	return desc
}

// generateServiceDescription creates a human-readable description of a service
func (s *NmapScanner) generateServiceDescription(port Port) string {
	desc := fmt.Sprintf("Service detected on port %d/%s.", port.PortID, port.Protocol)

	if port.Service.Product != "" {
		desc += fmt.Sprintf("\nProduct identified as %s", port.Service.Product)
		if port.Service.Version != "" {
			desc += fmt.Sprintf(" version %s", port.Service.Version)
		}
	}

	if port.Service.ExtraInfo != "" {
		desc += fmt.Sprintf("\nAdditional info: %s", port.Service.ExtraInfo)
	}

	return desc
}

// determineSeverityForPort sets severity based on the port and service
func (s *NmapScanner) determineSeverityForPort(port Port) string {
	// Define high-risk ports
	highRiskPorts := map[int]bool{
		21:   true, // FTP
		22:   true, // SSH
		23:   true, // Telnet
		25:   true, // SMTP
		53:   true, // DNS
		110:  true, // POP3
		135:  true, // MSRPC
		137:  true, // NetBIOS
		138:  true, // NetBIOS
		139:  true, // NetBIOS
		445:  true, // SMB
		1433: true, // MSSQL
		1521: true, // Oracle
		3306: true, // MySQL
		3389: true, // RDP
		5432: true, // PostgreSQL
		5900: true, // VNC
		6379: true, // Redis
	}

	// Check if it's a known high-risk port
	if highRiskPorts[port.PortID] {
		return models.SeverityMedium
	}

	// Check for sensitive services regardless of port
	sensitiveServices := map[string]bool{
		"ftp":        true,
		"ssh":        true,
		"telnet":     true,
		"smtp":       true,
		"dns":        true,
		"http":       true,
		"https":      true,
		"mysql":      true,
		"postgresql": true,
		"redis":      true,
		"mongodb":    true,
		"memcached":  true,
		"rdp":        true,
		"vnc":        true,
	}

	if sensitiveServices[strings.ToLower(port.Service.Name)] {
		return models.SeverityMedium
	}

	// Default to low severity for other open ports
	return models.SeverityLow
}
