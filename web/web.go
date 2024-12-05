package web

import (
	"embed"
	"io/fs"
)

//go:embed go-auth-admin/auth-admin/assets/*
var fsAuthAdminAssetsFS embed.FS

func MustAuthAdminAssetsFS() fs.FS {
	res, err := fs.Sub(fsAuthAdminAssetsFS, "go-auth-admin/auth-admin/assets")
	if err != nil {
		panic(err)
	}
	return res
}

//go:embed  go-auth-admin/index.html
var fsAuthAdminIndexHTML embed.FS

func MustAuthAdminIndexHTML() string {

	data, err := fsAuthAdminIndexHTML.ReadFile("go-auth-admin/index.html")
	if err != nil {
		panic(err)
	}

	return string(data)
}
