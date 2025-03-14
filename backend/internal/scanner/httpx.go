// internal/scanner/httpx.go
package scanner

import (
	"backend/internal/models"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TODO: Write down a spec on what we want out of this scanner, what types of assets shoudl it create
// and also if this scanner should be able to generate findings and
// what is counted as a finding from this tool.

// NOTES:
// Identified web apps or generally all apps running on this target or service
// generally because of different applicaton libarys being used and fraemworks,,,

// HTTPXScanner implements the Scanner interface for performing HTTPX scans
type HTTPXScanner struct {
	binPath string
	timeout int
}

// HTTPXResult represents the JSON output from HTTPX
type HTTPXResult struct {
	Timestamp          string            `json:"timestamp"`
	URL                string            `json:"url"`
	Input              string            `json:"input"`
	StatusCode         int               `json:"status_code"`
	Title              string            `json:"title"`
	WebServer          string            `json:"webserver"`
	ContentType        string            `json:"content_type"`
	Method             string            `json:"method"`
	Host               string            `json:"host"`
	ContentLength      int               `json:"content_length"`
	TLSData            *TLSData          `json:"tls,omitempty"`
	ResponseTime       string            `json:"time"`
	Technologies       []string          `json:"technologies,omitempty"`
	ServerPort         string            `json:"port"`
	Chain              []string          `json:"chain,omitempty"`
	Lines              int               `json:"lines,omitempty"`
	Words              int               `json:"words,omitempty"`
	FaviconMmh3        string            `json:"favicon_mmh3,omitempty"`
	FaviconMD5         string            `json:"favicon_md5,omitempty"`
	ResponseHash       string            `json:"hash,omitempty"`
	A                  []string          `json:"a,omitempty"`
	CNAMEs             []string          `json:"cname,omitempty"`
	CDNName            string            `json:"cdn_name,omitempty"`
	HTTP2              bool              `json:"http2,omitempty"`
	Pipeline           bool              `json:"pipeline,omitempty"`
	SecurityHeaders    map[string]string `json:"security_headers,omitempty"`
	WebApplicationInfo map[string]string `json:"webapplication_info,omitempty"`
}

// TLSData contains TLS/SSL certificate information
type TLSData struct {
	Cipher                string `json:"cipher"`
	Version               string `json:"version"`
	CertificateNotBefore  string `json:"certificate_not_before"`
	CertificateNotAfter   string `json:"certificate_not_after"`
	Expired               bool   `json:"expired"`
	CertificateIssuer     string `json:"certificate_issuer"`
	CertificateSubject    string `json:"certificate_subject"`
	CertificateCommonName string `json:"certificate_common_name"`
	JarmHash              string `json:"jarm,omitempty"`
}

// NewHTTPXScanner creates a new HTTPX scanner
func NewHTTPXScanner() *HTTPXScanner {
	return &HTTPXScanner{
		binPath: "httpx",
		timeout: 30, // Default timeout in seconds
	}
}

// Initialize checks if httpx is available
func (s *HTTPXScanner) Initialize(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, s.binPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("httpx command not available: %w", err)
	}

	// Verify it's actually HTTPX
	if !strings.Contains(string(output), "projectdiscovery.io") {
		return fmt.Errorf("command '%s' does not appear to be httpx", s.binPath)
	}

	return nil
}

// ConvertTarget converts a Target to a format suitable for httpx
func (s *HTTPXScanner) ConvertTarget(target models.Target) interface{} {
	return target.Value
}

