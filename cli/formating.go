package cli

import (
	"text/tabwriter"
	"os"
	"fmt"
)

const perPage = 25

type Tabler interface {
	Row(n int) []interface{}
}

func Table(format string, headers []interface{}, separators []interface{}, rows []Tabler) int {
	table := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Println()
	fmt.Fprintf(table, format, headers...)
	fmt.Fprintf(table, format, separators...)
	for index := 0; index < len(rows); index++ {
		row := rows[index]
		fmt.Fprintf(table, format, row.Row(index+1)...)

		if (index+1)%perPage == 0 {
			table.Flush()
			selected := Unswer(len(rows), "Press enter to show next objects or select object...")
			if selected != 0 {
				return selected-1
			}
		}
	}
	table.Flush()
	fmt.Println()
	return MustUnswer(len(rows), "Select object: ")-1
}
