package utils

func GetNatsSubject(organizationID, stackID string) string {
	return organizationID + "." + stackID
}
