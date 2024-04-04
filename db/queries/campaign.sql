-- name: CreateCampaign :exec
insert into campaigns (name) values (?);

-- name: GetCampaignsWithoutSources :many
select c.name, c.id from campaigns c
left join campaigns_sources cs on c.id = cs.campaign_id
where cs.source_id is null;