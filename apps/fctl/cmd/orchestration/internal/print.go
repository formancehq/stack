package internal

import (
	"io"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
)

func PrintWorkflowInstance(out io.Writer, w shared.Workflow, instance shared.WorkflowInstance) error {
	fctl.Section.WithWriter(out).Println("Stages")

	ind := 0
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(out).
		WithData(
			fctl.Prepend(
				fctl.Map(instance.Status,
					func(src shared.StageStatus) []string {
						stage := w.Config.Stages[ind]
						var name string
						for name = range stage {
						}
						ind = ind + 1
						return []string{
							name,
							src.StartedAt.Format(time.RFC3339),
							func() string {
								if src.TerminatedAt != nil {
									return src.TerminatedAt.Format(time.RFC3339)
								}
								return ""
							}(),
							func() string {
								if src.Error != nil {
									return *src.Error
								}
								return ""
							}(),
						}
					}),
				[]string{"Name", "Started at", "Terminated at", "Error"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}
	return nil
}
