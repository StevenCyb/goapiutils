# GoHttpUtils

[![GitHub release badge](https://badgen.net/github/release/StevenCyb/goapiutils/latest?label=Latest&logo=GitHub)](https://github.com/StevenCyb/goapiutils/releases/latest)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/StevenCyb/goapiutils/ci-test?label=Tests&logo=GitHub)
![GitHub](https://img.shields.io/github/license/StevenCyb/goapiutils)

This is a collection of query parsers and extractors that can be used for for API's.

## Query parser
- MongoDB
  - [RSQL parser for MongoDB find queries](parser/mongo/rsql/README.md)
  - [Parser for MongoDB sort options](parser/mongo/sort/README.md)
  - [JSON Patch for MongoDB](parser/mongo/jsonpatch/README.md)
- Object
  - [Parser for subset query](parser/object/subset/README.md)

## Extractor
- HTTP-Request Parameter
  - [Extract query value](extractor/http/request/parameter/README.md#query-parameter)
  - [Extract patch value](extractor/http/request/parameter/README.md#path-parameter)
