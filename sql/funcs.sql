
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
