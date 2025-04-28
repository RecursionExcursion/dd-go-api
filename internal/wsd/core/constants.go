package core

/* Supported os and arch */

const win64 = "win64"
const linux64 = "linux64"
const mac64 = "mac64"

type compliationPair struct {
	os   string
	arch string
}

var SupportedArchitecture = map[string]compliationPair{
	win64:   {os: "windows", arch: "amd64"},
	linux64: {os: "linux", arch: "arm64"},
	mac64:   {os: "darwin", arch: "arm64"},
}
