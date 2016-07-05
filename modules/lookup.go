package zdns

import (
	"flag"
	"os"
	"fmt"
	"strings"
	"github.com/zmap/zdns"
	"github.com/zmap/zdns/modules"
)

type Lookup interface {

	DoLookup(name string) (interface {}, error)

}

type GenericLookup struct {

	NameServers *[]string
	Timeout int
}


type LookupFactory interface {
	// expected to add any necessary commandline flags if being
	// run as a standalone scanner
	AddFlags(flags *flag.FlagSet)
	// global initialization. Gets called once globally
	// This is called after command line flags have been parsed
	Initialize(conf *conf.GlobalConf) bool
	// We can't set variables on an interface, so write functions
	// that define any settings for the factory
	AllowStdIn() (bool)
	// Return a single scanner which will scan a single host
	MakeLookup() (Lookup)

}

type GenericLookupFactory struct {

	NameServers *[]string
	Timeout int
}

func (l GenericLookupFactory) Initialize(c *conf.GlobalConf) bool {
	return true
}

func (s GenericLookupFactory) AllowStdIn() bool {
	return true
}


// keep a mapping from name to factory
var Lookups map[string]LookupFactory;

func RegisterLookup(name string, s LookupFactory) {
	if Lookups == nil {
		Lookups = make(map[string]LookupFactory, 100)
	}
	Lookups[name] = s
}

func ValidLookupsString() string {

	valid := make([]string, len(Lookups))
	for k, _ := range Lookups {
		valid = append(valid, k)
		fmt.Println("loaded module:", k)
	}
	return strings.Join(valid, ", ")
}

func GetLookup(name string) *LookupFactory {

	factory, ok := Lookups[name]
	if !ok {
		fmt.Println("[error] Invalid module:", os.Args[1], "Valid modules:", ValidLookupsString())
		os.Exit(1)
	}
	return &factory
}
