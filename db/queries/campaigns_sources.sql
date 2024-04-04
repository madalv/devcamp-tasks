-- name: CreateCampaignSourceLink :exec
insert into campaigns_sources(source_id, campaign_id) values (?, ?);