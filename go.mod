module github.com/amorenoz/ovnmodel

go 1.15

require (
	github.com/abiosoft/readline v0.0.0-20180607040430-155bce2042db // indirect
	github.com/cenk/hub v1.0.1 // indirect
	github.com/ebay/go-ovn v0.0.0-00010101000000-000000000000
	github.com/ebay/libovsdb v0.2.0
	github.com/fatih/color v1.10.0
	github.com/flynn-archive/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/k0kubun/pp v2.4.0+incompatible
	github.com/olekukonko/tablewriter v0.0.5
	gopkg.in/abiosoft/ishell.v2 v2.0.0
)

replace (
	github.com/ebay/go-ovn => github.com/amorenoz/go-ovn v0.1.1-0.20210309183059-6707d7ad0208
	github.com/ebay/libovsdb => github.com/amorenoz/libovsdb v0.2.1-0.20210305093614-80f588f3e2db
)
