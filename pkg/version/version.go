package version

import "github.com/coreos/go-semver/semver"

var (
	GitRepository string
	GitCommit     string
	VersionMajor  int64 = 0
	VersionMinor  int64 = 5
	VersionPatch  int64 = 0
	VersionPre          = ""
	VersionDev    string
)

var Version = semver.Version{
	Major:      VersionMajor,
	Minor:      VersionMinor,
	Patch:      VersionPatch,
	PreRelease: semver.PreRelease(VersionPre),
	Metadata:   VersionDev,
}
