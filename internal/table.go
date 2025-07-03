package internal

import (
	"fmt"
	"strings"
)

type Table struct {
	Headers []string
	Rows    [][]string
	widths  []int
}

func NewTable(headers []string) *Table {
	return &Table{
		Headers: headers,
		Rows:    make([][]string, 0),
		widths:  make([]int, len(headers)),
	}
}

func (t *Table) AddRow(row []string) {
	if len(row) != len(t.Headers) {
		panic(fmt.Sprintf("row length %d does not match header length %d", len(row), len(t.Headers)))
	}
	t.Rows = append(t.Rows, row)
}

func (t *Table) calculateWidths() {
	// Initialize with header widths
	for i, header := range t.Headers {
		t.widths[i] = len(header)
	}

	// Check all rows for maximum width
	for _, row := range t.Rows {
		for i, cell := range row {
			if len(cell) > t.widths[i] {
				t.widths[i] = len(cell)
			}
		}
	}
}

func (t *Table) Print() {
	if len(t.Rows) == 0 {
		return
	}

	t.calculateWidths()

	// Ensure minimum column spacing
	spacing := 4

	// Print header
	for i, header := range t.Headers {
		fmt.Printf("%-*s", t.widths[i], header)
		if i < len(t.Headers)-1 {
			fmt.Print(strings.Repeat(" ", spacing))
		}
	}
	fmt.Println()

	// Print separator
	totalWidth := 0
	for i, width := range t.widths {
		totalWidth += width
		if i < len(t.widths)-1 {
			totalWidth += spacing
		}
	}
	fmt.Println(strings.Repeat("-", totalWidth))

	// Print rows
	for _, row := range t.Rows {
		for i, cell := range row {
			fmt.Printf("%-*s", t.widths[i], cell)
			if i < len(row)-1 {
				fmt.Print(strings.Repeat(" ", spacing))
			}
		}
		fmt.Println()
	}
}