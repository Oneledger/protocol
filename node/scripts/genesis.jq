.chain_id = "OneLedger-Root" | 
.validators[0].name = "David-Node" |
.validators[1].name = "Alice-Node" |
.validators[2].name = "Bob-Node" |
.validators[3].name = "Carol-Node" |
. + {app_state: { account : "Zero", states : ""} } |
.app_state.states = [{amount: "1000000000000", coin: "OLT" }, {amount: "10000", coin: "VT"}]
