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

CREATE TABLE IF NOT EXISTS key_value
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   resource_id  INTEGER,
   key          VARCHAR(32),
   value        TEXT,
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   modified_time TEXT GENERATED ALWAYS AS (DATETIME(modified, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   modified      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   UNIQUE (resource_id, key)
)
;

CREATE TRIGGER IF NOT EXISTS update_key_value_modified
   AFTER UPDATE
   ON key_value
   FOR EACH ROW
BEGIN
   UPDATE key_value SET modified = CASE WHEN old.value=new.value THEN modified
      ELSE STRFTIME('%s', 'now') END
   WHERE id = old.id;
END
;



CREATE TABLE IF NOT EXISTS group_type
(
   type        CHAR(1) PRIMARY KEY,
   sort        INTEGER,
   name        VARCHAR(32),
   description VARCHAR(128)
)
;

DELETE FROM group_type WHERE 1=1
;

DELETE FROM sqlite_sequence WHERE name = 'group_type'
;

INSERT
INTO group_type
   (sort, type, name, description)
VALUES
   (1, 'G', 'TabGroup', 'Browser''s name for the containing Tab Group'),
   (2, 'T', 'Tag', 'Human-specified keywords for resource content'),
   (3, 'C', 'Category', 'AI-generated top-level categorization'),
   (4, 'K', 'Keyword', 'AI-generated notable tags for resource content'),
   (5, 'I', 'Invalid', 'Unspecified or not a valid group type')
;

CREATE TABLE IF NOT EXISTS `group`
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   name         VARCHAR(32),
   type         CHAR(1),
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   latest_time  TEXT GENERATED ALWAYS AS (DATETIME(latest, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   latest       INTEGER DEFAULT (STRFTIME('%s', 'now')),
   UNIQUE (name, type)
)
;
DROP INDEX IF EXISTS idx_group__type;
CREATE INDEX idx_group__type ON `group`(type);

CREATE TRIGGER IF NOT EXISTS update_group_latest
   AFTER UPDATE
   ON `group`
   FOR EACH ROW
BEGIN
   UPDATE `group` SET latest = STRFTIME('%s', 'now') WHERE id = old.id;
END
;

CREATE TABLE IF NOT EXISTS resource_group
(
   group_id    INTEGER,
   resource_id INTEGER,
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   latest_time  TEXT GENERATED ALWAYS AS (DATETIME(latest, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   latest       INTEGER DEFAULT (STRFTIME('%s', 'now')),
   UNIQUE (group_id, resource_id)
)
;

CREATE TRIGGER IF NOT EXISTS update_resource_group_latest
   AFTER UPDATE
   ON resource_group
   FOR EACH ROW
BEGIN
   UPDATE resource_group
   SET latest = STRFTIME('%s', 'now')
   WHERE 1=1
      AND group_id = old.group_id
      AND resource_id = old.resource_id;
END
;

CREATE TABLE IF NOT EXISTS resource
(
   id           INTEGER PRIMARY KEY AUTOINCREMENT,
   url          VARCHAR(256) UNIQUE,
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   visited_time TEXT GENERATED ALWAYS AS (DATETIME(visited, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   visited      INTEGER DEFAULT (STRFTIME('%s', 'now'))
)
;

CREATE TRIGGER IF NOT EXISTS update_resource_visited
   AFTER UPDATE
   ON resource
   FOR EACH ROW
BEGIN
   UPDATE resource SET visited = STRFTIME('%s', 'now') WHERE id = old.id;
END
;

--=======================================================================--
-- VIEWS BELOW
--=======================================================================--
DROP VIEW IF EXISTS groups_with_counts;

CREATE VIEW groups_with_counts AS
SELECT
   g.id,
   COUNT(*) AS resource_count,
   g.name,
   g.type,
   gt.name AS type_name
FROM `group` g
  JOIN resource_group rg ON rg.group_id=g.id
  JOIN group_type gt ON gt.type=g.type
GROUP BY
   g.id,
   g.name,
   gt.name
ORDER BY
   gt.name,
   iif(g.name='<none>',0,1),
   g.name;
