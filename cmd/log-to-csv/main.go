package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/satriawinarah/nginx-log-parser/internal/log-to-csv"
)

func main() {
	// Define command line flags
	var (
		inputPath  = flag.String("input", "", "Input nginx log file path or directory containing .log files")
		outputPath = flag.String("output", "", "Output CSV file path or directory for CSV files")
		help       = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	// Show help if requested
	if *help {
		fmt.Println("Nginx Log to CSV Converter")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  log-to-csv -input <log-file> -output <csv-file>")
		fmt.Println("  log-to-csv -input <log-directory> -output <csv-directory>")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  log-to-csv -input input/access.log -output output/access.csv")
		fmt.Println("  log-to-csv -input input/ -output output/")
		fmt.Println("")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		return
	}

	// Validate required flags
	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Error: Both -input and -output flags are required")
		fmt.Println("Use -help for usage information")
		os.Exit(1)
	}

	// Check if input path exists
	inputInfo, err := os.Stat(*inputPath)
	if err != nil {
		log.Fatalf("Error accessing input path: %v", err)
	}

	// Process based on whether input is a file or directory
	if inputInfo.IsDir() {
		// Process directory
		fmt.Printf("Processing log files in directory: %s\n", *inputPath)
		fmt.Printf("Output directory: %s\n", *outputPath)

		if err := logtocsv.ProcessLogDirectory(*inputPath, *outputPath); err != nil {
			log.Fatalf("Error processing log directory: %v", err)
		}
	} else {
		// Process single file
		fmt.Printf("Processing single log file: %s\n", *inputPath)
		fmt.Printf("Output file: %s\n", *outputPath)

		if err := logtocsv.ConvertLogToCSV(*inputPath, *outputPath); err != nil {
			log.Fatalf("Error converting log file: %v", err)
		}
	}

	fmt.Println("\nConversion completed successfully!")
}
