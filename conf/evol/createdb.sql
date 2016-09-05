CREATE TABLE IF NOT EXISTS ipsums(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    uri TEXT NOT NULL,
    desc TEXT,
    adminEmail TEXT NOT NULL,
    adminKey TEXT NOT NULL,
    resetToken TEXT,
    resetTS INTEGER,
    deleteToken TEXT,
    deleteTS INTEGER,
    created INTEGER
);

CREATE UNIQUE INDEX if not exists idx_ipsums_uri on ipsums (uri);

CREATE TABLE IF NOT EXISTS ipsumtext(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    ipsum_id INTEGER NOT NULL,
    data TEXT NOT NULL,
    created INTEGER,
    FOREIGN KEY(ipsum_id) REFERENCES ipsums(id)

);

CREATE INDEX if not exists idx_ipsumtext_ipsum_id on ipsumtext (ipsum_id);
