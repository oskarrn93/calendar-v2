package testdata

import "embed"

//go:embed basketball/*
//go:embed esport/*
//go:embed football/*
//go:embed nba/*
var Content embed.FS
