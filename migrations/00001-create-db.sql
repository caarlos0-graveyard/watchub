create table if not exists tokens(
	id bigint primary key,
	token text not null,
	next timestamp,
	disabled boolean not null default false,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp
);

create index if not exists idx_tokens_id on tokens(id);
create index if not exists idx_tokens_next on tokens(next);

create table if not exists users(
	id integer primary key autoincrement,
	login text not null unique,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp
);

create table if not exists repositories(
	id bigint primary key,
	name text not null unique,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp
);

create table if not exists starred_repositories(
	id integer primary key autoincrement,
	user_id bigint not null,
	repository_id bigint not null,
	stargazer_id integer not null,

	foreign key (user_id) references tokens (id) on delete cascade on update no action,
	foreign key (repository_id) references repositories (id) on delete cascade on update no action,
	foreign key (stargazer_id) references users (id) on delete cascade on update no action
);

create table if not exists user_followers(
	id integer primary key autoincrement,
	user_id bigint not null,
	follower_id integer not null,

	foreign key (user_id) references tokens (id) on delete cascade on update no action,
	foreign key (follower_id) references users (id) on delete cascade on update no action
);
