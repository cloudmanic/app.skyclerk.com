# MySQL to SQLite Migration Guide

This document outlines the complete process for migrating a MySQL database dump to SQLite for the Skyclerk application.

## Overview

The migration process converts a MySQL dump file to SQLite-compatible format and imports the data into the existing SQLite database structure created by the Go application's GORM migrations.

## Prerequisites

- Python 3.x installed
- SQLite3 command-line tool
- Existing SQLite database with proper schema (created by running the Go application)
- MySQL dump file to migrate

## Step-by-Step Process

### 1. Create the Migration Script

Create `migrate_mysql_to_sqlite.py`:

```python
#!/usr/bin/env python3
"""
MySQL to SQLite Migration Script for Skyclerk
Converts MySQL dump to SQLite-compatible format
"""

import re
import sys
import os

def convert_mysql_to_sqlite(input_file, output_file):
    """Convert MySQL dump to SQLite format"""
    
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Remove MySQL-specific SET statements and comments
    content = re.sub(r'/\*!\d+.*?\*/;', '', content, flags=re.MULTILINE)
    content = re.sub(r'/\*!\d+.*?\*/', '', content, flags=re.MULTILINE)
    
    # Convert data types
    # MySQL int(N) to SQLite INTEGER
    content = re.sub(r'\bint\(\d+\)\s+unsigned\s+NOT\s+NULL\s+AUTO_INCREMENT\b', 'INTEGER PRIMARY KEY AUTOINCREMENT', content)
    content = re.sub(r'\bint\(\d+\)\s+unsigned\b', 'INTEGER', content)
    content = re.sub(r'\bint\(\d+\)\b', 'INTEGER', content)
    
    # MySQL varchar to TEXT (SQLite handles varchar as TEXT anyway)
    content = re.sub(r'\bvarchar\((\d+)\)', r'TEXT', content)
    
    # MySQL text to TEXT
    content = re.sub(r'\btext\b', 'TEXT', content)
    
    # MySQL ENUM to TEXT (SQLite doesn't support ENUM)
    content = re.sub(r'\benum\([^)]+\)', 'TEXT', content)
    
    # MySQL timestamp with current_timestamp() to DATETIME with CURRENT_TIMESTAMP
    content = re.sub(r'\btimestamp\s+NOT\s+NULL\s+DEFAULT\s+current_timestamp\(\)\s+ON\s+UPDATE\s+current_timestamp\(\)', 'DATETIME DEFAULT CURRENT_TIMESTAMP', content)
    content = re.sub(r'\btimestamp\s+NOT\s+NULL\s+DEFAULT\s+\'0000-00-00 00:00:00\'', 'DATETIME DEFAULT CURRENT_TIMESTAMP', content)
    content = re.sub(r'\btimestamp\b', 'DATETIME', content)
    
    # Fix apostrophes in data - escape single quotes for SQL
    # This is a bit tricky as we need to escape quotes in INSERT data but not in schema
    lines = content.split('\n')
    for i, line in enumerate(lines):
        if line.strip().startswith('INSERT INTO'):
            # Escape single quotes in INSERT statements
            lines[i] = line.replace("\\'", "''")  # MySQL escaping to SQLite escaping
            
    content = '\n'.join(lines)
    
    # Remove MySQL ENGINE and other table options
    content = re.sub(r'\)\s+ENGINE=\w+[^;]*;', ');', content)
    
    # Remove MySQL AUTO_INCREMENT from table definition
    content = re.sub(r'\s+AUTO_INCREMENT=\d+', '', content)
    
    # Remove MySQL charset and collation
    content = re.sub(r'\s+CHARACTER\s+SET\s+\w+', '', content)
    content = re.sub(r'\s+COLLATE\s+\w+', '', content)
    content = re.sub(r'\s+DEFAULT\s+CHARSET=\w+', '', content)
    
    # Remove MySQL specific KEY definitions that aren't PRIMARY KEY
    # Keep PRIMARY KEY but remove other KEY definitions
    lines = content.split('\n')
    filtered_lines = []
    in_table_def = False
    
    for line in lines:
        # Track if we're inside a table definition
        if line.strip().startswith('CREATE TABLE'):
            in_table_def = True
            filtered_lines.append(line)
        elif in_table_def and line.strip() == ');':
            in_table_def = False
            filtered_lines.append(line)
        elif in_table_def:
            # Skip KEY definitions that aren't PRIMARY KEY
            if (line.strip().startswith('KEY ') or 
                line.strip().startswith('UNIQUE KEY ') or
                line.strip().startswith('INDEX ')):
                continue
            else:
                # Remove trailing comma if this was the last field before a KEY
                if (len(filtered_lines) > 0 and 
                    filtered_lines[-1].strip().endswith(',') and
                    not line.strip().startswith('PRIMARY KEY')):
                    # Check if the next non-empty line after this is the closing paren
                    remaining_lines = [l for l in lines[lines.index(line):] if l.strip()]
                    if remaining_lines and remaining_lines[0].strip() == ');':
                        # Remove the trailing comma from the previous line
                        filtered_lines[-1] = filtered_lines[-1].rstrip().rstrip(',')
                
                filtered_lines.append(line)
        else:
            filtered_lines.append(line)
    
    content = '\n'.join(filtered_lines)
    
    # Remove duplicate PRIMARY KEY definitions (in case AUTO_INCREMENT conversion created them)
    content = re.sub(r',\s*PRIMARY KEY \([^)]+\)', '', content)
    
    # Clean up trailing commas before closing parentheses in table definitions
    content = re.sub(r',(\s*\n\s*\);)', r'\1', content)
    
    # Add SQLite-specific pragmas at the beginning
    sqlite_header = """-- SQLite migration from MySQL dump
PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;

"""
    
    # Add commit at the end
    sqlite_footer = """
COMMIT;
PRAGMA foreign_keys = ON;
"""
    
    content = sqlite_header + content + sqlite_footer
    
    # Write the converted content
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"Successfully converted MySQL dump to SQLite format")
    print(f"Input: {input_file}")
    print(f"Output: {output_file}")

def main():
    if len(sys.argv) != 3:
        print("Usage: python3 migrate_mysql_to_sqlite.py <input_mysql_dump> <output_sqlite_file>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    if not os.path.exists(input_file):
        print(f"Error: Input file '{input_file}' not found")
        sys.exit(1)
    
    try:
        convert_mysql_to_sqlite(input_file, output_file)
    except Exception as e:
        print(f"Error during conversion: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
```

