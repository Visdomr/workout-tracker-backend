#!/bin/bash
echo "=== Registered Users ==="
sqlite3 workout_tracker.db -header -column "SELECT id, username, email, created_at FROM users;"
echo ""
echo "Total users: $(sqlite3 workout_tracker.db "SELECT COUNT(*) FROM users;")"
