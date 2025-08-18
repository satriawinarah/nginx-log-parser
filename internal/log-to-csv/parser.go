package logtocsv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// NginxLogEntry represents a parsed nginx log entry
type NginxLogEntry struct {
	IP             string
	Timestamp      string
	RequestMethod  string
	RequestURI     string
	HTTPVersion    string
	ResponseStatus string
	ResponseSize   string
	Referer        string
	UserAgent      string
	RemoteUser     string
}

// ParseNginxLog parses a single nginx log line
func ParseNginxLog(logLine string) (*NginxLogEntry, error) {
	// Standard nginx log format:
	// IP - - [timestamp] "method uri version" status size "referer" "user-agent"

	// Remove any leading/trailing whitespace
	logLine = strings.TrimSpace(logLine)
	if logLine == "" {
		return nil, fmt.Errorf("empty log line")
	}

	entry := &NginxLogEntry{}

	// Use regex to parse the log line more reliably
	// This handles cases where fields might contain spaces or special characters
	re := regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+\[([^\]]+)\]\s+"([^"]*)"\s+(\S+)\s+(\S+)\s+"([^"]*)"\s+"([^"]*)"`)
	matches := re.FindStringSubmatch(logLine)

	if len(matches) < 10 {
		// Try alternative format for logs without referer or user-agent
		reAlt := regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+\[([^\]]+)\]\s+"([^"]*)"\s+(\S+)\s+(\S+)`)
		matches = reAlt.FindStringSubmatch(logLine)
		if len(matches) < 7 {
			return nil, fmt.Errorf("unable to parse log line: %s", logLine)
		}
	}

	entry.IP = matches[1]
	entry.RemoteUser = matches[2]
	entry.Timestamp = matches[4]

	// Parse the request line (method, URI, version)
	requestParts := strings.Split(matches[5], " ")
	if len(requestParts) >= 3 {
		entry.RequestMethod = requestParts[0]
		entry.RequestURI = requestParts[1]
		entry.HTTPVersion = requestParts[2]
	} else if len(requestParts) >= 2 {
		entry.RequestMethod = requestParts[0]
		entry.RequestURI = requestParts[1]
		entry.HTTPVersion = ""
	} else if len(requestParts) >= 1 {
		entry.RequestMethod = requestParts[0]
		entry.RequestURI = ""
		entry.HTTPVersion = ""
	}

	entry.ResponseStatus = matches[6]
	entry.ResponseSize = matches[7]

	if len(matches) > 8 {
		entry.Referer = matches[8]
	}
	if len(matches) > 9 {
		entry.UserAgent = matches[9]
	}

	return entry, nil
}

// ConvertLogToCSV converts nginx log file to CSV format
func ConvertLogToCSV(inputPath, outputPath string) error {
	// Open input log file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output CSV file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Create CSV writer
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Write CSV header
	header := []string{
		"IP",
		"Timestamp",
		"RequestMethod",
		"RequestURI",
		"HTTPVersion",
		"ResponseStatus",
		"ResponseSize",
		"Referer",
		"UserAgent",
		"RemoteUser",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Process log file line by line
	scanner := bufio.NewScanner(inputFile)
	lineNumber := 0
	successCount := 0
	errorCount := 0

	for scanner.Scan() {
		lineNumber++
		logLine := scanner.Text()

		if strings.TrimSpace(logLine) == "" {
			continue
		}

		entry, err := ParseNginxLog(logLine)
		if err != nil {
			errorCount++
			fmt.Printf("Warning: Failed to parse line %d: %v\n", lineNumber, err)
			continue
		}

		// Write entry to CSV
		row := []string{
			entry.IP,
			entry.Timestamp,
			entry.RequestMethod,
			entry.RequestURI,
			entry.HTTPVersion,
			entry.ResponseStatus,
			entry.ResponseSize,
			entry.Referer,
			entry.UserAgent,
			entry.RemoteUser,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}

		successCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	fmt.Printf("Conversion completed successfully!\n")
	fmt.Printf("Input file: %s\n", inputPath)
	fmt.Printf("Output file: %s\n", outputPath)
	fmt.Printf("Successfully parsed: %d lines\n", successCount)
	fmt.Printf("Failed to parse: %d lines\n", errorCount)
	fmt.Printf("Total lines processed: %d\n", lineNumber)

	return nil
}

// ProcessLogDirectory processes all .log files in a directory and converts them to CSV
func ProcessLogDirectory(inputDir, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Find all .log files in input directory
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("failed to read input directory: %w", err)
	}

	var logFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".log") {
			logFiles = append(logFiles, entry.Name())
		}
	}

	if len(logFiles) == 0 {
		return fmt.Errorf("no .log files found in directory: %s", inputDir)
	}

	fmt.Printf("Found %d log files to process:\n", len(logFiles))
	for _, file := range logFiles {
		fmt.Printf("  - %s\n", file)
	}

	// Process each log file
	for _, logFile := range logFiles {
		inputPath := filepath.Join(inputDir, logFile)
		outputPath := filepath.Join(outputDir, strings.TrimSuffix(logFile, ".log")+".csv")

		fmt.Printf("\nProcessing: %s\n", logFile)
		if err := ConvertLogToCSV(inputPath, outputPath); err != nil {
			fmt.Printf("Error processing %s: %v\n", logFile, err)
			continue
		}
	}

	return nil
}
