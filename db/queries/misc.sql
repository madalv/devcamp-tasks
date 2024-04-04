-- name: GetNamesOfCampaignsAndSources :many
select c.name from campaigns c
union
select s.name from sources s;