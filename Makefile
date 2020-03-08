cmd/jcquery :
	CGO_ENABLED=1 go build -o cmd/jcquery cmd/jcquery.go


cmd/jcquery.exe : cmd/jcquery.go
	CGO_ENABLED=1 go build -o cmd/jcquery.exe cmd/jcquery.go
