package main

import (
	"go-auth-admin/internal/cmd"
	"go-auth-admin/internal/config"

	"go-auth-admin/internal/config/consts"
	xlog "go-auth-admin/internal/util/utillog"
)

//nolint:gochecknoglobals
var (
	Version     = "" //  "1.0.0"
	ShortCommit = "" // "1a2b3c4"
	Commit      = "" // "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0"
	Date        = ""
)

func main() {

	xlog.Info("Build info: [name: %v] [version: %v] [date: %v] [short-commit: %v]", consts.AppName, Version, Date, ShortCommit)

	config.AppVersion, config.AppCommit, config.AppDate, config.ShortCommit = Version, Commit, Date, ShortCommit

	config.ReadFlags()
	//
	x := cmd.Command{}

	x.Exec()
}
