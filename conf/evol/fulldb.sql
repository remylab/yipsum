CREATE TABLE IF NOT EXISTS ipsums(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    uri TEXT NOT NULL,
    desc TEXT,
    adminKey TEXT NOT NULL,
    newAdminKey,
    adminEmail TEXT NOT NULL,
    newAdminEmail,
    created INTEGER
);

CREATE UNIQUE INDEX if not exists idx_ipsums_uri on ipsums (uri);
