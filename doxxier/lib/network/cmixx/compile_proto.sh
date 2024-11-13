#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR || return

protoc --go_out=. --go_opt=paths=source_relative ./direct_message.proto