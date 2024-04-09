alter table campaigns
add column blacklist varchar(255) default "",
add column whitelist varchar(255) default "";