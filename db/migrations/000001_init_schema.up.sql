create table campaigns (
  id bigint unsigned auto_increment primary key,
  name varchar(255) not null
);

create table sources (
    id bigint unsigned auto_increment primary key, 
    name varchar(255) not null
);

create table campaigns_sources (
  campaign_id bigint unsigned not null,
  source_id bigint unsigned not null,
  foreign key (campaign_id) references campaigns(id) on delete cascade,
  foreign key (source_id) references sources(id) on delete cascade
);

-- ensure that a source can be used in a maximum of 10 campaigns
create trigger trg_chk_campaign_limit
before insert on campaigns_sources
for each row
begin
  set @cnt = (select count(*) from campaigns_sources where source_id = new.source_id);
    if @cnt >= 10 then
      signal sqlstate "45000" set message_text = "Campaigns limit reached";
    end if;
 end;