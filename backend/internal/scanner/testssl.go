// internal/scanner/testssl.go
package scanner

import (
	"backend/internal/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	//"net"
	"os"
	"os/exec"
	"strings"
	//"time"
	//"github.com/google/uuid"
)

type TestSSLScanner struct {
	resolverTimeout int
}

type TestSSLOutput struct {
	Id string `json:"id"`
	Ip string `json:"ip"`
	//port string `json:"port"`
	//severity *string `json:"severity"`
	//cve      *string `json:"cve"`
	//cwe      *string `json:"cwe"`
	Finding string `json:"finding"`
}

// NewTestSSLScanner creates a new TestSSL resolver scanner
func NewTestSSLScanner() *TestSSLScanner {
	return &TestSSLScanner{
		resolverTimeout: 5, // Default timeout in seconds
	}
}

// Initialize checks if TestSSL resolution is available
func (s *TestSSLScanner) Initialize(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "testssl.sh", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil || strings.TrimSpace(string(output)) == "" {
		return fmt.Errorf("testssl.sh not available: %w", err)
	}

	return nil
}

// ConvertTarget converts a Target to a format suitable for TestSSL resolution
func (s *TestSSLScanner) ConvertTarget(target models.Target) interface{} {
	switch target.TargetType {
	case models.TargetTypeDomain:
		return target.Value
	case models.TargetTypeIP:
		return target.Value
	default:
		return nil
	}
}

func (s *TestSSLScanner) ConvertService(service models.Service) interface{} {
	// TODO: Implement this...
	return nil
}

func generateRandomID(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Scan performs TestSSL resolution against the target
func (s *TestSSLScanner) Scan(ctx context.Context, target interface{}, params models.JSONB) (*models.ScanResults, error) {
	if target == nil {
		return nil, fmt.Errorf("invalid target for TestSSL resolver")
	}

	targetValue := target.(string)
	scanResults := &models.ScanResults{
		Certificates: []models.Certificate{},
		Findings:     []models.Finding{},
	}

	var certificates []models.Certificate
	var findings []models.Finding

	randomBytes, err := generateRandomID(16)
	if err != nil {
		fmt.Println("Error generating random ID:", err)
	}
	randPrefix := fmt.Sprintf("tempdata_%d_%s_%d",
		time.Now().UnixNano(),
		randomBytes,
		os.Getpid())

	fileName := os.TempDir() + "/" + randPrefix + ".json"
	tmpFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Erroring creating file:", err)
	}
	cmd := exec.CommandContext(ctx, "testssl.sh", "-q", "--color", "0", "--jsonfile", fileName, targetValue)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running testssl.sh: ", err)
		fmt.Println("Command output:", output)
	}

	err = tmpFile.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
	}

	jsonFile, err := os.Open(fileName)
	if err == nil {
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		var outputJson []TestSSLOutput
		err = json.Unmarshal(byteValue, &outputJson)
		if err == nil {
			var certificate models.Certificate
			layout := "2006-01-02 15:04"
			for _, f := range outputJson {
				switch f.Id {
				case "cert_subjectAltName":
					certificate.Domain = f.Finding
				case "cert_notBefore":
					parsedTime, err := time.Parse(layout, f.Finding)
					if err == nil {
						certificate.IssuedAt = parsedTime
					} else {
						fmt.Printf("Error parsing time: %v \n", err)
					}
				case "cert_notAfter":
					parsedTime, err := time.Parse(layout, f.Finding)
					if err == nil {
						certificate.ExpiresAt = parsedTime
					} else {
						fmt.Printf("Error parsing time: %v \n", err)
					}
				case "cert_caIssuers":
					certificate.Issuer = f.Finding
				default:
					continue
				}
			}
			certificates = append(certificates, certificate)
		}
	}

	err = os.Remove(fileName)
	if err != nil {
		fmt.Println("Error removing file:", err)
	}

	scanResults.Findings = findings
	scanResults.Certificates = certificates
	return scanResults, nil
}

// Type returns the scanner type identifier
func (s *TestSSLScanner) Type() string {
	return "dns"
}

// SupportsTargetType indicates whether this scanner can handle the specified target type
func (s *TestSSLScanner) SupportsTargetType(targetType string) bool {
	switch targetType {
	case models.TargetTypeDomain, models.TargetTypeIP:
		return true
	default:
		return false
	}
}

// SupportsServices indicates whether this scanner can scan services
func (s *TestSSLScanner) SupportsServices() bool {
	return false
}