// ConvertService converts a Service to a format suitable for httpx
func (s *HTTPXScanner) ConvertService(service models.Service) interface{} {
	// Format as URL with protocol, host, and port
	targetValue, ok := service.RawInfo["target_value"].(string)
	if !ok {
		// If target value not found in RawInfo, we can't convert this service
		return nil
	}

	protocol := "http"
	if service.ServiceName == "https" || service.Port == 443 {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s:%d", protocol, targetValue, service.Port)
}

// Scan performs an httpx scan against the target
func (s *HTTPXScanner) Scan(ctx context.Context, target interface{}, params models.JSONB) (*models.ScanResults, error) {
	targetValue, ok := target.(string)
	if !ok {
		return nil, fmt.Errorf("invalid target format for httpx scanner")
	}

	scanResults := &models.ScanResults{
		Findings:        []models.Finding{},
		NewTargets:      []models.Target{},
		TargetRelations: []models.TargetRelation{},
		Services:        []models.Service{},
	}

	// Parse parameters or use defaults
	timeout := s.timeout
	threads := 50
	followRedirects := true
	techDetect := true
	statusCode := true
	title := true
	webServer := true
	contentType := true
	tls := true
	favicon := true
	jarm := false
	probe := true
	ports := ""
	http2 := true
	securityHeaders := true
	extractCNAME := true

	// Override with provided parameters if available
	if val, ok := params["timeout"].(float64); ok {
		timeout = int(val)
	}
	if val, ok := params["threads"].(float64); ok {
		threads = int(val)
	}
	if val, ok := params["follow_redirects"].(bool); ok {
		followRedirects = val
	}
	if val, ok := params["tech_detect"].(bool); ok {
		techDetect = val
	}
	if val, ok := params["status_code"].(bool); ok {
		statusCode = val
	}
	if val, ok := params["title"].(bool); ok {
		title = val
	}
	if val, ok := params["web_server"].(bool); ok {
		webServer = val
	}
	if val, ok := params["content_type"].(bool); ok {
		contentType = val
	}
	if val, ok := params["tls"].(bool); ok {
		tls = val
	}
	if val, ok := params["favicon"].(bool); ok {
		favicon = val
	}
	if val, ok := params["jarm"].(bool); ok {
		jarm = val
	}
	if val, ok := params["probe"].(bool); ok {
		probe = val
	}
	if val, ok := params["ports"].(string); ok {
		ports = val
	}
	if val, ok := params["http2"].(bool); ok {
		http2 = val
	}
	if val, ok := params["security_headers"].(bool); ok {
		securityHeaders = val
	}
	if val, ok := params["extract_cname"].(bool); ok {
		extractCNAME = val
	}

	// Create a temporary file for output
	outputFile, err := os.CreateTemp("", "httpx-output-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary output file: %w", err)
	}
	defer os.Remove(outputFile.Name())
	defer outputFile.Close()

	// Start building the command
	args := []string{
		"-json",
		"-o", outputFile.Name(),
		"-timeout", strconv.Itoa(timeout),
		"-threads", strconv.Itoa(threads),
	}

	// Add target
	if strings.Contains(targetValue, "://") {
		// If the target is already a URL, use it directly
		args = append(args, "-u", targetValue)
	} else {
		// Otherwise, treat it as a host
		args = append(args, "-u", targetValue)
	}

	// Add optional parameters
	if followRedirects {
		args = append(args, "-follow-redirects")
	}
	if techDetect {
		args = append(args, "-tech-detect")
	}
	if statusCode {
		args = append(args, "-status-code")
	}
	if title {
		args = append(args, "-title")
	}
	if webServer {
		args = append(args, "-web-server")
	}
	if contentType {
		args = append(args, "-content-type")
	}
	if tls {
		args = append(args, "-tls-probe")
	}
	if favicon {
		args = append(args, "-favicon")
	}
	if jarm {
		args = append(args, "-jarm")
	}
	if probe {
		args = append(args, "-probe")
	}
	if ports != "" {
		args = append(args, "-ports", ports)
	}
	if http2 {
		args = append(args, "-http2")
	}
	if securityHeaders {
		args = append(args, "-include-response-header")
	}
	if extractCNAME {
		args = append(args, "-cname")
	}

	// Run httpx command
	cmd := exec.CommandContext(ctx, s.binPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start httpx command: %w", err)
	}

	// Read stdout for progress updates
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// Log progress but don't process - results will be in the JSON file
		fmt.Println("HTTPX progress:", line)
	}

	// Wait for completion
	if err := cmd.Wait(); err != nil {
		// Don't treat this as a fatal error - the command might exit with non-zero status
		// if no targets were found, but it may still have output
		fmt.Printf("HTTPX command exited with error: %v\n", err)
	}

	// Parse results from the JSON file
	results, err := s.parseHTTPXOutput(outputFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to parse httpx output: %w", err)
	}

	// Process results to create findings, targets, and services
	for _, result := range results {
		// Create findings, new targets, and services based on HTTPX results
		s.processHTTPXResult(result, scanResults)
	}

	// If no results were found, create a "no web servers found" finding
	if len(results) == 0 {
		finding := models.Finding{
			Title:       fmt.Sprintf("No web servers found for target %s", targetValue),
			Description: fmt.Sprintf("HTTPX did not discover any web servers for the target %s with the current scan parameters.", targetValue),
			Severity:    models.SeverityLow,
			FindingType: "no_web_servers",
			Details: models.JSONB{
				"target": targetValue,
			},
		}
		scanResults.Findings = append(scanResults.Findings, finding)
	}

	return scanResults, nil
}

