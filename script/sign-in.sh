curl -v --location --request POST 'http://localhost:8080/signin' \
    --header 'Content-Type: application/json' \
    --data-raw '{ "username": "admin", "password": "password" }'
