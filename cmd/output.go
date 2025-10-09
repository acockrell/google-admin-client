package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// OutputFormat represents supported output formats
type OutputFormat string

const (
	OutputFormatJSON  OutputFormat = "json"
	OutputFormatCSV   OutputFormat = "csv"
	OutputFormatYAML  OutputFormat = "yaml"
	OutputFormatTable OutputFormat = "table"
	OutputFormatPlain OutputFormat = "plain"
)

// Global output settings (set by root command flags)
var (
	outputFormat OutputFormat = OutputFormatPlain
	quietMode    bool         = false
)

// ValidateOutputFormat checks if the given format string is valid
func ValidateOutputFormat(format string) error {
	switch OutputFormat(format) {
	case OutputFormatJSON, OutputFormatCSV, OutputFormatYAML, OutputFormatTable, OutputFormatPlain:
		return nil
	default:
		return fmt.Errorf("invalid output format: %s (must be json, csv, yaml, table, or plain)", format)
	}
}

// FormatOutput formats and outputs data according to the global output settings
func FormatOutput(data interface{}, headers []string) error {
	return FormatOutputWithWriter(os.Stdout, data, headers)
}

// FormatOutputWithWriter formats and outputs data to the specified writer
// This is useful for testing and for writing to different outputs
func FormatOutputWithWriter(w io.Writer, data interface{}, headers []string) error {
	if data == nil {
		if !quietMode {
			if _, err := fmt.Fprintln(w, "No data to display"); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		}
		return nil
	}

	switch outputFormat {
	case OutputFormatJSON:
		return formatJSON(w, data)
	case OutputFormatCSV:
		return formatCSV(w, data, headers)
	case OutputFormatYAML:
		return formatYAML(w, data)
	case OutputFormatTable:
		return formatTable(w, data, headers)
	case OutputFormatPlain:
		return formatPlain(w, data)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

// formatJSON outputs data as pretty-printed JSON
func formatJSON(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// formatYAML outputs data as YAML
func formatYAML(w io.Writer, data interface{}) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	defer func() {
		if err := encoder.Close(); err != nil {
			Logger.Warn().Err(err).Msg("Failed to close YAML encoder")
		}
	}()
	return encoder.Encode(data)
}

// formatCSV outputs data as CSV with headers
func formatCSV(w io.Writer, data interface{}, headers []string) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write headers
	if len(headers) > 0 && !quietMode {
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write CSV headers: %w", err)
		}
	}

	// Convert data to rows
	rows, err := convertToRows(data, headers)
	if err != nil {
		return err
	}

	// Write rows
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return writer.Error()
}

// formatTable outputs data as a formatted table
func formatTable(w io.Writer, data interface{}, headers []string) error {
	table := tablewriter.NewWriter(w)

	// Set headers
	if len(headers) > 0 && !quietMode {
		// Convert []string to []any for Header method
		headerAny := make([]any, len(headers))
		for i, h := range headers {
			headerAny[i] = h
		}
		table.Header(headerAny...)
	}

	// Convert data to rows
	rows, err := convertToRows(data, headers)
	if err != nil {
		return err
	}

	// Add rows to table
	for _, row := range rows {
		// Convert []string to []any for Append method
		rowAny := make([]any, len(row))
		for i, cell := range row {
			rowAny[i] = cell
		}
		if err := table.Append(rowAny); err != nil {
			return fmt.Errorf("failed to append row: %w", err)
		}
	}

	return table.Render()
}

// formatPlain outputs data in a simple, human-readable plain text format
func formatPlain(w io.Writer, data interface{}) error {
	// For plain format, we'll use reflection to output fields
	v := reflect.ValueOf(data)

	// Handle pointer
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Handle slice/array
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if err := formatPlainItem(w, v.Index(i).Interface()); err != nil {
				return err
			}
		}
		return nil
	}

	// Handle single item
	return formatPlainItem(w, data)
}

