module github.com/Luet-lab/extensions/extensions/package-browser

go 1.14

require (
	github.com/Sabayon/pkgs-checker v0.7.3-0.20201029211214-b71c01e603ee // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/mudler/luet v0.0.0-20201004175813-2cb0f3ab5ddf
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	gopkg.in/macaron.v1 v1.3.9
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/docker/docker => github.com/Luet-lab/moby v17.12.0-ce-rc1.0.20200605210607-749178b8f80d+incompatible
