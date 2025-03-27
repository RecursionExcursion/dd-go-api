package wsd

type compliationPair struct {
	os   string
	arch string
}

var compilationPairs = map[string]compliationPair{
	"win": {os: "windows", arch: "amd64"},
	"lin": {os: "linux", arch: "arm64"},
	"dar": {os: "darwin", arch: "arm64"},
}
