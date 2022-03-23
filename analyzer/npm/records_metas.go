package npm

import "time"

type Records struct {
	Version         string            `bson:"version,omitempty"`
	Name            string            `bson:"name,omitempty"`
	// DevDependencies map[string]string `bson:"devDependencies,omitempty"`
	Dependencies    map[string]string `bson:"dependencies,,omitempty"`
	Licence         string            `bson:"license,omitempty"`
}

func NewRecords() *Records {
	return &Records{
		Version: "",
		Name: "",
		Licence: "",
		// DevDependencies: map[string]string{},
		Dependencies:    map[string]string{},
	}
}

type Metas struct {
	Version string    `bson:"version,omitempty"`
	Name    string    `bson:"name,omitempty"`
	Time    time.Time `bson:"time,omitempty"`
}
