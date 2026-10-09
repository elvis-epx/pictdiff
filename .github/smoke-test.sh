#!/bin/bash
# Runs a built pictdiff binary against test_data and checks the diff metrics.
# Usage: .github/smoke-test.sh path/to/binary

BIN=$1
TMP=$(mktemp -d)

check() {
	METRIC=$("$BIN" "test_data/$1.png" "test_data/$2.png" "$TMP/out.png") || exit 1
	if [ "$METRIC" != "$3" ]; then
		echo "$BIN $1 $2: expected metric $3, got $METRIC"
		exit 1
	fi
	echo "$BIN $1 $2: ok"
}

check old new 2089964
check olda newa 11180000
check oldb newb 17711226

rm -rf "$TMP"
