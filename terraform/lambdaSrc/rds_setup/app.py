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
    conn = pymysql.connect(host=rds_host, user=name, passwd=password, db=db_name, connect_timeout=5)
except:
    logger.error(f"ERROR: Unexpected error: Could not connect to MySql instance., rds_host is {rds_host}, db_name is {db_name}")
    sys.exit()

logger.info("SUCCESS: Connection to RDS mysql instance succeeded")
def handler_setup(event, context):

    with conn.cursor() as cur:
        
        for x in allOp:
            cur.execute(x)    
        
    return {"statusCode": 200,
        "headers": {"Content-Type": "application/json"},
        "body": "\"DB setted correctly\""         }