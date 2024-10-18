export JWT_SECRET=eUbP9shywUygMx7u
export MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" 
export MONGO_DATABASE=demo

PWD=$(pwd)
cd ../api;
go test -v *.go
if [ $? == 0 ]; then
  GIN_MODE=release go test -v -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
fi

cd $PWD
