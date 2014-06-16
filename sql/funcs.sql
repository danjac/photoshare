
-- Function: add_tag(character varying)

-- DROP FUNCTION add_tag(character varying);

CREATE OR REPLACE FUNCTION add_tag(name character varying)
  RETURNS bigint AS
WITH s AS (
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
  LANGUAGE sql VOLATILE
  COST 100;
ALTER FUNCTION add_tag(character varying)
  OWNER TO postgres;


-- Function: add_tags(bigint, character varying[])

-- DROP FUNCTION add_tags(bigint, character varying[]);

CREATE OR REPLACE FUNCTION add_tags(IN pid bigint, VARIADIC names character varying[])
  RETURNS void AS
$BODY$DECLARE 
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
END;$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION add_tags(bigint, character varying[])
  OWNER TO postgres;
