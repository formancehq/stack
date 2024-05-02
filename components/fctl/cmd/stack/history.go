package stack

import (
	"errors"
	"strings"

	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	"github.com/formancehq/fctl/pkg/printer"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

const (
	pageSizeFlag = "page-size"
	cursorFlag   = "cursor"

	actionFlag      = "action"
	dataFlag        = "data"
	userIdFlag      = "user-id"
	displayDataFlag = "display-data"
)

type HistoryStore struct {
	Cursor *membershipclient.LogCursorData `json:"cursor"`
}

type HistoryController struct {
	store *HistoryStore
}

var _ fctl.Controller[*HistoryStore] = (*HistoryController)(nil)

func NewDefaultHistoryStore() *HistoryStore {
	return &HistoryStore{
		Cursor: &membershipclient.LogCursorData{},
	}
}
func NewHistoryController() *HistoryController {
	return &HistoryController{
		store: NewDefaultHistoryStore(),
	}
}

func NewHistoryCommand() *cobra.Command {
	return fctl.NewMembershipCommand("history [id]",
		fctl.WithShortDescription("Query stack history"),
		fctl.WithAliases("hist"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag(actionFlag, "", "Filter on Action"),
		fctl.WithStringFlag(userIdFlag, "", "Filter on UserId, use SYSTEM to filter on system logs"),
		fctl.WithStringFlag(dataFlag, "", "Filter on modified Data with --data key=value, key is a jsonb text path"),

		fctl.WithBoolFlag(displayDataFlag, false, "Display data"),

		fctl.WithStringFlag(cursorFlag, "", "Cursor"),
		fctl.WithIntFlag(pageSizeFlag, 10, "Page size"),
		fctl.WithPreRunE(func(cmd *cobra.Command, args []string) error {
			return fctl.CheckMembershipVersion("v0.29.0")(cmd, args)
		}),
		fctl.WithController(NewHistoryController()),
	)
}
func (c *HistoryController) GetStore() *HistoryStore {
	return c.store
}

func (c *HistoryController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := store.GetStore(cmd.Context())
	pageSize := fctl.GetInt(cmd, pageSizeFlag)
	stackID := args[0]
	req := store.Client().ListLogs(cmd.Context(), store.OrganizationId()).PageSize(int32(pageSize)).StackId(stackID)

	cursor := fctl.GetString(cmd, cursorFlag)
	userID := fctl.GetString(cmd, userIdFlag)
	action := fctl.GetString(cmd, actionFlag)
	data := fctl.GetString(cmd, dataFlag)
	if cursor != "" {
		if userID != "" || action != "" || data != "" || stackID != "" {
			return nil, errors.New("cursor can't be used with other flags")
		}

		req = req.Cursor(cursor)
	}

	if stackID == "" && cursor == "" {
		return nil, errors.New("stack-id or cursor is required")
	}

	if userID != "" {
		req = req.UserId(userID)
	}

	if action != "" {
		if !strings.Contains(action, "stacks") {
			return nil, errors.New("stacks history are scoped to 'stacks.*' actions")
		}
		req = req.Action(membershipclient.Action(action))
	}

	if data != "" {
		keyVal := strings.Split(data, "=")
		if len(keyVal) != 2 {
			return nil, errors.New("data filter must be in the form key=value")
		}

		req = req.Key(keyVal[0]).Value(keyVal[1])
	}

	log, res, err := req.Execute()
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		return nil, errors.New("error listing stack logs")
	}

	c.store.Cursor = &log.Data
	return c, nil
}

func (c *HistoryController) Render(cmd *cobra.Command, args []string) error {
	return printer.LogCursor(cmd.OutOrStdout(), c.store.Cursor, fctl.GetBool(cmd, displayDataFlag))
}
