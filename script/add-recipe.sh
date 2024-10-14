curl -v --location --request POST 'http://localhost:8080/recipes' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "name": "Homemade Pizza",
        "ingredients": ["..."],
        "instructions": ["..."],
        "tags": ["dinner", "fastfood"]
    }'
