
create table campaigns (
  id bigint unsigned auto_increment primary key,
  name varchar(255) not null
)

create table sources (
    id bigint unsigned auto_increment primary key, 
    name varchar(255) not null
)

create table campaigns_sources (
  campaign_id bigint unsigned not null,
  source_id bigint unsigned not null,
  foreign key (campaign_id) references campaigns(id),
  foreign key (source_id) references sources(id)
)

-- ensure that a source can be used in a maximum of 10 campaigns
create trigger trg_chk_campaign_limit
before insert on campaigns_sources
for each row
begin
  set @cnt = (select count(*) from campaigns_sources where source_id = new.source_id);
  
  if @cnt >= 10 then
    signal sqlstate "45000" set message_text = "Campaigns limit reached";
  end if;
end

-- Queries

-- top 5 sources by number of campaigns
select s.name, count(cs.campaign_id) as nr_campaigns from campaigns_sources cs
join sources s on s.id = cs.source_id
group by cs.source_id
order by nr_campaigns desc
limit 5;

-- campaigns without sources
select c.name, c.id from campaigns c
left join campaigns_sources cs on c.id = cs.campaign_id
where cs.source_id is null;

-- all names of campaigns + sources
select c.name from campaigns c
union
select s.name from sources s;


