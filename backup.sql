--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3 (Debian 13.3-1.pgdg100+1)
-- Dumped by pg_dump version 13.3 (Debian 13.3-1.pgdg100+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: persondb; Type: TABLE; Schema: public; Owner: personuser
--

CREATE TABLE public.persondb (
    id uuid NOT NULL,
    salary integer,
    married boolean,
    profession character varying(30)
);


ALTER TABLE public.persondb OWNER TO personuser;

--
-- Name: users; Type: TABLE; Schema: public; Owner: personuser
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    username character varying,
    password character varying,
    refreshtoken character varying
);


ALTER TABLE public.users OWNER TO personuser;

--
-- Data for Name: persondb; Type: TABLE DATA; Schema: public; Owner: personuser
--

COPY public.persondb (id, salary, married, profession) FROM stdin;
aa715651-d339-45bb-83de-d2701b449094	100	t	football player
ebd74a3d-2dc3-4642-aa2e-384499a98d1c	2000	t	trainer
a9b8115f-2272-4187-8a0c-b9b638900f32	500	t	hockey player
c623b2d8-cd5d-4273-803e-213ea29ee4eb	5000	t	waiter
c4f21c39-2975-42f2-9686-abf4c2cf4750	600	t	troll
e03a7548-633f-45fd-b32e-96aabbb3141d	3500	t	gamer
8c0fcb46-7ed5-4da0-8dfd-a256edc24adc	3500	t	gamer
2a5fad93-0a41-49d2-9e54-5fc0b849c15e	3500	t	gamer
38cb5e4b-a83b-4af6-9008-aa91f11a2701	20000	t	restorator
51963429-fb48-4784-8773-ef25ff5f8b70	8888	f	trainer
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: personuser
--

COPY public.users (id, username, password, refreshtoken) FROM stdin;
cac9e029-618d-419d-b28a-9fb3226325c4	oleg	$2a$14$40Ozo0zDDsQ.vyaerB8tKeJ23tCSJLm9SvTVrJJDSDSnAKuGae8KG	$2a$14$1CD2hZ6bDw/NZScIS7tipuI.Zop.nYELGv0R56z9NcwTUdnwZ084W
994621f0-5486-4967-b7fc-b48959ec133f	Jake	$2a$14$1ph1bgddaLiPMT0.4a3NTuJeymm0hXuWDWSgEqrcR3PA0xqFqWXKq	$2a$14$x6bEWfpe9pnR2uSGoMqCWuwut1QWmAFkEaJl/1EY8ZNOPUy4vX8Bi
0755887c-a0f8-4ee7-9aaa-9264eb05cecb	SWAGGER	$2a$14$bnBBxYPI11vT5gN/BUu8Y.BiBefmGh4l2hQ82ClQp1VwyoPjocpla	$2a$14$5cM6oGH8nv.aoBPB16X2MeSB.jonkAhQlAfCMy4zmuBFwKjLr9oj2
9a8c1f2b-4868-4825-80f5-16f4869a3554	Mfffn	$2a$14$jkm4.8a8oXBK/rFaiHYH2OJjDIS3zKvn1ym8A0MrRdDur4E4t0lQm	$2a$14$bUoNGoTG.ArTTwpFKjp1QufZmXjQMhNPSDxYcKlTic0idIShUsiIy
073ab5e1-6c91-4236-887f-7c93bef3b5ac	youarewelcome	$2a$14$MX5JbHlQV6oHB9xJCINmke3MxtIu5Uakij.fucMlBt2gidxEsbCXq	$2a$14$NCS/y45cZgdHQzvRotBWxuW0DpEmIHbKoHSIhRNkTXJrId6Y1INCK
c9f72a53-2056-4324-a64f-6b057d3343f8	eugene	$2a$14$acmK/gTggOpBpkDRJpljj.9ZlvLHpOxH641xj94iXhE0UyIC/JcJK	$2a$14$EmwzlYqwSTmU/OpYkLzt9eoXqNzwU/1sqL3uhjAdMc29HyE7WDmE6
\.


--
-- Name: persondb persondb_pkey; Type: CONSTRAINT; Schema: public; Owner: personuser
--

ALTER TABLE ONLY public.persondb
    ADD CONSTRAINT persondb_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: personuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