// parseHTTPXOutput parses the JSON output from HTTPX
func (s *HTTPXScanner) parseHTTPXOutput(outputFile string) ([]HTTPXResult, error) {
	var results []HTTPXResult

	// Open the file
	file, err := os.Open(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open httpx output file: %w", err)
	}
	defer file.Close()

	// Read line by line as each line is a separate JSON object
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var result HTTPXResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			fmt.Printf("Error parsing HTTPX result line: %v\n", err)
			continue
		}

		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return results, fmt.Errorf("error reading httpx output file: %w", err)
	}

	return results, nil
}

// processHTTPXResult converts an HTTPX result to findings, targets, and services
func (s *HTTPXScanner) processHTTPXResult(result HTTPXResult, scanResults *models.ScanResults) {
	// Extract hostname and port from the result
	hostname := result.Host
	if hostname == "" {
		// Try to extract hostname from URL
		hostname = s.extractHostFromURL(result.URL)
	}

	// Parse the port
	port, _ := strconv.Atoi(result.ServerPort)
	if port == 0 {
		if strings.HasPrefix(result.URL, "https://") {
			port = 443
		} else {
			port = 80
		}
	}

	// Create a service for the web server
	service := models.Service{
		ID:          uuid.New(),
		Port:        port,
		Protocol:    "tcp",
		ServiceName: s.determineServiceName(result),
		Version:     result.WebServer,
		Title:       result.Title,
		Description: s.generateServiceDescription(result),
		Banner:      result.WebServer,
		RawInfo: models.JSONB{
			"url":           result.URL,
			"status_code":   result.StatusCode,
			"content_type":  result.ContentType,
			"response_time": result.ResponseTime,
			"technologies":  result.Technologies,
			"host":          hostname,
			"discovered_at": time.Now().Format(time.RFC3339),
		},
	}

	// Add more details to the service info
	if result.TLSData != nil {
		service.RawInfo["tls"] = map[string]interface{}{
			"cipher":                  result.TLSData.Cipher,
			"version":                 result.TLSData.Version,
			"certificate_issuer":      result.TLSData.CertificateIssuer,
			"certificate_subject":     result.TLSData.CertificateSubject,
			"certificate_common_name": result.TLSData.CertificateCommonName,
			"certificate_not_before":  result.TLSData.CertificateNotBefore,
			"certificate_not_after":   result.TLSData.CertificateNotAfter,
			"expired":                 result.TLSData.Expired,
		}
	}

	if result.SecurityHeaders != nil && len(result.SecurityHeaders) > 0 {
		service.RawInfo["security_headers"] = result.SecurityHeaders
	}

	// Add favicon hashes if available
	if result.FaviconMmh3 != "" {
		service.RawInfo["favicon_mmh3"] = result.FaviconMmh3
	}
	if result.FaviconMD5 != "" {
		service.RawInfo["favicon_md5"] = result.FaviconMD5
	}

	// Add the service to the results
	scanResults.Services = append(scanResults.Services, service)

	// Create findings based on the HTTPX results
	findings := s.createFindings(result, service.ID)
	scanResults.Findings = append(scanResults.Findings, findings...)

	// Check if we need to create a new target for the hostname
	if hostname != "" && hostname != result.Input {
		// Create a new target for the hostname
		target := models.Target{
			ID:         uuid.New(),
			TargetType: models.TargetTypeDomain,
			Value:      hostname,
			Metadata: models.JSONB{
				"discovered_from": result.Input,
				"discovery_scan":  "httpx",
				"discovered_at":   time.Now().Format(time.RFC3339),
			},
		}
		scanResults.NewTargets = append(scanResults.NewTargets, target)

		// Create a relation between input and hostname
		relation := models.TargetRelation{
			ID:            uuid.New(),
			SourceID:      uuid.Nil, // Will be set by worker to the original target ID
			DestinationID: target.ID,
			RelationType:  "resolves_to",
			Metadata: models.JSONB{
				"discovered_at": time.Now().Format(time.RFC3339),
			},
		}
		scanResults.TargetRelations = append(scanResults.TargetRelations, relation)
	}

	// Process CNAMEs if available
	if result.CNAMEs != nil && len(result.CNAMEs) > 0 {
		for _, cname := range result.CNAMEs {
			// Create a new target for each CNAME
			cnameTarget := models.Target{
				ID:         uuid.New(),
				TargetType: models.TargetTypeDomain,
				Value:      cname,
				Metadata: models.JSONB{
					"discovered_from": result.Input,
					"discovery_scan":  "httpx",
					"discovered_at":   time.Now().Format(time.RFC3339),
				},
			}
			scanResults.NewTargets = append(scanResults.NewTargets, cnameTarget)

			// Create a relation between input and CNAME
			relation := models.TargetRelation{
				ID:            uuid.New(),
				SourceID:      uuid.Nil, // Will be set by worker to the original target ID
				DestinationID: cnameTarget.ID,
				RelationType:  "cname_of",
				Metadata: models.JSONB{
					"discovered_at": time.Now().Format(time.RFC3339),
				},
			}
			scanResults.TargetRelations = append(scanResults.TargetRelations, relation)
		}
	}

	// Process A records if available
	if result.A != nil && len(result.A) > 0 {
		for _, ip := range result.A {
			// Create a new target for each IP
			ipTarget := models.Target{
				ID:         uuid.New(),
				TargetType: models.TargetTypeIP,
				Value:      ip,
				Metadata: models.JSONB{
					"discovered_from": result.Input,
					"discovery_scan":  "httpx",
					"discovered_at":   time.Now().Format(time.RFC3339),
				},
			}
			scanResults.NewTargets = append(scanResults.NewTargets, ipTarget)

			// Create a relation between input and IP
			relation := models.TargetRelation{
				ID:            uuid.New(),
				SourceID:      uuid.Nil, // Will be set by worker to the original target ID
				DestinationID: ipTarget.ID,
				RelationType:  "resolves_to",
				Metadata: models.JSONB{
					"discovered_at": time.Now().Format(time.RFC3339),
				},
			}
			scanResults.TargetRelations = append(scanResults.TargetRelations, relation)
		}
	}
}

