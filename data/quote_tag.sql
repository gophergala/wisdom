CREATE TABLE IF NOT EXISTS quote_tags (
    quote_id integer REFERENCES quote(id),
    tag_id integer REFERENCES tag(id),
    PRIMARY KEY (quote_id, tag_id)
);