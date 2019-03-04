create schema if not exists sumwhere collate utf8mb4_0900_ai_ci;

create table if not exists advertisement
(
	id bigint auto_increment
		primary key,
	image varchar(255) null,
	create_at datetime null,
	update_at datetime null
);

create table if not exists banner
(
	id bigint auto_increment
		primary key,
	image_url varchar(255) null,
	create_at datetime null,
	update_at datetime null
);

create table if not exists `character`
(
	id bigint auto_increment
		primary key,
	type_name varchar(255) null
)
collate=utf8mb4_unicode_ci;

create table if not exists chat_room
(
	id bigint auto_increment
		primary key,
	create_at datetime null,
	update_at datetime null
);

create table if not exists country
(
	id bigint auto_increment
		primary key,
	name varchar(255) null,
	image_url varchar(255) null
);

create table if not exists event
(
	id bigint auto_increment
		primary key,
	image_url varchar(255) null,
	title varchar(255) null,
	text text null,
	start_at datetime null,
	end_at datetime null
);

create table if not exists interest
(
	id bigint auto_increment
		primary key,
	type_name varchar(255) null
)
collate=utf8mb4_unicode_ci;

create table if not exists match_type
(
	id bigint auto_increment
		primary key,
	title varchar(255) null,
	sub_title varchar(255) null,
	image_url varchar(255) null,
	is_enable tinyint(1) null
);

create table if not exists notice
(
	id bigint auto_increment
		primary key,
	title varchar(255) null,
	text text null,
	create_at datetime null,
	update_at datetime null
);

create table if not exists profile
(
	id bigint auto_increment
		primary key,
	user_id bigint null,
	age int null,
	job varchar(255) null,
	character_type text null,
	trip_style_type text null,
	image1 varchar(255) null,
	image2 varchar(255) null,
	image3 varchar(255) null,
	image4 varchar(255) null,
	create_at datetime null,
	update_at datetime null,
	constraint profile_user_id_fk
		unique (user_id),
	constraint profile_user_id_fk
		foreign key (user_id) references ??? ()
			on delete cascade
);

create table if not exists purchase_product
(
	show_name varchar(255) null,
	product_name varchar(255) null,
	increase float null,
	price float null
);

create table if not exists report_type
(
	id bigint auto_increment
		primary key,
	name varchar(255) not null
);

create table if not exists trip_place
(
	id bigint auto_increment
		primary key,
	trip varchar(255) null,
	image_url varchar(255) null,
	country_id bigint null,
	discription varchar(255) null,
	constraint trip_place_country_id_fk
		foreign key (country_id) references country (id)
			on delete cascade
)
collate=utf8mb4_unicode_ci;

create table if not exists trip_style
(
	id bigint auto_increment
		primary key,
	type varchar(255) not null,
	name varchar(255) not null
)
collate=utf8mb4_unicode_ci;

create table if not exists user
(
	id bigint auto_increment
		primary key,
	email varchar(50) null,
	nickname varchar(10) null,
	username varchar(255) null,
	gender varchar(255) null,
	age int(10) null,
	token varchar(255) null,
	main_profile_image varchar(255) null,
	join_type varchar(50) not null,
	sns_id varchar(255) null,
	password varchar(255) null,
	point bigint default 1 null,
	admin tinyint(1) null,
	has_profile tinyint(1) default 0 null,
	created_at datetime null,
	updated_at datetime null,
	kakao_token varchar(255) null,
	facebook_token varchar(255) null,
	kakao_id int null,
	facebook_id varchar(255) null,
	deleted_at datetime null,
	fcm_token varchar(255) null
)
collate=utf8mb4_unicode_ci;

create table if not exists chat_member
(
	id bigint auto_increment
		primary key,
	room_id bigint null,
	user_id bigint null,
	constraint chat_member_chat_room_id_fk
		foreign key (room_id) references chat_room (id)
			on delete cascade,
	constraint chat_member_user_id_fk
		foreign key (user_id) references user (id)
			on delete cascade
);

create table if not exists `match`
(
	id bigint auto_increment
		primary key,
	maker_id bigint null,
	title varchar(100) not null,
	destination varchar(50) not null,
	start_date datetime null,
	end_date datetime null,
	image varchar(100) null,
	constraint match_user_id_fk
		foreign key (maker_id) references user (id)
			on delete cascade
)
collate=utf8mb4_unicode_ci;

create table if not exists match_member
(
	id bigint auto_increment
		primary key,
	match_id bigint null,
	user_id bigint null,
	join_date datetime null,
	constraint match_member_match_id_fk
		foreign key (match_id) references `match` (id)
			on delete cascade
)
collate=utf8mb4_unicode_ci;

create table if not exists purchase_history
(
	id bigint auto_increment
		primary key,
	user_id bigint null,
	message varchar(255) null,
	positive_value tinyint(1) null,
	`key` int null,
	create_at datetime null,
	constraint purchase_history_user_id_fk
		foreign key (user_id) references user (id)
			on delete cascade
);

create table if not exists push
(
	id bigint auto_increment
		primary key,
	user_id bigint null,
	fcm_token varchar(255) null,
	match_alert tinyint(1) default 1 null,
	friend_alert tinyint(1) default 1 null,
	chat_alert tinyint(1) default 1 null,
	event_alert tinyint(1) default 1 null,
	constraint push_user_id_fk
		unique (user_id),
	constraint push_user_id_fk
		foreign key (user_id) references user (id)
			on delete cascade
);

create table if not exists report
(
	id bigint auto_increment
		primary key,
	user_id bigint null,
	target_user_id bigint null,
	report_type bigint null,
	comment text null,
	create_at datetime null,
	update_at datetime null,
	constraint report_report_type_id_fk
		foreign key (report_type) references report_type (id)
			on delete cascade,
	constraint report_user_id_fk
		foreign key (user_id) references user (id)
			on delete set null,
	constraint report_user_id_fk_2
		foreign key (target_user_id) references user (id)
			on delete set null
);

create table if not exists trip
(
	id bigint auto_increment
		primary key,
	user_id bigint not null,
	match_type_id bigint null,
	trip_place_id bigint null,
	region varchar(255) null,
	activity varchar(255) not null,
	gender_type varchar(20) null,
	start_date datetime null,
	end_date datetime null,
	create_at datetime null,
	delete_at datetime null,
	update_at datetime null,
	constraint travel_user_id_fk
		foreign key (user_id) references user (id)
			on delete cascade,
	constraint trip_match_type_id_fk
		foreign key (match_type_id) references match_type (id)
			on delete set null,
	constraint trip_trip_type_id_fk
		foreign key (trip_place_id) references trip_place (id)
			on delete cascade
)
collate=utf8mb4_unicode_ci;

create table if not exists match_request
(
	from_match_id bigint null,
	to_match_id bigint null,
	create_at datetime null,
	constraint match_request_trip_id_fk
		foreign key (from_match_id) references trip (id)
			on delete set null,
	constraint match_request_trip_id_fk_2
		foreign key (to_match_id) references trip (id)
			on delete set null
)
collate=utf8mb4_unicode_ci;

create table if not exists tripmatch_history
(
	user_id bigint null,
	trip_id bigint null,
	created_at datetime null,
	constraint tripmatch_history_trip_id_fk
		foreign key (trip_id) references trip (id)
			on delete cascade,
	constraint tripmatch_history_user_id_fk
		foreign key (user_id) references user (id)
			on delete cascade
)
collate=utf8mb4_unicode_ci;