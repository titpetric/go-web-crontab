module github.com/titpetric/go-web-crontab

go 1.27.0

require (
	github.com/apex/log v1.9.0
	github.com/go-bridget/mig v0.4.4
	github.com/jmoiron/sqlx v1.4.0
	github.com/namsral/flag v1.7.4-pre
	github.com/robfig/cron/v3 v3.0.1
	github.com/titpetric/pdo v0.2.1
	github.com/titpetric/phpscript v0.3.0
	github.com/titpetric/platform v0.6.0
	golang.org/x/net v0.58.0
	modernc.org/sqlite v1.57.0
)

require (
	github.com/a-h/templ v0.3.1020 // indirect
	github.com/dlclark/regexp2 v1.12.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/expr-lang/expr v1.17.8 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/go-chi/chi/v5 v5.3.1 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/titpetric/oida v0.1.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
	modernc.org/libc v1.75.3 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.0 // indirect
)

replace github.com/titpetric/platform => /root/workspace/github/platform

replace github.com/titpetric/phpscript => /root/workspace/titpetric-1.27/phpscript

replace github.com/titpetric/oida => /root/workspace/titpetric-1.27/oida

replace github.com/titpetric/pdo => /root/workspace/titpetric-1.27/pdo
