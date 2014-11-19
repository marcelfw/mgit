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

func FormatRow(name string, value string) interface{} {
	lines := strings.Split(value, "\n")

	switch {
	case len(lines) == 0 || (len(lines) == 1 && value == ""):
		columns := make([]string, 2, 2)
		columns[0] = name
		columns[1] = "<no output>"
		return columns
	case len(lines) == 1:
		columns := make([]string, 2, 2)
		columns[0] = name
		columns[1] = value
		return columns
	default:
		rows := make([][]string, 0, len(lines))
		for idx, line := range lines {
			columns := make([]string, 2, 2)
			var pre string // pre is used to hopefully make it easier to see the lines belong together
			switch {
			case idx == 0:
				columns[0] = name
				pre = "   "
			case idx == len(lines)-1:
				pre = "\\_ "
			default:
				pre = "|  "
			}
			columns[1] = pre + strings.TrimSpace(line)
			rows = append(rows, columns)
		}
		return rows
	}

	return nil
}
