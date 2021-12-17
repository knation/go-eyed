// Used to create, manage, and parse IDs
// Kirk Morales

package eyed

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"

	"github.com/segmentio/ksuid"
)

// EYEDs will be 28 characters (including the underscore) plus the
// length of the prefix.
type EYED struct {
	idType *EYEDType
	ksuid  ksuid.KSUID
}

// EYEDType is used to differentiate between different ID types
// with different prefixes.
type EYEDType struct {
	name   string
	prefix string
	check  *regexp.Regexp
}

// Holds a map of all prefixes to their associated types
var registeredTypes = make(map[string](EYEDType))

// Registers a new distinct type of EYED
func RegisterType(name string, prefix string) *EYEDType {
	idType := EYEDType{name, prefix, regexp.MustCompile("^" + prefix + "_[a-zA-Z0-9]{27}")}
	registeredTypes[prefix] = idType

	return &idType
}

// Gets the EYEDType for the given ID string
func GetType(id string) (*EYEDType, bool) {
	idParts := strings.Split(id, "_")
	if len(idParts) != 2 {
		return nil, false
	}

	idType := registeredTypes[idParts[0]]
	if idType.name != "" {
		return &idType, true
	} else {
		return nil, false
	}
}

// Parses the given ID string and returns an EYED object
func Parse(id string) (EYED, bool) {
	newEYED := EYED{}

	idParts := strings.Split(id, "_")
	if len(idParts) != 2 {
		return newEYED, false
	}

	idType := registeredTypes[idParts[0]]
	if idType.name == "" {
		return newEYED, false
	}

	parsedKsuid, err := ksuid.Parse(idParts[1])
	if err != nil {
		return newEYED, false
	}

	newEYED.idType = &idType
	newEYED.ksuid = parsedKsuid

	return newEYED, true
}

// Creates a new EYED of the give type
func (idType *EYEDType) New() EYED {
	return EYED{idType, ksuid.New()}
}

// Returns the name of an EYEDType
func (idType *EYEDType) Name() string {
	return idType.name
}

// Returns the prefix of an EYEDType
func (idType *EYEDType) Prefix() string {
	return idType.prefix
}

// Checks if the given ID matches the type
func (idType *EYEDType) Is(id string) bool {
	return idType.check.Match([]byte(id))
}

func (id *EYEDType) Value() (driver.Value, error) {
	matched := id.Is([]byte(p))
	if !matched {
		return driver.Value(""), fmt.Errorf("Number '%s' not a valid PhoneNumber format.", p)
	}
	return driver.Value(string(p)), nil
}

// Returns the string reprsentation of the ID
func (id EYED) String() string {
	return id.idType.prefix + "_" + id.ksuid.String()
}

// Returns the EYEDType for the given EYED
func (id EYED) Type() *EYEDType {
	return id.idType
}

// Returns the KSUID for the given EYED
func (id EYED) Ksuid() ksuid.KSUID {
	return id.ksuid
}
