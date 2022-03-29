package composer

import "time"

type Records struct {
	Version string `bson:"version,omitempty"`
	Name    string `bson:"name,omitempty"`
	// DevDependencies map[string]string `bson:"devDependencies,omitempty"`
	Require map[string]string `bson:"require,omitempty"`
	Licence []string            `bson:"license,omitempty"`
}

func NewRecords() *Records {
	return &Records{
		Version: "",
		Name:    "",
		Licence: []string{},
		// DevDependencies: map[string]string{},
		Require: map[string]string{},
	}
}

type Metas struct {
	Version string    `bson:"version,omitempty"`
	Name    string    `bson:"name,omitempty"`
	Time    time.Time `bson:"time,omitempty"`
}
