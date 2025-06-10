#!/bin/bash

# MySQL to SQLite Migration Script for Skyclerk
# Usage: ./migrate_db.sh <mysql_dump_file> [sqlite_db_path]

set -e  # Exit on any error

# Check arguments
if [ $# -lt 1 ]; then
    echo "Usage: $0 <mysql_dump_file> [sqlite_db_path]"
    echo "Example: $0 production_dump.sql backend/cache/sqlite/skyclerk.db"
    exit 1
fi

MYSQL_DUMP="$1"
SQLITE_DB="${2:-backend/cache/sqlite/sk_test_mj7zC9Al5M.db}"

# Check if input file exists
if [ ! -f "$MYSQL_DUMP" ]; then
    echo "Error: MySQL dump file '$MYSQL_DUMP' not found"
    exit 1
fi

# Check if SQLite database exists
if [ ! -f "$SQLITE_DB" ]; then
    echo "Error: SQLite database '$SQLITE_DB' not found"
    echo "Make sure to run the Go application first to create the database schema"
    exit 1
fi

echo "Starting MySQL to SQLite migration..."
echo "Source: $MYSQL_DUMP"
echo "Target: $SQLITE_DB"
echo

# Create backup
BACKUP_FILE="${SQLITE_DB}.backup.$(date +%Y%m%d_%H%M%S)"
echo "Creating backup: $BACKUP_FILE"
cp "$SQLITE_DB" "$BACKUP_FILE"

# Temporary file names
CONVERTED_DUMP="temp_converted_dump.sql"
INSERTS_ONLY="temp_inserts_only.sql" 
SINGLE_INSERTS="temp_single_inserts.sql"
FINAL_IMPORT="temp_final_import.sql"

# Cleanup function
cleanup() {
    echo "Cleaning up temporary files..."
    rm -f "$CONVERTED_DUMP" "$INSERTS_ONLY" "$SINGLE_INSERTS" "$FINAL_IMPORT"
}

# Set trap to cleanup on exit
trap cleanup EXIT

echo "Step 1/6: Converting MySQL dump to SQLite format..."
python3 - "$MYSQL_DUMP" "$CONVERTED_DUMP" << 'EOF'
import re
import sys
import os

def convert_mysql_to_sqlite(input_file, output_file):
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Remove MySQL-specific SET statements and comments
    content = re.sub(r'/\*!\d+.*?\*/;', '', content, flags=re.MULTILINE)
    content = re.sub(r'/\*!\d+.*?\*/', '', content, flags=re.MULTILINE)
    
    # Convert data types
    content = re.sub(r'\bint\(\d+\)\s+unsigned\s+NOT\s+NULL\s+AUTO_INCREMENT\b', 'INTEGER PRIMARY KEY AUTOINCREMENT', content)
    content = re.sub(r'\bint\(\d+\)\s+unsigned\b', 'INTEGER', content)
    content = re.sub(r'\bint\(\d+\)\b', 'INTEGER', content)
    content = re.sub(r'\bvarchar\((\d+)\)', r'TEXT', content)
    content = re.sub(r'\btext\b', 'TEXT', content)
    content = re.sub(r'\benum\([^)]+\)', 'TEXT', content)
    content = re.sub(r'\btimestamp\s+NOT\s+NULL\s+DEFAULT\s+current_timestamp\(\)\s+ON\s+UPDATE\s+current_timestamp\(\)', 'DATETIME DEFAULT CURRENT_TIMESTAMP', content)
    content = re.sub(r'\btimestamp\s+NOT\s+NULL\s+DEFAULT\s+\'0000-00-00 00:00:00\'', 'DATETIME DEFAULT CURRENT_TIMESTAMP', content)
    content = re.sub(r'\btimestamp\b', 'DATETIME', content)
    
    # Fix apostrophes in data
    lines = content.split('\n')
    for i, line in enumerate(lines):
        if line.strip().startswith('INSERT INTO'):
            lines[i] = line.replace("\\'", "''")
    content = '\n'.join(lines)
    
    # Remove MySQL-specific elements
    content = re.sub(r'\)\s+ENGINE=\w+[^;]*;', ');', content)
    content = re.sub(r'\s+AUTO_INCREMENT=\d+', '', content)
    content = re.sub(r'\s+CHARACTER\s+SET\s+\w+', '', content)
    content = re.sub(r'\s+COLLATE\s+\w+', '', content)
    content = re.sub(r'\s+DEFAULT\s+CHARSET=\w+', '', content)
    
    # Remove KEY definitions
    lines = content.split('\n')
    filtered_lines = []
    in_table_def = False
    
    for line in lines:
        if line.strip().startswith('CREATE TABLE'):
            in_table_def = True
            filtered_lines.append(line)
        elif in_table_def and line.strip() == ');':
            in_table_def = False
            filtered_lines.append(line)
        elif in_table_def:
            if (line.strip().startswith('KEY ') or 
                line.strip().startswith('UNIQUE KEY ') or
                line.strip().startswith('INDEX ')):
                continue
            else:
                if (len(filtered_lines) > 0 and 
                    filtered_lines[-1].strip().endswith(',') and
                    not line.strip().startswith('PRIMARY KEY')):
                    remaining_lines = [l for l in lines[lines.index(line):] if l.strip()]
                    if remaining_lines and remaining_lines[0].strip() == ');':
                        filtered_lines[-1] = filtered_lines[-1].rstrip().rstrip(',')
                filtered_lines.append(line)
        else:
            filtered_lines.append(line)
    
    content = '\n'.join(filtered_lines)
    content = re.sub(r',\s*PRIMARY KEY \([^)]+\)', '', content)
    content = re.sub(r',(\s*\n\s*\);)', r'\1', content)
    
    sqlite_header = """-- SQLite migration from MySQL dump
PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;

"""
    sqlite_footer = """
COMMIT;
PRAGMA foreign_keys = ON;
"""
    
    content = sqlite_header + content + sqlite_footer
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"Converted MySQL dump to SQLite format")

