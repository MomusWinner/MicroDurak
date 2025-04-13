-- +goose Up
create extension if not exists "uuid-ossp";

create table player (
    id uuid primary key default uuid_generate_v4(),
    name varchar(300) not null,
    age smallint not null,
    rating integer not null
);

create table player_auth (
  id uuid primary key default uuid_generate_v4(),
  player_id uuid references player on delete cascade not null,
  email varchar(300) not null,
  password varchar(300) not null
);

create table match_result (
    id uuid primary key default uuid_generate_v4(),
    player_count smallint not null
);

create table player_placement (
    match_result_id uuid references match_result on delete cascade,
    player_id uuid references player on delete cascade,
    player_place smallint not null,
    rank_change integer not null,
	primary key (match_result_id, player_id)
);

-- +goose Down
drop table cascade player, match_result, player_placement, player_auth;
