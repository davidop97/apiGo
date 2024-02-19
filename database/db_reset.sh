#!/bin/bash

# add route to database
cd "$(dirname "$0")"

# - database
mysql -u root -p < mysqlapigo_db.sql

# - user
# Check if the user exists before attempting to drop it
user_exist=$(mysql -u root -p -e "SELECT 1 FROM mysql.user WHERE user='mysql_apigo_user'" | grep -q 1 && echo "true" || echo "false")

if [ "$user_exist" == "true" ]; then
    echo "User 'mysql_apigo_user' exists. Dropping..."
    mysql -u root -p -e "DROP USER 'mysql_apigo_user'@'localhost';"
else
    echo "User 'mysql_apigo_user' does not exist. Skipping drop."
fi

# Recreate user
mysql -u root -p < mysqlapigo_db_user.sql
