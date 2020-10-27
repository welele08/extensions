package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	pack "github.com/mudler/luet/pkg/package"

	entropy "github.com/Sabayon/pkgs-checker/pkg/entropy"
	compiler "github.com/mudler/luet/pkg/compiler"
	"gopkg.in/yaml.v2"
)

// #cgo LDFLAGS: -lpthread

func main() {

	dbPath := os.Getenv("ENTROPY_DB")
	if dbPath == "" {
		dbPath = "/var/lib/entropy/client/database/amd64/equo.db"
	}

	luet := os.Getenv("LUET")
	if luet == "" {
		luet = "luet"
	}

	dir, err := ioutil.TempDir(os.TempDir(), "prefix")
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("DEBUG") != "true" {
		defer os.RemoveAll(dir)
	}

	// Retrive list of installed packages (return []EntropyPackage)
	pkgs, err := entropy.RetrieveRepoPackages(dbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Found", len(pkgs), "packages")
	if os.Getenv("DEBUG") == "true" {
		fmt.Println(pkgs)
	}

	brokenPkg := []string{}

	for idx, pkg := range pkgs {
		index := idx + 1
		fmt.Println(fmt.Sprintf("[ %3d / %3d ] Retrieving data for %s",
			index, len(pkgs), pkg))

		// Retrieve pkg detail (EntropyPackageDetail)
		detail, err := entropy.RetrievePackageData(pkg, dbPath)
		if err != nil {
			fmt.Println(err)
			brokenPkg = append(brokenPkg, fmt.Sprintf("%s: %s", pkg.String(), err.Error()))
			continue
		}
		// print list of files
		if os.Getenv("DEBUG") == "true" {
			fmt.Println("[", index, "/", len(pkgs), "]", "files: ", detail.Files)
		}

		var files []string

		for _, f := range detail.Files {
			files = append(files, strings.TrimPrefix(f, "/"))
		}

		// We use category with slot when slot != 0
		category := pkg.Category
		if pkg.Slot != "0" {
			slot := pkg.Slot
			// Ignore sub-slot
			if strings.Contains(pkg.Slot, "/") {
				slot = pkg.Slot[0:strings.Index(slot, "/")]
			}
			category = fmt.Sprintf("%s-%s", category, slot)
		}

		version := pkg.Version
		if pkg.VersionSuffix != "" {
			version += pkg.VersionSuffix
		}

		a := compiler.PackageArtifact{
			CompileSpec: &compiler.LuetCompilationSpec{Package: &pack.DefaultPackage{Name: pkg.Name, Category: category, Version: version}},
			Files:       files,
		}

		data, err := yaml.Marshal(a)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		metadata := filepath.Join(dir, a.GetCompileSpec().GetPackage().GetFingerPrint()+".metadata.yaml")
		fmt.Println(fmt.Sprintf("[ %3d / %3d ] Generating metadata for %s:%s (%d) at %s",
			index, len(pkgs), pkg, pkg.Slot, len(files), metadata))
		err = ioutil.WriteFile(metadata, data, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(fmt.Sprintf("[ %3d / %3d ] Creating db entry for %s",
			index, len(pkgs), pkg))

		cmd := exec.Command(luet, "database", "create", metadata)
		cmd.Env = os.Environ()
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Failed creating package", pkg, err.Error())
			fmt.Println(string(out))
			brokenPkg = append(brokenPkg, fmt.Sprintf("%s: %s - %s", pkg.String(), string(out), err.Error()))
			continue
		}

		fmt.Println(fmt.Sprintf("[ %3d / %3d ] %s",
			index, len(pkgs), string(out)))
	}

	if len(brokenPkg) > 0 {
		fmt.Println(fmt.Sprintf("%d Packages not imported:", len(brokenPkg)))
		for _, pkg := range brokenPkg {
			fmt.Println(pkg)
		}
		fmt.Println("Some packages are not loaded correctly. Please, share your experience to the community")
	}
	fmt.Println("All done.")
}
