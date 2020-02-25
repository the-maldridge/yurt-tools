package docker

import (
	"path"
	"strings"
)

// ParseIdentifier figures out what the different parts of the id mean
// and splits them apart.  This allows other parts of the system to
// cleanly assert that parts of the Image struct are set.
func ParseIdentifier(in string) Image {
	di := Image{}

	// First set the tag for the image.
	if !strings.Contains(in, ":") {
		di.Tag = "latest"
	} else {
		// Tag must be after a colon
		di.Tag = strings.FieldsFunc(in, func(c rune) bool { return c == ':' })[1]
	}

	// Split the image and owner.
	parts := strings.Split(in, "/")
	imgAndTag := parts[len(parts)-1]

	i := strings.SplitN(imgAndTag, ":", 2)
	di.Image = i[0]

	if len(parts) > 1 {
		// If there is more than one part, it must be from the
		// owner
		di.Owner = strings.Join(parts[:len(parts)-1], "/")
	}

	return di
}

func (i Image) URL() string {
	// If the owner is unset, it must be from the docker hub
	// library as that is where docker will pull unqualified image
	// names from.
	if i.Owner == "" {
		return "https://hub.docker.com/_/" + i.Image
	}

	// If there is exactly 1 slash in the name, it is likely from
	// hub.docker.com.  To be sure we could check for any dots in
	// the name, but this should be sufficient for most cases.  If
	// you're here because it wasn't, send a PR extending this
	// check.
	switch strings.Count(i.Owner, "/") {
	case 1:
		return "https://hub.docker.com/r/" + i.Owner + "/" + i.Image
	default:
		return "https://" + i.Owner + "/" + i.Image
	}
}

// Registry returns the registry that owns the requested image.
// This allows pulling registries from a map rather than trying to
// initialize a new registry for every context.
func (i Image) Registry() RegistryInfo {
	imgURL := i.URL()
	switch {
	case strings.Contains(imgURL, "quay.io"):
		return RegistryInfo{"quay-io", "https://quay.io/"}
	default:
		return RegistryInfo{"docker-hub", "https://registry-1.docker.io/"}
	}
}

// RepoStr returns the string that can be fed directly to docker to
// pull the image.
func (i Image) RepoStr() string {
	if i.Owner == "" {
		return i.Image + ":" + i.Tag
	}
	return path.Join(i.Owner, i.Image) + ":" + i.Tag
}
