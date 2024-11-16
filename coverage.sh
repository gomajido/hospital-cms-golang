#!/bin/bash
#
# Code coverage generation

COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
PKG_LIST=$(go list ./... |
grep -v /go-core-nabitu/cmd |
grep -v /go-core-nabitu/config |
grep -v /go-core-nabitu/testemail |
grep -v /go-core-nabitu/assets |
grep -v /go-core-nabitu/internal/router |
grep -v /go-core-nabitu/internal/response |
grep -v /go-core-nabitu/internal/dependency |
grep -v /go-core-nabitu/internal/constant |
grep -v /go-core-nabitu/internal/helper |
grep -v /go-core-nabitu/pkg 
)

# Create the coverage files directory
mkdir -p "$COVERAGE_DIR";

# Create a coverage file for each package
go test -covermode=count -coverprofile "${COVERAGE_DIR}/coverage.cov" ${PKG_LIST} ;

# Merge the coverage profile files
#echo 'mode: count' > "${COVERAGE_DIR}"/coverage.cov ;
#tail -q -n +2 "${COVERAGE_DIR}"/*.cov >> "${COVERAGE_DIR}"/coverage.cov ;

# Display the global code coverage
go tool cover -func="${COVERAGE_DIR}"/coverage.cov ;

# If needed, generate HTML report
if [ "$1" == "html" ]; then
    go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o "${COVERAGE_DIR}"/coverage.html;
    open "${COVERAGE_DIR}"/coverage.html;
fi

# Remove the coverage files directory
sleep 5;
rm -rf "$COVERAGE_DIR";
