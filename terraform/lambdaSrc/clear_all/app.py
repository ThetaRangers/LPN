import sys
import os
import logging
import pymysql
from statement import *

def handler_clear(event, context):
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
        cur.execute(clearCluster)
        cur.execute(clearNodes)

        json_data = { }
    conn.commit()
    conn.close()
    return json_data