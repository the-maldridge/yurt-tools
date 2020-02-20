package docker

// An Image is the complete specification that can be fed to the
// docker command line to define a single image to pull and work with.
type Image struct {
	Owner string
	Image string
	Tag   string
}
