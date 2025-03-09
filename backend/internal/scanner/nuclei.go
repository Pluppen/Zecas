package scanner

import (
	"backend/internal/models"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// NucleiScanner implements the Scanner interface for performing Nuclei scans
type NucleiScanner struct {
	binPath      string
	templatesDir string
	timeout      int
}

// NucleiResult represents the JSON output from Nuclei
type NucleiResult struct {
	Template         string           `json:"template"`
	TemplateID       string           `json:"template-id"`
	TemplatePath     string           `json:"template-path"`
	Info             NucleiResultInfo `json:"info"`
	Type             string           `json:"type"`
	Host             string           `json:"host"`
	Matcher          string           `json:"matcher_name,omitempty"`
	ExtractedResults []string         `json:"extracted_results,omitempty"`
	IP               string           `json:"ip,omitempty"`
	Timestamp        string           `json:"timestamp"`
	CurlCommand      string           `json:"curl_command,omitempty"`
	MatcherStatus    bool             `json:"matcher_status,omitempty"`
	MatchedAt        string           `json:"matched_at,omitempty"`
	Request          string           `json:"request,omitempty"`
	Response         string           `json:"response,omitempty"`
}

// NucleiResultInfo contains metadata about the vulnerability
type NucleiResultInfo struct {
	Name           string               `json:"name"`
	Author         []string             `json:"author"`
	Tags           []string             `json:"tags"`
	Description    string               `json:"description"`
	Reference      []string             `json:"reference,omitempty"`
	Severity       string               `json:"severity"`
	Classification NucleiClassification `json:"classification,omitempty"`
}

// NucleiClassification contains additional classification information
type NucleiClassification struct {
	CVEIDs     []string `json:"cve-id,omitempty"`
	CVSSScore  string   `json:"cvss-score,omitempty"`
	CVSSVector string   `json:"cvss-metrics,omitempty"`
	CWEIDs     []string `json:"cwe-id,omitempty"`
}

// NewNucleiScanner creates a new Nuclei scanner
func NewNucleiScanner() *NucleiScanner {
	// Default templates directory - can be overridden in scan parameters
	templatesDir := os.Getenv("NUCLEI_TEMPLATES_DIR")
	if templatesDir == "" {
		// Default to ~/.nuclei-templates if not specified
		homeDir, err := os.UserHomeDir()
		if err == nil {
			templatesDir = filepath.Join(homeDir, ".nuclei-templates")
		} else {
			templatesDir = "/opt/nuclei-templates" // Fallback
		}
	}

	return &NucleiScanner{
		binPath:      "nuclei",
		templatesDir: templatesDir,
		timeout:      300, // Default timeout in seconds (5 minutes)
	}
}

// Initialize checks if nuclei is available and templates are installed
func (s *NucleiScanner) Initialize(ctx context.Context) error {
	// Check if nuclei command is available
	cmd := exec.CommandContext(ctx, s.binPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nuclei command not available: %w", err)
	}

	// Ensure we got a proper version output
	if !strings.Contains(string(output), "Nuclei Engine Version") {
		return fmt.Errorf("could not verify nuclei installation")
	}

	// Check if templates directory exists or try to update templates
	if _, err := os.Stat(s.templatesDir); os.IsNotExist(err) {
		log.Println("Nuclei templates directory not found, attempting to update templates...")
		updateCmd := exec.CommandContext(ctx, s.binPath, "-update-templates")
		updateOutput, updateErr := updateCmd.CombinedOutput()
		if updateErr != nil {
			return fmt.Errorf("failed to update nuclei templates: %w (%s)", updateErr, string(updateOutput))
		}
	}

	return nil
}

// ConvertTarget converts a Target to a format suitable for nuclei
func (s *NucleiScanner) ConvertTarget(target models.Target) interface{} {
	return target.Value
}

// ConvertService converts a Service to a format suitable for nuclei
func (s *NucleiScanner) ConvertService(service models.Service) interface{} {
	// For a service, format as URL with protocol, host, and port
	targetDetails, err := s.getServiceTargetDetails(service)
	if err != nil {
		return nil
	}
	return targetDetails
}

// getServiceTargetDetails retrieves the target value and formats it for scanning a service
func (s *NucleiScanner) getServiceTargetDetails(service models.Service) (string, error) {
	// This would normally query the database to get the target details
	// For this implementation, we'll assume the target details are available in the service.RawInfo
	targetValue, ok := service.RawInfo["target_value"].(string)
	if !ok {
		return "", fmt.Errorf("target value not available in service raw info")
	}

	// Format based on service type
	protocol := "http"
	if service.Protocol == "tcp" {
		switch service.ServiceName {
		case "http":
			protocol = "http"
		case "https":
			protocol = "https"
		case "ssh", "ftp", "smtp", "imap", "pop3":
			// These protocols can be scanned directly
			protocol = service.ServiceName
		default:
			// Default to http for unknown services
			protocol = "http"
		}
	}

	// Create URL format for the service
	return fmt.Sprintf("%s://%s:%d", protocol, targetValue, service.Port), nil
}

// Scan performs a nuclei scan against the target
func (s *NucleiScanner) Scan(ctx context.Context, target interface{}, params models.JSONB) (*models.ScanResults, error) {
	targetValue, ok := target.(string)
	if !ok {
		return nil, fmt.Errorf("invalid target format for nuclei scanner")
	}

	scanResults := &models.ScanResults{
		Findings:        []models.Finding{},
		NewTargets:      []models.Target{},
		TargetRelations: []models.TargetRelation{},
		Services:        []models.Service{},
	}

	// Get scan parameters or use defaults
	templateTags := []string{"cve"}              // Default to CVE checks
	templatePaths := []string{}                  // Default to empty (will use template-tags instead)
	templateExclude := []string{"dos"}           // Default to exclude DoS templates
	severity := []string{"medium,high,critical"} // Default severities to scan for
	timeout := s.timeout
	rateLimit := 150               // Default requests per second
	bulkSize := 25                 // Default number of templates to run concurrently
	templatesDir := s.templatesDir // Use default templates dir
	headless := false              // Default to non-headless mode
	includeAll := false            // Don't include all templates by default

	// Override with provided parameters if available
	if val, ok := params["template_tags"].([]interface{}); ok && len(val) > 0 {
		templateTags = []string{}
		for _, tag := range val {
			if tagStr, ok := tag.(string); ok {
				templateTags = append(templateTags, tagStr)
			}
		}
	}

	if val, ok := params["template_paths"].([]interface{}); ok && len(val) > 0 {
		templatePaths = []string{}
		for _, path := range val {
			if pathStr, ok := path.(string); ok {
				templatePaths = append(templatePaths, pathStr)
			}
		}
	}

	if val, ok := params["template_exclude"].([]interface{}); ok && len(val) > 0 {
		templateExclude = []string{}
		for _, exclude := range val {
			if excludeStr, ok := exclude.(string); ok {
				templateExclude = append(templateExclude, excludeStr)
			}
		}
	}

	if val, ok := params["severity"].([]interface{}); ok && len(val) > 0 {
		severity = []string{}
		for _, sev := range val {
			if sevStr, ok := sev.(string); ok {
				severity = append(severity, sevStr)
			}
		}
	}

	if val, ok := params["timeout"].(float64); ok {
		timeout = int(val)
	}

	if val, ok := params["rate_limit"].(float64); ok {
		rateLimit = int(val)
	}

	if val, ok := params["bulk_size"].(float64); ok {
		bulkSize = int(val)
	}

	if val, ok := params["templates_dir"].(string); ok && val != "" {
		templatesDir = val
	}

	if val, ok := params["headless"].(bool); ok {
		headless = val
	}

	if val, ok := params["include_all"].(bool); ok {
		includeAll = val
	}

	// Create a temporary file for output
	outputFile, err := os.CreateTemp("", "nuclei-output-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary output file: %w", err)
	}
	defer os.Remove(outputFile.Name())
	defer outputFile.Close()

	// Construct the nuclei command
	args := []string{
		"-target", targetValue,
		"-j",                    // Output in JSON format
		"-o", outputFile.Name(), // Output file
		"-timeout", fmt.Sprintf("%d", timeout),
		"-rate-limit", fmt.Sprintf("%d", rateLimit),
		"-bulk-size", fmt.Sprintf("%d", bulkSize),
	}

	// Add template paths if provided
	for _, path := range templatePaths {
		args = append(args, "-t", path)
	}

	// Add template tags if provided and no specific templates
	if len(templatePaths) == 0 && len(templateTags) > 0 {
		args = append(args, "-tags", strings.Join(templateTags, ","))
	}

	// Add template exclusions
	if len(templateExclude) > 0 {
		args = append(args, "-exclude-tags", strings.Join(templateExclude, ","))
	}

	// Add severity levels
	if len(severity) > 0 {
		args = append(args, "-severity", strings.Join(severity, ","))
	}

	// Set templates directory if different from default
	if templatesDir != s.templatesDir {
		args = append(args, "-templates-dir", templatesDir)
	}

	// Enable headless scanning if requested
	if headless {
		args = append(args, "-headless")
	}

	// Include all templates if requested (overrides tags and severity)
	if includeAll {
		args = append(args, "-include-all")
	}

	// Run nuclei with a timeout
	cmd := exec.CommandContext(ctx, s.binPath, args...)
	cmd.Stderr = os.Stderr

	// Create a pipe for real-time processing
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start nuclei scan: %w", err)
	}

	// Start a goroutine to log progress
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Printf("Nuclei progress: %s", scanner.Text())
		}
	}()

	// Wait for the command to complete
	err = cmd.Wait()
	// We don't treat this as an error since nuclei will exit with status 1 if vulnerabilities are found
	if err != nil {
		log.Printf("Nuclei scan completed with status: %v", err)
	}

	// Parse the results
	findings, err := s.parseNucleiOutput(outputFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to parse nuclei output: %w", err)
	}

	// Process findings
	scanResults.Findings = findings

	// Return the scan results
	return scanResults, nil
}

