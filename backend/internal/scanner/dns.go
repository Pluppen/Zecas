// internal/scanner/dns.go
package scanner

import (
	"backend/internal/models"
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DNSScanner implements the Scanner interface for resolving DNS records
type DNSScanner struct {
	resolverTimeout int
}

// DNSRecordType defines the types of DNS records to check
type DNSRecordType string

// DNS record types
const (
	RecordTypeA     DNSRecordType = "A"
	RecordTypeAAAA  DNSRecordType = "AAAA"
	RecordTypeCNAME DNSRecordType = "CNAME"
	RecordTypeMX    DNSRecordType = "MX"
	RecordTypeTXT   DNSRecordType = "TXT"
	RecordTypeNS    DNSRecordType = "NS"
	RecordTypeSOA   DNSRecordType = "SOA"
	RecordTypePTR   DNSRecordType = "PTR"
)

// NewDNSScanner creates a new DNS resolver scanner
func NewDNSScanner() *DNSScanner {
	return &DNSScanner{
		resolverTimeout: 5, // Default timeout in seconds
	}
}

// Initialize checks if DNS resolution is available
func (s *DNSScanner) Initialize(ctx context.Context) error {
	// Try to resolve a known domain to check if DNS resolution works
	_, err := net.LookupHost("example.com")
	if err != nil {
		// Try an alternative check with the dig command
		cmd := exec.CommandContext(ctx, "dig", "+short", "example.com")
		output, err := cmd.CombinedOutput()
		if err != nil || strings.TrimSpace(string(output)) == "" {
			return fmt.Errorf("DNS resolution not available: %w", err)
		}
	}

	return nil
}

// ConvertTarget converts a Target to a format suitable for DNS resolution
func (s *DNSScanner) ConvertTarget(target models.Target) interface{} {
	switch target.TargetType {
	case models.TargetTypeDomain:
		return target.Value
	case models.TargetTypeIP:
		// For IP targets, we'll do reverse DNS lookups
		return target.Value
	default:
		return nil
	}
}

// ConvertService returns nil since DNS scanner doesn't scan services
func (s *DNSScanner) ConvertService(service models.Service) interface{} {
	return nil
}

// Scan performs DNS resolution against the target
func (s *DNSScanner) Scan(ctx context.Context, target interface{}, params models.JSONB) (*models.ScanResults, error) {
	if target == nil {
		return nil, fmt.Errorf("invalid target for DNS resolver")
	}

	targetValue := target.(string)
	scanResults := &models.ScanResults{
		Findings:        []models.Finding{},
		NewTargets:      []models.Target{},
		TargetRelations: []models.TargetRelation{},
		Services:        []models.Service{},
	}

	// Configure which record types to query
	recordTypes := []DNSRecordType{
		RecordTypeA,
		RecordTypeAAAA,
		RecordTypeCNAME,
		RecordTypeMX,
		RecordTypeTXT,
		RecordTypeNS,
	}

	// Check if this is an IP for reverse lookup
	isIP := net.ParseIP(targetValue) != nil

	// If specific record types are requested
	if types, ok := params["record_types"].([]interface{}); ok && len(types) > 0 {
		recordTypes = []DNSRecordType{}
		for _, t := range types {
			if typeStr, ok := t.(string); ok {
				recordTypes = append(recordTypes, DNSRecordType(typeStr))
			}
		}
	}

	var dnsRecords []models.DNSRecord
	var findings []models.Finding

	if isIP {
		// Perform reverse DNS lookup
		names, err := net.LookupAddr(targetValue)
		if err != nil {
			finding := models.Finding{
				Title:       fmt.Sprintf("No reverse DNS for %s", targetValue),
				Description: fmt.Sprintf("No PTR records found for IP %s", targetValue),
				Severity:    models.SeverityLow,
				FindingType: "dns_no_ptr",
				Details: models.JSONB{
					"ip":    targetValue,
					"error": err.Error(),
				},
			}
			findings = append(findings, finding)
		} else {
			// Create targets for PTR records
			for _, name := range names {
				// Standardize name format (remove trailing dot)
				name = strings.TrimSuffix(name, ".")

				newTarget := models.Target{
					ID:         uuid.New(),
					ProjectID:  uuid.Nil, // Will be set by worker
					TargetType: models.TargetTypeDomain,
					Value:      name,
					Metadata: models.JSONB{
						"discovered_from": targetValue,
						"discovery_scan":  "dns_ptr",
						"discovered_at":   time.Now().Format(time.RFC3339),
					},
				}
				scanResults.NewTargets = append(scanResults.NewTargets, newTarget)

				// Create relation between IP and hostname
				relation := models.TargetRelation{
					ID:            uuid.New(),
					SourceID:      uuid.Nil, // Will be set by worker
					DestinationID: newTarget.ID,
					RelationType:  models.RelationResolvesTo,
					Metadata: models.JSONB{
						"discovered_at": time.Now().Format(time.RFC3339),
						"record_type":   "PTR",
					},
				}
				scanResults.TargetRelations = append(scanResults.TargetRelations, relation)
			}

			finding := models.Finding{
				Title:       fmt.Sprintf("Reverse DNS for %s", targetValue),
				Description: fmt.Sprintf("IP %s resolves to %s", targetValue, strings.Join(names, ", ")),
				Severity:    models.SeverityInfo,
				FindingType: "dns_ptr",
				Details: models.JSONB{
					"ip":        targetValue,
					"hostnames": names,
				},
			}
			findings = append(findings, finding)
		}
	} else {
		// Process each requested record type
		for _, recordType := range recordTypes {
			var records []string
			var err error

			// Lookup appropriate record type
			switch recordType {
			case RecordTypeA:
				ips, err := net.LookupHost(targetValue)
				if err == nil {
					records = ips

					// Create targets for A records (IPs)
					for _, ip := range ips {
						// Only create targets for IPv4 addresses
						if net.ParseIP(ip).To4() != nil {
							newTarget := models.Target{
								ID:         uuid.New(),
								ProjectID:  uuid.Nil, // Will be set by worker
								TargetType: models.TargetTypeIP,
								Value:      ip,
								Metadata: models.JSONB{
									"discovered_from": targetValue,
									"discovery_scan":  "dns_a",
									"discovered_at":   time.Now().Format(time.RFC3339),
								},
							}
							scanResults.NewTargets = append(scanResults.NewTargets, newTarget)

							// Create relation between domain and IP
							relation := models.TargetRelation{
								ID:            uuid.New(),
								SourceID:      uuid.Nil, // Will be set by worker
								DestinationID: newTarget.ID,
								RelationType:  models.RelationResolvesTo,
								Metadata: models.JSONB{
									"discovered_at": time.Now().Format(time.RFC3339),
									"record_type":   "A",
								},
							}
							scanResults.TargetRelations = append(scanResults.TargetRelations, relation)

							dnsRecord := models.DNSRecord{
								RecordType:  (string)(recordType),
								RecordValue: ip,
							}
							dnsRecords = append(dnsRecords, dnsRecord)
						}
					}
				}
			case RecordTypeAAAA:
				// Filter only IPv6 addresses
				ips, err := net.LookupIP(targetValue)
				if err == nil {

					for _, ip := range ips {
						if ip.To4() == nil { // IPv6 addresses have To4() == nil
							dnsRecord := models.DNSRecord{
								RecordType:  (string)(recordType),
								RecordValue: ip.String(),
							}
							dnsRecords = append(dnsRecords, dnsRecord)
							records = append(records, ip.String())
						}
					}
				}
			case RecordTypeCNAME:
				cname, err := net.LookupCNAME(targetValue)
				if err == nil && cname != "" {
					cname = strings.TrimSuffix(cname, ".")
					records = []string{cname}

					// Create target for CNAME
					newTarget := models.Target{
						ID:         uuid.New(),
						ProjectID:  uuid.Nil, // Will be set by worker
						TargetType: models.TargetTypeDomain,
						Value:      cname,
						Metadata: models.JSONB{
							"discovered_from": targetValue,
							"discovery_scan":  "dns_cname",
							"discovered_at":   time.Now().Format(time.RFC3339),
						},
					}
					scanResults.NewTargets = append(scanResults.NewTargets, newTarget)

					// Create relation between domain and CNAME
					relation := models.TargetRelation{
						ID:            uuid.New(),
						SourceID:      uuid.Nil, // Will be set by worker
						DestinationID: newTarget.ID,
						RelationType:  models.RelationResolvesTo,
						Metadata: models.JSONB{
							"discovered_at": time.Now().Format(time.RFC3339),
							"record_type":   "CNAME",
						},
					}
					scanResults.TargetRelations = append(scanResults.TargetRelations, relation)

					dnsRecord := models.DNSRecord{
						RecordType:  (string)(recordType),
						RecordValue: cname,
					}
					dnsRecords = append(dnsRecords, dnsRecord)
				}
			case RecordTypeMX:
				mxs, err := net.LookupMX(targetValue)
				if err == nil {
					for _, mx := range mxs {
						mxHost := strings.TrimSuffix(mx.Host, ".")
						records = append(records, fmt.Sprintf("%s (priority: %d)", mxHost, mx.Pref))

						// Create target for MX hostname
						newTarget := models.Target{
							ID:         uuid.New(),
							ProjectID:  uuid.Nil, // Will be set by worker
							TargetType: models.TargetTypeDomain,
							Value:      mxHost,
							Metadata: models.JSONB{
								"discovered_from": targetValue,
								"discovery_scan":  "dns_mx",
								"discovered_at":   time.Now().Format(time.RFC3339),
								"mx_priority":     mx.Pref,
							},
						}
						scanResults.NewTargets = append(scanResults.NewTargets, newTarget)

						// Create relation between domain and MX host
						relation := models.TargetRelation{
							ID:            uuid.New(),
							SourceID:      uuid.Nil, // Will be set by worker
							DestinationID: newTarget.ID,
							RelationType:  models.RelationResolvesTo,
							Metadata: models.JSONB{
								"discovered_at": time.Now().Format(time.RFC3339),
								"record_type":   "MX",
								"priority":      mx.Pref,
							},
						}
						scanResults.TargetRelations = append(scanResults.TargetRelations, relation)

						dnsRecord := models.DNSRecord{
							RecordType:  (string)(recordType),
							RecordValue: fmt.Sprintf("%s (priority: %d)", mxHost, mx.Pref),
						}
						dnsRecords = append(dnsRecords, dnsRecord)
					}
				}
			case RecordTypeTXT:
				txts, err := net.LookupTXT(targetValue)
				if err == nil {
					records = txts
					for _, txt := range txts {
						dnsRecord := models.DNSRecord{
							RecordType:  (string)(recordType),
							RecordValue: txt,
						}
						dnsRecords = append(dnsRecords, dnsRecord)
					}
				}
			case RecordTypeNS:
				nss, err := net.LookupNS(targetValue)
				if err == nil {
					for _, ns := range nss {
						nsHost := strings.TrimSuffix(ns.Host, ".")
						records = append(records, nsHost)

						// Create target for NS hostname
						newTarget := models.Target{
							ID:         uuid.New(),
							ProjectID:  uuid.Nil, // Will be set by worker
							TargetType: models.TargetTypeDomain,
							Value:      nsHost,
							Metadata: models.JSONB{
								"discovered_from": targetValue,
								"discovery_scan":  "dns_ns",
								"discovered_at":   time.Now().Format(time.RFC3339),
							},
						}
						scanResults.NewTargets = append(scanResults.NewTargets, newTarget)

						// Create relation between domain and NS host
						relation := models.TargetRelation{
							ID:            uuid.New(),
							SourceID:      uuid.Nil, // Will be set by worker
							DestinationID: newTarget.ID,
							RelationType:  models.RelationResolvesTo,
							Metadata: models.JSONB{
								"discovered_at": time.Now().Format(time.RFC3339),
								"record_type":   "NS",
							},
						}
						scanResults.TargetRelations = append(scanResults.TargetRelations, relation)

						dnsRecord := models.DNSRecord{
							RecordType:  (string)(recordType),
							RecordValue: nsHost,
						}
						dnsRecords = append(dnsRecords, dnsRecord)
					}
				}
			case RecordTypeSOA:
				// SOA records need to be queried using a specific DNS tool as Go's net package doesn't support it directly
				cmd := exec.CommandContext(ctx, "dig", "+short", "SOA", targetValue)
				output, err := cmd.CombinedOutput()
				if err == nil {
					outputStr := strings.TrimSpace(string(output))
					if outputStr != "" {
						records = []string{outputStr}
						dnsRecord := models.DNSRecord{
							RecordType:  (string)(recordType),
							RecordValue: outputStr,
						}
						dnsRecords = append(dnsRecords, dnsRecord)
					}
				}
			}

			if err != nil || len(records) == 0 {
				// No records found
				finding := models.Finding{
					Title:       fmt.Sprintf("No %s records for %s", recordType, targetValue),
					Description: fmt.Sprintf("No %s DNS records were found for %s", recordType, targetValue),
					Severity:    models.SeverityInfo,
					FindingType: "dns_no_records",
					Details: models.JSONB{
						"domain":      targetValue,
						"record_type": string(recordType),
						"error":       err != nil,
					},
				}
				findings = append(findings, finding)
			} else {
				// Records found
				finding := models.Finding{
					Title:       fmt.Sprintf("%s records for %s", recordType, targetValue),
					Description: s.generateRecordDescription(targetValue, recordType, records),
					Severity:    s.determineSeverityForRecords(recordType, records),
					FindingType: "dns_records",
					Details: models.JSONB{
						"domain":      targetValue,
						"record_type": string(recordType),
						"records":     records,
					},
				}
				findings = append(findings, finding)
			}
		}
	}
	scanResults.Findings = findings
	scanResults.DNSRecords = dnsRecords
	return scanResults, nil
}

// generateRecordDescription creates a human-readable description of DNS records
func (s *DNSScanner) generateRecordDescription(domain string, recordType DNSRecordType, records []string) string {
	desc := fmt.Sprintf("The following %s records were found for %s:", recordType, domain)

	for _, record := range records {
		desc += fmt.Sprintf("\n• %s", record)
	}

	switch recordType {
	case RecordTypeA:
		desc += "\n\nThese IP addresses are the direct hosts for this domain."
	case RecordTypeAAAA:
		desc += "\n\nThese are IPv6 addresses for this domain."
	case RecordTypeCNAME:
		desc += "\n\nThis domain is an alias pointing to another canonical name."
	case RecordTypeMX:
		desc += "\n\nThese servers handle email for this domain (lower priority values are preferred)."
	case RecordTypeTXT:
		desc += "\n\nThese text records may contain SPF, DKIM, or other domain verification information."
	case RecordTypeNS:
		desc += "\n\nThese nameservers are authoritative for this domain."
	case RecordTypePTR:
		desc += "\n\nThis IP address has reverse DNS pointing to these hostnames."
	}

	return desc
}

// generateSummaryDescription creates a summary of all DNS findings
func (s *DNSScanner) generateSummaryDescription(target string, records map[DNSRecordType][]string, isIP bool) string {
	if isIP {
		if ptr, ok := records[RecordTypePTR]; ok && len(ptr) > 0 {
			return fmt.Sprintf("IP %s has reverse DNS records pointing to %s", target, strings.Join(ptr, ", "))
		}
		return fmt.Sprintf("IP %s has no reverse DNS records", target)
	}

	desc := fmt.Sprintf("DNS resolution results for %s:", target)

	for recordType, typeRecords := range records {
		if len(typeRecords) > 0 {
			desc += fmt.Sprintf("\n\n%s Records:", recordType)
			for _, record := range typeRecords {
				desc += fmt.Sprintf("\n• %s", record)
			}
		}
	}

	if len(records) == 0 {
		desc += "\n\nNo DNS records were found for this target."
	}

	return desc
}

// determineSeverityForRecords sets severity based on record type and content
func (s *DNSScanner) determineSeverityForRecords(recordType DNSRecordType, records []string) string {
	// Most DNS records are informational
	severity := models.SeverityInfo

	// Check for specific security-related records
	if recordType == RecordTypeTXT {
		for _, record := range records {
			// Missing SPF record might be a security issue
			if strings.Contains(record, "v=spf1") {
				if strings.Contains(record, "all") {
					severity = models.SeverityInfo // Has SPF record with proper configuration
				} else {
					severity = models.SeverityLow // Has SPF but might be misconfigured
				}
			}

			// DMARC record
			if strings.Contains(record, "v=DMARC1") {
				severity = models.SeverityInfo
			}
		}
	}

	return severity
}

// Type returns the scanner type identifier
func (s *DNSScanner) Type() string {
	return "dns"
}

// SupportsTargetType indicates whether this scanner can handle the specified target type
func (s *DNSScanner) SupportsTargetType(targetType string) bool {
	switch targetType {
	case models.TargetTypeDomain, models.TargetTypeIP:
		return true
	default:
		return false
	}
}

// SupportsServices indicates whether this scanner can scan services
func (s *DNSScanner) SupportsServices() bool {
	return false
}
