
// 创建用户表
CREATE TABLE public.t_user (
	id serial NOT NULL,
	username varchar(150) NOT NULL,
	"password" varchar(128) NOT NULL,
	email varchar(254) NOT NULL,
	phone varchar(254) NOT NULL,
	CONSTRAINT t_user_key PRIMARY KEY (id),
	CONSTRAINT t_user_username_key UNIQUE (username)
);

// 创建服务端ip端口记录表
CREATE TABLE public.t_server (
	id serial NOT NULL,
	ipaddr varchar(128) not NULL,
	port varchar(16) not null,
	userid int8 not null,
	CONSTRAINT pk_t_server PRIMARY KEY (id)
);
ALTER TABLE public.t_server ADD CONSTRAINT fk_t_user_reference_t_server 
FOREIGN KEY (userid) REFERENCES t_user(id) ON UPDATE RESTRICT ON DELETE RESTRICT;

// 创建客户端ip端口记录表
CREATE TABLE public.t_client (
	id serial NOT NULL,
	ipaddr varchar(128) not NULL,
	port varchar(16) not null,
	userid int8 not null,
	CONSTRAINT pk_t_client PRIMARY KEY (id)
);
ALTER TABLE public.t_client ADD CONSTRAINT fk_t_user_reference_t_client 
FOREIGN KEY (userid) REFERENCES t_user(id) ON UPDATE RESTRICT ON DELETE RESTRICT;