// internal/scanner/subdomain.go
package scanner

// TODO: Write down a spec on what we want out of this scanner, what types of assets shoudl it create
// and also if this scanner should be able to generate findings and
// what is counted as a finding from this tool.

import (
	"backend/internal/models"
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SubdomainScanner implements the Scanner interface for discovering subdomains
type SubdomainScanner struct {
	binPath         string
	wordlistPath    string
	resolverTimeout int
}

// NewSubdomainScanner creates a new subdomain scanner
func NewSubdomainScanner() *SubdomainScanner {
	return &SubdomainScanner{
		binPath:         "subfinder",      // Default tool, subfinder
		wordlistPath:    "subdomains.txt", // Default wordlist
		resolverTimeout: 5,                // Default timeout in seconds
	}
}

// Initialize checks if the scanner tools are available
func (s *SubdomainScanner) Initialize(ctx context.Context) error {
	// Check if subfinder is available
	cmd := exec.CommandContext(ctx, s.binPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try to fall back to alternative tools or methods
		// For example, we could check if we have basic DNS tools
		cmd = exec.CommandContext(ctx, "dig", "-version")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("subdomain scanning tools not available: %w", err)
		}
		// Set binPath to use dig-based approach
		s.binPath = "dig"
	} else if !strings.Contains(string(output), "subfinder") {
		return fmt.Errorf("command '%s' does not appear to be subfinder", s.binPath)
	}

	if _, err := os.Stat(s.wordlistPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("wordlist %s does not exists", s.wordlistPath)
	}

	return nil
}

// ConvertTarget converts a Target to a format suitable for subdomain scanning
func (s *SubdomainScanner) ConvertTarget(target models.Target) interface{} {
	// Only domain targets are supported
	if target.TargetType == models.TargetTypeDomain {
		return target.Value
	}
	return nil
}

// ConvertService returns nil since subdomain scanner doesn't scan services
func (s *SubdomainScanner) ConvertService(service models.Service) interface{} {
	return nil
}

