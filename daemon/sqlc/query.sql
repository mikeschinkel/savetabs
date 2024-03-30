-- name: LoadGroup :one
SELECT * FROM `group` WHERE id = ? LIMIT 1;

-- name: ListGroup :one
SELECT * FROM `group` ORDER BY latest DESC;

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

-- name: ListGroupType :one
SELECT * FROM group_type ORDER BY name DESC;

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