### 2. Create INSERT Extraction Script

Create `extract_inserts.py`:

```python
#!/usr/bin/env python3
"""
Extract only INSERT statements from MySQL dump and fix table names for SQLite
"""

import re
import sys

def extract_inserts(input_file, output_file):
    """Extract INSERT statements and fix table names"""
    
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Table name mappings from MySQL dump to SQLite
    table_mappings = {
        'Categories': 'Categories',
        'Contacts': 'Contacts', 
        'Files': 'Files',
        'FilesToLedger': 'FilesToLedger',
        'Labels': 'Labels',
        'LabelsToLedger': 'LabelsToLedger',
        'Ledger': 'Ledger',
        'SnapClerk': 'SnapClerk',
        'accounts': 'accounts',
        'acct_to_users': 'acct_to_users',
        'activities': 'activities',
        'applications': 'applications',
        'billings': 'billings',
        'connected_accounts': 'connected_accounts',
        'forgot_passwords': 'forgot_passwords',
        'invites': 'invites',
        'sessions': 'sessions',
        'users': 'users'
    }
    
    lines = content.split('\n')
    insert_lines = []
    
    # Add SQLite header
    insert_lines.append("-- SQLite data import")
    insert_lines.append("PRAGMA foreign_keys = OFF;")
    insert_lines.append("BEGIN TRANSACTION;")
    insert_lines.append("")
    
    i = 0
    while i < len(lines):
        line = lines[i]
        # Only keep INSERT statements and their data
        if line.strip().startswith('INSERT INTO'):
            # Extract table name and fix it
            match = re.match(r'INSERT INTO `([^`]+)`', line)
            if match:
                old_table = match.group(1)
                if old_table in table_mappings:
                    new_table = table_mappings[old_table]
                    # Replace the table name in the line
                    line = line.replace(f'`{old_table}`', f'`{new_table}`')
                    insert_lines.append(line)
                    
                    # Add the values line(s) that follow
                    i += 1
                    while i < len(lines) and not lines[i].strip().startswith('INSERT INTO') and lines[i].strip() != '':
                        if lines[i].strip().startswith('(') or lines[i].strip() == ';':
                            insert_lines.append(lines[i])
                        i += 1
                    i -= 1  # Back up one since the loop will increment
                else:
                    print(f"Warning: Unknown table {old_table}")
        i += 1
    
    # Add SQLite footer
    insert_lines.append("")
    insert_lines.append("COMMIT;")
    insert_lines.append("PRAGMA foreign_keys = ON;")
    
    content = '\n'.join(insert_lines)
    
    # Write the extracted content
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"Extracted INSERT statements")
    print(f"Input: {input_file}")
    print(f"Output: {output_file}")
    print(f"Found {len([l for l in insert_lines if l.startswith('INSERT')])} INSERT statements")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 extract_inserts.py <input_sql> <output_sql>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    extract_inserts(input_file, output_file)
```

### 3. Create Single INSERT Converter

Create `convert_to_single_inserts.py`:

```python
#!/usr/bin/env python3
"""
Convert multi-row INSERT statements to single-row INSERT statements for SQLite
"""

import re
import sys

def convert_to_single_inserts(input_file, output_file):
    """Convert multi-row INSERTs to single-row INSERTs"""
    
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    lines = content.split('\n')
    output_lines = []
    
    current_insert_header = None
    
    for line in lines:
        if line.strip().startswith('INSERT INTO'):
            # Store the INSERT header
            current_insert_header = line.strip()
        elif line.strip().startswith('(') and current_insert_header:
            # This is a values line, convert to complete INSERT
            # Remove trailing comma and/or semicolon
            values_line = line.strip().rstrip(',').rstrip(';')
            complete_insert = f"{current_insert_header} {values_line};"
            output_lines.append(complete_insert)
        elif not line.strip().startswith('('):
            # Not a values line, add as-is and reset header if needed
            output_lines.append(line)
            if line.strip() and not line.strip().startswith('INSERT INTO'):
                current_insert_header = None
    
    content = '\n'.join(output_lines)
    
    # Write the converted content
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"Converted to single-row INSERT statements")
    print(f"Input: {input_file}")
    print(f"Output: {output_file}")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 convert_to_single_inserts.py <input_sql> <output_sql>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    convert_to_single_inserts(input_file, output_file)
