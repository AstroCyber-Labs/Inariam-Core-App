package main

import (
	"gitea/pcp-inariam/inariam/core/cli"
	_ "gitea/pcp-inariam/inariam/core/services/api/docs"
	_ "gitea/pcp-inariam/inariam/pkgs/log"
)

func main() {
	cli.Execute()
}
