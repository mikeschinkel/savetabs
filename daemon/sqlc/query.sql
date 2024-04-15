-- noinspection SqlResolveForFile @ any/"sqlc"

-- name: LoadGroup :one
SELECT * FROM `group` WHERE id = ? LIMIT 1;

-- name: LoadGroupType :one
SELECT * FROM group_type WHERE type = ? LIMIT 1;

-- name: ListGroupsByType :many
SELECT * FROM `group` WHERE type = ? ORDER BY name;

-- name: LoadGroupsBySlug :one
SELECT * FROM `group` WHERE slug = ? LIMIT 1;

-- name: ListGroupsWithCounts :many
SELECT * FROM groups_with_counts_view;

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
   CAST(CASE WHEN g.ID IS NULL THEN 0
      ELSE COUNT(DISTINCT rg.link_id) END AS INTEGER) AS link_count,
   gt.sort
FROM group_type gt
   LEFT JOIN `group` g ON gt.type=g.type
   LEFT JOIN link_group rg ON g.id=rg.group_id
GROUP BY
   gt.sort,
   gt.type,
   gt.name
ORDER BY
   gt.sort;

-- name: LoadLink :one
SELECT * FROM link WHERE id = ? LIMIT 1;

-- name: ListLinks :many
SELECT * FROM link ORDER BY url LIMIT 100;

-- name: ListFilteredLinks :many
SELECT *
FROM link
WHERE id IN (sqlc.slice('ids'));

-- name: ListLinkIdsByGroupSlugs :many
SELECT CAST(rg.link_id AS INTEGER) AS link_id
FROM link_group rg
JOIN `group` g ON g.id=rg.group_id
WHERE g.slug IN (sqlc.slice('slugs'));

-- name: ListLinkIdsByKeyValues :many
SELECT CAST(link_id AS INTEGER) AS link_id
FROM key_value
WHERE kv_pair IN (sqlc.slice('pairs'));

-- name: ListLinksForGroup :many
SELECT DISTINCT
   id,
   link_id,
   url,
   group_id,
   cast(group_name AS VARCHAR(32)) AS group_name,
   cast(group_slug AS VARCHAR(32)) AS group_slug,
   cast(group_type AS VARCHAR(1))  AS group_type,
   cast(type_name AS VARCHAR(32))  AS type_name,
   domain,
   group_ids,
   group_types,
   group_names,
   quoted_group_types,
   quoted_group_slugs,
   quoted_group_names
FROM
   links_view
WHERE true
   AND group_type = ?
   AND group_slug = ?
ORDER BY
   url;



-- name: UpsertLinksFromVarJSON :exec
INSERT INTO link (url)
SELECT r.value AS url
FROM var
   JOIN json_each( var.value ) r ON var.key='json'
WHERE var.id = ?
ON CONFLICT (url)
   DO UPDATE
            SET visited = strftime('%s','now');

-- name: UpsertLinkGroupsFromVarJSON :exec
INSERT INTO link_group (group_id, link_id)
SELECT g.id, r.id
FROM var
   JOIN json_each( var.value ) j ON var.key='json'
   JOIN link r ON r.url=json_extract(j.value,'$.link_url')
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

-- name: ListKeyValues :many
SELECT * FROM key_value ORDER BY link_id,key DESC;

-- name: UpsertKeyValuesFromVarJSON :exec
INSERT INTO key_value (link_id, key, value)
SELECT
   r.id,
   json_extract(kv.value,'$.key'),
   json_extract(kv.value,'$.value')
FROM var
   JOIN json_each( var.value ) kv ON var.key='json'
   JOIN link r ON r.url=json_extract(kv.value,'$.url')
WHERE var.id = ?
   ON CONFLICT (link_id,key)
   DO UPDATE
      SET value = excluded.value;

