BEGIN {
	# while(( getline line<"accounts.csv") > 0 ) {
	#      print line
	#   }
	FS=","
	first=1
	position=1
}

first == 1 {
	count=1
	for (i=1;i<=NF;i++) {
		columns[$i]=count
		count++
	}
	first = 0
	next
}

{
	for (i=1;i<=NF;i++) {
		data[i]=$i
	}

	name=data[columns["Identity"]]
	identity[name]=position
	position++

	for (col in columns) {
		variable[name,col]=data[columns[col]]
	}
}

END {
	if (Attribute == "Peers") {
		init=0
		for (name in identity) {
			nodeName=variable[name,"NodeName"]
			commandName="./getNodeId " nodeName
			commandName | getline id
			entry=id "@" Prefix variable[name,"P2PAddress"]
			if (init == 0) {
				buffer=entry
				init=1
			} else {
				last = buffer
				buffer=last "," entry
			}
		}
		print buffer
	} else {
		print Prefix variable[Identity,Attribute]
	}
}
