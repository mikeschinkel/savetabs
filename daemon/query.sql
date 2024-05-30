-- name: LoadGroup :one
SELECT * FROM `group`
WHERE true
   AND id = ?
   AND archived IN (sqlc.slice('groups_archived'))
   AND deleted IN (sqlc.slice('groups_deleted'))
LIMIT 1;

-- name: LoadGroupName :one
SELECT name FROM `group` WHERE id = ?;

-- name: LoadGroupTypeAndName :one
SELECT type,name FROM `group` WHERE id = ?;

-- name: MergeLinksGroups :exec
UPDATE OR IGNORE link_group
SET group_id = ?
WHERE group_id IN (sqlc.slice('group_ids'));

-- name: MarkGroupsDeleted :exec
UPDATE `group` SET deleted = 1 WHERE id IN (sqlc.slice('group_ids'));

-- name: LoadAltGroupIdsByName :many
SELECT id FROM `group` WHERE id <> ? AND name = ?;

-- name: UpdateGroupName :exec
UPDATE `group` SET name = ?, slug = ?
WHERE id = ?;

-- name: LoadGroupType :one
SELECT * FROM group_type WHERE type = ? LIMIT 1;

-- name: ListGroupsByType :many
SELECT
   g.*,
   CAST(gt.name AS TEXT) AS type_name
FROM `group` g
   JOIN group_type gt ON gt.type = g.type
WHERE true
   AND g.type = ?
   AND g.archived IN (sqlc.slice('groups_archived'))
   AND g.deleted IN (sqlc.slice('groups_deleted'))
ORDER BY g.name;

-- name: LoadGroupsBySlug :one
SELECT * FROM `group`
WHERE true
   AND slug = ?
   AND archived IN (sqlc.slice('groups_archived'))
   AND deleted IN (sqlc.slice('groups_deleted'))
LIMIT 1;

-- name: UpsertGroupsFromVarJSON :exec
INSERT INTO `group` (name,type,slug)
SELECT
   json_extract(r.value,'$.name') AS name,
   json_extract(r.value,'$.type') AS type,
   json_extract(r.value,'$.slug') AS slug
FROM var
   JOIN json_each( var.value ) r ON var.key='json'
WHERE var.id = ?
    ON CONFLICT (name,type)
        DO UPDATE
            SET latest = strftime('%s','now');

-- name: ListGroupsType :many
SELECT * FROM groups_type;

-- name: LoadGroupTypeWithStats :one
SELECT * FROM groups_type WHERE type = ? LIMIT 1;

-- name: LoadLinkIdByUrl :one
SELECT id FROM link WHERE original_url = ? LIMIT 1;

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
FROM link
WHERE true
   AND id = ?
LIMIT 1
;

-- name: UpsertLink :one
INSERT INTO link
   (original_url,title,visited)
VALUES
   (?,?,strftime('%s','now'))
ON CONFLICT (original_url)
   DO UPDATE
      SET
         title = excluded.title,
         visited = strftime('%s','now')
RETURNING id
;

-- name: UpsertLinksFromVarJSON :exec
INSERT INTO link (original_url,title,visited)
SELECT
   json_extract(r.value,'$.original_url'),
   json_extract(r.value,'$.title'),
   strftime('%s','now')
FROM var
   JOIN json_each( var.value ) r ON var.key='json'
WHERE var.id = ?
ON CONFLICT (original_url)
   DO UPDATE
   SET
      title = excluded.title,
      visited = strftime('%s','now')
;

-- name: UpsertLinkGroupsFromVarJSON :exec
INSERT INTO link_group (group_id, link_id)
SELECT g.id, r.id
FROM var
   JOIN json_each( var.value ) j ON var.key='json'
   JOIN link r ON r.original_url=json_extract(j.value,'$.link_url')
   JOIN `group` g ON true
      AND g.name=json_extract(j.value,'$.group_name')
      AND g.type=json_extract(j.value,'$.group_type')
WHERE var.id = ?
ON CONFLICT (group_id, link_id)
   DO UPDATE
      SET latest = strftime('%s','now')
;

-- name: UpsertLinkMetaFromVarJSON :exec
INSERT INTO meta (link_id,key,value)
SELECT
   CAST(json_extract(r.value,'$.link_id') AS INTEGER),
   CAST(json_extract(r.value,'$.key') AS TEXT),
   CAST(json_extract(r.value,'$.value') AS TEXT)
FROM var
   JOIN json_each( var.value ) r ON var.key='json'
WHERE var.id = ?
ON CONFLICT (link_id,key)
   DO UPDATE
   SET modified = strftime('%s','now')
