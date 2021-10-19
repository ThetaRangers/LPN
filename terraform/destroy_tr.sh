#!/bin/bash
if [[ $# != 1 ]]  
then
	echo "Var file needed"
else
	terraform plan -destroy -var-file="$1" -out=tfdestroy 
	terraform apply tfdestroy 
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -f terraform.tfstate
	rm -f terraform.tfstate.backup
	rm -f tfdestroy
	rm -f tfplan
	find ./ -name \*.zip -delete
fi