CREATE TABLE booking_items (
    id BIGSERIAL PRIMARY KEY,

    booking_id BIGINT NOT NULL,

    seat_id BIGINT NOT NULL,

    price BIGINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_booking_items_booking
        FOREIGN KEY (booking_id)
        REFERENCES bookings(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_booking_items_seat
        FOREIGN KEY (seat_id)
        REFERENCES seats(id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_booking_seat
        UNIQUE (booking_id, seat_id)
);