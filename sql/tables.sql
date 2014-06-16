
-- Table: photo_tags

-- DROP TABLE photo_tags;

CREATE TABLE photo_tags
(
  photo_id bigint NOT NULL,
  tag_id bigint NOT NULL,
  CONSTRAINT photo_tags_pkey PRIMARY KEY (tag_id, photo_id),
  CONSTRAINT photo_id_fkey FOREIGN KEY (photo_id)
      REFERENCES photos (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);
ALTER TABLE photo_tags
  OWNER TO postgres;


-- Table: photos

-- DROP TABLE photos;

CREATE TABLE photos
(
  id serial NOT NULL,
  owner_id integer,
  created_at timestamp with time zone,
  title text,
  photo text,
  CONSTRAINT photos_pkey PRIMARY KEY (id),
  CONSTRAINT owner_id_fkey FOREIGN KEY (owner_id)
      REFERENCES users (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);
ALTER TABLE photos
  OWNER TO postgres;


-- Table: photos

-- DROP TABLE photos;

CREATE TABLE photos
(
  id serial NOT NULL,
  owner_id integer,
  created_at timestamp with time zone,
  title text,
  photo text,
  CONSTRAINT photos_pkey PRIMARY KEY (id),
  CONSTRAINT owner_id_fkey FOREIGN KEY (owner_id)
      REFERENCES users (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);
ALTER TABLE photos
  OWNER TO postgres;


-- Table: photo_tags

-- DROP TABLE photo_tags;

CREATE TABLE photo_tags
(
  photo_id bigint NOT NULL,
  tag_id bigint NOT NULL,
  CONSTRAINT photo_tags_pkey PRIMARY KEY (tag_id, photo_id),
  CONSTRAINT photo_id_fkey FOREIGN KEY (photo_id)
      REFERENCES photos (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);
ALTER TABLE photo_tags
  OWNER TO postgres;

-- Index: fki_photo_id_fkey

-- DROP INDEX fki_photo_id_fkey;

CREATE INDEX fki_photo_id_fkey
  ON photo_tags
  USING btree
  (photo_id);



-- Table: tags

-- DROP TABLE tags;

CREATE TABLE tags
(
  id bigserial NOT NULL,
  name text,
  CONSTRAINT tags_pkey PRIMARY KEY (id),
  CONSTRAINT tags_name_key UNIQUE (name)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE tags
  OWNER TO postgres;


-- Table: users

-- DROP TABLE users;

CREATE TABLE users
(
  id serial NOT NULL,
  created_at timestamp with time zone,
  name text,
  password text,
  email text,
  admin boolean,
  active boolean,
  CONSTRAINT users_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE users
  OWNER TO postgres;
