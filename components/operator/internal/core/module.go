package core

func ValidateInstalledVersion(module Module) bool {
	labels := module.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	if labels["formance.com/installed-version"] == module.GetVersion() {
		return false
	}
	labels["formance.com/installed-version"] = module.GetVersion()
	module.SetLabels(labels)

	return true
}

func GetInstalledVersion(module Module) string {
	return module.GetLabels()["formance.com/installed-version"]
}
