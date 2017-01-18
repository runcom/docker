package registeredtransports

import (
	"fmt"

	"github.com/containers/image/types"
)

// KnownTransports is a registry of known ImageTransport instances.
var KnownTransports map[string]types.ImageTransport

func init() {
	KnownTransports = make(map[string]types.ImageTransport)
}

// Register TODO(runcom)
func Register(t types.ImageTransport) {
	name := t.Name()
	if _, ok := KnownTransports[name]; ok {
		panic(fmt.Sprintf("Duplicate image transport name %s", name))
	}
	KnownTransports[name] = t
}

// ImageName converts a types.ImageReference into an URL-like image name, which MUST be such that
// ParseImageName(ImageName(reference)) returns an equivalent reference.
//
// This is the generally recommended way to refer to images in the UI.
//
// NOTE: The returned string is not promised to be equal to the original input to ParseImageName;
// e.g. default attribute values omitted by the user may be filled in in the return value, or vice versa.
func ImageName(ref types.ImageReference) string {
	return ref.Transport().Name() + ":" + ref.StringWithinTransport()
}
