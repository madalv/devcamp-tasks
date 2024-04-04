-- name: CreateSource :exec
insert into sources (name) values (?);

-- name: GetSourcesWithMostCampaigns :many
select s.name, count(cs.campaign_id) as nr_campaigns from campaigns_sources cs
join sources s on s.id = cs.source_id
group by cs.source_id
order by nr_campaigns desc
limit 5;