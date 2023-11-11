curl -H 'Content-Type: application/json' \
     -d '{"title":"Doom Guy","author":"John Romero"}' \
     -X PATCH  \
     localhost:8080/books/2
