module github.com/AndreasSko/go-jwlm

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.2.12
	github.com/AndreasSko/jwpub-snippets v0.0.0-20210704143853-4c0b907b6dd3
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Jeffail/gabs v1.4.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Netflix/go-expect v0.0.0-20200312175327-da48e75238e2
	github.com/buger/goterm v1.0.0
	github.com/cavaliercoder/grab v1.0.1-0.20201108051000-98a5bfe305ec
	github.com/codeclysm/extract/v3 v3.0.2
	github.com/davecgh/go-spew v1.1.1
	github.com/hinshun/vt10x v0.0.0-20180809195222-d55458df857c
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/klauspost/compress v1.13.1
	github.com/kr/pty v1.1.8 // indirect
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20210611083646-a4fc73990273 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/AndreasSko/jwpub-snippets v0.0.0-20210627134355-912640c387c2 => ./publication/jwpub-snippets
