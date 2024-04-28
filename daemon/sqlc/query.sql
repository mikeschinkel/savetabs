-- noinspection SqlResolveForFile @ any/"sqlc"

-- name: LoadGroup :one
SELECT * FROM `group`
WHERE true
   AND id = ?
   AND archived IN (sqlc.slice('groups_archived'))
   AND deleted IN (sqlc.slice('groups_deleted'))
LIMIT 1;

-- name: LoadGroupType :one
SELECT * FROM group_type WHERE type = ? LIMIT 1;

-- name: ListGroupsByType :many
SELECT *
FROM `group`
WHERE true
   AND type = ?
   AND archived IN (sqlc.slice('groups_archived'))
   AND deleted IN (sqlc.slice('groups_deleted'))
ORDER BY name;

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

SELECT * FROM group_type WHERE type = ? LIMIT 1;

-- name: ListGroupsType :many
SELECT DISTINCT
   gt.type,
   gt.name,
   gt.plural,
   COUNT(DISTINCT g.id) AS group_count,
   COUNT(DISTINCT g.archived=1) AS groups_archived,
   COUNT(DISTINCT g.deleted=1) AS groups_deleted,
   CAST(CASE WHEN g.ID IS NULL THEN 0
      ELSE COUNT(DISTINCT rg.link_id) END AS INTEGER) AS link_count,
   COUNT(DISTINCT l.archived=1) AS links_archived,
   COUNT(DISTINCT l.deleted=1) AS links_deleted,
   gt.sort
FROM group_type gt
   LEFT JOIN `group` g ON gt.type=g.type
   LEFT JOIN link_group rg ON g.id=rg.group_id
   LEFT JOIN link l ON l.id=rg.link_id
GROUP BY
   gt.sort,
   gt.type,
   gt.name
ORDER BY
   gt.sort;

-- name: LoadLink :one
SELECT *
FROM link
WHERE true
   AND id = ?
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
LIMIT 1
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
SELECT CAST(ifnull(url,'<invalid>') AS TEXT) FROM link
WHERE id IN (sqlc.slice('link_ids'))
;

-- name: ArchiveLinks :exec
UPDATE link
SET archived=1
WHERE id IN (sqlc.slice('link_ids'))
;

-- name: DeleteLinks :exec
UPDATE link
SET deleted=1
WHERE id IN (sqlc.slice('link_ids'))
;

-- name: ListFilteredLinks :many
SELECT
   l.id,
   l.original_url,
   l.created_time,
   l.visited_time,
   c.id AS content_id,
   c.title,
   c.body
FROM
   link l
    LEFT JOIN content c ON c.link_id=l.id
WHERE true
   AND l.id IN (sqlc.slice('ids'))
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
GROUP BY
   l.id,
   l.original_url,
   l.created_time,
   l.visited_time,
   c.id,
   c.title,
   c.body
HAVING 1=0
   OR c.created=max(c.created)
   OR count(c.id)=0
ORDER BY
   l.original_url;

-- name: ListLatestUnparsedLinkURLs :many
SELECT
   id,
   original_url
FROM link
WHERE true
   AND sld == ''
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY
   id DESC
LIMIT 8; -- LIMIT was chosen as slice len == slice cap for 8

-- name: UpdateLinkParts :exec
UPDATE link
SET
   scheme = ?,
   subdomain = ?,
   sld = ?,
   tld = ?,
   port = ?,
   path = ?,
   query = ?,
   fragment = ?
WHERE
   original_url = ?;


-- name: ListLinkIdsByGroupSlugs :many
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

-- name: ListLinkIdsByMetadata :many
SELECT CAST(m.link_id AS INTEGER) AS link_id
FROM metadata m
   JOIN link l ON l.id=m.link_id
WHERE true
   AND m.kv_pair IN (sqlc.slice('kv_pairs'))
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

-- name: UpsertLinksFromVarJSON :exec
INSERT INTO link (original_url)
SELECT r.value AS url
FROM var
   JOIN json_each( var.value ) r ON var.key='json'
WHERE var.id = ?
ON CONFLICT (original_url)
   DO UPDATE
            SET visited = strftime('%s','now');

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
            SET latest = strftime('%s','now');

-- name: UpsertVar :execlastid
INSERT INTO var (key,value) VALUES (?,?)
ON CONFLICT (key) DO UPDATE SET value = excluded.value;

-- name: DeleteVar :exec
DELETE FROM var WHERE id = ?;

-- name: ListMetadata :many
SELECT *
FROM metadata m
   JOIN link l ON m.link_id = l.id
WHERE true
   AND archived IN (sqlc.slice('links_archived'))
   AND deleted IN (sqlc.slice('links_deleted'))
ORDER BY link_id,key DESC;

-- name: UpsertMetadataFromVarJSON :exec
INSERT INTO metadata (link_id, key, value)
SELECT
   r.id,
   json_extract(kv.value,'$.key'),
   json_extract(kv.value,'$.value')
FROM var
   JOIN json_each( var.value ) kv ON var.key='json'
   JOIN link r ON r.original_url=json_extract(kv.value,'$.url')
WHERE var.id = ?
   ON CONFLICT (link_id,key)
   DO UPDATE
      SET value = excluded.value;