// createFindings generates findings based on the HTTPX result
func (s *HTTPXScanner) createFindings(result HTTPXResult, serviceID uuid.UUID) []models.Finding {
	var findings []models.Finding

	// Create a basic web server finding
	webServerFinding := models.Finding{
		Title:       fmt.Sprintf("Web server detected: %s", result.Host),
		Description: s.generateWebServerDescription(result),
		Severity:    models.SeverityInfo,
		FindingType: "web_server_detected",
		ServiceID:   &serviceID,
		Details: models.JSONB{
			"url":          result.URL,
			"status_code":  result.StatusCode,
			"web_server":   result.WebServer,
			"content_type": result.ContentType,
			"title":        result.Title,
		},
		DiscoveredAt: time.Now(),
		Verified:     true,
		Fixed:        false,
		Manual:       false,
	}
	findings = append(findings, webServerFinding)

	// Check for technology stack findings
	if result.Technologies != nil && len(result.Technologies) > 0 {
		techFinding := models.Finding{
			Title:       fmt.Sprintf("Technology stack identified on %s", result.Host),
			Description: s.generateTechnologiesDescription(result),
			Severity:    models.SeverityInfo,
			FindingType: "technology_detection",
			ServiceID:   &serviceID,
			Details: models.JSONB{
				"url":          result.URL,
				"technologies": result.Technologies,
			},
			DiscoveredAt: time.Now(),
			Verified:     true,
			Fixed:        false,
			Manual:       false,
		}
		findings = append(findings, techFinding)
	}

	// Check for TLS/SSL findings
	if result.TLSData != nil {
		severity := models.SeverityInfo
		tlsDesc := s.generateTLSDescription(result)

		// Check for expired certificates
		if result.TLSData.Expired {
			severity = models.SeverityHigh
			tlsDesc = "❌ " + tlsDesc + "\n\nThe SSL certificate has expired, which can lead to security warnings for users and reduced trust. It's recommended to renew the certificate as soon as possible."
		}

		tlsFinding := models.Finding{
			Title:       fmt.Sprintf("SSL/TLS configuration for %s", result.Host),
			Description: tlsDesc,
			Severity:    severity,
			FindingType: "tls_configuration",
			ServiceID:   &serviceID,
			Details: models.JSONB{
				"url":                   result.URL,
				"tls_version":           result.TLSData.Version,
				"tls_cipher":            result.TLSData.Cipher,
				"certificate_issuer":    result.TLSData.CertificateIssuer,
				"certificate_subject":   result.TLSData.CertificateSubject,
				"certificate_not_after": result.TLSData.CertificateNotAfter,
				"expired":               result.TLSData.Expired,
			},
			DiscoveredAt: time.Now(),
			Verified:     true,
			Fixed:        false,
			Manual:       false,
		}
		findings = append(findings, tlsFinding)
	}

	// Check for security headers
	if result.SecurityHeaders != nil && len(result.SecurityHeaders) > 0 {
		securityIssues := s.analyzeSecurityHeaders(result.SecurityHeaders)
		if len(securityIssues) > 0 {
			securityFinding := models.Finding{
				Title:       fmt.Sprintf("Security header analysis for %s", result.Host),
				Description: s.generateSecurityHeadersDescription(result, securityIssues),
				Severity:    models.SeverityMedium,
				FindingType: "security_headers",
				ServiceID:   &serviceID,
				Details: models.JSONB{
					"url":              result.URL,
					"security_headers": result.SecurityHeaders,
					"issues":           securityIssues,
				},
				DiscoveredAt: time.Now(),
				Verified:     true,
				Fixed:        false,
				Manual:       false,
			}
			findings = append(findings, securityFinding)
		}
	}

	return findings
}

