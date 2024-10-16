(cd ../auth0; \
AUTH0_DOMAIN="dev-ed8p2huh8zf16ycg.us.auth0.com" \
AUTH0_API_IDENTIFIER="https://api.recipes.io" \
MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
MONGO_DATABASE=demo \
go run main.go)
