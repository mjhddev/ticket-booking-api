CREATE TYPE seat_status AS ENUM (
    'available',
    'reserved',
    'sold'
);

CREATE TABLE seats (
    id BIGSERIAL PRIMARY KEY,

    event_id BIGINT NOT NULL,

    seat_number VARCHAR(20) NOT NULL,

    status seat_status NOT NULL DEFAULT 'available',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_seats_event
        FOREIGN KEY (event_id)
        REFERENCES events(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_seat_event_number
        UNIQUE (event_id, seat_number)
);