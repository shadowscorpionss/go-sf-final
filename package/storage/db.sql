drop table if exists news_comments;
DROP TABLE IF EXISTS news;

CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    link TEXT NOT NULL UNIQUE
);

create table news_comments (
	id serial primary key,
	parent_id bigint,
	post_id bigint not null,
	content text not null,
	pub_time integer default 0,	
	author_name varchar(100) not null,

	CONSTRAINT fk_parent_comment
      FOREIGN KEY(parent_id) 
        REFERENCES news_comments(id),
	constraint fk_post_comment
	  foreign key (post_id)
		references news(id)
	
);