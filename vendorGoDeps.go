package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	command string
	errs    []string
	deps    Godeps
)

// Dependency represents a godeps dependency.
type Dependency struct {
	ImportPath string `json:"ImportPath"`
	Revision   string `json:"Rev"`
}

func (dep *Dependency) extractRepoName() string {
	path := strings.Split(dep.ImportPath, "/")
	if len(path) <= 3 {
		return dep.ImportPath
	}
	return strings.Join(path[:3], "/")
}

// Godeps represents a set of godeps dependencies.
type Godeps struct {
	Deps      []Dependency `json:"Deps"`
	GoVersion string       `json:"GoVersion"`
}

func init() {
	command = "git"
}

func main() {
	fileContent, err := ioutil.ReadFile("Godeps/godeps.json")

	if err != nil {
		log.Fatalf("Unable to read godeps.json file: %v\n", err)
	}

	if err := json.Unmarshal([]byte(fileContent), &deps); err != nil {
		log.Fatal(err)
	}

	var submodules = make(map[string]string)

	for _, dep := range deps.Deps {
		if dep.ImportPath != "" {
			repoName := dep.extractRepoName()
			submodules[repoName] = dep.Revision
		}
	}

	fmt.Println("vendoring submodules")
	fmt.Println("please run the git checkout commmand below to use the right versions in GoDeps")

	for repoPath, hash := range submodules {
		time.Sleep(1 * time.Second)
		vendorPath := "vendor/" + repoPath
		args := []string{"submodule", "add", "http://" + repoPath, vendorPath}
		fmt.Println("pushd", vendorPath, " && git checkout ", hash, " && popd")
		if err := exec.Command(command, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			errs = append(errs, repoPath)
			continue
		}
	}

	if len(errs) > 0 {
		for _, skipped := range errs {
			fmt.Println("unable to vendor", skipped)
		}
	}

	if len(errs) == 0 {
		fmt.Println("successfully vendored all dependencies")
		fmt.Println("deleting Godeps folder")
		if err := exec.Command(command, "rm", "-r", "Godeps").Run(); err != nil {
			log.Fatalf("error removing Godeps folder: %v\n", err)
		}
		if err := exec.Command(command, "add", ".").Run(); err != nil {
			log.Fatalf("error staging deleted Godeps folder: %v\n", err)
		}
		if err := exec.Command(command, "commit", "-m", "deleted Godeps folder and added gitmodules file").Run(); err != nil {
			log.Fatalf("error committing deleted Godeps folder: %v\n", err)
		}
	}
}
