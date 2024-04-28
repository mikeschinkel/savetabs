CREATE TABLE IF NOT EXISTS var
(
   id    INTEGER PRIMARY KEY AUTOINCREMENT,
   key   VARCHAR(10) NOT NULL,
   value VARCHAR(128),
   UNIQUE (key)
)
;

CREATE TABLE IF NOT EXISTS action
(
   id   INTEGER PRIMARY KEY AUTOINCREMENT,
   name VARCHAR(16)
)
;

CREATE TABLE IF NOT EXISTS history
(
   id          INTEGER PRIMARY KEY AUTOINCREMENT,
   foreign_id  INTEGER,
   date_time   INTEGER,
   action_id   INTEGER,
   action_data BLOB -- JSONB
)
;

CREATE TABLE IF NOT EXISTS metadata
(
   id            INTEGER PRIMARY KEY AUTOINCREMENT,
   link_id       INTEGER      NOT NULL,
   key           VARCHAR(32)  NOT NULL,
   value         VARCHAR(512) NOT NULL,
   kv_pair       VARCHAR(544) NOT NULL GENERATED ALWAYS AS (key || '=' || value) VIRTUAL,
   created_time  TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   modified_time TEXT GENERATED ALWAYS AS (DATETIME(modified, 'unixepoch')) VIRTUAL,
   created       INTEGER DEFAULT (STRFTIME('%s', 'now')),
   modified      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   UNIQUE (link_id, key)
)
;

DROP INDEX IF EXISTS idx_metadata__kv_pair
;

CREATE INDEX idx_metadata__kv_pair ON metadata (kv_pair)
;

CREATE TRIGGER IF NOT EXISTS update_metadata_modified
   AFTER UPDATE
   ON metadata
   FOR EACH ROW
BEGIN
   UPDATE metadata
   SET
      modified = CASE
                    WHEN old.value = new.value THEN modified
                    ELSE STRFTIME('%s', 'now') END
   WHERE
      id = old.id;
END
;

CREATE TABLE IF NOT EXISTS content
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   link_id      INTEGER     NOT NULL,
   title        VARCHAR(64) NOT NULL,
   body         TEXT        NOT NULL,
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now'))
)
;

CREATE TRIGGER IF NOT EXISTS update_content_modified
   AFTER UPDATE
   ON content
   FOR EACH ROW
BEGIN
   UPDATE content
   SET
      created = CASE
                   WHEN old.created = new.created THEN created
                   ELSE STRFTIME('%s', 'now') END
   WHERE
      id = old.id;
END
;

CREATE TABLE IF NOT EXISTS group_type
(
   type        CHAR(1) PRIMARY KEY,
   sort        INTEGER,
   name        VARCHAR(32),
   plural      VARCHAR(32),
   description VARCHAR(128)
)
;


DELETE FROM group_type WHERE TRUE
;

DELETE FROM sqlite_sequence WHERE name = 'group_type'
;

INSERT
   INTO group_type
   (sort, type, name,       plural,       description                                 )
VALUES
   (1,    'G',  'TabGroup', 'TabGroups',  'Browser''s name for containing Tab Group'  ),
   (2,    'C',  'Category', 'Categories', 'AI-generated top-level categorization'     ),
   (3,    'T',  'Tag',      'Tags',       'Human-specified keywords for link content' ),
   (4,    'K',  'Keyword',  'Keywords',   'AI-generated notable tags for link content'),
   (5,    'B',  'Bookmark', 'Bookmarks',  'Browser''s name for saved links'           ),
   (6,    'I',  'Invalid',  'Invalids',   'Unspecified or not a valid group type'     )
;

CREATE TABLE IF NOT EXISTS `group`
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   name         VARCHAR(32) NOT NULL,
   type         CHAR(1)     NOT NULL,
   slug         VARCHAR(32) NOT NULL,
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   latest_time  TEXT GENERATED ALWAYS AS (DATETIME(latest, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   latest       INTEGER DEFAULT (STRFTIME('%s', 'now')),
   archived     INTEGER NOT NULL DEFAULT 0,
   deleted      INTEGER NOT NULL DEFAULT 0,
   UNIQUE (name, type),
   UNIQUE (slug)
)
;

DROP INDEX IF EXISTS idx_group__type
;

CREATE INDEX idx_group__type ON `group` (type)
;

CREATE TRIGGER IF NOT EXISTS update_group_latest
   AFTER UPDATE
   ON `group`
   FOR EACH ROW
BEGIN
   UPDATE `group` SET latest = STRFTIME('%s', 'now') WHERE id = old.id;
END
;

CREATE TABLE IF NOT EXISTS link_group
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   group_id     INTEGER NOT NULL,
   link_id      INTEGER NOT NULL,
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   latest_time  TEXT GENERATED ALWAYS AS (DATETIME(latest, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   latest       INTEGER DEFAULT (STRFTIME('%s', 'now')),
   UNIQUE (group_id, link_id)
)
;

DROP TRIGGER IF EXISTS update_link_group_latest
;

CREATE TRIGGER IF NOT EXISTS update_link_group_latest
   AFTER UPDATE
   ON link_group
   FOR EACH ROW
BEGIN
   UPDATE link_group
   SET
      latest = STRFTIME('%s', 'now')
   WHERE
      TRUE
      AND group_id = old.group_id
      AND link_id = old.link_id;
END
;

CREATE TABLE IF NOT EXISTS link
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   title        VARCHAR(128) NOT NULL DEFAULT '',
   scheme       VARCHAR(5)  NOT NULL DEFAULT '',
   subdomain    VARCHAR(32) NOT NULL DEFAULT '',
   sld          VARCHAR(32) NOT NULL DEFAULT '',
   tld          VARCHAR(10) NOT NULL DEFAULT '',
   port         VARCHAR(6) NOT NULL DEFAULT '',
   path         VARCHAR(64) NOT NULL DEFAULT '',
   query        VARCHAR(64) NOT NULL DEFAULT '',
   fragment     VARCHAR(64) NOT NULL DEFAULT '',
   original_url VARCHAR(256) NOT NULL DEFAULT '',
   url          VARCHAR(256) GENERATED ALWAYS AS (scheme||'://'
      ||CASE WHEN sld='' THEN sld ELSE sld||'.'||tld END
      ||CASE WHEN port='' THEN '' ELSE ':'||port END
      ||path
      ||CASE WHEN query='' THEN '' ELSE '?'||query END
      ||CASE WHEN fragment='' THEN '' ELSE '#'||fragment END),
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   visited_time TEXT GENERATED ALWAYS AS (DATETIME(visited, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   visited      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   archived     INTEGER NOT NULL DEFAULT 0,
   deleted      INTEGER NOT NULL DEFAULT 0

)
;

DROP INDEX IF EXISTS idx_link__original_url;
CREATE UNIQUE INDEX idx_link__original_url ON link(original_url);

--DROP TRIGGER IF  EXISTS update_link_visited
CREATE TRIGGER IF NOT EXISTS update_link_visited
   AFTER UPDATE
   ON link
   FOR EACH ROW
BEGIN
   UPDATE link
   SET visited = STRFTIME('%s', 'now'),
      title = CASE WHEN title='' THEN url ELSE title END
   WHERE id = old.id;
END
;

--=======================================================================--
-- VIEWS BELOW
--=======================================================================--
