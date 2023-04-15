package buildinfo

import "runtime/debug"

var BuildInfo *debug.BuildInfo

func init() {
	var ok bool
	BuildInfo, ok = debug.ReadBuildInfo()
	if !ok {
		panic("cannot read build info")
	}
}

func GetSentryVersion() string {
	for _, module := range BuildInfo.Deps {
		if module.Path == "" {
			return module.Version
		}
	}
	return "not found"
}
