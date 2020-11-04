package main

// e.g. FORMAT=json TREE="$HOME/_git/portage-tree/multi-arch $HOME/_git/portage-tree/amd64" go run main.go layers/system-x layers/mate layers/gnome

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	helpers "github.com/mudler/luet/cmd/helpers"
	pkg "github.com/mudler/luet/pkg/package"
	"github.com/mudler/luet/pkg/solver"
	"github.com/mudler/luet/pkg/tree"

	"github.com/mudler/luet/pkg/compiler"
	"github.com/mudler/luet/pkg/compiler/backend"
)

func rankMapStringInt(values map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range values {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	ranked := make([]string, len(values))
	for i, kv := range ss {
		ranked[i] = kv.Key
	}
	return ranked
}

type BuildTree struct {
	order map[string]int
}

func (bt *BuildTree) Increase(s string) {
	if bt.order == nil {
		bt.order = make(map[string]int)
	}
	if _, ok := bt.order[s]; !ok {
		bt.order[s] = 0
	}

	bt.order[s]++
}

func (bt *BuildTree) Reset(s string) {
	if bt.order == nil {
		bt.order = make(map[string]int)
	}
	bt.order[s] = 0
}

func (bt *BuildTree) Level(s string) int {
	return bt.order[s]
}

func ints(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func (bt *BuildTree) AllInLevel(l int) []string {
	var all []string
	for k, v := range bt.order {
		if v == l {
			all = append(all, k)
		}
	}
	return all
}

func (bt *BuildTree) Order(compilationTree map[string]map[string]interface{}) {
	sentinel := false
	for !sentinel {
		sentinel = true

	LEVEL:
		for _, l := range bt.AllLevels() {

			for _, j := range bt.AllInLevel(l) {
				for _, k := range bt.AllInLevel(l) {
					if j == k {
						continue
					}
					if _, ok := compilationTree[j][k]; ok {
						if bt.Level(j) == bt.Level(k) {
							bt.Increase(j)
							sentinel = false
							break LEVEL
						}
					}
				}
			}
		}
	}
}

func (bt *BuildTree) AllLevels() []int {
	var all []int
	for _, v := range bt.order {
		all = append(all, v)
	}

	sort.Sort(sort.IntSlice(all))

	return ints(all)
}

func (bt *BuildTree) JSON() (string, error) {
	type buildjob struct {
		Jobs []string `json:"packages"`
	}

	result := []buildjob{}
	for _, l := range bt.AllLevels() {
		result = append(result, buildjob{Jobs: bt.AllInLevel(l)})
	}
	dat, err := json.Marshal(&result)
	return string(dat), err
}

func main() {
	opts := compiler.NewDefaultCompilerOptions()
	compilerSpecs := compiler.NewLuetCompilationspecs()
	solverOpts := solver.Options{Type: solver.SingleCoreSimple, Concurrency: 1}
	compilerBackend := backend.NewSimpleDockerBackend()

	db := pkg.NewInMemoryDatabase(false)
	bt := &BuildTree{}
	generalRecipe := tree.NewCompilerRecipe(db)
	luetCompiler := compiler.NewLuetCompiler(compilerBackend, generalRecipe.GetDatabase(), opts, solverOpts)

	for _, src := range strings.Split(os.Getenv("TREE"), " ") {
		err := generalRecipe.Load(src)
		if err != nil {
			log.Fatal("Error: " + err.Error())
		}
	}

	for _, a := range os.Args[1:] {

		pack, err := helpers.ParsePackageStr(a)
		if err != nil {
			log.Fatal("Invalid package string ", a, ": ", err.Error())
		}

		spec, err := luetCompiler.FromPackage(pack)
		if err != nil {
			log.Fatal("Error: " + err.Error())
		}
		compilerSpecs.Add(spec)
	}

	//toCalculate, _ := luetCompiler.ComputeMinimumCompilableSet(compilerSpecs.All()...)

	compilationOrder := map[string]int{}
	compilationTree := map[string]map[string]interface{}{}

	for _, sp := range compilerSpecs.All() {
		//for _, sp := range toCalculate {
		//		fmt.Println(sp.GetPackage().HumanReadableString())

		ass, err := luetCompiler.ComputeDepTree(sp)
		if err != nil {
			log.Fatal("Error: " + err.Error())
		}

		bt.Reset(fmt.Sprintf("%s/%s", sp.GetPackage().GetCategory(), sp.GetPackage().GetName()))

		for _, p := range ass {

			bt.Reset(fmt.Sprintf("%s/%s", p.Package.GetCategory(), p.Package.GetName()))

			//	fmt.Println(ass.Package.HumanReadableString())
			compilationOrder[fmt.Sprintf("%s/%s", p.Package.GetCategory(), p.Package.GetName())]++

			spec, err := luetCompiler.FromPackage(p.Package)
			if err != nil {
				log.Fatal("Error: " + err.Error())
			}
			ass, err := luetCompiler.ComputeDepTree(spec)
			if err != nil {
				log.Fatal("Error: " + err.Error())
			}
			for _, r := range ass {
				if compilationTree[fmt.Sprintf("%s/%s", p.Package.GetCategory(), p.Package.GetName())] == nil {
					compilationTree[fmt.Sprintf("%s/%s", p.Package.GetCategory(), p.Package.GetName())] = make(map[string]interface{})
				}
				compilationTree[fmt.Sprintf("%s/%s", p.Package.GetCategory(), p.Package.GetName())][fmt.Sprintf("%s/%s", r.Package.GetCategory(), r.Package.GetName())] = nil
			}
			if compilationTree[fmt.Sprintf("%s/%s", sp.GetPackage().GetCategory(), sp.GetPackage().GetName())] == nil {
				compilationTree[fmt.Sprintf("%s/%s", sp.GetPackage().GetCategory(), sp.GetPackage().GetName())] = make(map[string]interface{})
			}
			compilationTree[fmt.Sprintf("%s/%s", sp.GetPackage().GetCategory(), sp.GetPackage().GetName())][fmt.Sprintf("%s/%s", p.Package.GetCategory(), p.Package.GetName())] = nil
		}
		//	compilationOrder[sp.GetPackage().HumanReadableString()]++
	}

	bt.Order(compilationTree)

	if os.Getenv("FORMAT") == "json" {
		data, err := bt.JSON()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(data)
	} else {
		for _, l := range bt.AllLevels() {
			fmt.Println(strings.Join(bt.AllInLevel(l), " "))
		}
	}

}