// parseNucleiOutput reads the Nuclei JSON output and converts it to findings
func (s *NucleiScanner) parseNucleiOutput(outputFile string) ([]models.Finding, error) {
	var findings []models.Finding

	// Open the output file
	file, err := os.Open(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open nuclei output file: %w", err)
	}
	defer file.Close()

	// Create a scanner to read line by line (Nuclei outputs each result as a separate JSON object)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var result NucleiResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			log.Printf("Failed to parse nuclei result: %v", err)
			continue
		}

		// Create a finding from the result
		finding := s.createFindingFromNucleiResult(result)
		findings = append(findings, finding)
	}

	if err := scanner.Err(); err != nil {
		return findings, fmt.Errorf("error reading nuclei output file: %w", err)
	}

	return findings, nil
}

// createFindingFromNucleiResult converts a Nuclei result to a finding
func (s *NucleiScanner) createFindingFromNucleiResult(result NucleiResult) models.Finding {
	// Determine severity based on Nuclei's severity
	severity := models.SeverityMedium // Default
	switch strings.ToLower(result.Info.Severity) {
	case "critical":
		severity = models.SeverityCritical
	case "high":
		severity = models.SeverityHigh
	case "medium":
		severity = models.SeverityMedium
	case "low":
		severity = models.SeverityLow
	case "info":
		severity = models.SeverityInfo
	}

	// Format the description
	description := result.Info.Description
	if description == "" {
		description = fmt.Sprintf("Nuclei found a %s issue using template %s", result.Info.Severity, result.Template)
	}

	// Add references if available
	if len(result.Info.Reference) > 0 {
		description += "\n\nReferences:\n"
		for _, ref := range result.Info.Reference {
			description += fmt.Sprintf("- %s\n", ref)
		}
	}

	// Add CVSS information if available
	if result.Info.Classification.CVSSScore != "" {
		description += fmt.Sprintf("\nCVSS Score: %s", result.Info.Classification.CVSSScore)
		if result.Info.Classification.CVSSVector != "" {
			description += fmt.Sprintf("\nCVSS Vector: %s", result.Info.Classification.CVSSVector)
		}
	}

	// Add matched data if available
	if result.MatchedAt != "" {
		description += fmt.Sprintf("\n\nMatched at: %s", result.MatchedAt)
	}

	// Add extracted results if available
	if len(result.ExtractedResults) > 0 {
		description += "\n\nExtracted data:\n"
		for _, extracted := range result.ExtractedResults {
			description += fmt.Sprintf("- %s\n", extracted)
		}
	}

	// Create metadata for the details field
	details := models.JSONB{
		"template_id":   result.TemplateID,
		"template_name": result.Template,
		"host":          result.Host,
		"matcher":       result.Matcher,
		"timestamp":     result.Timestamp,
		"type":          result.Type,
	}

	// Add CVE IDs if available
	if len(result.Info.Classification.CVEIDs) > 0 {
		details["cves"] = result.Info.Classification.CVEIDs
	}

	// Add CWE IDs if available
	if len(result.Info.Classification.CWEIDs) > 0 {
		details["cwes"] = result.Info.Classification.CWEIDs
	}

	// Add CVSS information if available
	if result.Info.Classification.CVSSScore != "" {
		details["cvss_score"] = result.Info.Classification.CVSSScore
	}
	if result.Info.Classification.CVSSVector != "" {
		details["cvss_vector"] = result.Info.Classification.CVSSVector
	}

	// Add tags if available
	if len(result.Info.Tags) > 0 {
		details["tags"] = result.Info.Tags
	}

	// Add curl command if available (useful for reproduction)
	if result.CurlCommand != "" {
		details["curl_command"] = result.CurlCommand
	}

	// Add request/response if available (can be helpful for verification)
	if result.Request != "" {
		details["request"] = result.Request
	}
	if result.Response != "" {
		details["response"] = result.Response
	}

	// Create the finding
	finding := models.Finding{
		Title:        result.Info.Name,
		Description:  description,
		Severity:     severity,
		FindingType:  s.determineFindingType(result),
		Details:      details,
		DiscoveredAt: time.Now(),
		Verified:     false, // Requires manual verification
		Fixed:        false,
		Manual:       false,
	}

	return finding
}

