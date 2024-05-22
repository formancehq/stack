package webhookrunner

import (
	_ "errors"
	"testing"

	"github.com/formancehq/webhooks/internal/commons"
	_ "github.com/formancehq/webhooks/internal/commons"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
)

func TestRunTosaveProcess(t *testing.T){

	t.Run("ToSaveProcess", tosaveProcess)

}


func tosaveProcess(t *testing.T){

	newHook := commons.NewHook("Test", []string{"testevent"}, "/localhost/", "", false)
	err := Database.SaveHook(*newHook)
	require.NoError(t, err)

	attempt1 := commons.NewSharedAttempt(newHook.ID, newHook.Name, newHook.Endpoint, "eventtest", "payload")
	attempt2 := commons.NewSharedAttempt(newHook.ID, newHook.Name, newHook.Endpoint, "eventtest", "payload")
	attempt3 := commons.NewSharedAttempt(newHook.ID, newHook.Name, newHook.Endpoint, "eventtest", "payload")
	
	Database.SaveAttempt(*attempt1.Val)
	Database.SaveAttempt(*attempt2.Val)
	Database.SaveAttempt(*attempt3.Val)

	commons.SetSuccesStatus(attempt1.Val)
	commons.SetAbortDisableHook(attempt2.Val)
	commons.SetAbortNoRetryModeStatus(attempt3.Val)

	WebhookRunner.State.ToSaveAttempts.Add(attempt1)
	WebhookRunner.State.ToSaveAttempts.Add(attempt2)
	WebhookRunner.State.ToSaveAttempts.Add(attempt3)

	WebhookRunner.ToSaveProcess()

	ra1, err := Database.GetAttempt(attempt1.Val.ID)
	require.NoError(t, err)
	ra2, err := Database.GetAttempt(attempt2.Val.ID)
	require.NoError(t, err)
	ra3, err := Database.GetAttempt(attempt3.Val.ID)
	require.NoError(t, err)

	require.Equal(t, commons.SuccessStatus, ra1.Status)
	require.Equal(t, commons.AbortStatus, ra2.Status)
	require.Equal(t, commons.AbortDisabledHook, ra2.Comment)
	require.Equal(t, commons.AbortStatus, ra3.Status)
	require.Equal(t, commons.AbortNoRetryMode, ra3.Comment)
}