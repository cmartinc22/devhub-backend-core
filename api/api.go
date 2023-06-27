package api

import "embed"

//go:embed */schemas/*
//go:embed api.html
var ApiContent embed.FS
