#!/bin/bash -ex
cd "$(dirname "$0")"
go install tool
mage build coverage crosscompile
cat report.out
