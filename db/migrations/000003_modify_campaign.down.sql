alter table campaigns
    drop column domain_list,
    drop column list_type,
    add column blacklist varchar(255) default "",
    add column whitelist varchar(255) default "";