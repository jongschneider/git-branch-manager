package internal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Table struct {
	table *table.Table
}

func NewTable(headers []string) *Table {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(TableBorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return TableHeaderStyle
			}
			return TableCellStyle
		}).
		Headers(headers...)

	return &Table{
		table: t,
	}
}

func (t *Table) AddRow(row []string) {
	t.table.Row(row...)
}

func (t *Table) Print() {
	if t.table == nil {
		return
	}
	fmt.Println(t.table.String())
}

func (t *Table) String() string {
	if t.table == nil {
		return ""
	}
	return t.table.String()
}

