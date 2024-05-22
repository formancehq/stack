package storage

import (
	"fmt"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
)

type Table struct {
	Name    string
	Columns map[string]string
}

func (t *Table) ListColumns() string {
	return strings.Join(collectionutils.Values(t.Columns), ",")
}

func (t *Table) ExcludedRow() string {
	str := ""
	for key, column := range t.Columns {
		if key == "ID" {
			continue
		}
		str += fmt.Sprintf("%s = EXCLUDED.%s,", column, column)
	}
	str = str[:len(str)-1] + ";"
	return str

}

func StrArray(arr []string) string {
	var sb strings.Builder

	for i, str := range arr {
		sb.WriteString("")
		sb.WriteString(str)
		sb.WriteString("")
		if i < len(arr)-1 {
			sb.WriteString(",")
		}
	}

	return fmt.Sprintf("{%s}", sb.String())
}
