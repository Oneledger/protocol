#!/bin/bash
  
rm -r New-Node/*.log
rm -r New-Node/olfullnode*
rm -r New-Node/consensus/data
rm -r New-Node/consensus/config/addrbook.json

sed -i '/\"last\_/d' New-Node/consensus/config/priv_validator.json
