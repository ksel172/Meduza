#!/bin/bash

QUERY_PATH="/home/queries"
TMP_QUERY_PATH="/tmp/queries"

mkdir $TMP_QUERY_PATH
cp $QUERY_PATH/* $TMP_QUERY_PATH

# Connect to storage
echo "Connecting to database '$POSTGRES_DB'..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -c "\\connect $POSTGRES_DB";
echo "Finished connecting to database '$POSTGRES_DB' with status code $?"

# Set up Schema
echo "Replacing environmental variables in SQL files..."
FILE="$TMP_QUERY_PATH/schema.sql"

echo "Replacing schema in '$FILE'..."
cat "$FILE" | sed -i -e "s/{POSTGRES_SCHEMA}/$POSTGRES_SCHEMA/g" "$FILE"
echo "Finished replacing schema in '$FILE' with status code $?"
cat "$FILE"

echo "Creating schema using PSQL..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$FILE"
echo "Finished creating schema using PSQL with status code $?"

# Set up Tables
echo "Creating tables..."
tables=("users" "listeners" "agents" "agent_config" "agent_info" "agent_task" "agent_command" "payloads" "modules" "teams" "team_members" "certificates" )

for t in "${tables[@]}";do
    current_file="$TMP_QUERY_PATH/$t.sql"
    echo "Replacing environmental variables in '$current_file'..."
    cat "$current_file" | sed -i -e "s/{POSTGRES_SCHEMA}/$POSTGRES_SCHEMA/g" -e "s/{TABLE_NAME}/$t/g" "$current_file"
    echo "Finished replacing schema in '$current_file' with status code $?"
    cat $current_file
    echo "Creating table '$t' using PSQL..."
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$current_file"
    echo "Creating table '$t' in PSQL with status code $?"
done

echo "Finished creating tables"

# Enable pgcrypto extension
echo "Enabling pgcrypto extension..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS pgcrypto;
EOSQL
echo "Finished enabling pgcrypto extension"

# Insert default admin
echo "Inserting default admin..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    INSERT INTO ${POSTGRES_SCHEMA}.users (username, pw_hash, role) VALUES ('Meduza', crypt('Meduza', gen_salt('bf')), 'admin');
EOSQL
echo "Finished inserting default admin"