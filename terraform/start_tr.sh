#!/bin/bash
if [[ $# != 2 ]]  
then
	echo "Usage ./start_tr.sh varfile lambda_function_name"
else
	zip "$2.zip" "$2"
	terraform init -input=false
	terraform plan -input=false -var-file="$1" -out=tfplan
	terraform apply -auto-approve tfplan
fi
