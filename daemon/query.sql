-- name: LoadGroup :one
SELECT
   *
FROM
   `group`
WHERE
   TRUE
   AND id = ?
   AND archived IN (sqlc.slice('groups_archived'))
   AND deleted IN (sqlc.slice('groups_deleted'))
LIMIT 1
;

-- name: LoadGroupName :one
SELECT name FROM `group` WHERE id = ?
;

-- name: ValidateGroups :many
SELECT
   id
FROM
   `group`
WHERE
   id IN (sqlc.slice('group_ids'))
GROUP BY
   id
;

-- name: ListLinkGroupMoveExceptions :many
SELECT
   exception,
   from_group_id,
   from_group_name,
   link_id,
   link_url,
   to_group_id,
   to_group_name
FROM
   link_group_move_exceptions
WHERE true
   AND from_id IN (sqlc.slice('from_ids'))
   AND to_id IN (sqlc.slice('to_ids'))
;

-- name: LoadLinkGroupMoveStatus :one
SELECT
   IFNULL(all_links_missing+no_to_groups_found,0)<>0 AS error_exception,
   *
FROM (
   SELECT
      IFNULL(COUNT(*) - SUM(link_found), 0) = 0                         AS all_links_missing,
      IFNULL(SUM(to_group_found), 0) = 0                                AS no_to_groups_found,
      IFNULL(COUNT(*), 0)                                               AS link_count,
      IFNULL(COUNT(*) - SUM(link_found), 0)                             AS link_exceptions,
      IFNULL(COUNT(*) - SUM(from_group_found), 0) <> 0                  AS from_group_exception,
      IFNULL(COUNT(*) - SUM(to_group_found), 0) <> 0                    AS to_group_exception,
      CAST(IFNULL(COUNT(*) - SUM(from_link_group_found), 0) AS INTEGER) AS from_link_group_exceptions,
      CAST(IFNULL(SUM(to_link_group_found), 0) AS INTEGER)              AS to_link_group_exceptions
   --SELECT *
   FROM
      link_group_move_from_to
   WHERE
      TRUE
      AND from_id IN (sqlc.slice('from_ids'))
      AND to_id IN (sqlc.slice('to_ids'))
      -- AND from_id IN (461, 462, 463, 464, 465)
      -- AND to_id IN (466, 467, 468, 469, 470)
)x

;

-- SELECT * FROM link_group_move_from_to
-- WHERE from_id IN (411,412,413,414,415,416,417,418,419,420) AND to_id IN (411,412,413,414,415,416,417,418,419,420);

-- name: DeleteLinkGroupMove :exec
DELETE
FROM
   link_group_move
WHERE
   id IN (sqlc.slice('link_group_ids'))
;

-- name: InsertLinkGroupMoveFromVarJSON :many
INSERT INTO link_group_move (role, group_id, link_id)
SELECT
   CAST(IFNULL(JSON_EXTRACT(j.value, '$.group_role'), '') AS TEXT) AS role,
   CAST(IFNULL(JSON_EXTRACT(j.value, '$.group_id'), 0) AS INTEGER) AS group_id,
   CAST(IFNULL(JSON_EXTRACT(j.value, '$.link_id'), 0) AS INTEGER)  AS link_id
FROM
   var
      JOIN JSON_EACH(var.value) j
         ON var.key = 'json'
WHERE
   var.id = ?
RETURNING role,id
;

-- name: ListLinkIds :many
SELECT
   id
FROM
   link
WHERE
   id IN (sqlc.slice('link_ids'))
GROUP BY
   id
;

-- name: ListGroupIds :many
SELECT
   id
FROM
   `group`
WHERE
   id IN (sqlc.slice('group_ids'))
GROUP BY
   id
;

-- name: LoadGroupIdBySlug :one
SELECT id FROM `group` WHERE slug = ?
;

-- name: LoadGroupTypeAndName :one
SELECT type, name FROM `group` WHERE id = ?
;

-- name: MergeLinksGroups :exec
UPDATE OR IGNORE link_group
SET
   group_id = ?
WHERE
   group_id IN (sqlc.slice('group_ids'))
;

-- name: MoveLinksToGroup :many
UPDATE link_group
SET
   group_id = ?,
   latest   = STRFTIME('%s', 'now')
WHERE
   TRUE
   AND group_id = ?
   AND link_id IN (sqlc.slice('link_ids'))
RETURNING id
;

-- name: MarkGroupsDeleted :exec
UPDATE `group` SET deleted = 1 WHERE id IN (sqlc.slice('group_ids'))
;

-- name: LoadAltGroupIdsByName :many
SELECT id FROM `group` WHERE id <> ? AND name = ?
;

-- name: UpdateGroupName :exec
UPDATE `group`
SET
   name = ?,
   slug = ?
WHERE
   id = ?
;

-- name: LoadGroupType :one
SELECT * FROM group_type WHERE type = ? LIMIT 1
;

-- name: ListGroupsByType :many
SELECT
   g.*,
   CAST(gt.name AS TEXT) AS type_name
FROM
   `group` g
      JOIN group_type gt
         ON gt.type = g.type
