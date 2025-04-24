CREATE TABLE nodes (
    id TEXT PRIMARY KEY,
    sequence TEXT DEFAULT 0,
    node_type TEXT NOT NULL,
    contents TEXT,
    children TEXT,
    deleted INTEGER DEFAULT 0,
    parent TEXT,
    created_at INTEGER,
    modified_at INTEGER,
    deleted_at INTEGER,
    FOREIGN KEY (parent) REFERENCES nodes (id)
);
