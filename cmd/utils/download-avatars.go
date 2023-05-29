package main

import (
	"fmt"
	"os/exec"

	"github.com/alecthomas/kong"
	"github.com/letsblockit/letsblockit/data"
)

const (
	outputFolder      = "data/assets/images/contributors"
	outputPathPattern = outputFolder + "/%s.webp"
)

type downloadAvatarsCmd struct {
	All bool `help:"Download all contributor avatars."`
}

func (c *downloadAvatarsCmd) Run(k *kong.Context) error {
	list, err := data.ParseContributors()
	k.FatalIfErrorf(err, "failed to parse contributors.json")
	contributors := list.GetAll()
	k.Printf("Downloading avatars for %d contributors", len(contributors))

	for _, contributor := range contributors {
		targetFile := fmt.Sprintf(outputPathPattern, contributor.Login)
		execOrFatal(k, "magick", contributor.AvatarUrl,
			"-quality", "80", "-define", "webp:image-hint=picture",
			"-define", "webp:method=6", "-define", "webp:alpha-filtering=2",
			"-resize", "96x96", targetFile)
	}
	execOrFatal(k, "git", "add", outputFolder)

	return nil
}

func execOrFatal(k *kong.Context, name string, args ...string) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	k.FatalIfErrorf(err, "Command %s failed:\n:%s", cmd, string(output))
}
