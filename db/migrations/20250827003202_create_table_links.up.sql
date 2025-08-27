create table links
(
    id            varchar(100)  not null,
    user_id       varchar(100)  not null,
    title         varchar(50)   not null,
    short_url     varchar(50)   UNIQUE,
    long_url      varchar(100)  UNIQUE,
    is_active     BOOLEAN,
    created_at    bigint        not null,
    updated_at    bigint        not null,
    primary key (id),
    foreign key fk_links_user_id (user_id) references users (id)
) engine = InnoDB;