#!/bin/bash
if [[ $# != 1 ]]  
then
	echo "Needed var file"
else
	terraform plan -destroy -var-file="$1" -out=tfdestroy 
	terraform apply tfdestroy 
fi