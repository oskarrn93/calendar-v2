package testdata

import "embed"

//go:embed nba/*
var Content embed.FS
