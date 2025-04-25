

author="Jorge Luis Borges"
book_title="Fictions"

echo "author: ${author}"
echo "title: ${book_title}"

echo "Pinging the API..."
curl localhost:8080/api/v1/ping | jq


echo ""
echo "Inserting first book..."
curl -X POST localhost:8080/api/v1/books \
    -H 'Content-Type: application/json' \
    -d "{\"Author\":\"${author}\",\"Title\":\"${book_title}\"}" | jq

echo ""
echo "Inserting second book..."
curl -X POST localhost:8080/api/v1/books \
    -H 'Content-Type: application/json' \
    -d "{\"Author\":\"${author}\",\"Title\":\"The Aleph\"}" | jq

echo ""
echo "Searching for book title..."
curl "localhost:8080/api/v1/books/title/?title=${book_title}" | jq

