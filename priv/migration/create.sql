create table if not exists user(
    id text primary key,
    username text unique not null,
    key blob unique not null
);

create table if not exists message(
    id text primary key,
    type integer default 0,
    user text not null,
    conversation text not null,
    datetime text not null,
    foreign key(user) references user(id) on delete cascade,
    foreign key(conversation) references conversation(id) on delete cascade
);

create table if not exists conversation(
    id text primary key,
    users text not null
);
