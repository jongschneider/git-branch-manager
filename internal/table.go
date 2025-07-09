package internal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Table struct {
	table           *table.Table
	headers         []string
	rows            [][]string
	terminalWidth   int
	minColumnWidths map[string]int
}

func NewTable(headers []string) *Table {
	termWidth := GetTerminalWidth()

	// Set minimum column widths for responsive design
	minWidths := map[string]int{
		"WORKTREE":     12,
		"BRANCH":       30, // Branch names can be long
		"GIT STATUS":   12,
		"SYNC STATUS":  12,
		"PATH":         30, // Paths are long
		"STATUS":       12,
		"ENV VARIABLE": 15,
		"ENV VAR":      12,
		"ISSUES":       20,
	}

	return &Table{
		headers:         headers,
		rows:            make([][]string, 0),
		terminalWidth:   termWidth,
		minColumnWidths: minWidths,
	}
}

// NewTestTable creates a table with a specific terminal width for testing
func NewTestTable(headers []string, termWidth int) *Table {
	// Set minimum column widths for responsive design
	minWidths := map[string]int{
		"WORKTREE":     12,
		"BRANCH":       30, // Branch names can be long
		"GIT STATUS":   12,
		"SYNC STATUS":  12,
		"PATH":         30, // Paths are long
		"STATUS":       12,
		"ENV VARIABLE": 15,
		"ENV VAR":      12,
		"ISSUES":       20,
	}

	return &Table{
		headers:         headers,
		rows:            make([][]string, 0),
		terminalWidth:   termWidth,
		minColumnWidths: minWidths,
	}
}

func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

func (t *Table) Print() {
	t.buildTable()
	if t.table == nil {
		return
	}
	fmt.Println(t.table.String())
}

func (t *Table) String() string {
	t.buildTable()
	if t.table == nil {
		return ""
	}
	return t.table.String()
}

// buildTable creates the responsive table based on terminal width
func (t *Table) buildTable() {
	// Calculate responsive headers and data
	responsiveHeaders, columnIndices := t.getResponsiveHeaders()

	// Build the table with responsive headers
	t.table = table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(TableBorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return TableHeaderStyle
			}
			return TableCellStyle
		}).
		Headers(responsiveHeaders...)

	// Add rows with only the columns that fit
	for _, row := range t.rows {
		responsiveRow := make([]string, len(columnIndices))
		for i, colIndex := range columnIndices {
			if colIndex < len(row) {
				responsiveRow[i] = row[colIndex]
			}
		}
		t.table.Row(responsiveRow...)
	}
}

// getResponsiveHeaders returns headers and column indices that fit in the terminal
func (t *Table) getResponsiveHeaders() ([]string, []int) {
	// Calculate estimated width needed for each column
	estimatedWidth := t.calculateEstimatedWidth()

	// If we have enough width for all columns, show everything
	// But be more conservative - require extra margin for readability
	if estimatedWidth <= t.terminalWidth-20 {
		indices := make([]int, len(t.headers))
		for i := range indices {
			indices[i] = i
		}
		return t.headers, indices
	}

	// Narrow terminal - prioritize columns and omit PATH if present
	prioritizedHeaders := make([]string, 0)
	columnIndices := make([]int, 0)

	// Always include these columns first (in priority order)
	priorityOrder := []string{"ENV VARIABLE", "WORKTREE", "BRANCH", "GIT STATUS", "SYNC STATUS", "STATUS"}

	for _, priorityHeader := range priorityOrder {
		for i, header := range t.headers {
			if header == priorityHeader {
				prioritizedHeaders = append(prioritizedHeaders, header)
				columnIndices = append(columnIndices, i)
				break
			}
		}
	}

	// Calculate width without PATH column
	widthWithoutPath := t.calculateEstimatedWidthForHeaders(prioritizedHeaders)

	// If there's still room and PATH exists, add it
	if widthWithoutPath+t.minColumnWidths["PATH"]+10 <= t.terminalWidth { // 10 for borders and spacing
		for i, header := range t.headers {
			if header == "PATH" {
				prioritizedHeaders = append(prioritizedHeaders, header)
				columnIndices = append(columnIndices, i)
				break
			}
		}
	}

	return prioritizedHeaders, columnIndices
}

// calculateEstimatedWidth estimates the total width needed for all columns
func (t *Table) calculateEstimatedWidth() int {
	totalWidth := 0
	for i, header := range t.headers {
		// Calculate the maximum width needed for this column
		maxWidth := len(header) // Header width

		// Check content width in all rows
		for _, row := range t.rows {
			if i < len(row) {
				contentWidth := len(row[i])
				if contentWidth > maxWidth {
					maxWidth = contentWidth
				}
			}
		}

		// Use minimum width if content is smaller
		if minWidth, exists := t.minColumnWidths[header]; exists && minWidth > maxWidth {
			maxWidth = minWidth
		} else if !exists && maxWidth < 15 {
			maxWidth = 15 // default column width
		}

		totalWidth += maxWidth
	}
	// Add space for borders and padding (approximately 3 chars per column for borders and spacing)
	// Plus 4 for the outer borders
	return totalWidth + len(t.headers)*3 + 4
}

// calculateEstimatedWidthForHeaders estimates width for specific headers
func (t *Table) calculateEstimatedWidthForHeaders(headers []string) int {
	totalWidth := 0
	for _, header := range headers {
		// Find the column index for this header
		colIndex := -1
		for i, h := range t.headers {
			if h == header {
				colIndex = i
				break
			}
		}

		if colIndex == -1 {
			continue // Header not found
		}

		// Calculate the maximum width needed for this column
		maxWidth := len(header) // Header width

		// Check content width in all rows
		for _, row := range t.rows {
			if colIndex < len(row) {
				contentWidth := len(row[colIndex])
				if contentWidth > maxWidth {
					maxWidth = contentWidth
				}
			}
		}

		// Use minimum width if content is smaller
		if minWidth, exists := t.minColumnWidths[header]; exists && minWidth > maxWidth {
			maxWidth = minWidth
		} else if !exists && maxWidth < 15 {
			maxWidth = 15 // default column width
		}

		totalWidth += maxWidth
	}
	// Add space for borders and padding
	return totalWidth + len(headers)*3 + 4
}
