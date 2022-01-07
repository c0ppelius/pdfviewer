package src

import "embed"

//go:embed *
var Assets embed.FS

// separate package for embedded assets as embed provides no native facility for moving them up the directory
