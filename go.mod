module github.com/amorenoz/ovnmodel

go 1.15

require (
	github.com/cenk/hub v1.0.1 // indirect
	github.com/ebay/go-ovn v0.0.0-00010101000000-000000000000
	github.com/ebay/libovsdb v0.2.0
	github.com/fatih/color v1.10.0
	github.com/k0kubun/pp v2.4.0+incompatible
)

replace (
	github.com/ebay/go-ovn => github.com/amorenoz/go-ovn v0.1.1-0.20210305093751-39f4443b448d
	github.com/ebay/libovsdb => github.com/amorenoz/libovsdb v0.2.1-0.20210305093614-80f588f3e2db
)
