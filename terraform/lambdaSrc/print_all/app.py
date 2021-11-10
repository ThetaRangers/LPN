import sys
import os
import logging
import pymysql
from statement import *

def handler_print(event, context):
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

    with conn.cursor() as cur:
        print("NodesIP\tIPstring\n")
        cur.execute(getNodes)

        rows = cur.fetchall()
        for row in rows:
            ip = f'{row[0]}'
            ipStr = f'{row[1]}'
            print("%s\t%s\n" %(ip, ipStr))

        print("ClusterID\tNodeIP\tAttached\n")
        cur.execute(getReplicaOf)

        rows = cur.fetchall()
        for row in rows:
            clusterId = f'{row[0]}'
            nodeIp = f'{row[1]}'
            attached = row[2]
            print("%s\t%s\t%d\n" %(clusterId, nodeIp, attached))


        json_data = { }
    conn.close()
    return json_data