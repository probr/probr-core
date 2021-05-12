package flags

import (
	"flag"
	"path/filepath"

	"github.com/citihub/probr-core/internal/core"
)

// Core relates to the primary probr execution
var Core *flag.FlagSet

// List allows users to view their installed binaries
var List *flag.FlagSet

// Version displays information about this probr core installation
var Version *flag.FlagSet

func init() {
	Core = flag.NewFlagSet("probr", flag.ExitOnError)
	addBinariesFlag(Core)
	addAllFlag(Core)

	List = flag.NewFlagSet("probr list", flag.ExitOnError)
	addBinariesFlag(List)
	addAllFlag(List)

	Version = flag.NewFlagSet("probr version", flag.ExitOnError)
	addVerboseFlag(Version)
}

func addBinariesFlag(flagSet *flag.FlagSet) {
	core.BinariesPath = flagSet.String("binaries-path", filepath.Join(core.UserHomeDir(), "probr", "binaries"), "Location for service pack binaries. If not provided, default value is: [UserHomeDir]/probr/binaries")
}

func addVerboseFlag(flagSet *flag.FlagSet) {
	core.Verbose = flagSet.Bool("v", false, "Display extended version information")
}

func addAllFlag(flagSet *flag.FlagSet) {
	core.AllPacks = flagSet.Bool("all", false, "Include all installed packs, not just those specified within config.yml")
}
