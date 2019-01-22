// +build ignore

package main

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/32leaves/ruruku/cmd"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
date: %s
title: "%s"
series: commands
menu: %s
hideFromIndex: true
---
`

func main() {
	filePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		title := strings.Replace(base, "_", " ", -1)
		inmenu := "false"
		if base == "ruruku" {
			title = "Command Line Interface"
			inmenu = "true"
		}
		return fmt.Sprintf(fmTemplate, now, title, inmenu)
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/cli/" + strings.ToLower(base) + "/"
	}

	err := doc.GenMarkdownTreeCustom(cmd.GetRoot(), "www/content/cli", filePrepender, linkHandler)
	if err != nil {
		panic(err)
	}
}
