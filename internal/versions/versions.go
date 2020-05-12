package versions

import (
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

type VersionInfo struct {
	Current       string
	Available     []string
	UpToDate      bool
	NonComparable bool
}

func Compare(c string, a []string) VersionInfo {
	vi := VersionInfo{Current: c, UpToDate: true}

	have, err := version.NewVersion(c)
	if err != nil {
		vi.NonComparable = true
		return vi
	}

	// The practice of having git hashes as your release candidate makes version
	// comparison difficult.  Instead we will ignore these hashes.
	// This could break comparison over very long release candidates, e.g. rc100000
	// But this kind of value is unlikely to exist in the Prerelease
	// We check 8-digit endings starting with `20` to catch most dates
	gitHash := regexp.MustCompile("[a-f0-9]{7,}$")
	date := regexp.MustCompile("20[0-9]{6}$")
	if gitHash.MatchString(have.Prerelease()) && !date.MatchString(have.Prerelease()) {
		noPreVersion := strings.Replace(c, have.Prerelease(), "", -1)
		newHave, err := version.NewVersion(noPreVersion)
		if err == nil {
			have = newHave
		}
	}

	versions := []*version.Version{}
	for i := range a {
		v, err := version.NewVersion(a[i])
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}
	sort.Sort(sort.Reverse(version.Collection(versions)))

	for _, v := range versions {
		if have.LessThan(v) {
			vi.Available = append(vi.Available, v.Original())
			vi.UpToDate = false
		}
		if len(vi.Available) > 5 {
			break
		}
	}

	return vi
}
