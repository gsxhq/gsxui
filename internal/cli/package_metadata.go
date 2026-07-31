package cli

import (
	"fmt"
	"os"
)

func packageMetadataArtifacts(root string) ([]artifact, error) {
	artifacts := make([]artifact, 0, 2)
	for _, relative := range []string{"package.json", "package-lock.json"} {
		target, err := artifactPath(root, relative)
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(target)
		if err != nil {
			return nil, fmt.Errorf("stage %s for transactional npm integration: %w", relative, err)
		}
		artifacts = append(artifacts, artifact{
			RelativePath:       relative,
			Content:            content,
			AlwaysTrack:        true,
			MutableAfterCommit: true,
		})
	}
	return artifacts, nil
}