```

## Migration Commands

Here are the exact commands to run for a complete migration:

```bash
# 1. Backup existing SQLite database
cp backend/cache/sqlite/sk_test_mj7zC9Al5M.db backend/cache/sqlite/sk_test_mj7zC9Al5M.db.backup

# 2. Convert MySQL dump to SQLite format
python3 migrate_mysql_to_sqlite.py your_mysql_dump.sql converted_dump.sql

# 3. Extract only INSERT statements with proper table mapping
python3 extract_inserts.py converted_dump.sql inserts_only.sql

# 4. Convert to single-row INSERT statements
python3 convert_to_single_inserts.py inserts_only.sql single_inserts.sql

# 5. Fix quote escaping for SQLite compatibility
sed "s/\\\\'/\\'\\'/g" single_inserts.sql > final_import.sql

# 6. Import data into SQLite database
sqlite3 backend/cache/sqlite/sk_test_mj7zC9Al5M.db ".read final_import.sql"

# 7. Verify import success
sqlite3 backend/cache/sqlite/sk_test_mj7zC9Al5M.db "SELECT COUNT(*) FROM Categories;"
sqlite3 backend/cache/sqlite/sk_test_mj7zC9Al5M.db "SELECT COUNT(*) FROM Contacts;"
sqlite3 backend/cache/sqlite/sk_test_mj7zC9Al5M.db "SELECT COUNT(*) FROM Ledger;"
sqlite3 backend/cache/sqlite/sk_test_mj7zC9Al5M.db "SELECT COUNT(*) FROM accounts;"
sqlite3 backend/cache/sqlite/sk_test_mj7zC9Al5M.db "SELECT COUNT(*) FROM users;"

