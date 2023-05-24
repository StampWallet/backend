package managers

//go:generate $GOPATH/bin/mockgen --destination mocks/mocks.go --build_flags=--mod=mod . AuthManager,BusinessManager,ItemDefinitionManager,LocalCardManager,TransactionManager,VirtualCardManager
