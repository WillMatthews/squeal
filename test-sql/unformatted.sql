CREATE TABLE foo_bar_baz (
    id SERIAL PRIMARY KEY,
    foo VARCHAR(255) NOT NULL,
    bar INTEGER NOT NULL,
    baz DATE,
    qux BOOLEAN DEFAULT FALSE
);

SELECT foo, bar, baz
FROM foo_bar_baz
WHERE qux = TRUE
ORDER BY baz DESC;

-- ALTER TABLE foo_bar_baz ADD COLUMN quux TEXT;

-- DROP TABLE foo_bar_baz;

-- TRUNCATE TABLE foo_bar_baz;

-- UPDATE foo_bar_baz SET qux = FALSE WHERE qux = TRUE;

-- CREATE OR REPLACE FUNCTION foo_bar_trigger() RETURNS TRIGGER AS $$ BEGIN
--     IF NEW.qux = TRUE THEN
--         RAISE EXCEPTION 'qux cannot be true';
--     END IF;
--     RETURN NEW;
-- END; $$ LANGUAGE plpgsql;

-- CREATE TRIGGER foo_bar_trigger BEFORE INSERT OR UPDATE ON foo_bar_baz FOR EACH ROW EXECUTE FUNCTION foo_bar_trigger();

-- CREATE INDEX foo_bar_baz_foo_idx ON foo_bar_baz (foo);

-- CREATE UNIQUE INDEX foo_bar_baz_foo_bar_idx ON foo_bar_baz (foo, bar);