WHERE
   TRUE
   AND g.type = ?
   AND g.archived IN (sqlc.slice('groups_archived'))
   AND g.deleted IN (sqlc.slice('groups_deleted'))
ORDER BY
   g.name
;

-- name: LoadGroupsBySlug :one
SELECT
   *
FROM
   `group`
WHERE
   TRUE
   AND slug = ?
   AND archived IN (sqlc.slice('groups_archived'))
   AND deleted IN (sqlc.slice('groups_deleted'))
LIMIT 1
;

-- name: UpsertGroupsFromVarJSON :exec
INSERT INTO `group` (name, type, slug, archived, deleted)
SELECT
   JSON_EXTRACT(r.value, '$.name')     AS name,
   JSON_EXTRACT(r.value, '$.type')     AS type,
   JSON_EXTRACT(r.value, '$.slug')     AS slug,
   JSON_EXTRACT(r.value, '$.archived') AS archived,
   JSON_EXTRACT(r.value, '$.deleted')  AS deleted
FROM
   var
      JOIN JSON_EACH(var.value) r
         ON var.key = 'json'
WHERE
   var.id = ?
ON CONFLICT (slug)
   DO UPDATE
   SET
      archived = excluded.archived,
      deleted  = excluded.deleted,
      latest   = STRFTIME('%s', 'now')
;

-- name: ListGroupsType :many
SELECT * FROM groups_type
;

-- name: LoadGroupTypeWithStats :one
SELECT * FROM groups_type WHERE type = ? LIMIT 1
;

-- name: LoadLinkIdByUrl :one
SELECT id FROM link WHERE original_url = ? LIMIT 1
;

-- name: LoadLink :one
SELECT
   id,
   original_url,
   created_time,
   visited_time,
   scheme,
   subdomain,
   sld,
   tld,
   path,
   query,
   fragment,
   port,
   url,
   title
FROM
   link
WHERE
   TRUE
   AND id = ?
LIMIT 1
;

-- name: UpsertLink :one
INSERT INTO link
   (original_url, title, visited)
VALUES
   (?, ?, STRFTIME('%s', 'now'))
ON CONFLICT (original_url)
   DO UPDATE
   SET
      title   = excluded.title,
      visited = STRFTIME('%s', 'now')
RETURNING id
;

-- name: UpsertLinksFromVarJSON :exec
INSERT INTO link (original_url, title, archived, deleted, visited)
SELECT
   JSON_EXTRACT(r.value, '$.original_url'),
   JSON_EXTRACT(r.value, '$.title'),
   JSON_EXTRACT(r.value, '$.archived'),
   JSON_EXTRACT(r.value, '$.deleted'),
   STRFTIME('%s', 'now')
FROM
   var
      JOIN JSON_EACH(var.value) r
         ON var.key = 'json'
WHERE
   var.id = ?
ON CONFLICT (original_url)
   DO UPDATE
   SET
      title    = excluded.title,
      archived = excluded.archived,
      deleted  = excluded.deleted,
      visited  = STRFTIME('%s', 'now')
;

-- name: UpsertLinkGroupsFromVarJSON :exec
INSERT INTO link_group (group_id, link_id)
SELECT
   g.id,
   r.id
FROM
   var
      JOIN JSON_EACH(var.value) j
         ON var.key = 'json'
      JOIN link r
         ON r.original_url = JSON_EXTRACT(j.value, '$.link_url')
      JOIN `group` g
         ON TRUE
      AND g.name = JSON_EXTRACT(j.value, '$.group_name')
      AND g.type = JSON_EXTRACT(j.value, '$.group_type')
WHERE
   var.id = ?
ON CONFLICT (group_id, link_id)
   DO UPDATE
   SET latest = STRFTIME('%s', 'now')
;

-- name: UpsertLinkMetaFromVarJSON :exec
INSERT INTO meta (link_id, key, value)
SELECT
   CAST(JSON_EXTRACT(r.value, '$.link_id') AS INTEGER),
   CAST(JSON_EXTRACT(r.value, '$.key') AS TEXT),
   CAST(JSON_EXTRACT(r.value, '$.value') AS TEXT)
FROM
   var
      JOIN JSON_EACH(var.value) r
         ON var.key = 'json'
WHERE
   var.id = ?
ON CONFLICT (link_id,key)
   DO UPDATE
   SET modified = STRFTIME('%s', 'now')
;

-- name: ListLinksLite :many
SELECT
   id,
   original_url                        AS url,
   CAST(IFNULL(created, 0) AS INTEGER) AS created,
   CAST(IFNULL(visited, 0) AS INTEGER) AS visited
FROM
   link
WHERE
   TRUE
   AND id > ?
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY
   id
LIMIT 100
;

-- name: ListLinks :many
SELECT
   *
FROM
   link
WHERE
   TRUE
   AND id > ?
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY
   original_url
LIMIT 100
;

-- name: GetLinkURLs :many
SELECT
   CAST(IFNULL(url, '<invalid>') AS TEXT) AS url
FROM
   link
WHERE
   id IN (sqlc.slice('link_ids'))
;

-- name: ArchiveLinks :exec
UPDATE link
SET
   archived=1
