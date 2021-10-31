#!/bin/bash
if [[ $# != 1 ]]  
then
	echo "Usage ./start_tr.sh private_var_file"
else
    cd ./lambdaSrc
    for d in ./*/
    do 
        (cd "$d";
        zipName=$(echo "$d" | awk '{ print substr( $d, 3, length($d)-3) }');
        zip "./../../$zipName" -r ./*;)
    done
    cd ..;
	terraform init -input=false
	terraform plan -input=false -var-file="$1" -out=tfplan
	terraform apply -auto-approve tfplan
	python ./setupDB.py
fi
