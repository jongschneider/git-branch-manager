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
	termWidth := getTerminalWidth()

	// Set minimum column widths for responsive design
	minWidths := map[string]int{
		"WORKTREE":     12,
		"BRANCH":       15,
		"GIT STATUS":   12,
		"SYNC STATUS":  12,
		"PATH":         20,
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
	if t.terminalWidth >= 120 {
		// Wide terminal - show all columns
		indices := make([]int, len(t.headers))
		for i := range indices {
			indices[i] = i
		}
		return t.headers, indices
	}

	// Narrow terminal - prioritize columns and omit PATH if present
	prioritizedHeaders := make([]string, 0)
	columnIndices := make([]int, 0)

	// Always include these columns first
	for i, header := range t.headers {
		if header != "PATH" {
			prioritizedHeaders = append(prioritizedHeaders, header)
			columnIndices = append(columnIndices, i)
		}
	}

	// If there's still room and PATH exists, add it
	if t.terminalWidth >= 100 {
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
