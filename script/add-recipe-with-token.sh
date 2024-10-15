if [ "$#" != 1 ]; then
  echo "token needed"
  exit 1
fi

curl -v --location --request POST 'http://localhost:8080/recipes' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: '"$1" \
    --data-raw '{
        "name": "Homemade Pizza",
        "ingredients": ["..."],
        "instructions": ["..."],
        "tags": ["dinner", "fastfood"]
    }'