// Scan performs subdomain enumeration against the target domain
func (s *SubdomainScanner) Scan(ctx context.Context, target interface{}, params models.JSONB) (*models.ScanResults, error) {
	if target == nil {
		return nil, fmt.Errorf("invalid target for subdomain scanner")
	}

	domain := target.(string)
	scanResults := &models.ScanResults{
		Findings:        []models.Finding{},
		NewTargets:      []models.Target{},
		TargetRelations: []models.TargetRelation{},
		Services:        []models.Service{},
	}

	// Configure scan parameters
	recursive := false
	resolveIP := true
	wordlist := s.wordlistPath
	//timeout := s.resolverTimeout

	// Override with provided parameters if available
	if val, ok := params["recursive"].(bool); ok {
		recursive = val
	}
	if val, ok := params["resolve_ip"].(bool); ok {
		resolveIP = val
	}
	if val, ok := params["wordlist"].(string); ok && val != "" {
		wordlist = val
	}
	if _, ok := params["timeout"].(float64); ok {
		//timeout = int(val)
	}

	var subdomains []string
	var err error

	// Use the appropriate method based on available tools
	if s.binPath == "subfinder" {
		subdomains, err = s.runSubfinder(ctx, domain, recursive)
	} else {
		// Fallback to wordlist-based approach using basic DNS tools
		subdomains, err = s.runWordlistEnumeration(ctx, domain, wordlist)
	}

	if err != nil {
		return nil, fmt.Errorf("subdomain scan failed: %w", err)
	}

	// Process discovered subdomains
	discoveredCount := 0
	resolvedCount := 0

	for _, subdomain := range subdomains {
		subdomain = strings.TrimSpace(subdomain)
		if subdomain == "" || subdomain == domain {
			continue
		}

		discoveredCount++

		// Create new target for the subdomain
		newTarget := models.Target{
			ID:         uuid.New(),
			ProjectID:  uuid.Nil,                // Will be set by worker
			TargetType: models.TargetTypeDomain, // Use domain type for subdomains
			Value:      subdomain,
			Metadata: models.JSONB{
				"discovered_from": domain,
				"discovery_scan":  "subdomain",
				"discovered_at":   time.Now().Format(time.RFC3339),
			},
		}
		scanResults.NewTargets = append(scanResults.NewTargets, newTarget)

		// Create relation between domain and subdomain
		relation := models.TargetRelation{
			ID:            uuid.New(),
			SourceID:      uuid.Nil, // Will be set by worker
			DestinationID: newTarget.ID,
			RelationType:  models.RelationParentOf,
			Metadata: models.JSONB{
				"discovered_at": time.Now().Format(time.RFC3339),
			},
		}
		scanResults.TargetRelations = append(scanResults.TargetRelations, relation)

		// Resolve IP addresses if requested
		if resolveIP {
			ips, err := net.LookupHost(subdomain)
			if err == nil && len(ips) > 0 {
				resolvedCount++

				// Add each IP as a target and create relations
				for _, ip := range ips {
					// Create new target for each IP
					ipTarget := models.Target{
						ID:         uuid.New(),
						ProjectID:  uuid.Nil, // Will be set by worker
						TargetType: models.TargetTypeIP,
						Value:      ip,
						Metadata: models.JSONB{
							"discovered_from": subdomain,
							"discovery_scan":  "subdomain_resolution",
							"discovered_at":   time.Now().Format(time.RFC3339),
						},
					}
					scanResults.NewTargets = append(scanResults.NewTargets, ipTarget)

					// Create relation between subdomain and IP
					ipRelation := models.TargetRelation{
						ID:            uuid.New(),
						SourceID:      newTarget.ID,
						DestinationID: ipTarget.ID,
						RelationType:  models.RelationResolvesTo,
						Metadata: models.JSONB{
							"discovered_at": time.Now().Format(time.RFC3339),
						},
					}
					scanResults.TargetRelations = append(scanResults.TargetRelations, ipRelation)
				}

				// Create finding for each resolved subdomain
				finding := models.Finding{
					Title:       fmt.Sprintf("Subdomain discovered: %s", subdomain),
					Description: fmt.Sprintf("Subdomain %s resolves to %s", subdomain, strings.Join(ips, ", ")),
					Severity:    models.SeverityInfo,
					FindingType: "subdomain_discovered",
					Details: models.JSONB{
						"parent_domain": domain,
						"subdomain":     subdomain,
						"ip_addresses":  ips,
					},
				}
				scanResults.Findings = append(scanResults.Findings, finding)
			} else {
				// Create finding for unresolved subdomain
				finding := models.Finding{
					Title:       fmt.Sprintf("Unresolved subdomain discovered: %s", subdomain),
					Description: fmt.Sprintf("Subdomain %s was discovered but does not resolve to an IP address", subdomain),
					Severity:    models.SeverityLow,
					FindingType: "subdomain_unresolved",
					Details: models.JSONB{
						"parent_domain": domain,
						"subdomain":     subdomain,
					},
				}
				scanResults.Findings = append(scanResults.Findings, finding)
			}
		} else {
			// Create basic finding for discovered subdomain
			finding := models.Finding{
				Title:       fmt.Sprintf("Subdomain discovered: %s", subdomain),
				Description: fmt.Sprintf("Subdomain %s was discovered for parent domain %s", subdomain, domain),
				Severity:    models.SeverityInfo,
				FindingType: "subdomain_discovered",
				Details: models.JSONB{
					"parent_domain": domain,
					"subdomain":     subdomain,
				},
			}
			scanResults.Findings = append(scanResults.Findings, finding)
		}
	}

	// Create summary finding
	finding := models.Finding{
		Title:       fmt.Sprintf("Subdomain enumeration for %s", domain),
		Description: fmt.Sprintf("Discovered %d subdomains for %s, %d of which resolved to IP addresses", discoveredCount, domain, resolvedCount),
		Severity:    models.SeverityInfo,
		FindingType: "subdomain_summary",
		Details: models.JSONB{
			"domain":          domain,
			"total_found":     discoveredCount,
			"total_resolved":  resolvedCount,
			"scan_parameters": params,
		},
	}
	scanResults.Findings = append(scanResults.Findings, finding)

	return scanResults, nil
}

// runSubfinder runs the subfinder tool to enumerate subdomains
func (s *SubdomainScanner) runSubfinder(ctx context.Context, domain string, recursive bool) ([]string, error) {
	args := []string{"-d", domain, "-silent"}

	if recursive {
		args = append(args, "-recursive")
	}

	cmd := exec.CommandContext(ctx, s.binPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("subfinder execution failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var results []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			results = append(results, line)
		}
	}

	return results, nil
}

// runWordlistEnumeration uses a wordlist and basic DNS tools to find subdomains
func (s *SubdomainScanner) runWordlistEnumeration(ctx context.Context, domain string, wordlistPath string) ([]string, error) {
	// Try to read wordlist
	catCmd := exec.CommandContext(ctx, "cat", wordlistPath)
	catOutput, err := catCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to setup wordlist pipe: %w", err)
	}

	if err := catCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start wordlist reader: %w", err)
	}

	var results []string

	// Read wordlist line by line and try to resolve each potential subdomain
	scanner := bufio.NewScanner(catOutput)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" || strings.HasPrefix(word, "#") {
			continue
		}

		subdomain := fmt.Sprintf("%s.%s", word, domain)

		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
			// Try to resolve the subdomain
			_, err := net.LookupHost(subdomain)
			if err == nil {
				results = append(results, subdomain)
			}
		}
	}

	if err := catCmd.Wait(); err != nil {
		return nil, fmt.Errorf("wordlist reader failed: %w", err)
	}

	return results, nil
}

// Type returns the scanner type identifier
func (s *SubdomainScanner) Type() string {
	return "subdomain"
}

// SupportsTargetType indicates whether this scanner can handle the specified target type
func (s *SubdomainScanner) SupportsTargetType(targetType string) bool {
	return targetType == models.TargetTypeDomain
}

// SupportsServices indicates whether this scanner can scan services
func (s *SubdomainScanner) SupportsServices() bool {
	return false
}
