#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.

# to start mongodb database
# from /Users/angeldionisio/Documents/mongodb-macos-x86_64-4.2.0
# bin/mongod --dbpath data/db