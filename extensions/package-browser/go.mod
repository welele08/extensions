module github.com/Luet-lab/extensions/extensions/package-browser

go 1.14

require (
	github.com/Sabayon/pkgs-checker v0.7.3-0.20201029211214-b71c01e603ee // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/mudler/luet v0.0.0-20210210080702-b93357e36c74
	github.com/narqo/go-badge v0.0.0-20190124110329-d9415e4e1e9f
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	gopkg.in/macaron.v1 v1.3.9
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

replace github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.0-rc9.0.20200221051241-688cf6d43cc4

replace github.com/docker/docker => github.com/Luet-lab/moby v17.12.0-ce-rc1.0.20200605210607-749178b8f80d+incompatible
