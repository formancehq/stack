package controllerutils

func JustError[T any](t T, err error) error {
	return err
}
