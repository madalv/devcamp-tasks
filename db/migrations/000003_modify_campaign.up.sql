alter table campaigns
    add column domain_list varchar(255) default "",
    add column list_type varchar(255) default "black",
    drop column blacklist,
    drop column whitelist;
