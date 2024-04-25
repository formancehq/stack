package core

import "fmt"

func ShellScript(cmd string, args ...any) []string {
	return []string{"sh", "-c",
		fmt.Sprintf(`/bin/sh <<'EOF'
		set -x
		`+cmd+`
 						    
EOF`, args...)}
}
