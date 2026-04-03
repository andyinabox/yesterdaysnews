package youtubedownloader

import (
	"context"
	"fmt"
	"strings"
)

var requiredVersion binVersion

func init() {
	requiredVersion = NewVersion(2026, 03, 03)
}

func (c *Client) CheckVersion(ctx context.Context) error {
	result, err := c.shell.Execute(ctx, fmt.Sprintf("%s --version", c.binPath))
	if err != nil {
		return fmt.Errorf("error determining yt-dlp version: %w", err)
	}

	versionString := strings.TrimSpace(string(result))

	installedVersion, err := ParseVersion(versionString)
	if err != nil {
		return fmt.Errorf("error parsing yt-dlp version %q: %w", versionString, err)
	}

	if !installedVersion.IsEqualOrGreaterThan(requiredVersion) {
		return fmt.Errorf("incompatible yt-dlp version. have %s, need %s or higher", installedVersion, requiredVersion)
	}

	return nil
}