// formatPlainItem formats a single item in plain text
func formatPlainItem(w io.Writer, item interface{}) error {
	v := reflect.ValueOf(item)

	// Handle pointer
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	// Handle struct
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			// Get field name (use json tag if available)
			fieldName := field.Name
			if jsonTag := field.Tag.Get("json"); jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "" && parts[0] != "-" {
					fieldName = parts[0]
				}
			}

			if _, err := fmt.Fprintf(w, "%s: %v\n", fieldName, value.Interface()); err != nil {
				return fmt.Errorf("failed to write field: %w", err)
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
		return nil
	}

	// Handle map
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			if _, err := fmt.Fprintf(w, "%v: %v\n", key.Interface(), v.MapIndex(key).Interface()); err != nil {
				return fmt.Errorf("failed to write map entry: %w", err)
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
		return nil
	}

	// Fallback: just print the value
	if _, err := fmt.Fprintln(w, item); err != nil {
		return fmt.Errorf("failed to write value: %w", err)
	}
	return nil
}

// convertToRows converts data to a slice of string slices for CSV/table output
func convertToRows(data interface{}, headers []string) ([][]string, error) {
	var rows [][]string

	v := reflect.ValueOf(data)

	// Handle pointer
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Handle slice/array
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			row, err := convertItemToRow(v.Index(i).Interface(), headers)
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
		}
		return rows, nil
	}

	// Handle single item
	row, err := convertItemToRow(data, headers)
	if err != nil {
		return nil, err
	}
	return [][]string{row}, nil
}

// convertItemToRow converts a single item to a row of strings
func convertItemToRow(item interface{}, headers []string) ([]string, error) {
	var row []string

	v := reflect.ValueOf(item)

	// Handle pointer
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return make([]string, len(headers)), nil
		}
		v = v.Elem()
	}

	// Handle struct
	if v.Kind() == reflect.Struct {
		// If headers are provided, extract values in header order
		if len(headers) > 0 {
			row = make([]string, len(headers))
			t := v.Type()

			// Build a map of json tag -> field index for quick lookup
			fieldMap := make(map[string]int)
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				fieldName := field.Name

				// Use json tag if available
				if jsonTag := field.Tag.Get("json"); jsonTag != "" {
					parts := strings.Split(jsonTag, ",")
					if parts[0] != "" && parts[0] != "-" {
						fieldName = parts[0]
					}
				}

				fieldMap[strings.ToLower(fieldName)] = i
			}

			// Extract values in header order
			for i, header := range headers {
				headerLower := strings.ToLower(header)
				if fieldIdx, ok := fieldMap[headerLower]; ok {
					value := v.Field(fieldIdx)
					row[i] = fmt.Sprintf("%v", value.Interface())
				} else {
					row[i] = ""
				}
			}
		} else {
			// No headers provided, just iterate over fields
			t := v.Type()
			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				if !field.IsExported() {
					continue
				}
				value := v.Field(i)
				row = append(row, fmt.Sprintf("%v", value.Interface()))
			}
		}
		return row, nil
	}

	// Handle map
	if v.Kind() == reflect.Map {
		if len(headers) > 0 {
			row = make([]string, len(headers))
			for i, header := range headers {
				// Try to find the value for this header in the map
				for _, key := range v.MapKeys() {
					if strings.EqualFold(fmt.Sprintf("%v", key.Interface()), header) {
						value := v.MapIndex(key)
						row[i] = fmt.Sprintf("%v", value.Interface())
						break
					}
				}
			}
		} else {
			// No headers, just dump all map values
			for _, key := range v.MapKeys() {
				value := v.MapIndex(key)
				row = append(row, fmt.Sprintf("%v", value.Interface()))
			}
		}
		return row, nil
	}

	// Fallback: convert to single string
	return []string{fmt.Sprintf("%v", item)}, nil
}

// QuietPrintf prints output only if not in quiet mode
func QuietPrintf(format string, args ...interface{}) {
	if !quietMode {
		fmt.Printf(format, args...)
	}
}

// QuietPrintln prints output only if not in quiet mode
func QuietPrintln(args ...interface{}) {
	if !quietMode {
		fmt.Println(args...)
	}
}
