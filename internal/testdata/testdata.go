package testdata

import "embed"

//go:embed football/*
//go:embed nba/*
var Content embed.FS
