#!/bin/bash -x

OLD1=test_data/old.png
NEW1=test_data/new.png
OLD2=test_data/olda.png
NEW2=test_data/newa.png
REF1=test_data/reference_diff.png
REF2=test_data/reference_diffa.png

# ./pictdiff.py $OLD $NEW $REF
# ./pictdiff.py $OLD2 $NEW2 $REF2

REFMETRIC1=2089964
REFMETRIC2=11180000

npm install
go build pictdiff.go || exit 1
cargo build --release || exit 1

for version in "./pictdiff.py" "target/release/pictdiff" "./pictdiff" "node pictdiff.js"; do
	for sample in 1 2; do
		OLD="OLD$sample"
		NEW="NEW$sample"
		REF="REF$sample"
		REFMETRIC="REFMETRIC$sample"

		$version ${!OLD} ${!NEW} tmp.png || exit 1
		METRIC=$($version ${!OLD} ${!NEW} tmp.png)
		if [ "$METRIC" -ne "${!REFMETRIC}" ]; then
			echo "Command produced an unexpected diff metric"
			exit 1
		fi
	
		./pictdiff.py ${!REF} tmp.png tmp2.png || exit 1
		MAPMETRIC=$(./pictdiff.py ${!REF} tmp.png tmp2.png)
		if [ "$MAPMETRIC" -ne 0 ]; then
			echo "Command produced an unexpected diff map"
			exit 1
		fi
	done
done

rm tmp.png tmp2.png
