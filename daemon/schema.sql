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

--DROP TABLE link_group_move;
CREATE TABLE IF NOT EXISTS link_group_move
(
   id            INTEGER PRIMARY KEY AUTOINCREMENT,
   role          CHAR(1) NOT NULL DEFAULT '?',
   group_id      INTEGER NOT NULL,
   link_id       INTEGER NOT NULL
)
;

CREATE TABLE IF NOT EXISTS meta
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

DROP INDEX IF EXISTS idx_meta__kv_pair
;

CREATE INDEX idx_meta__kv_pair ON meta (kv_pair)
;

CREATE TRIGGER IF NOT EXISTS update_meta_modified
   AFTER UPDATE
   ON meta
   FOR EACH ROW
BEGIN
   UPDATE meta
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
   body         TEXT        NOT NULL DEFAULT '',
   head         TEXT        NOT NULL DEFAULT '',
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   created      INTEGER              DEFAULT (STRFTIME('%s', 'now'))
)
;

DROP TRIGGER IF EXISTS update_content_prevent
;

-- Trigger to prevent updates
CREATE TRIGGER IF NOT EXISTS update_content_prevent
   BEFORE UPDATE
   ON content
BEGIN
   SELECT RAISE(FAIL, 'Updates are not allowed on content table');
END
;

-- -- Trigger to prevent deletes
-- CREATE TRIGGER delete_content_prevent
--    BEFORE DELETE ON content
-- BEGIN
--    SELECT RAISE(FAIL, 'Deletes are not allowed on content table');
-- END;


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
   created      INTEGER              DEFAULT (STRFTIME('%s', 'now')),
   latest       INTEGER              DEFAULT (STRFTIME('%s', 'now')),
   archived     INTEGER     NOT NULL DEFAULT 0,
   deleted      INTEGER     NOT NULL DEFAULT 0,
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
   scheme       VARCHAR(5)   NOT NULL DEFAULT '',
   subdomain    VARCHAR(32)  NOT NULL DEFAULT '',
   sld          VARCHAR(32)  NOT NULL DEFAULT '',
   tld          VARCHAR(10)  NOT NULL DEFAULT '',
   port         VARCHAR(6)   NOT NULL DEFAULT '',
   path         VARCHAR(64)  NOT NULL DEFAULT '',
   query        VARCHAR(64)  NOT NULL DEFAULT '',
   fragment     VARCHAR(64)  NOT NULL DEFAULT '',
   original_url VARCHAR(256) NOT NULL DEFAULT '',
   host         VARCHAR(64)  NOT NULL GENERATED ALWAYS AS (CAST(
         CASE WHEN TRIM(subdomain) = '' THEN '' ELSE subdomain || '.' END ||
         sld ||
         CASE WHEN TRIM(tld) = '' THEN '' ELSE '.' || tld END ||
         CASE WHEN TRIM(port) = '' THEN '' ELSE ':' || tld END AS TEXT)),
   url          VARCHAR(256) NOT NULL GENERATED ALWAYS AS (CAST(scheme || '://'
      || CASE WHEN sld = '' THEN sld ELSE sld || '.' || tld END
      || CASE WHEN port = '' THEN '' ELSE ':' || port END
      || path
      || CASE WHEN query = '' THEN '' ELSE '?' || query END
      || CASE WHEN fragment = '' THEN '' ELSE '#' || fragment END AS TEXT)),
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   visited_time TEXT GENERATED ALWAYS AS (DATETIME(visited, 'unixepoch')) VIRTUAL,
   created      INTEGER               DEFAULT (STRFTIME('%s', 'now')),
   visited      INTEGER               DEFAULT (STRFTIME('%s', 'now')),
   archived     INTEGER      NOT NULL DEFAULT 0,
   deleted      INTEGER      NOT NULL DEFAULT 0,
   parsed       INTEGER      NOT NULL DEFAULT 0
)
;

DROP INDEX IF EXISTS idx_link__original_url
;

CREATE UNIQUE INDEX idx_link__original_url ON link (original_url)
;

--DROP TRIGGER IF  EXISTS update_link_visited
CREATE TRIGGER IF NOT EXISTS update_link_visited
   AFTER UPDATE
   ON link
   FOR EACH ROW
BEGIN
   UPDATE link
   SET
      visited = STRFTIME('%s', 'now'),
      title   = CASE WHEN title = '' THEN url ELSE title END
   WHERE
      id = old.id;
END
;

--=======================================================================--
-- VIEWS BELOW
--=======================================================================--

DROP VIEW IF EXISTS groups_type
;

CREATE VIEW groups_type AS
SELECT DISTINCT
   gt.type,
   gt.name,
   gt.plural,
   COUNT(DISTINCT g.id)                                    AS group_count,
   COUNT(DISTINCT g.archived = 1)                          AS groups_archived,
   COUNT(DISTINCT g.deleted = 1)                           AS groups_deleted,
   CAST(CASE
           WHEN g.id IS NULL THEN 0
           ELSE COUNT(DISTINCT rg.link_id) END AS INTEGER) AS link_count,
   COUNT(DISTINCT l.archived = 1)                          AS links_archived,
   COUNT(DISTINCT l.deleted = 1)                           AS links_deleted,
   gt.sort
