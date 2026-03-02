package airesource

import "fmt"

const (
	APIVersionDraft = "ai-resource/draft"
)

var supportedVersions = []string{APIVersionDraft}

func IsSupportedVersion(version string) bool {
	for _, v := range supportedVersions {
		if v == version {
			return true
		}
	}
	return false
}

func UnsupportedVersionError(version string) error {
	return fmt.Errorf("unsupported apiVersion: %s (supported: %v)", version, supportedVersions)
}
