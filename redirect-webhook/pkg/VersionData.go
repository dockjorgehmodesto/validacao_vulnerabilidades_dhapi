package pkg

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const FileName = ".version"

type VersionFile struct {
	Version VersionData `json:"version"`
}

type VersionData struct {
	Current    string `json:"current"`
	Hash       string `json:"hash"`
	LastCommit string `json:"lastcommit"`
}

func ParseFileFromBranch(branch Branch) VersionFile {
	return VersionFile{
		Version: VersionData{
			Current:    branch.getVersion(),
			Hash:       branch.hash,
			LastCommit: "",
		},
	}
}

func GetCurrentVersionFile(workdir string) (VersionFile, error) {
	bytes, e := ioutil.ReadFile(workdir + FileName)
	if e != nil {
		return VersionFile{}, e
	}

	return ParseFile(bytes)
}

func ParseFile(bytes []byte) (VersionFile, error) {
	currentFile := VersionFile{}
	e := yaml.Unmarshal(bytes, &currentFile)
	if e != nil {
		return VersionFile{}, errors.Wrap(e, "error unmarshaling yaml data")
	}
	return currentFile, nil
}

func (v VersionFile) ShouldUpdate(workdir string) (bool, error) {
	currentFile, e := GetCurrentVersionFile(workdir)
	if e != nil {
		return false, e
	}

	return v.Version.Current != currentFile.Version.Current, nil
}

func (v VersionFile) Update(workdir string) error {
	out, _ := yaml.Marshal(v)
	e := ioutil.WriteFile(workdir+FileName, out, 0644)
	return e
}
