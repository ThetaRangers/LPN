import sys
import os
import logging
import pymysql
from statement import *

#rds settings
rds_host, rds_port  = os.environ['rds_endpoint'].split(":")
name = os.environ['db_username']
password = os.environ['db_password']
db_name = os.environ['db_name']


logger = logging.getLogger()
logger.setLevel(logging.INFO)

try:
    conn = pymysql.connect(host=rds_host, user=name, passwd=password, db=db_name, connect_timeout=50)
except:
    logger.error(f"ERROR: Unexpected error: Could not connect to MySql instance., rds_host is {rds_host}, db_name is {db_name}")
    sys.exit()

logger.info("SUCCESS: Connection to RDS mysql instance succeeded")
def handler_registry(event, context):

    #params
    ip = event['Ip']
    n = event['N']

    ipList = []

    with conn.cursor() as cur:
        args = (ip, n, 0, 0)
        cur.callproc('replicaSet',args)

        rows = cur.fetchall()
        for row in rows:
            ipList.append(f'{row[0]}')

        cur.execute(retParams)
        outParams = cur.fetchone()
        valid = outParams[0]
        crashed = outParams[1]
        
        print(crashed, valid)

        json_data = { "crashed": crashed,
                       "valid": valid,
                       "ipList": ipList}
    return json_data