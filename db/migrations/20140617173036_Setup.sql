-- +goose Up
-- +goose StatementBegin
--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;

--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: add_tag(character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION add_tag(name character varying) RETURNS bigint
    LANGUAGE sql
    AS $_$WITH s AS (
    SELECT id
    FROM tags
    WHERE name=$1
 ), i AS (
INSERT INTO tags (name) 
SELECT $1
WHERE NOT EXISTS (
    (SELECT 1 FROM s)
    )
    RETURNING id
    )
    SELECT id 
    FROM i 
    UNION ALL
    SELECT id
    FROM s;
$_$;


ALTER FUNCTION public.add_tag(name character varying) OWNER TO postgres;

--
-- Name: add_tags(bigint, character varying[]); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION add_tags(pid bigint, VARIADIC names character varying[]) RETURNS void
    LANGUAGE plpgsql
    AS $$DECLARE 
tag VARCHAR(200);
tid BIGINT;
BEGIN
DELETE FROM photo_tags WHERE photo_id=pid;
FOREACH tag IN ARRAY names
LOOP
     tid := add_tag(tag);
        
        IF (SELECT 1 FROM photo_tags WHERE photo_id=pid AND tag_id=tid) IS NULL THEN
      
		INSERT INTO photo_tags(photo_id, tag_id) VALUES(pid, tid);
        END IF;
END LOOP;
RETURN;
END;$$;


ALTER FUNCTION public.add_tags(pid bigint, VARIADIC names character varying[]) OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: photo_tags; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE photo_tags (
    photo_id bigint NOT NULL,
    tag_id bigint NOT NULL
);


ALTER TABLE public.photo_tags OWNER TO postgres;

--
-- Name: photos; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE photos (
    id integer NOT NULL,
    owner_id integer,
    created_at timestamp with time zone,
    title text,
    photo text
);


ALTER TABLE public.photos OWNER TO postgres;

--
-- Name: photos_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE photos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.photos_id_seq OWNER TO postgres;

--
-- Name: photos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE photos_id_seq OWNED BY photos.id;


--
-- Name: tag_counts; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE tag_counts (
    id bigint,
    name text,
    num_photos bigint,
    photo text
);


ALTER TABLE public.tag_counts OWNER TO postgres;

--
-- Name: tags; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE tags (
    id bigint NOT NULL,
    name text
);


ALTER TABLE public.tags OWNER TO postgres;

--
-- Name: tags_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tags_id_seq OWNER TO postgres;

--
-- Name: tags_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE tags_id_seq OWNED BY tags.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE users (
    id integer NOT NULL,
    created_at timestamp with time zone,
    name text,
    password text,
    email text,
    admin boolean,
    active boolean
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY photos ALTER COLUMN id SET DEFAULT nextval('photos_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY tags ALTER COLUMN id SET DEFAULT nextval('tags_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Name: photo_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY photo_tags
    ADD CONSTRAINT photo_tags_pkey PRIMARY KEY (tag_id, photo_id);


--
-- Name: photos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY photos
    ADD CONSTRAINT photos_pkey PRIMARY KEY (id);


--
-- Name: tags_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY tags
    ADD CONSTRAINT tags_name_key UNIQUE (name);


--
-- Name: tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: fki_photo_id_fkey; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX fki_photo_id_fkey ON photo_tags USING btree (photo_id);


--
-- Name: _RETURN; Type: RULE; Schema: public; Owner: postgres
--

CREATE RULE "_RETURN" AS ON SELECT TO tag_counts DO INSTEAD SELECT t.id, t.name, (SELECT count(*) AS count FROM photo_tags pt WHERE (t.id = pt.tag_id)) AS num_photos, (SELECT p.photo FROM (photos p JOIN photo_tags pt ON ((pt.photo_id = p.id))) WHERE (pt.tag_id = t.id) ORDER BY p.created_at DESC LIMIT 1) AS photo FROM tags t GROUP BY t.id HAVING ((SELECT count(*) AS count FROM photo_tags pt WHERE (t.id = pt.tag_id)) > 0) ORDER BY (SELECT count(*) AS count FROM photo_tags pt WHERE (t.id = pt.tag_id)) DESC;


--
-- Name: owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY photos
    ADD CONSTRAINT owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id);


--
-- Name: photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY photo_tags
    ADD CONSTRAINT photo_id_fkey FOREIGN KEY (photo_id) REFERENCES photos(id) ON DELETE CASCADE;


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--
-- +goose StatementEnd
