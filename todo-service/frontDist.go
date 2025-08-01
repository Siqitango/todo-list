package todoservice

import (
	"embed"
)

//go:embed dist/*
var FrontDist embed.FS
