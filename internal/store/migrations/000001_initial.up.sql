CREATE TABLE mods (
    name             TEXT PRIMARY KEY NOT NULL,
    archive_filename TEXT NOT NULL DEFAULT '',
    downloaded_md5   TEXT NOT NULL DEFAULT '',
    installed_md5    TEXT NOT NULL DEFAULT ''
);

CREATE TABLE settings (
    key   TEXT PRIMARY KEY NOT NULL,
    value TEXT NOT NULL
);

CREATE TABLE git_revisions (
    url      TEXT PRIMARY KEY NOT NULL,
    revision TEXT NOT NULL
);
