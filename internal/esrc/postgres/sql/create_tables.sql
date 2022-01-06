CREATE OR REPLACE FUNCTION update_modified_column() 
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.modified = now();
        RETURN NEW; 
    END;
    $$ language 'plpgsql';

CREATE EXTENSION "pgcrypto";

CREATE TABLE IF NOT EXISTS event_streams(
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    type VARCHAR (50) NOT NULL,
    version INT NOT NULL,
    created timestamp default current_timestamp,
    modified timestamp);

CREATE TRIGGER update_event_stream_modtime BEFORE UPDATE ON event_streams FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();

CREATE TABLE IF NOT EXISTS events(
    id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version INT NOT NULL,
    event_source_id uuid NOT NULL,
    data VARCHAR NOT NULL,
    published BOOLEAN NOT NULL DEFAULT false,
    created timestamp default current_timestamp,
    modified timestamp,
    CONSTRAINT fk_event_source
        FOREIGN KEY(event_source_id) 
            REFERENCES event_streams(id)
    );

CREATE TRIGGER update_event_modtime BEFORE UPDATE ON events FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE IF NOT EXISTS snapshots(
    event_source_id INT NOT NULL,
    version INT NOT NULL,
    data VARCHAR NOT NULL,
    created timestamp default current_timestamp,
    modified timestamp);

CREATE TRIGGER update_snapshot_modtime BEFORE UPDATE ON snapshots FOR EACH ROW EXECUTE PROCEDURE update_modified_column();