# 8. Clean up temporary files (optional)
rm converted_dump.sql inserts_only.sql single_inserts.sql
```

## Table Mappings

The migration handles these table mappings between MySQL and SQLite:

| MySQL Table Name | SQLite Table Name | Notes |
|------------------|-------------------|-------|
| Categories | Categories | Mixed case preserved |
| Contacts | Contacts | Mixed case preserved |
| Files | Files | Mixed case preserved |
| FilesToLedger | FilesToLedger | Mixed case preserved |
| Labels | Labels | Mixed case preserved |
| LabelsToLedger | LabelsToLedger | Mixed case preserved |
| Ledger | Ledger | Mixed case preserved |
| SnapClerk | SnapClerk | Mixed case preserved |
| accounts | accounts | Lowercase |
| acct_to_users | acct_to_users | Lowercase |
| activities | activities | Lowercase |
| applications | applications | Lowercase |
| billings | billings | Lowercase |
| connected_accounts | connected_accounts | Lowercase |
| forgot_passwords | forgot_passwords | Lowercase |
| invites | invites | Lowercase |
| sessions | sessions | Lowercase |
| users | users | Lowercase |

## Data Type Conversions

| MySQL Type | SQLite Type | Notes |
|------------|-------------|-------|
| int(N) unsigned NOT NULL AUTO_INCREMENT | INTEGER PRIMARY KEY AUTOINCREMENT | Primary key conversion |
| int(N) unsigned | INTEGER | Standard integer |
| int(N) | INTEGER | Standard integer |
| varchar(N) | TEXT | SQLite treats varchar as TEXT |
| text | TEXT | Direct mapping |
| enum(...) | TEXT | SQLite doesn't support ENUM |
| timestamp | DATETIME | With appropriate default handling |

## Common Issues and Solutions

### 1. Quote Escaping
**Problem**: MySQL uses `\'` for escaping single quotes, SQLite uses `''`
**Solution**: The `sed` command in step 5 handles this conversion

### 2. Table Name Case Sensitivity
**Problem**: MySQL dump may have different case than Go GORM models
**Solution**: The extract_inserts.py script maps table names correctly

### 3. MySQL-Specific Syntax
**Problem**: MySQL dump contains ENGINE, AUTO_INCREMENT, charset specifications
**Solution**: The migrate_mysql_to_sqlite.py script removes these

### 4. ENUM Data Types
**Problem**: SQLite doesn't support ENUM
**Solution**: Converted to TEXT, existing enum values work as strings

## Verification Steps

After migration, verify data integrity:

```bash
# Check row counts match expectations
sqlite3 your_database.db "SELECT name, COUNT(*) as count FROM sqlite_master WHERE type='table' GROUP BY name;"

# Test sample data with special characters
sqlite3 your_database.db "SELECT * FROM Categories WHERE CategoriesName LIKE '%''%' LIMIT 5;"

# Verify foreign key relationships still work
sqlite3 your_database.db "SELECT c.ContactsName, COUNT(l.LedgerId) FROM Contacts c LEFT JOIN Ledger l ON c.ContactsId = l.LedgerContactId GROUP BY c.ContactsId LIMIT 5;"
```

## Rollback Plan

If migration fails:

```bash
# Restore from backup
cp backend/cache/sqlite/sk_test_mj7zC9Al5M.db.backup backend/cache/sqlite/sk_test_mj7zC9Al5M.db
```

## Performance Notes

- Large dumps (100MB+) may take several minutes to process
- The quote escaping step (sed command) can be slow on very large files
- Consider splitting very large dumps if memory becomes an issue
- SQLite import is much faster than MySQL dump processing

## Troubleshooting

### Script Errors
- Ensure Python 3 is used (not Python 2)
- Check file permissions on input/output files
- Verify input file is valid MySQL dump format

### SQLite Import Errors
- Check for remaining quote escaping issues
- Verify table names match exactly (case sensitive)
- Ensure foreign key constraints are disabled during import

### Data Integrity Issues
- Compare row counts between source and target
- Spot check records with special characters
- Test application functionality after migration