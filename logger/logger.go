// logger/logger.go
package logger

import (
	"fmt"
	"strings"
)

// Log styles
const (
	InfoPrefix  = "[INFO] "
	WarnPrefix  = "[WARN] "
	ErrorPrefix = "[ERROR] "
)

func Info(message string) {
	fmt.Println(InfoPrefix + message)
	fmt.Println("--------------------------------")
}

// InfoTable logs slices of data in a table format
func InfoTable(header []string, rows [][]string) {
	// Find the maximum width for each column
	maxWidth := make([]int, len(header))
	for i, h := range header {
		maxWidth[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > maxWidth[i] {
				maxWidth[i] = len(cell)
			}
		}
	}

	// Create the format string for each row
	var formatBuilder strings.Builder
	formatBuilder.WriteString("[INFO] ")
	for _, width := range maxWidth {
		formatBuilder.WriteString(fmt.Sprintf("%%-%ds ", width))
	}
	format := formatBuilder.String()

	// Print the header
	fmt.Println(fmt.Sprintf(format, interfaceSlice(header)...))

	// Print each row
	for _, row := range rows {
		fmt.Println(fmt.Sprintf(format, interfaceSlice(row)...))
	}
	fmt.Println("--------------------------------")
}

// Helper function to convert a slice of strings to a slice of empty interfaces
func interfaceSlice(slice []string) []interface{} {
	interfaces := make([]interface{}, len(slice))
	for i, d := range slice {
		interfaces[i] = d
	}
	return interfaces
}

func Warn(message string) {
	fmt.Println(WarnPrefix + message)
	fmt.Println("--------------------------------")
}

func Error(message string) {
	fmt.Println(ErrorPrefix + message)
	fmt.Println("--------------------------------")
}