// extractHostFromURL extracts the hostname from a URL
func (s *HTTPXScanner) extractHostFromURL(url string) string {
	// Strip protocol
	hostPart := url
	if strings.Contains(hostPart, "://") {
		hostPart = strings.Split(hostPart, "://")[1]
	}

	// Strip path and query
	if strings.Contains(hostPart, "/") {
		hostPart = strings.Split(hostPart, "/")[0]
	}

	// Strip port
	if strings.Contains(hostPart, ":") {
		hostPart = strings.Split(hostPart, ":")[0]
	}

	return hostPart
}

// determineServiceName determines the service name based on the HTTPX result
func (s *HTTPXScanner) determineServiceName(result HTTPXResult) string {
	if strings.HasPrefix(result.URL, "https://") {
		return "https"
	}
	return "http"
}

// generateServiceDescription creates a human-readable description of the web service
func (s *HTTPXScanner) generateServiceDescription(result HTTPXResult) string {
	desc := fmt.Sprintf("Web service detected on %s at port %s.", result.Host, result.ServerPort)

	if result.Title != "" {
		desc += fmt.Sprintf("\nPage Title: %s", result.Title)
	}

	if result.WebServer != "" {
		desc += fmt.Sprintf("\nWeb Server: %s", result.WebServer)
	}

	if result.ContentType != "" {
		desc += fmt.Sprintf("\nContent Type: %s", result.ContentType)
	}

	if result.Technologies != nil && len(result.Technologies) > 0 {
		desc += fmt.Sprintf("\nTechnologies: %s", strings.Join(result.Technologies, ", "))
	}

	return desc
}

