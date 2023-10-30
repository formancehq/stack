package storage_test

import "testing"

func TestPaymentsStorage(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)

	testUninstallConnectors(t, store)
	testAccountsDeletedAfterConnectorUninstall(t, store)
}
