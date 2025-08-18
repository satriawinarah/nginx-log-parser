package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	nginx_parser "github.com/satriawinarah/nginx-log-parser/internal"
)

func parseLogLine(line string) (nginx_parser.Visitor, error) {
	visitor := nginx_parser.Visitor{}

	parts := strings.Fields(line)
	if len(parts) < 9 {
		return nginx_parser.Visitor{}, fmt.Errorf("invalid log line: %s", line)
	}

	ip := parts[0]
	timestamp := strings.Join(parts[3:5], " ")
	request := strings.Join(parts[5:8], " ")
	status := parts[8]
	userAgent := strings.Join(parts[9:], " ")

	visitor = *nginx_parser.ParseRequest(request, &visitor)
	visitor.IP = ip
	visitor.RequestTime = timestamp
	visitor.ResponseStatus = status
	visitor.UserAgent = userAgent

	return visitor, nil
}

func main() {
	file, err := os.Open("input/access.log")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	visitors := []nginx_parser.Visitor{}
	for scanner.Scan() {
		line := scanner.Text()
		log, err := parseLogLine(line)
		if err != nil {
			fmt.Println("Error parsing log line:", err)
			continue
		}
		visitors = append(visitors, log)
	}

	err = nginx_parser.SaveVisitorToDB(visitors)
	if err != nil {
		fmt.Println("Error saving visitors to DB:", err)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
