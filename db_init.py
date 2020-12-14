#!/usr/bin/env python3
from __future__ import print_function

import os
import mysql.connector
from mysql.connector import errorcode

TABLES = {}
TABLES["users"] = (
    "CREATE TABLE users ("
    "   id BIGINT UNSIGNED NOT NULL,"
    "   username VARCHAR(20) NOT NULL,"
    "   email VARCHAR(255) NOT NULL,"
    "   primary key (id))"
)

if __name__ == "__main__":
    DB_NAME = os.environ.get('DB_NAME')
    cnx = mysql.connector.connect(
        user = os.environ.get('DB_USER'), 
        password = os.environ.get('DB_PASS'), 
        host = os.environ.get('DB_HOST')
    )
    cursor = cnx.cursor()
    try:
        cursor.execute("CREATE DATABASE {} DEFAULT CHARACTER SET 'utf8'".format(DB_NAME))
        cursor.execute("USE {}".format(DB_NAME))
        for table_name in TABLES:
            table_descripton = TABLES[table_name]
            print("creating table {}: ".format(table_name), end='')
            cursor.execute(table_descripton)
    except mysql.connector.Error as err:
        if err.errno == errorcode.ER_TABLE_EXISTS_ERROR:
            print("already exists.")
        else:
            print(err.msg)
    else:
        print("OK.")
    finally:
        cursor.close()
        cnx.close()