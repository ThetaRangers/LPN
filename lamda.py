from __future__ import print_function
import json
import boto3
import operator

def lambda_handler(event, context):
    '''
        -ip: node ip requested to join
        -network: network to join
        -bootstrap: if the node is bootstrap
        -n: replication coefficient
    '''
    ip = event['ip']
    n = event['n']
    bootstrap = event['bootstrap']

    if 'network' in event:
        dynamo = boto3.resource('dynamodb').Table(event['network'])

    if not ip:
        raise ValueError('Empty ip string')

    x = dynamo.scan()["Items"]
    result = []
    list = []

    min = 0

    for item in x:
        list.append([item["Key"], item["Value"]])


    list = sorted(list, key=operator.itemgetter(1))[:n+1]

    for i in list:
        #Update values
        if i[0] != ip:
            new_value = i[1] + 1
            dynamo.put_item(
                Item = {
                "Key":i[0],
                "Value":new_value
            })

            result.append(i[0])

    result = result[:n]

    dynamo.put_item(
        Item = {
            "Key":ip,
            "Value":0,
            #"isBootstrap":bootstrap
        }
    )

    return json.dumps(result)