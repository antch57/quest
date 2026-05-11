// Package tablefmt provides helpers for rendering CLI tables.
package tablefmt

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// Render writes a rounded text table to w using the provided headers and rows.
//
// Each row is appended as-is. If row widths differ from the header width,
// go-pretty handles the layout during rendering.
func Render(w io.Writer, headers []string, rows [][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault

	headerRow := make(table.Row, len(headers))
	for i, header := range headers {
		headerRow[i] = header
	}
	t.AppendHeader(headerRow)

	for _, row := range rows {
		tableRow := make(table.Row, len(row))
		for i, cell := range row {
			tableRow[i] = cell
		}
		t.AppendRow(tableRow)
	}

	t.Render()
}
