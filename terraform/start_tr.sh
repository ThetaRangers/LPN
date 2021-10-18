#!/bin/bash
if [[ $# != 1 ]]  
then
	echo "Needed var file"
else
	terraform init -input=false
	terraform plan -input=false -var-file="$1" -out=tfplan
	terraform apply -auto-approve tfplan
fi
