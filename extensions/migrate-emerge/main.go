package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	compiler "github.com/mudler/luet/pkg/compiler"
	pack "github.com/mudler/luet/pkg/package"
	"gopkg.in/yaml.v2"
)

func getPackagefiles(p *pack.DefaultPackage) ([]string, error) {
	//  /var/db/pkg/app-misc/c_rehash-1.7-r1/CONTENTS
	// dir /usr
	// dir /usr/bin
	// obj /usr/bin/c_rehash ebf2b33b613c94eda4ef5c457c751a91 1608616028

	files := []string{}
	contentsB, err := ioutil.ReadFile(filepath.Join("/var/db/pkg", p.GetCategory(), fmt.Sprintf("%s-%s", p.GetName(), p.GetVersion()), "CONTENTS"))
	if err != nil {
		return files, err
	}
	lines := strings.Split(string(contentsB), "\n")
	for _, l := range lines {
		data := strings.Split(l, " ")
		if data[0] == "obj" {
			files = append(files, strings.TrimPrefix(data[1], "/"))
		}
	}
	return files, nil
}

func getGentooInstalledPackages() ([]*pack.DefaultPackage, error) {
	var packages []*pack.DefaultPackage

	cmd := exec.Command("qlist", "-I", "--format", "%{CATEGORY}/%{PN}@%{PVR}")
	out, err := cmd.Output()
	if err != nil {
		return packages, err
	}

	packs := strings.Split(string(out), "\n")
	for _, p := range packs {
		if p == "" {
			continue
		}
		pieces := strings.Split(p, "/")
		if len(pieces) != 2 {
			return packages, errors.New("Invalid count" + p)
		}
		data := strings.Split(pieces[1], "@")
		packages = append(packages, &pack.DefaultPackage{
			Name:     data[0],
			Version:  data[1],
			Category: pieces[0],
		})
	}
	return packages, nil
}

func main() {

	luet := os.Getenv("LUET")
	if luet == "" {
		luet = "luet"
	}
	var dir string

	if os.Getenv("MIGRATE_OUTPUT") != "" {
		dir = os.Getenv("MIGRATE_OUTPUT")
	} else {
		tempdir, err := ioutil.TempDir(os.TempDir(), "migrate-emerge")
		if err != nil {
			log.Fatal(err)
		}
		dir = tempdir
		if os.Getenv("DEBUG") != "true" {
			defer os.RemoveAll(dir)
		}
	}

	// qlist -I --format "%{CATEGORY}/%{PN}@%{PV}"
	pkgs, err := getGentooInstalledPackages()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Found", len(pkgs), "packages")
	fmt.Println(pkgs[0])
	if os.Getenv("DEBUG") == "true" {
		fmt.Println(pkgs)
	}

	brokenPkg := []string{}

	for i, pkg := range pkgs {
		index := i + 1
		fmt.Printf("[ %3d / %3d ] Retrieving data for %s\n",
			index, len(pkgs), pkg)

		files, err := getPackagefiles(pkg)
		if err != nil {
			fmt.Println("Failed getting package files", pkg, err.Error())
			brokenPkg = append(brokenPkg, fmt.Sprintf("%s: %s", pkg.String(), err.Error()))
			continue
		}
		a := compiler.PackageArtifact{
			CompileSpec: &compiler.LuetCompilationSpec{
				Package: pkg,
			},
			Files: files,
		}

		data, err := yaml.Marshal(a)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		metadata := filepath.Join(dir, a.GetCompileSpec().GetPackage().GetFingerPrint()+".metadata.yaml")
		fmt.Printf("[ %3d / %3d ] Generating metadata for %s (%d) at %s\n",
			index, len(pkgs), pkg, len(files), metadata)
		err = ioutil.WriteFile(metadata, data, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("[ %3d / %3d ] Creating db entry for %s\n",
			index, len(pkgs), pkg)

		cmd := exec.Command(luet, "database", "create", metadata)
		cmd.Env = os.Environ()
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Failed creating package", pkg, err.Error())
			fmt.Println(string(out))
			brokenPkg = append(brokenPkg, fmt.Sprintf("%s: %s - %s", pkg.String(), string(out), err.Error()))
			continue
		}

		fmt.Printf("[ %3d / %3d ] %s\n",
			index, len(pkgs), string(out))
	}

	if len(brokenPkg) > 0 {
		fmt.Printf("%d Packages not imported:\n", len(brokenPkg))
		for _, pkg := range brokenPkg {
			fmt.Println(pkg)
		}
		fmt.Println("Some packages are not loaded correctly. Please, share your experience to the community")
	}
	fmt.Println("All done.")
}