;

-- name: ListLinks :many
SELECT *
FROM link
WHERE true
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY original_url
LIMIT 100;

-- name: GetLinkURLs :many
SELECT CAST(ifnull(url,'<invalid>') AS TEXT) AS url
FROM link
WHERE id IN (sqlc.slice('link_ids'))
;

-- name: ArchiveLinks :exec
UPDATE link
SET archived=1
WHERE id IN (sqlc.slice('link_ids'))
;

-- name: MarkLinksDeleted :exec
UPDATE link
SET deleted=1
WHERE id IN (sqlc.slice('link_ids'))
;

-- name: MarkLinksDeletedByGroupIds :exec
UPDATE link
SET deleted=1
WHERE true
   AND id NOT IN (SELECT lg.link_id FROM link_group lg WHERE lg.group_id=?)
   AND id IN (SELECT lg.link_id FROM link_group lg WHERE lg.group_id IN (sqlc.slice('group_ids')))
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
WHERE true
   AND id IN (sqlc.slice('ids'))
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY
   original_url;

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
   created=max(created);

-- name: InsertContent :exec
INSERT INTO content
   (link_id,head,body)
VALUES
   (?,?,?)
;
-- name: ListLatestUnparsedLinkURLs :many
SELECT
   id,
   original_url
FROM link
WHERE true
   AND parsed = 0
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY
   id DESC
LIMIT 8; -- LIMIT was chosen as slice len == slice cap for 8

-- name: UpdateParsedLink :exec
UPDATE link
SET
   title = ?,
   scheme = ?,
   subdomain = ?,
   sld = ?,
   tld = ?,
   port = ?,
   path = ?,
   query = ?,
   fragment = ?,
   parsed = 1
WHERE
   original_url = ?;

-- name: ListLinkIdsByGroup :many
SELECT CAST(l.id AS INTEGER) AS link_id
FROM
   link l
   JOIN link_group rg ON l.id=rg.link_id
   JOIN `group` g ON g.id=rg.group_id
WHERE true
   AND g.slug IN (sqlc.slice('slugs'))
   AND l.archived IN (sqlc.slice('links_archived'))
   AND l.deleted IN (sqlc.slice('links_deleted'))
;

-- name: ListLinkIdsByMeta :many
SELECT CAST(m.link_id AS INTEGER) AS link_id
FROM meta m
   JOIN link l ON l.id=m.link_id
WHERE true
   AND m.kv_pair IN (sqlc.slice('kv_pairs'))
   AND m.key IN (sqlc.slice('keys'))
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'));

-- name: ListLinkIdsByGroupType :many
SELECT CAST(link_id AS INTEGER) AS link_id
FROM link_group lg
        JOIN `group` g ON lg.group_id = g.id
        JOIN link l ON l.id=lg.link_id
WHERE true
   AND g.type IN (sqlc.slice('groupTypes'))
   AND l.archived IN (sqlc.slice('links_archived'))
   AND l.deleted IN (sqlc.slice('links_deleted'));

-- name: ListLinkIdsNotInGroupType :many
SELECT CAST(l.id AS INTEGER) AS link_id
FROM link l
WHERE TRUE
   AND l.archived IN (sqlc.slice('links_archived'))
   AND l.deleted IN (sqlc.slice('links_deleted'))
   AND l.id NOT IN (
      SELECT lg.link_id
      FROM link_group lg
        JOIN `group` g ON lg.group_id = g.id
      WHERE g.type IN (sqlc.slice('groupTypes'))
   );

-- name: UpsertVar :one
INSERT INTO var (key,value) VALUES (?,?)
ON CONFLICT (key) DO UPDATE SET value = excluded.value
RETURNING id;

-- name: DeleteVar :exec
DELETE FROM var WHERE id = ?;

-- name: ListLinkMeta :many
SELECT m.*
FROM meta m
   JOIN link l ON m.link_id = l.id
WHERE true
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY link_id,key DESC
;

-- name: ListLinkMetaForLinkId :many
SELECT m.key,m.value
FROM meta m
   JOIN link l ON m.link_id = l.id
WHERE link_id = ?
;

-- name: UpsertMetaFromVarJSON :exec
INSERT INTO meta (link_id, key, value)
SELECT
   l.id,
   json_extract(kv.value,'$.key'),
   json_extract(kv.value,'$.value')
FROM var
   JOIN json_each( var.value ) kv ON var.key='json'
   JOIN link l ON l.original_url=json_extract(kv.value,'$.url')
WHERE var.id = ?
   ON CONFLICT (link_id,key)
   DO UPDATE
      SET value = excluded.value;

