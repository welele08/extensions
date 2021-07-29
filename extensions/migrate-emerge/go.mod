module github.com/Luet-lab/extensions/extensions/migrate-emerge

go 1.14

require (
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/mudler/luet v0.0.0-20201004175813-2cb0f3ab5ddf
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/docker/docker => github.com/Luet-lab/moby v17.12.0-ce-rc1.0.20200605210607-749178b8f80d+incompatible