FROM
   group_type gt
      LEFT JOIN `group` g
         ON gt.type = g.type
      LEFT JOIN link_group rg
         ON g.id = rg.group_id
      LEFT JOIN link l
         ON l.id = rg.link_id
GROUP BY
   gt.sort,
   gt.type,
   gt.name
ORDER BY
   gt.sort
;

DROP VIEW IF EXISTS link_group_move_half;

CREATE VIEW link_group_move_half AS
SELECT
   id,
   link_group_id,
   role,
   group_name,
   group_type,
   type_name,
   in_link_id         AS link_id,
   in_group_id        AS group_id,
   out_group_id <> 0  AS group_found,
   out_link_id <> 0   AS link_found,
   link_group_id <> 0 AS link_group_found,
   link_url
FROM
   (
      SELECT
         m.id,
         m.role,
         CAST(IFNULL(lg.id, 0) AS INTEGER)      AS link_group_id,
         CAST(IFNULL(g.name, '') AS TEXT)       AS group_name,
         CAST(IFNULL(gt.type, '') AS TEXT)      AS group_type,
         CAST(IFNULL(gt.name, '') AS TEXT)      AS type_name,
         CAST(IFNULL(m.link_id, 0) AS INTEGER)  AS in_link_id,
         CAST(IFNULL(m.group_id, 0) AS INTEGER) AS in_group_id,
         CAST(IFNULL(g.id, 0) AS INTEGER)       AS out_group_id,
         CAST(IFNULL(l.id, 0) AS INTEGER)       AS out_link_id,
         CAST(IFNULL(l.url, '') AS TEXT)        AS link_url
      FROM
         link_group_move m
         LEFT JOIN link_group lg ON m.group_id=lg.group_id AND m.link_id=lg.link_id
            LEFT JOIN link l ON m.link_id=l.id
         LEFT JOIN `group` g ON m.group_id=g.id
         LEFT JOIN group_type gt ON gt.type = g.type
      ) x

;

DROP VIEW IF EXISTS link_group_move_from_to;

CREATE VIEW link_group_move_from_to AS
SELECT
   f.type_name,
   f.id AS from_id,
   t.id AS to_id,
   f.link_group_id AS from_link_group_id,
   f.link_group_id<>0 AS from_link_group_found,
   f.group_id AS from_group_id,
   f.link_id,
   t.group_id AS to_group_id,
   f.group_name AS from_group_name,
   f.link_url,
   t.link_group_id AS to_link_group_id,
   t.link_group_id<>0 AS to_link_group_found,
   t.group_name AS to_group_name,
   f.group_found AS from_group_found,
   f.link_found AS link_found,
   t.group_found AS to_group_found
FROM
   link_group_move_half f
      JOIN link_group_move_half AS t ON f.link_id=t.link_id
WHERE true
   AND f.role='from'
   AND t.role='to'
;


DROP VIEW IF EXISTS link_group_move_exceptions;

CREATE VIEW link_group_move_exceptions AS
SELECT
   'from_link_group_missing' AS exception, --ssot[from_link_group_missing]: '([^']+)'
   from_group_id,
   from_group_name,
   link_id,
   link_url,
   to_group_id,
   to_group_name,
   from_id,
   to_id
FROM
   link_group_move_from_to
WHERE
   from_link_group_id = 0
   AND link_url <> ''
   AND to_group_id <> 0
UNION
SELECT
   'link_missing' AS exception, --ssot[link_missing]: '([^']+)'
   from_group_id,
   from_group_name,
   link_id,
   link_url,
   to_group_id,
   to_group_name,
   from_id,
   to_id
FROM
   link_group_move_from_to
WHERE
   link_found = 0
UNION
SELECT
   'from_group_missing' AS exception, --ssot[from_group_missing]: '([^']+)'
   from_group_id,
   from_group_name,
   link_id,
   link_url,
   to_group_id,
   to_group_name,
   from_id,
   to_id
FROM
   link_group_move_from_to
WHERE
   from_group_found = 0
UNION
SELECT
   'to_group_found' AS exception, --ssot[to_group_found]: '([^']+)'
   from_group_id,
   from_group_name,
   link_id,
   link_url,
   to_group_id,
   to_group_name,
   from_id,
   to_id
FROM
   link_group_move_from_to
WHERE
   to_group_found = 0
UNION
SELECT
   'to_link_group_found' AS exception, --ssot[to_link_group_found]: '([^']+)'
   from_group_id,
   from_group_name,
   link_id,
   link_url,
   to_group_id,
   to_group_name,
   from_id,
   to_id
FROM
   link_group_move_from_to
WHERE
   to_link_group_id <> 0;
;
