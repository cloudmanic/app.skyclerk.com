# MySQL to SQLite Migration - Quick Reference

## One-Command Migration

```bash
./migrate_db.sh your_mysql_dump.sql [path/to/sqlite/db]
```

**Example:**
```bash
./migrate_db.sh production_dump.sql backend/cache/sqlite/skyclerk.db
```

## Manual Step-by-Step (if automation fails)

```bash
# 1. Backup
cp backend/cache/sqlite/sk_test_mj7zC9Al5M.db backend/cache/sqlite/sk_test_mj7zC9Al5M.db.backup

# 2. Convert
python3 migrate_mysql_to_sqlite.py mysql_dump.sql converted.sql

# 3. Extract INSERTs
python3 extract_inserts.py converted.sql inserts.sql

# 4. Single-row INSERTs
python3 convert_to_single_inserts.py inserts.sql single.sql

# 5. Fix quotes
sed "s/\\\\'/\\'\\'/g" single.sql > final.sql

# 6. Import
sqlite3 backend/cache/sqlite/sk_test_mj7zC9Al5M.db ".read final.sql"

# 7. Verify
sqlite3 backend/cache/sqlite/sk_test_mj7zC9Al5M.db "SELECT COUNT(*) FROM Categories;"
```

## Files Created

- `MIGRATION.md` - Complete documentation
- `migrate_db.sh` - Automated migration script
- `mysql_data_imported.sql` - Example of final working import file

## Verification Commands

```bash
# Check all table counts
sqlite3 your_db.db "SELECT name, (SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=t.name) as count FROM (SELECT DISTINCT name as name FROM sqlite_master WHERE type='table') t;"

# Test special characters
sqlite3 your_db.db "SELECT * FROM Categories WHERE CategoriesName LIKE '%''%' LIMIT 3;"
```

## Rollback

```bash
cp your_db.db.backup.TIMESTAMP your_db.db
```

## Common Issues

1. **Quote errors**: Run the sed command from step 5
2. **Table not found**: Check table mappings in extract_inserts.py
3. **Permission errors**: Check file permissions and paths
4. **Large files**: May need to split very large dumps (>500MB)