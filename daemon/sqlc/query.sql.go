// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
	"strings"
)

const archiveLinks = `-- name: ArchiveLinks :exec
;

UPDATE link
SET archived=1
WHERE id IN (/*SLICE:link_ids*/?)
`

func (q *Queries) ArchiveLinks(ctx context.Context, linkIds []int64) error {
	query := archiveLinks
	var queryParams []interface{}
	if len(linkIds) > 0 {
		for _, v := range linkIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:link_ids*/?", strings.Repeat(",?", len(linkIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:link_ids*/?", "NULL", 1)
	}
	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}

const deleteLinks = `-- name: DeleteLinks :exec
;

UPDATE link
SET deleted=1
WHERE id IN (/*SLICE:link_ids*/?)
`

func (q *Queries) DeleteLinks(ctx context.Context, linkIds []int64) error {
	query := deleteLinks
	var queryParams []interface{}
	if len(linkIds) > 0 {
		for _, v := range linkIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:link_ids*/?", strings.Repeat(",?", len(linkIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:link_ids*/?", "NULL", 1)
	}
	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}

const deleteVar = `-- name: DeleteVar :exec
DELETE FROM var WHERE id = ?
`

func (q *Queries) DeleteVar(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteVar, id)
	return err
}

const getLinkURLs = `-- name: GetLinkURLs :many
SELECT CAST(ifnull(url,'<invalid>') AS TEXT) AS url
FROM link
WHERE id IN (/*SLICE:link_ids*/?)
`

