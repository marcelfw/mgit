// Copyright (c) 2014 Marcel Wouters

// Package engine implements the engine.
// This source generates the nice output.
package engine

import (
	"bytes"
	"strings"
)

// Output an text string table.
func ReturnTextTable(header []string, rows [][]string) string {
	var buffer bytes.Buffer

	// Storage for column widths and line.
	var column_width []int
	var line_columns []string

	// Init column width header columns.
	if header != nil {
		column_width = make([]int, len(header))
		line_columns = make([]string, len(header))

		for idx, column := range header {
			column_width[idx] = len(column)
		}
	}

	// Determine column widths.
	for _, row := range rows {
		if len(column_width) == 0 {
			column_width = make([]int, len(row))
			line_columns = make([]string, len(row))
		}

		for idx, column := range row {
			if len(column) > column_width[idx] {
				column_width[idx] = len(column)
			}
		}
	}

	if header != nil {
		// Fill line columns.
		for idx, _ := range header {
			line_columns[idx] = strings.Repeat("-", column_width[idx])
		}

		// Inserts header and lines into rows.
		rows = append(rows, header, header)
		copy(rows[2:], rows[0:len(rows)-1])
		rows[0] = header
		rows[1] = line_columns
	}

	// Write actual columns.
	no_of_columns := len(column_width)
	for _, row := range rows {
		for idx, column := range row {
			if idx > 0 {
				buffer.WriteString("  ")
			}

			buffer.WriteString(column)

			if idx < (no_of_columns-1) && len(column) < column_width[idx] {
				buffer.WriteString(strings.Repeat(" ", column_width[idx]-len(column)))
			}
		}

		buffer.WriteString("\n")
	}

	return buffer.String()
}
