package flags

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/probr/probr/internal/core"
)

// Run flags relate to the primary probr execution
var Run *flag.FlagSet

// List flags manage the view of installed binaries
var List *flag.FlagSet

// Version flags relate to the version information for this probr installation
var Version *flag.FlagSet

func init() {
	Run = flag.NewFlagSet("probr", flag.ExitOnError)
	core.ConfigPath = Run.String("config-file", defaultConfigPath(), "Location for service pack binaries.")
	addBinariesFlag(Run)
	addAllFlag(Run)

	List = flag.NewFlagSet("probr list", flag.ExitOnError)
	addBinariesFlag(List)
	addAllFlag(List)

	Version = flag.NewFlagSet("probr version", flag.ExitOnError)
	addVerboseFlag(Version)
}

func addBinariesFlag(flagSet *flag.FlagSet) {
	core.BinariesPath = flagSet.String("binaries-path", filepath.Join(core.UserHomeDir(), "probr", "binaries"), "Location for service pack binaries.")
}

func addVerboseFlag(flagSet *flag.FlagSet) {
	core.Verbose = flagSet.Bool("v", false, "Display extended version information")
}

func addAllFlag(flagSet *flag.FlagSet) {
	core.AllPacks = flagSet.Bool("all", false, "Include all installed packs, not just those specified within the provided config")
}

func defaultConfigPath() string {
	workDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(workDir, "config.yml")
}