convert_mysql_to_sqlite(sys.argv[1], sys.argv[2])
EOF

echo "Step 2/6: Extracting INSERT statements..."
python3 - "$CONVERTED_DUMP" "$INSERTS_ONLY" << 'EOF'
import re
import sys

def extract_inserts(input_file, output_file):
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    table_mappings = {
        'Categories': 'Categories', 'Contacts': 'Contacts', 'Files': 'Files',
        'FilesToLedger': 'FilesToLedger', 'Labels': 'Labels', 'LabelsToLedger': 'LabelsToLedger',
        'Ledger': 'Ledger', 'SnapClerk': 'SnapClerk', 'accounts': 'accounts',
        'acct_to_users': 'acct_to_users', 'activities': 'activities', 'applications': 'applications',
        'billings': 'billings', 'connected_accounts': 'connected_accounts', 
        'forgot_passwords': 'forgot_passwords', 'invites': 'invites', 'sessions': 'sessions', 'users': 'users'
    }
    
    lines = content.split('\n')
    insert_lines = ["-- SQLite data import", "PRAGMA foreign_keys = OFF;", "BEGIN TRANSACTION;", ""]
    
    i = 0
    while i < len(lines):
        line = lines[i]
        if line.strip().startswith('INSERT INTO'):
            match = re.match(r'INSERT INTO `([^`]+)`', line)
            if match:
                old_table = match.group(1)
                if old_table in table_mappings:
                    new_table = table_mappings[old_table]
                    line = line.replace(f'`{old_table}`', f'`{new_table}`')
                    insert_lines.append(line)
                    
                    i += 1
                    while i < len(lines) and not lines[i].strip().startswith('INSERT INTO') and lines[i].strip() != '':
                        if lines[i].strip().startswith('(') or lines[i].strip() == ';':
                            insert_lines.append(lines[i])
                        i += 1
                    i -= 1
        i += 1
    
    insert_lines.extend(["", "COMMIT;", "PRAGMA foreign_keys = ON;"])
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write('\n'.join(insert_lines))
    
    print(f"Extracted INSERT statements")

extract_inserts(sys.argv[1], sys.argv[2])
EOF

echo "Step 3/6: Converting to single-row INSERT statements..."
python3 - "$INSERTS_ONLY" "$SINGLE_INSERTS" << 'EOF'
import sys

def convert_to_single_inserts(input_file, output_file):
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    lines = content.split('\n')
    output_lines = []
    current_insert_header = None
    
    for line in lines:
        if line.strip().startswith('INSERT INTO'):
            current_insert_header = line.strip()
        elif line.strip().startswith('(') and current_insert_header:
            values_line = line.strip().rstrip(',').rstrip(';')
            complete_insert = f"{current_insert_header} {values_line};"
            output_lines.append(complete_insert)
        elif not line.strip().startswith('('):
            output_lines.append(line)
            if line.strip() and not line.strip().startswith('INSERT INTO'):
                current_insert_header = None
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write('\n'.join(output_lines))
    
    print(f"Converted to single-row INSERT statements")

convert_to_single_inserts(sys.argv[1], sys.argv[2])
EOF

echo "Step 4/6: Fixing quote escaping..."
sed "s/\\\\'/\\'\\'/g" "$SINGLE_INSERTS" > "$FINAL_IMPORT"

echo "Step 5/6: Importing data into SQLite database..."
sqlite3 "$SQLITE_DB" ".read $FINAL_IMPORT"

echo "Step 6/6: Verifying import..."
echo "Data verification:"
echo -n "Categories: "
sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM Categories;"
echo -n "Contacts: "
sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM Contacts;"
echo -n "Ledger: "
sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM Ledger;"
echo -n "Accounts: "
sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM accounts;"
echo -n "Users: "
sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM users;"

echo
echo "Migration completed successfully!"
echo "Backup created at: $BACKUP_FILE"
echo "To rollback: cp $BACKUP_FILE $SQLITE_DB"