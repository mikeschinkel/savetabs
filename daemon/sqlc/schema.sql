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

DROP INDEX IF EXISTS idx_key_value__kv_pair
;

CREATE INDEX idx_key_value__kv_pair ON key_value (kv_pair)
;

CREATE TRIGGER IF NOT EXISTS update_key_value_modified
   AFTER UPDATE
   ON key_value
   FOR EACH ROW
BEGIN
   UPDATE key_value
   SET
      modified = CASE
                    WHEN old.value = new.value THEN modified
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
   url          VARCHAR(256) UNIQUE,
   created_time TEXT GENERATED ALWAYS AS (DATETIME(created, 'unixepoch')) VIRTUAL,
   visited_time TEXT GENERATED ALWAYS AS (DATETIME(visited, 'unixepoch')) VIRTUAL,
   created      INTEGER DEFAULT (STRFTIME('%s', 'now')),
   visited      INTEGER DEFAULT (STRFTIME('%s', 'now'))
)
;

CREATE TRIGGER IF NOT EXISTS update_link_visited
   AFTER UPDATE
   ON link
   FOR EACH ROW
BEGIN
   UPDATE link SET visited = STRFTIME('%s', 'now') WHERE id = old.id;
END
;

--=======================================================================--
-- VIEWS BELOW
--=======================================================================--
DROP VIEW IF EXISTS groups_with_counts_view
;

-- SELECT * FROM groups_with_counts_view;
CREATE VIEW groups_with_counts_view AS
SELECT
   g.id,
   COUNT(DISTINCT rg.link_id) AS link_count,
   g.name,
   g.type,
   g.slug,
   gt.name                    AS type_name,
   gt.plural                  AS type_plural
FROM
   `group` g
      JOIN group_type gt
         ON gt.type = g.type
      LEFT JOIN link_group rg
         ON rg.group_id = g.id
GROUP BY
   g.id,
   g.name,
   gt.name
ORDER BY
   gt.name,
   g.name
;


DROP VIEW IF EXISTS link_group_ids_view
;

-- SELECT * FROM link_group_ids_view;
CREATE VIEW IF NOT EXISTS link_group_ids_view AS
SELECT
   link_id,
   group_ids,
   group_types,
   group_slugs,
   group_names,
   '''' || REPLACE(group_types, ',', ''',''') || '''' AS quoted_group_types,
   '''' || REPLACE(group_names, ',', ''',''') || '''' AS quoted_group_names,
   '''' || REPLACE(group_slugs, ',', ''',''') || '''' AS quoted_group_slugs
FROM
   (
      SELECT DISTINCT
         rg.link_id,
         GROUP_CONCAT(DISTINCT g.id)                           AS group_ids,
         GROUP_CONCAT(DISTINCT g.type)                         AS group_types,
         GROUP_CONCAT(DISTINCT g.type || ':' || g.name)        AS group_names,
         GROUP_CONCAT(DISTINCT LOWER(g.type) || ':' || g.slug) AS group_slugs
      FROM
         `group` g
            JOIN link_group rg
               ON g.id = rg.group_id
      GROUP BY
         rg.link_id
      ) x
;

DROP VIEW IF EXISTS links_view
;

-- SELECT * FROM links_view;
CREATE VIEW IF NOT EXISTS links_view AS
SELECT
   r.id,
   r.id      AS link_id,
   r.url,
   g.id      AS group_id,
   g.name    AS group_name,
   g.slug    AS group_slug,
   g.type    AS group_type,
   gt.name   AS type_name,
   sld.value AS domain,
   gs.group_ids,
   gs.group_types,
   gs.group_names,
   gs.quoted_group_types,
   gs.quoted_group_slugs,
   gs.quoted_group_names
FROM
   `group` g
      JOIN group_type gt
         ON gt.type = g.type
      LEFT JOIN link_group rg
         ON g.id = rg.group_id
      LEFT JOIN link r
         ON r.id = rg.link_id
      LEFT JOIN key_value sld
         ON sld.key = 'sld' AND sld.link_id = r.id
      LEFT JOIN link_group_ids_view gs
         ON r.id = gs.link_id
ORDER BY
   g.name,
   r.url
;





