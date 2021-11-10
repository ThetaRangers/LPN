import sys
import os
import logging
import pymysql
from statement import *

def handler_registry(event, context):
    #rds settings
    rds_host, rds_port  = os.environ['rds_endpoint'].split(":")
    name = os.environ['db_username']
    password = os.environ['db_password']
    db_name = os.environ['db_name']


    logger = logging.getLogger()
    logger.setLevel(logging.INFO)

    try:
        conn = pymysql.connect(host=rds_host, user=name, passwd=password, db=db_name, connect_timeout=5)
    except:
        logger.error(f"ERROR: Unexpected error: Could not connect to MySql instance., rds_host is {rds_host}, db_name is {db_name}")
        sys.exit()

    logger.info("SUCCESS: Connection to RDS mysql instance succeeded")

    #params
    ip = event['ip']
    n = event['n']
    ipStr = event['ipStr']

    ipList = []

    with conn.cursor() as cur:
        args = (ip, ipStr, n, 0, 0)
        cur.callproc('clusterSet',args)

        rows = cur.fetchall()
        for row in rows:
            ip = f'{row[0]}'
            ipStr = f'{row[1]}'
            jsonField = { "ip": ip,
                          "strIp": ipStr }
            ipList.append(jsonField)

        cur.execute(retParams)
        outParams = cur.fetchone()
        valid = outParams[0]
        crashed = outParams[1]

        json_data = { "crashed": crashed,
                       "valid": valid,
                       "ipList": ipList}
    conn.close()
    return json_data