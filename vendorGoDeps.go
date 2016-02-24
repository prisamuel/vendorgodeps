package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	importPath *regexp.Regexp
	command    string
	errs       []string
)

func init() {
	command = "git"
}

func extractRepoName(importPath string) string {
	path := strings.Split(importPath, "/")
	if len(path) <= 3 {
		return importPath
	}
	return strings.Join(path[:3], "/")
}

func main() {
	fileContent, err := ioutil.ReadFile("Godeps/godeps.json")

	if err != nil {
		fmt.Println("Unable to read file")
		panic(err)
	}

	//Dependencies ...
	type Dependencies struct {
		ImportPath string `json:"ImportPath"`
		Revision   string `json:"Rev"`
	}

	type Godeps struct {
		Deps      []Dependencies `json:"Deps"`
		GoVersion string         `json:"GoVersion"`
	}

	var deps Godeps

	if err := json.Unmarshal([]byte(fileContent), &deps); err != nil {
		panic(err)
	}

	var submodules = make(map[string]string)

	for _, dep := range deps.Deps {
		if dep.ImportPath != "" {
			repoName := extractRepoName(dep.ImportPath)
			submodules[repoName] = dep.Revision
		}
	}

	fmt.Println("vendoring submodules")
	fmt.Println("Please run the git checkout commmand below to use the right versions in GoDeps")

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

	if len(errs) > 1 {
		for _, skipped := range errs {
			fmt.Println("unable to vendor", skipped)
		}
	}

	if len(errs) == 0 {
		fmt.Println("successfully vendored all dependencies")
		fmt.Println("deleting Godeps folder")
		if err := exec.Command(command, "rm", "-r", "Godeps").Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if err := exec.Command(command, "add", ".").Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if err := exec.Command(command, "commit", "-m", "deleted Godeps folder and added gitmodules file").Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}