func (q *Queries) GetLinkURLs(ctx context.Context, linkIds []int64) ([]string, error) {
	query := getLinkURLs
	var queryParams []interface{}
	if len(linkIds) > 0 {
		for _, v := range linkIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:link_ids*/?", strings.Repeat(",?", len(linkIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:link_ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		items = append(items, url)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertContent = `-- name: InsertContent :exec
INSERT INTO content
   (link_id,head,body)
VALUES
   (?,?,?)
`

type InsertContentParams struct {
	LinkID int64  `json:"link_id"`
	Head   string `json:"head"`
	Body   string `json:"body"`
}

func (q *Queries) InsertContent(ctx context.Context, arg InsertContentParams) error {
	_, err := q.db.ExecContext(ctx, insertContent, arg.LinkID, arg.Head, arg.Body)
	return err
}

const listFilteredLinks = `-- name: ListFilteredLinks :many
;

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
   AND id IN (/*SLICE:ids*/?)
   AND archived IN (/*SLICE:links_archived*/?)
   AND deleted IN (/*SLICE:links_deleted*/?)
ORDER BY
   original_url
`

type ListFilteredLinksParams struct {
	Ids           []int64 `json:"ids"`
	LinksArchived []int64 `json:"links_archived"`
	LinksDeleted  []int64 `json:"links_deleted"`
}

type ListFilteredLinksRow struct {
	ID          int64          `json:"id"`
	OriginalUrl string         `json:"original_url"`
	CreatedTime sql.NullString `json:"created_time"`
	VisitedTime sql.NullString `json:"visited_time"`
	Title       string         `json:"title"`
	Scheme      string         `json:"scheme"`
	Subdomain   string         `json:"subdomain"`
	Sld         string         `json:"sld"`
	Tld         string         `json:"tld"`
	Path        string         `json:"path"`
	Query       string         `json:"query"`
	Fragment    string         `json:"fragment"`
	Port        string         `json:"port"`
	Archived    int64          `json:"archived"`
	Deleted     int64          `json:"deleted"`
}

func (q *Queries) ListFilteredLinks(ctx context.Context, arg ListFilteredLinksParams) ([]ListFilteredLinksRow, error) {
	query := listFilteredLinks
	var queryParams []interface{}
	if len(arg.Ids) > 0 {
		for _, v := range arg.Ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(arg.Ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	if len(arg.LinksArchived) > 0 {
		for _, v := range arg.LinksArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_archived*/?", strings.Repeat(",?", len(arg.LinksArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_archived*/?", "NULL", 1)
	}
	if len(arg.LinksDeleted) > 0 {
		for _, v := range arg.LinksDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", strings.Repeat(",?", len(arg.LinksDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListFilteredLinksRow
	for rows.Next() {
		var i ListFilteredLinksRow
		if err := rows.Scan(
			&i.ID,
			&i.OriginalUrl,
			&i.CreatedTime,
			&i.VisitedTime,
			&i.Title,
			&i.Scheme,
			&i.Subdomain,
			&i.Sld,
			&i.Tld,
			&i.Path,
			&i.Query,
			&i.Fragment,
			&i.Port,
			&i.Archived,
			&i.Deleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listGroupsByType = `-- name: ListGroupsByType :many
SELECT
   g.id, g.name, g.type, g.slug, g.created_time, g.latest_time, g.created, g.latest, g.archived, g.deleted,
   CAST(gt.name AS TEXT) AS type_name
FROM ` + "`" + `group` + "`" + ` g
   JOIN group_type gt ON gt.type = g.type
WHERE true
   AND g.type = ?
   AND g.archived IN (/*SLICE:groups_archived*/?)
   AND g.deleted IN (/*SLICE:groups_deleted*/?)
ORDER BY g.name
`

type ListGroupsByTypeParams struct {
	Type           string  `json:"type"`
	GroupsArchived []int64 `json:"groups_archived"`
	GroupsDeleted  []int64 `json:"groups_deleted"`
}

type ListGroupsByTypeRow struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Slug        string         `json:"slug"`
	CreatedTime sql.NullString `json:"created_time"`
	LatestTime  sql.NullString `json:"latest_time"`
	Created     sql.NullInt64  `json:"-"`
	Latest      sql.NullInt64  `json:"latest"`
	Archived    int64          `json:"archived"`
	Deleted     int64          `json:"deleted"`
	TypeName    string         `json:"type_name"`
}

func (q *Queries) ListGroupsByType(ctx context.Context, arg ListGroupsByTypeParams) ([]ListGroupsByTypeRow, error) {
	query := listGroupsByType
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Type)
	if len(arg.GroupsArchived) > 0 {
		for _, v := range arg.GroupsArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:groups_archived*/?", strings.Repeat(",?", len(arg.GroupsArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:groups_archived*/?", "NULL", 1)
	}
	if len(arg.GroupsDeleted) > 0 {
		for _, v := range arg.GroupsDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:groups_deleted*/?", strings.Repeat(",?", len(arg.GroupsDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:groups_deleted*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListGroupsByTypeRow
	for rows.Next() {
		var i ListGroupsByTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Slug,
			&i.CreatedTime,
			&i.LatestTime,
			&i.Created,
			&i.Latest,
			&i.Archived,
			&i.Deleted,
			&i.TypeName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listGroupsType = `-- name: ListGroupsType :many
SELECT type, name, plural, group_count, groups_archived, groups_deleted, link_count, links_archived, links_deleted, sort FROM groups_type
`

func (q *Queries) ListGroupsType(ctx context.Context) ([]GroupsType, error) {
	rows, err := q.db.QueryContext(ctx, listGroupsType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupsType
	for rows.Next() {
		var i GroupsType
		if err := rows.Scan(
			&i.Type,
			&i.Name,
			&i.Plural,
			&i.GroupCount,
			&i.GroupsArchived,
			&i.GroupsDeleted,
			&i.LinkCount,
			&i.LinksArchived,
			&i.LinksDeleted,
			&i.Sort,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLatestUnparsedLinkURLs = `-- name: ListLatestUnparsedLinkURLs :many
;
SELECT
   id,
   original_url
FROM link
WHERE true
   AND parsed = 0
   AND archived IN (/*SLICE:links_archived*/?)
   AND deleted IN (/*SLICE:links_deleted*/?)
ORDER BY
   id DESC
LIMIT 8
`

type ListLatestUnparsedLinkURLsParams struct {
	LinksArchived []int64 `json:"links_archived"`
	LinksDeleted  []int64 `json:"links_deleted"`
}

type ListLatestUnparsedLinkURLsRow struct {
	ID          int64  `json:"id"`
	OriginalUrl string `json:"original_url"`
}

func (q *Queries) ListLatestUnparsedLinkURLs(ctx context.Context, arg ListLatestUnparsedLinkURLsParams) ([]ListLatestUnparsedLinkURLsRow, error) {
	query := listLatestUnparsedLinkURLs
	var queryParams []interface{}
	if len(arg.LinksArchived) > 0 {
		for _, v := range arg.LinksArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_archived*/?", strings.Repeat(",?", len(arg.LinksArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_archived*/?", "NULL", 1)
	}
	if len(arg.LinksDeleted) > 0 {
		for _, v := range arg.LinksDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", strings.Repeat(",?", len(arg.LinksDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLatestUnparsedLinkURLsRow
	for rows.Next() {
		var i ListLatestUnparsedLinkURLsRow
		if err := rows.Scan(&i.ID, &i.OriginalUrl); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLinkIdsByGroup = `-- name: ListLinkIdsByGroup :many
SELECT CAST(l.id AS INTEGER) AS link_id
FROM
   link l
   JOIN link_group rg ON l.id=rg.link_id
   JOIN ` + "`" + `group` + "`" + ` g ON g.id=rg.group_id
WHERE true
   AND g.slug IN (/*SLICE:slugs*/?)
   AND l.archived IN (/*SLICE:links_archived*/?)
   AND l.deleted IN (/*SLICE:links_deleted*/?)
`

type ListLinkIdsByGroupParams struct {
	Slugs         []string `json:"slugs"`
	LinksArchived []int64  `json:"links_archived"`
	LinksDeleted  []int64  `json:"links_deleted"`
}

func (q *Queries) ListLinkIdsByGroup(ctx context.Context, arg ListLinkIdsByGroupParams) ([]int64, error) {
	query := listLinkIdsByGroup
	var queryParams []interface{}
	if len(arg.Slugs) > 0 {
		for _, v := range arg.Slugs {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:slugs*/?", strings.Repeat(",?", len(arg.Slugs))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:slugs*/?", "NULL", 1)
	}
	if len(arg.LinksArchived) > 0 {
		for _, v := range arg.LinksArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_archived*/?", strings.Repeat(",?", len(arg.LinksArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_archived*/?", "NULL", 1)
	}
	if len(arg.LinksDeleted) > 0 {
		for _, v := range arg.LinksDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", strings.Repeat(",?", len(arg.LinksDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var link_id int64
		if err := rows.Scan(&link_id); err != nil {
			return nil, err
		}
		items = append(items, link_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLinkIdsByGroupType = `-- name: ListLinkIdsByGroupType :many
SELECT CAST(link_id AS INTEGER) AS link_id
FROM link_group lg
        JOIN ` + "`" + `group` + "`" + ` g ON lg.group_id = g.id
        JOIN link l ON l.id=lg.link_id
WHERE true
   AND g.type IN (/*SLICE:groupTypes*/?)
   AND l.archived IN (/*SLICE:links_archived*/?)
   AND l.deleted IN (/*SLICE:links_deleted*/?)
`

type ListLinkIdsByGroupTypeParams struct {
	GroupTypes    []string `json:"groupTypes"`
	LinksArchived []int64  `json:"links_archived"`
	LinksDeleted  []int64  `json:"links_deleted"`
}

func (q *Queries) ListLinkIdsByGroupType(ctx context.Context, arg ListLinkIdsByGroupTypeParams) ([]int64, error) {
	query := listLinkIdsByGroupType
	var queryParams []interface{}
	if len(arg.GroupTypes) > 0 {
		for _, v := range arg.GroupTypes {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:groupTypes*/?", strings.Repeat(",?", len(arg.GroupTypes))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:groupTypes*/?", "NULL", 1)
	}
	if len(arg.LinksArchived) > 0 {
		for _, v := range arg.LinksArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_archived*/?", strings.Repeat(",?", len(arg.LinksArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_archived*/?", "NULL", 1)
	}
	if len(arg.LinksDeleted) > 0 {
		for _, v := range arg.LinksDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", strings.Repeat(",?", len(arg.LinksDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var link_id int64
		if err := rows.Scan(&link_id); err != nil {
			return nil, err
		}
		items = append(items, link_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLinkIdsByMeta = `-- name: ListLinkIdsByMeta :many
;

SELECT CAST(m.link_id AS INTEGER) AS link_id
FROM meta m
   JOIN link l ON l.id=m.link_id
WHERE true
   AND m.kv_pair IN (/*SLICE:kv_pairs*/?)
   AND m.key IN (/*SLICE:keys*/?)
   AND archived IN (/*SLICE:links_archived*/?)
   AND deleted IN (/*SLICE:links_deleted*/?)
`

type ListLinkIdsByMetaParams struct {
	KvPairs       []string `json:"kv_pairs"`
	Keys          []string `json:"keys"`
	LinksArchived []int64  `json:"links_archived"`
	LinksDeleted  []int64  `json:"links_deleted"`
}

func (q *Queries) ListLinkIdsByMeta(ctx context.Context, arg ListLinkIdsByMetaParams) ([]int64, error) {
	query := listLinkIdsByMeta
	var queryParams []interface{}
	if len(arg.KvPairs) > 0 {
		for _, v := range arg.KvPairs {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:kv_pairs*/?", strings.Repeat(",?", len(arg.KvPairs))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:kv_pairs*/?", "NULL", 1)
	}
	if len(arg.Keys) > 0 {
		for _, v := range arg.Keys {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:keys*/?", strings.Repeat(",?", len(arg.Keys))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:keys*/?", "NULL", 1)
	}
	if len(arg.LinksArchived) > 0 {
		for _, v := range arg.LinksArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_archived*/?", strings.Repeat(",?", len(arg.LinksArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_archived*/?", "NULL", 1)
	}
	if len(arg.LinksDeleted) > 0 {
		for _, v := range arg.LinksDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", strings.Repeat(",?", len(arg.LinksDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var link_id int64
		if err := rows.Scan(&link_id); err != nil {
			return nil, err
		}
		items = append(items, link_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLinkIdsNotInGroupType = `-- name: ListLinkIdsNotInGroupType :many
SELECT CAST(l.id AS INTEGER) AS link_id
FROM link l
WHERE TRUE
   AND l.archived IN (/*SLICE:links_archived*/?)
   AND l.deleted IN (/*SLICE:links_deleted*/?)
   AND l.id NOT IN (
      SELECT lg.link_id
      FROM link_group lg
        JOIN ` + "`" + `group` + "`" + ` g ON lg.group_id = g.id
      WHERE g.type IN (/*SLICE:groupTypes*/?)
   )
`

type ListLinkIdsNotInGroupTypeParams struct {
	LinksArchived []int64  `json:"links_archived"`
	LinksDeleted  []int64  `json:"links_deleted"`
	GroupTypes    []string `json:"groupTypes"`
}

func (q *Queries) ListLinkIdsNotInGroupType(ctx context.Context, arg ListLinkIdsNotInGroupTypeParams) ([]int64, error) {
	query := listLinkIdsNotInGroupType
	var queryParams []interface{}
	if len(arg.LinksArchived) > 0 {
		for _, v := range arg.LinksArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_archived*/?", strings.Repeat(",?", len(arg.LinksArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_archived*/?", "NULL", 1)
	}
	if len(arg.LinksDeleted) > 0 {
		for _, v := range arg.LinksDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", strings.Repeat(",?", len(arg.LinksDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", "NULL", 1)
	}
	if len(arg.GroupTypes) > 0 {
		for _, v := range arg.GroupTypes {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:groupTypes*/?", strings.Repeat(",?", len(arg.GroupTypes))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:groupTypes*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var link_id int64
		if err := rows.Scan(&link_id); err != nil {
			return nil, err
		}
		items = append(items, link_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLinkMeta = `-- name: ListLinkMeta :many
SELECT m.id, m.link_id, m."key", m.value, m.kv_pair, m.created_time, m.modified_time, m.created, m.modified
FROM meta m
   JOIN link l ON m.link_id = l.id
WHERE true
   AND archived IN (/*SLICE:links_archived*/?)
   AND deleted IN (/*SLICE:links_deleted*/?)
ORDER BY link_id,key DESC
`

type ListLinkMetaParams struct {
	LinksArchived []int64 `json:"links_archived"`
	LinksDeleted  []int64 `json:"links_deleted"`
}

func (q *Queries) ListLinkMeta(ctx context.Context, arg ListLinkMetaParams) ([]Meta, error) {
	query := listLinkMeta
	var queryParams []interface{}
	if len(arg.LinksArchived) > 0 {
		for _, v := range arg.LinksArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_archived*/?", strings.Repeat(",?", len(arg.LinksArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_archived*/?", "NULL", 1)
	}
	if len(arg.LinksDeleted) > 0 {
		for _, v := range arg.LinksDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", strings.Repeat(",?", len(arg.LinksDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Meta
	for rows.Next() {
		var i Meta
		if err := rows.Scan(
			&i.ID,
			&i.LinkID,
			&i.Key,
			&i.Value,
			&i.KvPair,
			&i.CreatedTime,
			&i.ModifiedTime,
			&i.Created,
			&i.Modified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLinkMetaForLinkId = `-- name: ListLinkMetaForLinkId :many
;

SELECT m.key,m.value
FROM meta m
   JOIN link l ON m.link_id = l.id
WHERE link_id = ?
`

type ListLinkMetaForLinkIdRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (q *Queries) ListLinkMetaForLinkId(ctx context.Context, linkID int64) ([]ListLinkMetaForLinkIdRow, error) {
	rows, err := q.db.QueryContext(ctx, listLinkMetaForLinkId, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLinkMetaForLinkIdRow
	for rows.Next() {
		var i ListLinkMetaForLinkIdRow
		if err := rows.Scan(&i.Key, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLinks = `-- name: ListLinks :many
;

SELECT id, title, scheme, subdomain, sld, tld, port, path, "query", fragment, original_url, host, url, created_time, visited_time, created, visited, archived, deleted, parsed
FROM link
WHERE true
   AND archived IN (/*SLICE:links_archived*/?)
   AND deleted IN (/*SLICE:links_deleted*/?)
ORDER BY original_url
LIMIT 100
`

type ListLinksParams struct {
	LinksArchived []int64 `json:"links_archived"`
	LinksDeleted  []int64 `json:"links_deleted"`
}

func (q *Queries) ListLinks(ctx context.Context, arg ListLinksParams) ([]Link, error) {
	query := listLinks
	var queryParams []interface{}
	if len(arg.LinksArchived) > 0 {
		for _, v := range arg.LinksArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_archived*/?", strings.Repeat(",?", len(arg.LinksArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_archived*/?", "NULL", 1)
	}
	if len(arg.LinksDeleted) > 0 {
		for _, v := range arg.LinksDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", strings.Repeat(",?", len(arg.LinksDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:links_deleted*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Link
	for rows.Next() {
		var i Link
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Scheme,
			&i.Subdomain,
			&i.Sld,
			&i.Tld,
			&i.Port,
			&i.Path,
			&i.Query,
			&i.Fragment,
			&i.OriginalUrl,
			&i.Host,
			&i.Url,
			&i.CreatedTime,
			&i.VisitedTime,
			&i.Created,
			&i.Visited,
			&i.Archived,
			&i.Deleted,
			&i.Parsed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const loadGroup = `-- name: LoadGroup :one
SELECT id, name, type, slug, created_time, latest_time, created, latest, archived, deleted FROM ` + "`" + `group` + "`" + `
WHERE true
   AND id = ?
   AND archived IN (/*SLICE:groups_archived*/?)
   AND deleted IN (/*SLICE:groups_deleted*/?)
LIMIT 1
`

type LoadGroupParams struct {
	ID             int64   `json:"id"`
	GroupsArchived []int64 `json:"groups_archived"`
	GroupsDeleted  []int64 `json:"groups_deleted"`
}

func (q *Queries) LoadGroup(ctx context.Context, arg LoadGroupParams) (Group, error) {
	query := loadGroup
	var queryParams []interface{}
	queryParams = append(queryParams, arg.ID)
	if len(arg.GroupsArchived) > 0 {
		for _, v := range arg.GroupsArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:groups_archived*/?", strings.Repeat(",?", len(arg.GroupsArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:groups_archived*/?", "NULL", 1)
	}
	if len(arg.GroupsDeleted) > 0 {
		for _, v := range arg.GroupsDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:groups_deleted*/?", strings.Repeat(",?", len(arg.GroupsDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:groups_deleted*/?", "NULL", 1)
	}
	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Slug,
		&i.CreatedTime,
		&i.LatestTime,
		&i.Created,
		&i.Latest,
		&i.Archived,
		&i.Deleted,
	)
	return i, err
}

const loadGroupName = `-- name: LoadGroupName :one
SELECT name FROM ` + "`" + `group` + "`" + ` WHERE id = ?
`

func (q *Queries) LoadGroupName(ctx context.Context, id int64) (string, error) {
	row := q.db.QueryRowContext(ctx, loadGroupName, id)
	var name string
	err := row.Scan(&name)
	return name, err
}

const loadGroupType = `-- name: LoadGroupType :one
SELECT type, sort, name, plural, description FROM group_type WHERE type = ? LIMIT 1
`

func (q *Queries) LoadGroupType(ctx context.Context, type_ string) (GroupType, error) {
	row := q.db.QueryRowContext(ctx, loadGroupType, type_)
	var i GroupType
	err := row.Scan(
		&i.Type,
		&i.Sort,
		&i.Name,
		&i.Plural,
		&i.Description,
	)
	return i, err
}

const loadGroupTypeWithStats = `-- name: LoadGroupTypeWithStats :one
SELECT type, name, plural, group_count, groups_archived, groups_deleted, link_count, links_archived, links_deleted, sort FROM groups_type WHERE type = ? LIMIT 1
`

func (q *Queries) LoadGroupTypeWithStats(ctx context.Context, type_ string) (GroupsType, error) {
	row := q.db.QueryRowContext(ctx, loadGroupTypeWithStats, type_)
	var i GroupsType
	err := row.Scan(
		&i.Type,
		&i.Name,
		&i.Plural,
		&i.GroupCount,
		&i.GroupsArchived,
		&i.GroupsDeleted,
		&i.LinkCount,
		&i.LinksArchived,
		&i.LinksDeleted,
		&i.Sort,
	)
	return i, err
}

const loadGroupsBySlug = `-- name: LoadGroupsBySlug :one
SELECT id, name, type, slug, created_time, latest_time, created, latest, archived, deleted FROM ` + "`" + `group` + "`" + `
WHERE true
   AND slug = ?
   AND archived IN (/*SLICE:groups_archived*/?)
   AND deleted IN (/*SLICE:groups_deleted*/?)
LIMIT 1
`

type LoadGroupsBySlugParams struct {
	Slug           string  `json:"slug"`
	GroupsArchived []int64 `json:"groups_archived"`
	GroupsDeleted  []int64 `json:"groups_deleted"`
}

func (q *Queries) LoadGroupsBySlug(ctx context.Context, arg LoadGroupsBySlugParams) (Group, error) {
	query := loadGroupsBySlug
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Slug)
	if len(arg.GroupsArchived) > 0 {
		for _, v := range arg.GroupsArchived {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:groups_archived*/?", strings.Repeat(",?", len(arg.GroupsArchived))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:groups_archived*/?", "NULL", 1)
	}
	if len(arg.GroupsDeleted) > 0 {
		for _, v := range arg.GroupsDeleted {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:groups_deleted*/?", strings.Repeat(",?", len(arg.GroupsDeleted))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:groups_deleted*/?", "NULL", 1)
	}
	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Slug,
		&i.CreatedTime,
		&i.LatestTime,
		&i.Created,
		&i.Latest,
		&i.Archived,
		&i.Deleted,
	)
	return i, err
}

const loadLatestContent = `-- name: LoadLatestContent :one
SELECT
   id, link_id, title, body, head, created_time, created
FROM
   content
WHERE
   link_id = ?
GROUP BY
   link_id,
   created
HAVING
   created=max(created)
`

// TODO: Untested, ensure query works
func (q *Queries) LoadLatestContent(ctx context.Context, linkID int64) (Content, error) {
	row := q.db.QueryRowContext(ctx, loadLatestContent, linkID)
	var i Content
	err := row.Scan(
		&i.ID,
		&i.LinkID,
		&i.Title,
		&i.Body,
		&i.Head,
		&i.CreatedTime,
		&i.Created,
	)
	return i, err
}

const loadLink = `-- name: LoadLink :one
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
`

type LoadLinkRow struct {
	ID          int64          `json:"id"`
	OriginalUrl string         `json:"original_url"`
	CreatedTime sql.NullString `json:"created_time"`
	VisitedTime sql.NullString `json:"visited_time"`
	Scheme      string         `json:"scheme"`
	Subdomain   string         `json:"subdomain"`
	Sld         string         `json:"sld"`
	Tld         string         `json:"tld"`
	Path        string         `json:"path"`
	Query       string         `json:"query"`
	Fragment    string         `json:"fragment"`
	Port        string         `json:"port"`
	Url         string         `json:"url"`
	Title       string         `json:"title"`
}

func (q *Queries) LoadLink(ctx context.Context, id int64) (LoadLinkRow, error) {
	row := q.db.QueryRowContext(ctx, loadLink, id)
	var i LoadLinkRow
	err := row.Scan(
		&i.ID,
		&i.OriginalUrl,
		&i.CreatedTime,
		&i.VisitedTime,
		&i.Scheme,
		&i.Subdomain,
		&i.Sld,
		&i.Tld,
		&i.Path,
		&i.Query,
		&i.Fragment,
		&i.Port,
		&i.Url,
		&i.Title,
	)
	return i, err
}

const loadLinkIdByUrl = `-- name: LoadLinkIdByUrl :one
SELECT id FROM link WHERE original_url = ? LIMIT 1
`

func (q *Queries) LoadLinkIdByUrl(ctx context.Context, originalUrl string) (int64, error) {
	row := q.db.QueryRowContext(ctx, loadLinkIdByUrl, originalUrl)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateParsedLink = `-- name: UpdateParsedLink :exec

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
   original_url = ?
`

type UpdateParsedLinkParams struct {
	Title       string `json:"title"`
	Scheme      string `json:"scheme"`
	Subdomain   string `json:"subdomain"`
	Sld         string `json:"sld"`
	Tld         string `json:"tld"`
	Port        string `json:"port"`
	Path        string `json:"path"`
	Query       string `json:"query"`
	Fragment    string `json:"fragment"`
	OriginalUrl string `json:"original_url"`
}

// LIMIT was chosen as slice len == slice cap for 8
func (q *Queries) UpdateParsedLink(ctx context.Context, arg UpdateParsedLinkParams) error {
	_, err := q.db.ExecContext(ctx, updateParsedLink,
		arg.Title,
		arg.Scheme,
		arg.Subdomain,
		arg.Sld,
		arg.Tld,
		arg.Port,
		arg.Path,
		arg.Query,
		arg.Fragment,
		arg.OriginalUrl,
	)
	return err
}

const upsertGroupsFromVarJSON = `-- name: UpsertGroupsFromVarJSON :exec
INSERT INTO ` + "`" + `group` + "`" + ` (name,type,slug)
SELECT
   json_extract(r.value,'$.name') AS name,
   json_extract(r.value,'$.type') AS type,
   json_extract(r.value,'$.slug') AS slug
FROM var
   JOIN json_each( var.value ) r ON var.key='json'
WHERE var.id = ?
    ON CONFLICT (name,type)
        DO UPDATE
            SET latest = strftime('%s','now')
`

func (q *Queries) UpsertGroupsFromVarJSON(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, upsertGroupsFromVarJSON, id)
	return err
}

const upsertLink = `-- name: UpsertLink :one
;

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
`

type UpsertLinkParams struct {
	OriginalUrl string `json:"original_url"`
	Title       string `json:"title"`
}

func (q *Queries) UpsertLink(ctx context.Context, arg UpsertLinkParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, upsertLink, arg.OriginalUrl, arg.Title)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const upsertLinkGroupsFromVarJSON = `-- name: UpsertLinkGroupsFromVarJSON :exec
;

INSERT INTO link_group (group_id, link_id)
SELECT g.id, r.id
FROM var
   JOIN json_each( var.value ) j ON var.key='json'
   JOIN link r ON r.original_url=json_extract(j.value,'$.link_url')
   JOIN ` + "`" + `group` + "`" + ` g ON true
      AND g.name=json_extract(j.value,'$.group_name')
      AND g.type=json_extract(j.value,'$.group_type')
WHERE var.id = ?
ON CONFLICT (group_id, link_id)
   DO UPDATE
      SET latest = strftime('%s','now')
`

func (q *Queries) UpsertLinkGroupsFromVarJSON(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, upsertLinkGroupsFromVarJSON, id)
	return err
}

const upsertLinkMetaFromVarJSON = `-- name: UpsertLinkMetaFromVarJSON :exec
;

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
`

func (q *Queries) UpsertLinkMetaFromVarJSON(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, upsertLinkMetaFromVarJSON, id)
	return err
}

const upsertLinksFromVarJSON = `-- name: UpsertLinksFromVarJSON :exec
;

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
`

func (q *Queries) UpsertLinksFromVarJSON(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, upsertLinksFromVarJSON, id)
	return err
}

const upsertMetaFromVarJSON = `-- name: UpsertMetaFromVarJSON :exec
;

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
      SET value = excluded.value
`

func (q *Queries) UpsertMetaFromVarJSON(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, upsertMetaFromVarJSON, id)
	return err
}

const upsertVar = `-- name: UpsertVar :one
INSERT INTO var (key,value) VALUES (?,?)
ON CONFLICT (key) DO UPDATE SET value = excluded.value
RETURNING id
`

type UpsertVarParams struct {
	Key   string         `json:"key"`
	Value sql.NullString `json:"value"`
}

func (q *Queries) UpsertVar(ctx context.Context, arg UpsertVarParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, upsertVar, arg.Key, arg.Value)
	var id int64
	err := row.Scan(&id)
	return id, err
}
