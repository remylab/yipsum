    CREATE TABLE IF NOT EXISTS ipsums(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        uri TEXT NOT NULL,
        desc TEXT,
        adminKey TEXT NOT NULL,
        newAdminKey,
        created INTEGER
    );
    