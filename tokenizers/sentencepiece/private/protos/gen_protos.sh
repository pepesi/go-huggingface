#!/bin/bash
set -e

# Find Go import path: presumably unique.
import_path="$(go list -f '{{.ImportPath}}')"

# Extract the domain and the rest of the path
domain=$(echo "$import_path" | cut -d '/' -f 1)
rest_of_path=$(echo "$import_path" | cut -d '/' -f 2-)

# Reverse the domain part (split by '.')
reversed_domain=$(echo "$domain" | awk -F '.' '{ for (i=NF; i>1; i--) printf "%s.", $i; print $1 }' | sed 's/\.$//')

# Combine the reversed domain with the rest of the path
tmp_link="$reversed_domain/$rest_of_path"
tmp_link=$(echo "$tmp_link" | tr '/.' '__')
rm -f "${tmp_link}"
ln -s . "${tmp_link}"
protoc --go_out=. --go_opt=paths=source_relative "./${tmp_link}/sentencepiece_model.proto"
rm -f "${tmp_link}"
