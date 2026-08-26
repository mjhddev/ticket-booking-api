CREATE TYPE event_status AS ENUM (
    'draft',
    'published',
    'cancelled'
);

CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,

    organizer_id BIGINT NOT NULL,

    title VARCHAR(150) NOT NULL,

    description TEXT,

    location VARCHAR(255) NOT NULL,

    start_time TIMESTAMP NOT NULL,

    end_time TIMESTAMP NOT NULL,

    status event_status NOT NULL DEFAULT 'draft',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_events_organizer
        FOREIGN KEY (organizer_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT check_event_time
        CHECK (end_time > start_time)
);