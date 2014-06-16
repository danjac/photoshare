
-- View: tag_counts

-- DROP VIEW tag_counts;

CREATE OR REPLACE VIEW tag_counts AS 
 SELECT t.id, t.name, ( SELECT count(*) AS count
           FROM photo_tags pt
          WHERE t.id = pt.tag_id) AS num_photos, ( SELECT p.photo
           FROM photos p
      JOIN photo_tags pt ON pt.photo_id = p.id
     WHERE pt.tag_id = t.id
     ORDER BY p.created_at DESC
    LIMIT 1) AS photo
   FROM tags t
  GROUP BY t.id
 HAVING (( SELECT count(*) AS count
           FROM photo_tags pt
          WHERE t.id = pt.tag_id)) > 0
  ORDER BY ( SELECT count(*) AS count
           FROM photo_tags pt
          WHERE t.id = pt.tag_id) DESC;

ALTER TABLE tag_counts
  OWNER TO postgres;