// generateWebServerDescription creates a detailed description of the web server finding
func (s *HTTPXScanner) generateWebServerDescription(result HTTPXResult) string {
	desc := fmt.Sprintf("A web server was detected at %s (status code: %d).", result.URL, result.StatusCode)

	if result.Title != "" {
		desc += fmt.Sprintf("\n\nPage Title: %s", result.Title)
	}

	if result.WebServer != "" {
		desc += fmt.Sprintf("\nServer: %s", result.WebServer)
	}

	if result.ContentType != "" {
		desc += fmt.Sprintf("\nContent Type: %s", result.ContentType)
	}

	if result.ContentLength > 0 {
		desc += fmt.Sprintf("\nContent Length: %d bytes", result.ContentLength)
	}

	if result.ResponseTime != "" {
		desc += fmt.Sprintf("\nResponse Time: %s", result.ResponseTime)
	}

	if result.HTTP2 {
		desc += "\nHTTP/2: Supported"
	}

	return desc
}

// generateTechnologiesDescription creates a description of the technologies finding
func (s *HTTPXScanner) generateTechnologiesDescription(result HTTPXResult) string {
	desc := fmt.Sprintf("The following technologies were detected on %s:", result.URL)

	for _, tech := range result.Technologies {
		desc += fmt.Sprintf("\n• %s", tech)
	}

	desc += "\n\nKnowing the technology stack can help identify potential vulnerabilities and attack vectors specific to these technologies."

	return desc
}

// generateTLSDescription creates a description of the TLS finding
func (s *HTTPXScanner) generateTLSDescription(result HTTPXResult) string {
	tls := result.TLSData
	desc := fmt.Sprintf("SSL/TLS analysis for %s:", result.URL)

	desc += fmt.Sprintf("\n\nVersion: %s", tls.Version)
	desc += fmt.Sprintf("\nCipher: %s", tls.Cipher)

	if tls.CertificateCommonName != "" {
		desc += fmt.Sprintf("\nCertificate Common Name: %s", tls.CertificateCommonName)
	}

	if tls.CertificateIssuer != "" {
		desc += fmt.Sprintf("\nIssuer: %s", tls.CertificateIssuer)
	}

	if tls.CertificateSubject != "" {
		desc += fmt.Sprintf("\nSubject: %s", tls.CertificateSubject)
	}

	if tls.CertificateNotBefore != "" && tls.CertificateNotAfter != "" {
		desc += fmt.Sprintf("\nValidity: %s to %s", tls.CertificateNotBefore, tls.CertificateNotAfter)
	}

	return desc
}

