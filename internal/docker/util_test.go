package docker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseIdentifierNormal(t *testing.T) {
	id := ParseIdentifier("testowner/testimage:testtag")
	expectedId := Image{
		Tag:   "testtag",
		Owner: "testowner",
		Image: "testimage",
	}
	assert.Equal(t, id, expectedId)
}

func TestParseIdentifierNoTag(t *testing.T) {
	id := ParseIdentifier("testowner/testimage")
	expectedId := Image{
		Tag:   "latest",
		Owner: "testowner",
		Image: "testimage",
	}
	assert.Equal(t, id, expectedId)
}

func TestParseIdentifierNoOwner(t *testing.T) {
	id := ParseIdentifier("testimage:testtag")
	expectedId := Image{
		Tag:   "testtag",
		Image: "testimage",
	}
	assert.Equal(t, id, expectedId)
}

func TestParseIdentifierNoOwnerNoTag(t *testing.T) {
	id := ParseIdentifier("testimage")
	expectedId := Image{
		Tag:   "latest",
		Image: "testimage",
	}
	assert.Equal(t, id, expectedId)
}

func TestParseIdentifierURL(t *testing.T) {
	id := ParseIdentifier("quay.io/testowner/testimage")
	expectedId := Image{
		Tag:   "latest",
		Owner: "quay.io/testowner",
		Image: "testimage",
	}
	assert.Equal(t, id, expectedId)
}

func TestURLNormal(t *testing.T) {
	image := Image{
		Tag:   "testtag",
		Owner: "testowner",
		Image: "testimage",
	}
	url := image.URL()
	assert.Equal(t, url, "https://hub.docker.com/r/testowner/testimage")
}

func TestURLNoOwner(t *testing.T) {
	image := Image{
		Tag:   "testtag",
		Image: "testimage",
	}
	url := image.URL()
	assert.Equal(t, url, "https://hub.docker.com/_/testimage")
}

func TestURLSlashOwner(t *testing.T) {
	image := Image{
		Tag:   "testtag",
		Owner: "quay.io/testowner",
		Image: "testimage",
	}
	url := image.URL()
	assert.Equal(t, url, "https://quay.io/testowner/testimage")
}

func TestRegistryInfoDocker(t *testing.T) {
	image := Image{
		Tag:   "testtag",
		Owner: "testowner",
		Image: "testimage",
	}
	registry := image.Registry()
	expectedRegistry := RegistryInfo{"docker-hub", "https://registry-1.docker.io/"}
	assert.Equal(t, registry, expectedRegistry)
}

func TestRegistryInfoQuay(t *testing.T) {
	image := Image{
		Tag:   "testtag",
		Owner: "quay.io/testowner",
		Image: "testimage",
	}
	registry := image.Registry()
	expectedRegistry := RegistryInfo{"quay-io", "https://quay.io/"}
	assert.Equal(t, registry, expectedRegistry)
}

func TestRepoStrNormal(t *testing.T) {
	image := Image{
		Tag:   "testtag",
		Owner: "testowner",
		Image: "testimage",
	}
	repo := image.RepoStr()
	assert.Equal(t, repo, "testowner/testimage:testtag")
}

func TestRepoStrNoOwner(t *testing.T) {
	image := Image{
		Tag:   "testtag",
		Image: "testimage",
	}
	repo := image.RepoStr()
	assert.Equal(t, repo, "testimage:testtag")
}
