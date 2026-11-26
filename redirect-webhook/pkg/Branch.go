package pkg

import "strings"

const (
	VersionBranchCandidate = "Candidate/"
	VersionBranchHotfix    = "Hotfix/"
	VersionBranchBugfix    = "Bugfix/"
)

type Branch struct {
	name string
	hash string
}

func (b Branch) isVersionBranch() bool {
	return strings.Index(b.name, VersionBranchCandidate) == 0 || strings.Index(b.name, VersionBranchHotfix) == 0 || strings.Index(b.name, VersionBranchBugfix) == 0
}

func (b Branch) getVersion() string {
	split := strings.Split(b.name, "/")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}
