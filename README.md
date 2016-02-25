### VENDORGODEPS

This tool was built to enable the easy migration from the [Godeps](https://github.com/tools/godep) dependency management tooling to the vendor folder approach mandated by the Go team. The way you get the code into that
folder is not mandated, but we use [git submodules](https://git-scm.com/book/en/v2/Git-Tools-Submodules), so this tool follows suit.

The tool performs the following steps:

* Reads the contents of the Godeps/godeps.json file.
* For each dependency, it gets the import path and commit hash.
* It then proceeds to add each one as a git submodule.
* If any errors are encountered, the program will list them and exit.
* If all dependencies have been vendored, the program proceeds to delete the Godeps folder.
* Bathe in the glory of mandated vendoring.

