(cd ../api; \
JWT_SECRET=eUbP9shywUygMx7u \
MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
MONGO_DATABASE=demo \
go test -v *.go)
