#!/bin/bash

# add route to database
cd "$(dirname "$0")"

# add data to database
mysql -u root -p <<EOF
$(cat mysqlapigo_db_data.sql)
EOF
