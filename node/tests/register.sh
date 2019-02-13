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
      bkey="b2x0ZXN0MDE="
      bpass="b2xwYXNzMDE="
      ekey="{'address':'d7858005867c3449f6673a91f6e4f719f10e12e5','crypto':{'cipher':'aes-128-ctr','ciphertext':'72f2b4ce4af1ff68f3b8c3d0c5b6b3c55e153441633df9ba5b615134f697c108','cipherparams':{'iv':'7839596b4f84b7b2c4cdf01e841a8cec'},'kdf':'scrypt','kdfparams':{'dklen':32,'n':262144,'p':1,'r':8,'salt':'00750dcae280f52db4950517334e02348488e638f2958dffb625ddead8d626f3'},'mac':'341af04c7105565890aac985d13e46b3362312fc4d488944d17d19f03c100386'},'id':'00888165-5663-4e71-beea-1b6055f91cbe','version':3}"
      epass="1234"
   elif [ $name = "Bob" ]
   then
      bkey="b2x0ZXN0MDI="
      bpass="b2xwYXNzMDI="
      ekey="{'address':'aafa2d8980a730b02195f9c8dfeafeb3e69a69ca','crypto':{'cipher':'aes-128-ctr','ciphertext':'cfc36b7deb503116482371b7d2596aa936758b8247279efce461cf0344ae4b31','cipherparams':{'iv':'fc200b937116258856dd0e5a085e011d'},'kdf':'scrypt','kdfparams':{'dklen':32,'n':262144,'p':1,'r':8,'salt':'7354c4523dfc70372c8c34616c15dc21448ac40617ffc3a7b3a9af7ee32c37e6'},'mac':'1b322fa3c5789cede87144783f2bd8c4588e5094e68f7880640ca9a5458b8aab'},'id':'87363a39-0171-4640-ba12-b5aacad7aed2','version':3}"
      epass="2345"
   elif [ $name = "Carol" ]
   then
      bkey="b2x0ZXN0MDM="
      bpass="b2xwYXNzMDM="
      ekey="{'address':'8a309f95de0e47edb61de8fa0cf8bdd722271789','crypto':{'cipher':'aes-128-ctr','ciphertext':'81becb7ca37be737af147aa0552b1639b770d76ba98fa82069325fe1ce6e1aa1','cipherparams':{'iv':'5be20f263a46d6cca53cb0ae490245fd'},'kdf':'scrypt','kdfparams':{'dklen':32,'n':262144,'p':1,'r':8,'salt':'12456c9a74778a06449596676cc90f2f046e306b5db74688600c04577529b9c2'},'mac':'6737c9dd93f0abc8e102590984790214c2b9dfc36ea6e2b769e80c19eb22e4e8'},'id':'fbaef12b-a667-4c4e-b4c7-7234ef37cbe9','version':3}"
      epass="3456"
   else
      bkey=""
      bpass=""
      ekey=""
      epass=""
   fi

	# Add the accounts, keys are generated internally
	olclient update --root $OLDATA/$name-Node --account "$name-OneLedger" --nodeaccount true

	#TODO: Need flag to set node account for the each external chain.
	olclient update --root $OLDATA/$name-Node --account "$name-BitCoin"  --chain "BitCoin"  --chainkey "$bkey/$bpass"
	olclient update --root $OLDATA/$name-Node --account "$name-Ethereum" --chain "Ethereum" --chainkey "$ekey/$epass"

	# Account must have money in order to pay the registration fee
	olclient testmint --root $OLDATA/$name-Node --party "$name-OneLedger" --amount 100000.1 --currency OLT

	# Register the identity across the chain
	olclient register --root $OLDATA/$name-Node --identity "$name" \
		--account "$name-OneLedger" --node "$name-Node" --fee 0.1
done

# Give it some time to get committed
olclient --root $OLDATA/Emma-Node wait --completed identity --identity $list
