
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
