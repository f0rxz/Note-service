#!/bin/bash

sqlite3 notes.db <<EOF
CREATE TABLE IF NOT EXISTS notes (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);
.exit
EOF

echo "Database and table was created successfully"
