package versions

import (
	"sort"

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
