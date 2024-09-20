package printer

import (
	"fmt"
	"io"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/go-libs/time"
	"github.com/pterm/pterm"
)

func LogCursor(writer io.Writer, cursor *membershipclient.LogCursorData, withData bool) error {
	header := []string{"Identifier", "User", "Date", "Action"}

	if withData {
		header = append(header, "Data")
	}
	tableData := fctl.Map(cursor.Data, func(log membershipclient.Log) []string {
		line := []string{
			log.Seq,
			func() string {
				if log.UserId == "" {
					return "SYSTEM"
				}
				return log.UserId
			}(),
			log.Date.Format(time.DateFormat),
			log.Action,
		}

		if withData {
			line = append(line, func() string {
				if log.Data == nil {
					return ""
				}
				return fmt.Sprintf("%v", log.Data)
			}())

		}

		return line
	})
	tableData = fctl.Prepend(tableData, header)

	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(writer).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return Cursor(writer, &membershipclient.Cursor{
		HasMore:  cursor.HasMore,
		PageSize: cursor.PageSize,
		Next:     cursor.Next,
		Previous: cursor.Previous,
	})
}
