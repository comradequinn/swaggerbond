# swaggerbond
A simple swagger file indexer which exposes a UI to search and view its contents

# installation
* clone the repo
* cd into the dir
* run `make demo`
* browse to `http://localhost:8080/`

To run in production mode, run `make start` rather than `make demo` and then add files, by whatever means, into the `swagger-files` directory. They will be promptly parsed and included in search results without restarting the service
