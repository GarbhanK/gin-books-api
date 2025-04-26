
ping_api() {
    echo "Pinging the API..."
    curl "localhost:8080/api/v1/ping" | jq
}

insert_book() {
    local author=$1
    local title=$2

    curl -X POST localhost:8080/api/v1/books \
        -H 'Content-Type: application/json' \
        -d "{\"Author\":\"${author}\",\"Title\":\"${title}\"}" | jq
}

get_book() {
    local title=$1

    curl "localhost:8080/api/v1/books/title/?title=${title}" | jq
}

delete_book() {
    local title=$1

    curl -X "DELETE" "localhost:8080/api/v1/books/?title=${title}" | jq
}

get_all_books() {
    curl "localhost:8080/api/v1/books/?table=books" | jq
}

echo "Running end-to-end test..."

ping_api
insert_book "Jorge Luis Borges" "Fictions"
insert_book "Jorge Luis Borges" "The Aleph"
get_book "Fictions"
delete_book "Fictions"
get_all_books
