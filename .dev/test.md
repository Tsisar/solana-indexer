```shell
curl -X POST http://localhost:8080/v1/graphql \
  -H "Content-Type: application/json" \
  -H "x-hasura-admin-secret: admin" \
  -d '{"query":"query MyQuery { _meta { deployment error_message has_indexing_errors } }"}'
```