package testdata

import "embed"

//go:embed esport/*
//go:embed football/*
//go:embed nba/*
var Content embed.FS
