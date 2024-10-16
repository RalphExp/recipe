curl -v --location --request POST 'http://localhost:8080/recipes'\
    --header 'Content-Type: application/json'\
    --header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InA3TVd0eERvejc3QlkzZE9Ca0VpdiJ9.eyJpc3MiOiJodHRwczovL2Rldi1lZDhwMmh1aDh6ZjE2eWNnLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJKbGpOR2x4MlE5dmVyMHluYTZyMG01dHh3MVJjYjVxekBjbGllbnRzIiwiYXVkIjoiaHR0cHM6Ly9hcGkucmVjaXBlcy5pbyIsImlhdCI6MTcyOTA0NzkxOSwiZXhwIjoxNzI5MTM0MzE5LCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMiLCJhenAiOiJKbGpOR2x4MlE5dmVyMHluYTZyMG01dHh3MVJjYjVxeiJ9.nuUEREcw8P5bpcXKlUv6otcU-n-20g9V_dGOctoPLupbBxoGYC3Ys5tY335-Ovup11pxC3IO2HhHBryY1uZBQ5oy67Nacb5lXnny29NJqOUGegN7k6UgmZgUCHgl5UDKlnZfY2wUIJYBfmZ1RelTvsS6SEqR6QXvjdLyQOyGT3zeWIfE1V5RBcafK2QLi92qmJJojJfx5NKncKpV5NB825NpnAls_3GK9MntBBjOWBF676tdoB5D2YPSPTwL1Cnfw5E4y0L5hvf204KwehrSTi1N_X4bw8c9GU3qqC50y2jhb2w6ZbkJXWAZWraykXwfeEkIIdyHdt32-o-4D_sj5Q'\
    --data '{
        "name": "Homemade Pizza",
        "ingredients": ["..."],
        "instructions": ["..."],
        "tags": ["dinner", "fastfood"]
    }'