WHERE
   id IN (sqlc.slice('link_ids'))
;

-- name: MarkLinksDeleted :exec
UPDATE link
SET
   deleted=1
WHERE
   id IN (sqlc.slice('link_ids'))
;

-- name: MarkLinksDeletedByGroupIds :exec
UPDATE link
SET
   deleted=1
WHERE
   TRUE
   AND id NOT IN (
      SELECT
         lg.link_id
      FROM
         link_group lg
      WHERE
         lg.group_id = ?
      )
   AND id IN (
      SELECT lg.link_id FROM link_group lg WHERE lg.group_id IN (sqlc.slice('group_ids'))
      )
;

-- name: ListFilteredLinks :many
SELECT
   id,
   original_url,
   created_time,
   visited_time,
   title,
   scheme,
   subdomain,
   sld,
   tld,
   path,
   query,
   fragment,
   port,
   archived,
   deleted
FROM
   link
WHERE
   TRUE
   AND id IN (sqlc.slice('ids'))
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY
   original_url
;

-- name: LoadLatestContent :one
-- TODO: Untested, ensure query works
SELECT
   *
FROM
   content
WHERE
   link_id = ?
GROUP BY
   link_id,
   created
HAVING
   created = MAX(created)
;

-- name: InsertContent :exec
INSERT INTO content
   (link_id, head, body)
VALUES
   (?, ?, ?)
;
-- name: ListLatestUnparsedLinkURLs :many
SELECT
   id,
   original_url
FROM
   link
WHERE
   TRUE
   AND parsed = 0
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY
   id DESC
LIMIT 8
;

-- LIMIT was chosen as slice len == slice cap for 8

-- name: UpdateParsedLink :exec
UPDATE link
SET
   title     = ?,
   scheme    = ?,
   subdomain = ?,
   sld       = ?,
   tld       = ?,
   port      = ?,
   path      = ?,
   query     = ?,
   fragment  = ?,
   parsed    = 1
WHERE
   original_url = ?
;

-- name: ListLinkIdsByGroup :many
SELECT
   CAST(l.id AS INTEGER) AS link_id
FROM
   link l
      JOIN link_group rg
         ON l.id = rg.link_id
      JOIN `group` g
         ON g.id = rg.group_id
WHERE
   TRUE
   AND g.slug IN (sqlc.slice('slugs'))
   AND l.archived IN (sqlc.slice('links_archived'))
   AND l.deleted IN (sqlc.slice('links_deleted'))
;

-- name: ListLinkIdsByMeta :many
SELECT
   CAST(m.link_id AS INTEGER) AS link_id
FROM
   meta m
      JOIN link l
         ON l.id = m.link_id
WHERE
   TRUE
   AND m.kv_pair IN (sqlc.slice('kv_pairs'))
   AND m.key IN (sqlc.slice('keys'))
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
;

-- name: ListLinkIdsByGroupType :many
SELECT
   CAST(link_id AS INTEGER) AS link_id
FROM
   link_group lg
      JOIN `group` g
         ON lg.group_id = g.id
      JOIN link l
         ON l.id = lg.link_id
WHERE
   TRUE
   AND g.type IN (sqlc.slice('groupTypes'))
   AND l.archived IN (sqlc.slice('links_archived'))
   AND l.deleted IN (sqlc.slice('links_deleted'))
;

-- name: ListLinkIdsNotInGroupType :many
SELECT
   CAST(l.id AS INTEGER) AS link_id
FROM
   link l
WHERE
   TRUE
   AND l.archived IN (sqlc.slice('links_archived'))
   AND l.deleted IN (sqlc.slice('links_deleted'))
   AND l.id NOT IN (
      SELECT
         lg.link_id
      FROM
         link_group lg
            JOIN `group` g
               ON lg.group_id = g.id
      WHERE
         g.type IN (sqlc.slice('groupTypes'))
      )
;

-- name: UpsertVar :one
INSERT INTO var (key, value)
VALUES (?, ?)
ON CONFLICT (key) DO UPDATE SET value = excluded.value
RETURNING id
;

-- name: DeleteVar :exec
DELETE FROM var WHERE id = ?
;

-- name: ListLinkMeta :many
SELECT
   m.*
FROM
   meta m
      JOIN link l
         ON m.link_id = l.id
WHERE
   TRUE
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY
   link_id,
   key DESC
;

-- name: ListLinkMetaForLinkId :many
SELECT
   m.key,
   m.value
FROM
   meta m
      JOIN link l
         ON m.link_id = l.id
WHERE
   link_id = ?
;

-- name: UpsertMetaFromVarJSON :exec
INSERT INTO meta (link_id, key, value)
SELECT
   l.id,
   JSON_EXTRACT(kv.value, '$.key'),
   JSON_EXTRACT(kv.value, '$.value')
FROM
   var
      JOIN JSON_EACH(var.value) kv
         ON var.key = 'json'
      JOIN link l
         ON l.original_url = JSON_EXTRACT(kv.value, '$.url')
WHERE
   var.id = ?
ON CONFLICT (link_id,key)
   DO UPDATE
   SET value = excluded.value
;

