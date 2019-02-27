#!/bin/bash
#
# Register all of the identities and accounts on OneLedger
#
# Need to test to see if this has already been done...
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts

list="David Alice Bob Carol Emma"

$CMD/startOneLedger

echo "=================== Test Registration ======================="
for name in $list
do

   if [ $name = "Alice" ]
   then
      bkey="eyJrZXkiOiJiMngwWlhOME1ERT0iLCJwYXNzIjoiYjJ4d1lYTnpNREU9IiwiYWRkcmVzcyI6IjJOQWdSd2sybm1tam04RTVjVjVONnR2U3Zxa1p5VkJhZkdIIiwicHJpdmtleSI6ImNUM1ZTTEVhMll4YlBuYjdoREFIcmV3WGYzbnk0TFo3d0ZwZzc4anBlOGY1Qkc3R1dId2oifQ=="
      ekey="eyJrZXkiOiJleUpoWkdSeVpYTnpJam9pWkRjNE5UZ3dNRFU0Tmpkak16UTBPV1kyTmpjellUa3haalpsTkdZM01UbG1NVEJsTVRKbE5TSXNJbU55ZVhCMGJ5STZleUpqYVhCb1pYSWlPaUpoWlhNdE1USTRMV04wY2lJc0ltTnBjR2hsY25SbGVIUWlPaUkzTW1ZeVlqUmpaVFJoWmpGbVpqWTRaak5pT0dNelpEQmpOV0kyWWpOak5UVmxNVFV6TkRReE5qTXpaR1k1WW1FMVlqWXhOVEV6TkdZMk9UZGpNVEE0SWl3aVkybHdhR1Z5Y0dGeVlXMXpJanA3SW1sMklqb2lOemd6T1RVNU5tSTBaamcwWWpkaU1tTTBZMlJtTURGbE9EUXhZVGhqWldNaWZTd2lhMlJtSWpvaWMyTnllWEIwSWl3aWEyUm1jR0Z5WVcxeklqcDdJbVJyYkdWdUlqb3pNaXdpYmlJNk1qWXlNVFEwTENKd0lqb3hMQ0p5SWpvNExDSnpZV3gwSWpvaU1EQTNOVEJrWTJGbE1qZ3daalV5WkdJME9UVXdOVEUzTXpNMFpUQXlNelE0TkRnNFpUWXpPR1l5T1RVNFpHWm1Zall5TldSa1pXRmtPR1EyTWpabU15SjlMQ0p0WVdNaU9pSXpOREZoWmpBMFl6Y3hNRFUxTmpVNE9UQmhZV001T0RWa01UTmxORFppTXpNMk1qTXhNbVpqTkdRME9EZzVORFJrTVRka01UbG1NRE5qTVRBd016ZzJJbjBzSW1sa0lqb2lNREE0T0RneE5qVXROVFkyTXkwMFpUY3hMV0psWldFdE1XSTJNRFUxWmpreFkySmxJaXdpZG1WeWMybHZiaUk2TTMwPSIsInBhc3N3b3JkIjoiMTIzNCJ9"
   elif [ $name = "Bob" ]
   then
      bkey="eyJrZXkiOiJiMngwWlhOME1EST0iLCJwYXNzIjoiYjJ4d1lYTnpNREk9IiwiYWRkcmVzcyI6IjJONWltZWdxREVZNXBocG9yYWFZUXFVRDkzV0xmR0pUajlRIiwicHJpdmtleSI6ImNWanY0N0p0Nkw5NkFyNlZSYUEzRlpGMlhic0NTQnUzZXJNN2phQjdYRWRyQ1BxV0RpQzYifQ=="
      ekey="eyJrZXkiOiJleUpoWkdSeVpYTnpJam9pWVdGbVlUSmtPRGs0TUdFM016QmlNREl4T1RWbU9XTTRaR1psWVdabFlqTmxOamxoTmpsallTSXNJbU55ZVhCMGJ5STZleUpqYVhCb1pYSWlPaUpoWlhNdE1USTRMV04wY2lJc0ltTnBjR2hsY25SbGVIUWlPaUpqWm1Nek5tSTNaR1ZpTlRBek1URTJORGd5TXpjeFlqZGtNalU1Tm1GaE9UTTJOelU0WWpneU5EY3lOemxsWm1ObE5EWXhZMll3TXpRMFlXVTBZak14SWl3aVkybHdhR1Z5Y0dGeVlXMXpJanA3SW1sMklqb2labU15TURCaU9UTTNNVEUyTWpVNE9EVTJaR1F3WlRWaE1EZzFaVEF4TVdRaWZTd2lhMlJtSWpvaWMyTnllWEIwSWl3aWEyUm1jR0Z5WVcxeklqcDdJbVJyYkdWdUlqb3pNaXdpYmlJNk1qWXlNVFEwTENKd0lqb3hMQ0p5SWpvNExDSnpZV3gwSWpvaU56TTFOR00wTlRJelpHWmpOekF6TnpKak9HTXpORFl4Tm1NeE5XUmpNakUwTkRoaFl6UXdOakUzWm1aak0yRTNZak5oT1dGbU4yVmxNekpqTXpkbE5pSjlMQ0p0WVdNaU9pSXhZak15TW1aaE0yTTFOemc1WTJWa1pUZzNNVFEwTnpnelpqSmlaRGhqTkRVNE9HVTFNRGswWlRZNFpqYzRPREEyTkRCallUbGhOVFExT0dJNFlXRmlJbjBzSW1sa0lqb2lPRGN6TmpOaE16a3RNREUzTVMwME5qUXdMV0poTVRJdFlqVmhZV05oWkRkaFpXUXlJaXdpZG1WeWMybHZiaUk2TTMwPSIsInBhc3N3b3JkIjoiMjM0NSJ9"
   elif [ $name = "Carol" ]
   then
      bkey="eyJrZXkiOiJiMngwWlhOME1ETT0iLCJwYXNzIjoiYjJ4d1lYTnpNRE09IiwiYWRkcmVzcyI6IjJORlRZV0pMNVQ2dEs1R2ZXWmJ5Q21jb2FNS0FyenhXRmp1IiwicHJpdmtleSI6ImNWRmlXRHoyOThYeHpNR0VZRzd3alV2YlV3TVdaZ2taQkpNdjlvOTkyb2VmS3hwanpUdnEifQ=="
      ekey="eyJrZXkiOiJleUpoWkdSeVpYTnpJam9pT0dFek1EbG1PVFZrWlRCbE5EZGxaR0kyTVdSbE9HWmhNR05tT0dKa1pEY3lNakkzTVRjNE9TSXNJbU55ZVhCMGJ5STZleUpqYVhCb1pYSWlPaUpoWlhNdE1USTRMV04wY2lJc0ltTnBjR2hsY25SbGVIUWlPaUk0TVdKbFkySTNZMkV6TjJKbE56TTNZV1l4TkRkaFlUQTFOVEppTVRZek9XSTNOekJrTnpaaVlUazRabUU0TWpBMk9UTXlOV1psTVdObE5tVXhZV0V4SWl3aVkybHdhR1Z5Y0dGeVlXMXpJanA3SW1sMklqb2lOV0psTWpCbU1qWXpZVFEyWkRaalkyRTFNMk5pTUdGbE5Ea3dNalExWm1RaWZTd2lhMlJtSWpvaWMyTnllWEIwSWl3aWEyUm1jR0Z5WVcxeklqcDdJbVJyYkdWdUlqb3pNaXdpYmlJNk1qWXlNVFEwTENKd0lqb3hMQ0p5SWpvNExDSnpZV3gwSWpvaU1USTBOVFpqT1dFM05EYzNPR0V3TmpRME9UVTVOalkzTm1Oak9UQm1NbVl3TkRabE16QTJZalZrWWpjME5qZzROakF3WXpBME5UYzNOVEk1WWpsak1pSjlMQ0p0WVdNaU9pSTJOek0zWXpsa1pEa3paakJoWW1NNFpURXdNalU1TURrNE5EYzVNREl4TkdNeVlqbGtabU16Tm1WaE5tVXlZamMyT1dVNE1HTXhPV1ZpTWpKbE5HVTRJbjBzSW1sa0lqb2labUpoWldZeE1tSXRZVFkyTnkwMFl6UmxMV0kwWXpjdE56SXpOR1ZtTXpkalltVTVJaXdpZG1WeWMybHZiaUk2TTMwPSIsInBhc3N3b3JkIjoiMzQ1NiJ9"
   else
      bkey=""
      ekey=""
   fi

	# Add the accounts, keys are generated internally
	olclient update --root $OLDATA/$name-Node --account "$name-OneLedger" --nodeaccount true

	#TODO: Need flag to set node account for the each external chain.
	olclient update --root $OLDATA/$name-Node --account "$name-BitCoin"  --chain "BitCoin"  --chainkey "$bkey"
	olclient update --root $OLDATA/$name-Node --account "$name-Ethereum" --chain "Ethereum" --chainkey "$ekey"

	# Account must have money in order to pay the registration fee
	olclient testmint --root $OLDATA/$name-Node --party "$name-OneLedger" --amount 100000.1 --currency OLT

	# Register the identity across the chain
	olclient register --root $OLDATA/$name-Node --identity "$name" \
		--account "$name-OneLedger" --node "$name-Node" --fee 0.1
done

# Give it some time to get committed
olclient --root $OLDATA/Emma-Node wait --completed identity --identity $list
