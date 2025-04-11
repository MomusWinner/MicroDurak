create extension if not exists "uuid-ossp";

create table user(
    id uuid primary key default uuid_generate_v4(),
    name varchar(300),
    password varchar(300),
    email varchar(300),
    age integer,
    rating integer
);

create table match_result(
    id uuid primary key default uuid_generate_v4(),
    user_count int
);

create table user_placement(
    match_resut_id uuid references match_result on delete cascade,
    user_id uuid references match_result on delete cascade,
    user_place integer,
    rank_change integer
);
