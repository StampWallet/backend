package database

//go:generate $GOPATH/bin/mockgen --destination mocks/mocks.go --build_flags=--mod=mod . BusinessAuthorizedAccessor,UserAuthorizedAccessor,AuthorizedTransactionAccessor
