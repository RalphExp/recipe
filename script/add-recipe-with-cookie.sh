curl -v --location -b cookies.txt --request POST 'http://localhost:8080/recipes' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "name": "Homemade Pizza",
        "ingredients": ["..."],
        "instructions": ["..."],
        "tags": ["dinner", "fastfood"]
    }'
