package modules

import (
	"fmt"
)

func ShellCommand(cmd string, args ...any) []string {
	return []string{"sh", "-c",
		fmt.Sprintf(`/bin/sh <<'EOF'

		`+cmd+`
 						    
EOF`, args...)}
}