// analyzeSecurityHeaders examines security headers for common issues
func (s *HTTPXScanner) analyzeSecurityHeaders(headers map[string]string) []string {
	var issues []string

	// Check for missing important security headers
	importantHeaders := map[string]string{
		"Content-Security-Policy":   "Missing Content-Security-Policy header, which helps prevent XSS and data injection attacks",
		"X-Frame-Options":           "Missing X-Frame-Options header, which can prevent clickjacking attacks",
		"X-Content-Type-Options":    "Missing X-Content-Type-Options header, which prevents MIME type sniffing",
		"Strict-Transport-Security": "Missing HSTS header, which enforces secure (HTTPS) connections",
		"X-XSS-Protection":          "Missing X-XSS-Protection header, which enables browser's XSS filtering",
		"Referrer-Policy":           "Missing Referrer-Policy header, which controls how much referrer information should be included with requests",
	}

	for header, issue := range importantHeaders {
		found := false
		for key := range headers {
			if strings.EqualFold(key, header) {
				found = true
				break
			}
		}
		if !found {
			issues = append(issues, issue)
		}
	}

	// Check for specific header issues
	for key, value := range headers {
		// Check for weak CSP
		if strings.EqualFold(key, "Content-Security-Policy") {
			if strings.Contains(value, "unsafe-inline") || strings.Contains(value, "unsafe-eval") {
				issues = append(issues, "Content-Security-Policy contains unsafe directives (unsafe-inline or unsafe-eval)")
			}
		}

		// Check for weak X-Frame-Options
		if strings.EqualFold(key, "X-Frame-Options") {
			if !strings.EqualFold(value, "DENY") && !strings.EqualFold(value, "SAMEORIGIN") {
				issues = append(issues, "X-Frame-Options has a potentially weak value: "+value)
			}
		}

		// Check for weak HSTS
		if strings.EqualFold(key, "Strict-Transport-Security") {
			if !strings.Contains(strings.ToLower(value), "max-age=") {
				issues = append(issues, "HSTS header is missing max-age directive")
			} else {
				// Extract max-age value
				parts := strings.Split(value, ";")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(strings.ToLower(part), "max-age=") {
						ageStr := strings.TrimPrefix(strings.ToLower(part), "max-age=")
						age, err := strconv.Atoi(ageStr)
						if err == nil && age < 10368000 { // Less than 120 days
							issues = append(issues, "HSTS max-age is too short (less than 120 days)")
						}
					}
				}
			}

			if !strings.Contains(strings.ToLower(value), "includesubdomains") {
				issues = append(issues, "HSTS header is missing includeSubDomains directive")
			}
		}
	}

	return issues
}

// generateSecurityHeadersDescription creates a description of security header findings
func (s *HTTPXScanner) generateSecurityHeadersDescription(result HTTPXResult, issues []string) string {
	desc := fmt.Sprintf("Security header analysis for %s identified %d issue(s):", result.URL, len(issues))

	// List all issues
	for _, issue := range issues {
		desc += fmt.Sprintf("\n• %s", issue)
	}

	// Add recommendations
	desc += "\n\nRecommendations:"
	desc += "\n• Implement all missing security headers"
	desc += "\n• Ensure Content-Security-Policy doesn't use unsafe directives"
	desc += "\n• Set HSTS max-age to at least 120 days and include the includeSubDomains directive"
	desc += "\n• Set X-Frame-Options to DENY or SAMEORIGIN"

	// Add existing headers
	if len(result.SecurityHeaders) > 0 {
		desc += "\n\nCurrent Security Headers:"
		for header, value := range result.SecurityHeaders {
			desc += fmt.Sprintf("\n%s: %s", header, value)
		}
	}

	return desc
}

// Type returns the scanner type identifier
func (s *HTTPXScanner) Type() string {
	return "httpx"
}

// SupportsTargetType indicates whether this scanner can handle the specified target type
func (s *HTTPXScanner) SupportsTargetType(targetType string) bool {
	switch targetType {
	case models.TargetTypeDomain, models.TargetTypeIP:
		return true
	case models.TargetTypeCIDR:
		return true // HTTPX can scan CIDR ranges
	default:
		return false
	}
}

// SupportsServices indicates whether this scanner can scan services
func (s *HTTPXScanner) SupportsServices() bool {
	return true
}
