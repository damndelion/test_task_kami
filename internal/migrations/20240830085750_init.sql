-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE rooms (
                       id text PRIMARY KEY,
                       room_name VARCHAR(255) NOT NULL
);

-- Bookings is many to one relation with rooms, (one room can have many bookings)
CREATE TABLE bookings (
                          id SERIAL PRIMARY KEY,
                          room_id text NOT NULL,
                          start_time TIMESTAMP NOT NULL,
                          end_time TIMESTAMP NOT NULL,
                          CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS rooms;
-- +goose StatementEnd
