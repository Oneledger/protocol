#!/bin/sh
#
# Produce all of the missing fee and gas errors
#

list="Junk Stuff"

for name in $list
do
	olclient update --root $OLDATA/Alice-Node --account "$name-Demo"

	olclient testmint --root $OLDATA/Alice-Node --party "$name-Demo" --amount 3 

	# Should get an error message about missing fees
	olclient register --root $OLDATA/Alice-Node --identity $name --account "$name-Demo" --node "Alice-Node" 

	olclient send --root $OLDATA/Alice-Node --party $name --counterparty Alice --amount 2.5

done

olclient send --root $OLDATA/Alice-Node --party Alice --counterparty Bob --amount 2.5 --fee 0.001

#olclient swap --root $OLDATA/Alice-Node --party Bob --counterparty Alice --amount 2.5
#olclient swap --root $OLDATA/Bob-Node --party Alice --counterparty Bob --amount 2.5
