-- name: LoadGroup :one
SELECT * FROM `group` WHERE id = ? LIMIT 1;

-- name: ListGroupsWithCounts :many
SELECT * FROM groups_with_counts;

-- name: ListGroupsWithCountsByGroupType :many
SELECT * FROM groups_with_counts WHERE type = ?;

-- name: UpsertGroupsFromVarJSON :exec
INSERT INTO `group` (name,type)
SELECT
   json_extract(r.value,'$.name') AS name,
   json_extract(r.value,'$.type') AS type
FROM var
   JOIN json_each( var.value ) r ON var.key='json'
WHERE var.id = ?
ON CONFLICT (name,type)
   DO UPDATE
   SET latest = strftime('%s','now');

-- name: LoadGroupType :one
SELECT * FROM group_type WHERE type = ? LIMIT 1;

-- name: ListGroupsType :many
SELECT DISTINCT
   gt.type,
   gt.name,
   gt.plural,
   COUNT(DISTINCT g.id) AS group_count,
   CAST(CASE WHEN g.ID IS NULL THEN 0
      ELSE COUNT(*) END AS INTEGER) AS resource_count,
   gt.sort
FROM group_type gt
   LEFT JOIN `group` g ON gt.type=g.type
   LEFT JOIN resource_group rg ON g.id=rg.group_id
GROUP BY
   gt.sort,
   gt.type,
   gt.name
ORDER BY
   gt.sort;

-- name: LoadResource :one
SELECT * FROM resource WHERE id = ? LIMIT 1;

-- name: ListResources :many
SELECT * FROM resource ORDER BY visited DESC;

-- name: UpsertResourcesFromVarJSON :exec
INSERT INTO resource (url)
SELECT r.value AS url
FROM var
   JOIN json_each( var.value ) r ON var.key='json'
WHERE var.id = ?
ON CONFLICT (url)
   DO UPDATE
            SET visited = strftime('%s','now');

-- name: UpsertResourceGroupsFromVarJSON :exec
INSERT INTO resource_group (group_id, resource_id)
SELECT g.id, r.id
FROM var
   JOIN json_each( var.value ) j ON var.key='json'
   JOIN resource r ON r.url=json_extract(j.value,'$.resource_url')
   JOIN `group` g ON 1=1
      AND g.name=json_extract(j.value,'$.group_name')
      AND g.type=json_extract(j.value,'$.group_type')
WHERE var.id = ?
ON CONFLICT (group_id, resource_id)
   DO UPDATE
            SET latest = strftime('%s','now');

-- name: UpsertVar :execlastid
INSERT INTO var (key,value) VALUES (?,?)
ON CONFLICT (key) DO UPDATE SET value = excluded.value;

-- name: DeleteVar :exec
DELETE FROM var WHERE id = ?;

-- name: ListKeyValues :many
SELECT * FROM key_value ORDER BY resource_id,key DESC;

-- name: UpsertKeyValuesFromVarJSON :exec
INSERT INTO key_value (resource_id, key, value)
SELECT
   r.id,
   json_extract(kv.value,'$.key'),
   json_extract(kv.value,'$.value')
FROM var
   JOIN json_each( var.value ) kv ON var.key='json'
   JOIN resource r ON r.url=json_extract(kv.value,'$.url')
WHERE var.id = ?
   ON CONFLICT (resource_id,key)
   DO UPDATE
      SET value = excluded.value;

