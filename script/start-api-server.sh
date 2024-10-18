(cd ../api; \
JWT_SECRET=eUbP9shywUygMx7u \
MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
MONGO_DATABASE=demo\
REDIS_URI=localhost:6379\
go run main.go)
