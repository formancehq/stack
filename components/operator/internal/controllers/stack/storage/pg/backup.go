package pg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/formancehq/operator/apis/stack/v1beta3"
)

const (
	compression = 5
)

func BackupDatabase(database string, conf v1beta3.PostgresConfig) ([]byte, error) {
	args := []string{
		"-h",
		conf.Host,
		"-p",
		fmt.Sprint(conf.Port),
		"-U",
		conf.Username,
		"-d",
		database,
		"--compress",
		fmt.Sprint(compression),
	}

	cmd := exec.Command("pg_dump", args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", conf.Password))

	// This should output as logger.Error only if we are in debug mode ?
	// cmd.Stderr = os.Stderr
	// stderr := cmd.StderrPipe()

	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	return cmd.Output()
}
