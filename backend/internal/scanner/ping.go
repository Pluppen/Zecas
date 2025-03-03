// internal/scanner/ping.go
package scanner

import (
	"backend/internal/models"
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// PingScanner implements the Scanner interface for pinging targets
type PingScanner struct {
	count   int
	timeout int
	binPath string
}

// NewPingScanner creates a new ping scanner
func NewPingScanner() *PingScanner {
	return &PingScanner{
		count:   3,      // Default to 3 pings
		timeout: 5,      // Default to 5 second timeout
		binPath: "ping", // Default ping command
	}
}

// Initialize checks if ping is available
func (s *PingScanner) Initialize(ctx context.Context) error {
	// Simple check to see if ping is available
	cmd := exec.CommandContext(ctx, s.binPath, "-c", "1", "-W", "1", "127.0.0.1")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ping command not available: %w", err)
	}

	return nil
}

// ConvertTarget converts a Target to a format suitable for ping
func (s *PingScanner) ConvertTarget(target models.Target) interface{} {
	// For ping scanner, we need to ensure we have a hostname or IP
	// Domain type is preferable, but IP works too

	// If target is a CIDR, we'll just use the network address
	if target.TargetType == models.TargetTypeCIDR {
		// Extract the network part before the /
		parts := strings.Split(target.Value, "/")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// Otherwise just use the value
	return target.Value
}

// Scan performs a ping scan against the target
func (s *PingScanner) Scan(ctx context.Context, target interface{}, params models.JSONB) ([]models.Finding, error) {
	targetValue := target.(string)

	// Configure ping parameters
	count := s.count
	timeout := s.timeout

	// Override with provided parameters if available
	if val, ok := params["count"].(float64); ok {
		count = int(val)
	}

	if val, ok := params["timeout"].(float64); ok {
		timeout = int(val)
	}

	// Build ping command
	// Different systems use different ping arguments
	args := []string{
		"-c", strconv.Itoa(count), // Count of pings
		"-W", strconv.Itoa(timeout), // Timeout in seconds
	}

	// Add target
	args = append(args, targetValue)

	// Execute ping command
	cmd := exec.CommandContext(ctx, s.binPath, args...)
	output, err := cmd.CombinedOutput()

	// Get output regardless of error (ping returns non-zero if host is unreachable)
	outputStr := string(output)

	// Create finding
	finding := models.Finding{
		ScanID:      nil,      // Will be set by worker
		TargetID:    uuid.Nil, // Will be set by worker
		Title:       fmt.Sprintf("Ping result for %s", targetValue),
		Description: s.generateDescription(targetValue, outputStr, err != nil),
		Severity:    s.determineSeverity(err != nil),
		FindingType: "ping",
		Details: models.JSONB{
			"target":     targetValue,
			"reachable":  err == nil,
			"output":     outputStr,
			"ping_count": count,
			"timeout":    timeout,
		},
	}

	// Attempt to resolve hostname
	if err == nil {
		// Try to get IP if target was a hostname
		ips, err := net.LookupHost(targetValue)
		if err == nil && len(ips) > 0 {
			finding.Details["resolved_ips"] = ips
		}

		// Parse stats if possible
		stats := s.parsePingStats(outputStr)
		for k, v := range stats {
			finding.Details[k] = v
		}
	}

	return []models.Finding{finding}, nil
}

// Type returns the scanner type identifier
func (s *PingScanner) Type() string {
	return "ping"
}

// generateDescription creates a human-readable description of ping results
func (s *PingScanner) generateDescription(target string, output string, failed bool) string {
	if failed {
		return fmt.Sprintf("Host %s is unreachable via ICMP ping. This could indicate that the host is down, blocks ICMP packets, or network connectivity issues.", target)
	}

	// Extract ping statistics if available
	stats := s.parsePingStats(output)

	description := fmt.Sprintf("Host %s is reachable via ICMP ping.", target)

	if val, ok := stats["min_rtt"]; ok {
		description += fmt.Sprintf("\nPing statistics: min/avg/max = %.2f/%.2f/%.2f ms",
			val, stats["avg_rtt"], stats["max_rtt"])
	}

	if val, ok := stats["packet_loss"]; ok {
		description += fmt.Sprintf("\nPacket Loss: %.1f%%", val)
	}

	return description
}

// determineSeverity sets severity based on ping results
func (s *PingScanner) determineSeverity(failed bool) string {
	if failed {
		return models.SeverityHigh
	}
	return models.SeverityInfo
}

// parsePingStats extracts metrics from ping output
func (s *PingScanner) parsePingStats(output string) map[string]float64 {
	stats := make(map[string]float64)

	// Find the statistics line that typically looks like:
	// rtt min/avg/max/mdev = 0.083/0.145/0.214/0.055 ms
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Look for packet loss info
		if strings.Contains(line, "packet loss") {
			parts := strings.Split(line, ",")
			for _, part := range parts {
				if strings.Contains(part, "packet loss") {
					lossStr := strings.TrimSpace(strings.Split(part, "packet loss")[0])
					lossStr = strings.TrimSuffix(lossStr, "%")
					if loss, err := strconv.ParseFloat(lossStr, 64); err == nil {
						stats["packet_loss"] = loss
					}
				}
			}
		}

		// Look for RTT stats
		if strings.Contains(line, "min/avg/max") {
			parts := strings.Split(line, "=")
			if len(parts) < 2 {
				continue
			}

			rtts := strings.Split(strings.TrimSpace(parts[1]), "/")
			if len(rtts) < 3 {
				continue
			}

			// Parse min/avg/max RTTs
			if min, err := strconv.ParseFloat(rtts[0], 64); err == nil {
				stats["min_rtt"] = min
			}

			if avg, err := strconv.ParseFloat(rtts[1], 64); err == nil {
				stats["avg_rtt"] = avg
			}

			if max, err := strconv.ParseFloat(rtts[2], 64); err == nil {
				stats["max_rtt"] = max
			}
		}
	}

	return stats
}
