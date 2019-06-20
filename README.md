# tir-library
Section of Tir's Library containing book entries and reviews from registered scholars.



### Build

To build the application binary just issue:

`go build -o bin/tir-library -v .`

### Run

In order to run the application you need to define your Postgres configuration in a `.env` at the root of the project. 
It is recommended to copy the file `sample.env` and replace the values. Testing environmnet deletes all DB entries
before running the application so be careful not to use it with a production / development DB.

Once you have the `.env` run with:

`./bin/tir-library`



### Tests

#### Postman

Currently tirlib is tested against postman requests collection. You can run them with the provided script:

`./run_tests.sh`

or manually if the application is already running:

`newman run tests/tir-library.postman_collection.json`

Requires newman installed: (`npm install -g newman`):


#### Golang

Golang tests are under construction. They must be improved to use the auth middleware as currently they bypass it. 
Moreover, the access to the DB from concurrent tests are not dealt with. 

### API Documentation

The REST API documentation can be found [here](tirlibrary.docs.apiary.io).