// determineFindingType categorizes the finding based on the template tags and info
func (s *NucleiScanner) determineFindingType(result NucleiResult) string {
	// Check if it's a CVE
	for _, cve := range result.Info.Classification.CVEIDs {
		if strings.HasPrefix(cve, "CVE-") {
			return "vulnerability"
		}
	}

	// Check tags for categorization
	for _, tag := range result.Info.Tags {
		switch strings.ToLower(tag) {
		case "cve", "vulnerability", "vuln":
			return "vulnerability"
		case "config", "misconfig", "misconfiguration":
			return "misconfiguration"
		case "exposure", "exposed", "disclosure":
			return "information_disclosure"
		case "intrusion", "takeover", "account":
			return "account_takeover"
		case "sqli", "xss", "ssrf", "csrf", "injection":
			return "injection"
		case "default", "default-login", "weak-password":
			return "default_credential"
		case "tech", "technology", "stack":
			return "technology_detection"
		}
	}

	// Default based on scan type
	switch strings.ToLower(result.Type) {
	case "http":
		return "web_vulnerability"
	case "dns":
		return "dns_issue"
	case "file":
		return "file_vulnerability"
	case "network":
		return "network_vulnerability"
	case "headless":
		return "browser_vulnerability"
	default:
		return "security_issue"
	}
}

// Type returns the scanner type identifier
func (s *NucleiScanner) Type() string {
	return "nuclei"
}

// SupportsTargetType indicates whether this scanner can handle the specified target type
func (s *NucleiScanner) SupportsTargetType(targetType string) bool {
	switch targetType {
	case models.TargetTypeDomain, models.TargetTypeIP:
		return true
	case models.TargetTypeCIDR:
		// Nuclei supports CIDR notation
		return true
	default:
		return false
	}
}

// SupportsServices indicates whether this scanner can scan services
func (s *NucleiScanner) SupportsServices() bool {
	return true
}
