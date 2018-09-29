package cmd

var usage = `
Usage:
	wikinote [options] <command> [<args>...]
	wikinote [options] serve [<args>...]
	wikinote [--token=<token>] [--addr=<addr>] (user|config) [<args>...]

Options:
	--version
	--help
	--token=<token>, -t=<token>  token
	--addr=<addr>                address [default: localhost:4000]
	-D, --debug                  debug options

Commands:
	serve   start server
	user    manage user
	config  config

Examples:
	wikinote -D serve
	wikinote --token=admin:admin config set auto-save=true
	wikinote user add user1 -p password --role admin
	wikinote user del user1
	wikinote user assign user2 as editor
	wikinote user list -f role=admin

`
var serveUsage = `
Usage:
	wikinote serve [--config=<config-file>] [--wiki-path=<wiki-path>] [--bind=<bind-addr>]

options:
	--config=<config-file>, -c=<config-file>  config file path     [default: $HOME/wiki/.app/config.yaml]
	--wiki-path=<wiki-path>, -w=<wiki-path>   wiki path            [default: $HOME/wiki]
	--bind=<bind-addr>                        bind address         [default: :4000]
`
var userUsage = `
Usage:
	wikinote user add <name> [-p=<password>]
	wikinote user set <name> [--email=<email>]
	wikinote user assign <name> as <role>
	wikinote user password <name>
	wikinote user list [--filter=<filter>]
`
var configUsage = `
Usage:
	wikinote config set <key> <value>
	wikinote config get <key>
`
