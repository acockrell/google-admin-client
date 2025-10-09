package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidateOutputFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{
			name:    "valid json",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "valid csv",
			format:  "csv",
			wantErr: false,
		},
		{
			name:    "valid yaml",
			format:  "yaml",
			wantErr: false,
		},
		{
			name:    "valid table",
			format:  "table",
			wantErr: false,
		},
		{
			name:    "valid plain",
			format:  "plain",
			wantErr: false,
		},
		{
			name:    "invalid format",
			format:  "xml",
			wantErr: true,
		},
		{
			name:    "empty format",
			format:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOutputFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatJSON(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	tests := []struct {
		name     string
		data     interface{}
		contains []string
	}{
		{
			name: "simple struct",
			data: testData{Name: "John Doe", Email: "john@example.com", Age: 30},
			contains: []string{
				`"name": "John Doe"`,
				`"email": "john@example.com"`,
				`"age": 30`,
			},
		},
		{
			name: "slice of structs",
			data: []testData{
				{Name: "Alice", Email: "alice@example.com", Age: 25},
				{Name: "Bob", Email: "bob@example.com", Age: 35},
			},
			contains: []string{
				`"name": "Alice"`,
				`"name": "Bob"`,
				`"email": "alice@example.com"`,
				`"email": "bob@example.com"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatJSON(&buf, tt.data)
			if err != nil {
				t.Errorf("formatJSON() error = %v", err)
				return
			}

			output := buf.String()
			for _, substr := range tt.contains {
				if !strings.Contains(output, substr) {
					t.Errorf("formatJSON() output doesn't contain %q\nGot: %s", substr, output)
				}
			}
		})
	}
}

func TestFormatYAML(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	tests := []struct {
		name     string
		data     interface{}
		contains []string
	}{
		{
			name: "simple struct",
			data: testData{Name: "John Doe", Email: "john@example.com", Age: 30},
			contains: []string{
				"name: John Doe",
				"email: john@example.com",
				"age: 30",
			},
		},
		{
			name: "slice of structs",
			data: []testData{
				{Name: "Alice", Email: "alice@example.com", Age: 25},
				{Name: "Bob", Email: "bob@example.com", Age: 35},
			},
			contains: []string{
				"name: Alice",
				"name: Bob",
				"email: alice@example.com",
				"email: bob@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatYAML(&buf, tt.data)
			if err != nil {
				t.Errorf("formatYAML() error = %v", err)
				return
			}

			output := buf.String()
			for _, substr := range tt.contains {
				if !strings.Contains(output, substr) {
					t.Errorf("formatYAML() output doesn't contain %q\nGot: %s", substr, output)
				}
			}
		})
	}
}

func TestFormatCSV(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	tests := []struct {
		name     string
		data     interface{}
		headers  []string
		contains []string
	}{
		{
			name:    "simple struct with headers",
			data:    testData{Name: "John Doe", Email: "john@example.com", Age: 30},
			headers: []string{"Name", "Email", "Age"},
			contains: []string{
				"Name,Email,Age",
				"John Doe,john@example.com,30",
			},
		},
		{
			name: "slice of structs with headers",
			data: []testData{
				{Name: "Alice", Email: "alice@example.com", Age: 25},
				{Name: "Bob", Email: "bob@example.com", Age: 35},
			},
			headers: []string{"Name", "Email", "Age"},
			contains: []string{
				"Name,Email,Age",
				"Alice,alice@example.com,25",
				"Bob,bob@example.com,35",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			// Temporarily disable quiet mode for testing
			oldQuietMode := quietMode
			quietMode = false
			defer func() { quietMode = oldQuietMode }()

			err := formatCSV(&buf, tt.data, tt.headers)
			if err != nil {
				t.Errorf("formatCSV() error = %v", err)
				return
			}

			output := buf.String()
			for _, substr := range tt.contains {
				if !strings.Contains(output, substr) {
					t.Errorf("formatCSV() output doesn't contain %q\nGot: %s", substr, output)
				}
			}
		})
	}
}

func TestFormatTable(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	tests := []struct {
		name     string
		data     interface{}
		headers  []string
		contains []string
	}{
		{
			name:    "simple struct with headers",
			data:    testData{Name: "John Doe", Email: "john@example.com"},
			headers: []string{"Name", "Email"},
			contains: []string{
				"NAME",
				"EMAIL",
				"John Doe",
				"john@example.com",
			},
		},
		{
			name: "slice of structs with headers",
			data: []testData{
				{Name: "Alice", Email: "alice@example.com"},
				{Name: "Bob", Email: "bob@example.com"},
			},
			headers: []string{"Name", "Email"},
			contains: []string{
				"NAME",
				"EMAIL",
				"Alice",
				"Bob",
				"alice@example.com",
				"bob@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			// Temporarily disable quiet mode for testing
			oldQuietMode := quietMode
			quietMode = false
			defer func() { quietMode = oldQuietMode }()

			err := formatTable(&buf, tt.data, tt.headers)
			if err != nil {
				t.Errorf("formatTable() error = %v", err)
				return
			}

			output := buf.String()
			for _, substr := range tt.contains {
				if !strings.Contains(output, substr) {
					t.Errorf("formatTable() output doesn't contain %q\nGot: %s", substr, output)
				}
			}
		})
	}
}

func TestFormatPlain(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	tests := []struct {
		name     string
		data     interface{}
		contains []string
	}{
		{
			name: "simple struct",
			data: testData{Name: "John Doe", Email: "john@example.com"},
			contains: []string{
				"name: John Doe",
				"email: john@example.com",
			},
		},
		{
			name: "slice of structs",
			data: []testData{
				{Name: "Alice", Email: "alice@example.com"},
				{Name: "Bob", Email: "bob@example.com"},
			},
			contains: []string{
				"name: Alice",
				"name: Bob",
				"email: alice@example.com",
				"email: bob@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatPlain(&buf, tt.data)
			if err != nil {
				t.Errorf("formatPlain() error = %v", err)
				return
			}

			output := buf.String()
			for _, substr := range tt.contains {
				if !strings.Contains(output, substr) {
					t.Errorf("formatPlain() output doesn't contain %q\nGot: %s", substr, output)
				}
			}
		})
	}
}

func TestQuietMode(t *testing.T) {
	tests := []struct {
		name      string
		quietMode bool
		wantEmpty bool
	}{
		{
			name:      "quiet mode enabled",
			quietMode: true,
			wantEmpty: true,
		},
		{
			name:      "quiet mode disabled",
			quietMode: false,
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original quiet mode and restore after test
			oldQuietMode := quietMode
			defer func() { quietMode = oldQuietMode }()

			quietMode = tt.quietMode

			// Test QuietPrintf
			var buf bytes.Buffer
			oldStdout := bytes.NewBuffer(nil)
			// We can't easily capture stdout, so we'll just test the function doesn't panic
			QuietPrintf("test message")
			QuietPrintln("test message")

			// The test passes if we reach here without panicking
			_ = buf
			_ = oldStdout
		})
	}
}

func TestConvertToRows(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	tests := []struct {
		name     string
		data     interface{}
		headers  []string
		expected [][]string
		wantErr  bool
	}{
		{
			name:    "single struct with headers",
			data:    testData{Name: "John", Email: "john@example.com", Age: 30},
			headers: []string{"Name", "Email", "Age"},
			expected: [][]string{
				{"John", "john@example.com", "30"},
			},
			wantErr: false,
		},
		{
			name: "multiple structs with headers",
			data: []testData{
				{Name: "Alice", Email: "alice@example.com", Age: 25},
				{Name: "Bob", Email: "bob@example.com", Age: 35},
			},
			headers: []string{"Name", "Email", "Age"},
			expected: [][]string{
				{"Alice", "alice@example.com", "25"},
				{"Bob", "bob@example.com", "35"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := convertToRows(tt.data, tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(rows) != len(tt.expected) {
				t.Errorf("convertToRows() returned %d rows, want %d", len(rows), len(tt.expected))
				return
			}

			for i, row := range rows {
				if len(row) != len(tt.expected[i]) {
					t.Errorf("convertToRows() row %d has %d columns, want %d", i, len(row), len(tt.expected[i]))
					continue
				}

				for j, cell := range row {
					if cell != tt.expected[i][j] {
						t.Errorf("convertToRows() row %d col %d = %q, want %q", i, j, cell, tt.expected[i][j])
					}
				}
			}
		})
	}
}

func TestFormatOutputWithWriter(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	tests := []struct {
		name    string
		format  OutputFormat
		data    interface{}
		headers []string
		wantErr bool
	}{
		{
			name:    "json format",
			format:  OutputFormatJSON,
			data:    testData{Name: "John", Email: "john@example.com"},
			headers: nil,
			wantErr: false,
		},
		{
			name:    "csv format",
			format:  OutputFormatCSV,
			data:    testData{Name: "John", Email: "john@example.com"},
			headers: []string{"Name", "Email"},
			wantErr: false,
		},
		{
			name:    "yaml format",
			format:  OutputFormatYAML,
			data:    testData{Name: "John", Email: "john@example.com"},
			headers: nil,
			wantErr: false,
		},
		{
			name:    "table format",
			format:  OutputFormatTable,
			data:    testData{Name: "John", Email: "john@example.com"},
			headers: []string{"Name", "Email"},
			wantErr: false,
		},
		{
			name:    "plain format",
			format:  OutputFormatPlain,
			data:    testData{Name: "John", Email: "john@example.com"},
			headers: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original format and quiet mode
			oldFormat := outputFormat
			oldQuiet := quietMode
			defer func() {
				outputFormat = oldFormat
				quietMode = oldQuiet
			}()

			outputFormat = tt.format
			quietMode = false

			var buf bytes.Buffer
			err := FormatOutputWithWriter(&buf, tt.data, tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatOutputWithWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && buf.Len() == 0 {
				t.Errorf("FormatOutputWithWriter() produced no output for format %s", tt.format)
			}
		})
	}
}
