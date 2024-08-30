-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO rooms(id, room_name) values
    ('A-101', 'Reading room 1'),
    ('B-101', 'Reading room 2'),
    ('C-101', 'Reading room 3'),
    ('D-101', 'Reading room 4'),
    ('E-101', 'Reading room 5'),
    ('F-101', 'Reading room 6');

INSERT INTO bookings(room_id, start_time, end_time) VALUES
    ('A-101', '2024-09-01 10:00:00', '2024-09-01 12:00:00'),
    ('B-101', '2024-09-02 14:00:00', '2024-09-02 16:00:00'),
    ('C-101', '2024-09-03 09:00:00', '2024-09-03 11:00:00'),
    ('D-101', '2024-09-04 13:00:00', '2024-09-04 15:00:00'),
    ('E-101', '2024-09-05 08:00:00', '2024-09-05 10:00:00'),
    ('F-101', '2024-09-06 11:00:00', '2024-09-06 13:00:00');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM bookings
WHERE room_id IN ('A-101', 'B-101', 'C-101', 'D-101', 'E-101', 'F-101')
  AND start_time IN (
                     '2024-09-01 10:00:00',
                     '2024-09-02 14:00:00',
                     '2024-09-03 09:00:00',
                     '2024-09-04 13:00:00',
                     '2024-09-05 08:00:00',
                     '2024-09-06 11:00:00'
    )
  AND end_time IN (
                   '2024-09-01 12:00:00',
                   '2024-09-02 16:00:00',
                   '2024-09-03 11:00:00',
                   '2024-09-04 15:00:00',
                   '2024-09-05 10:00:00',
                   '2024-09-06 13:00:00'
    );

DELETE FROM rooms
WHERE id IN ('A-101', 'B-101', 'C-101', 'D-101', 'E-101', 'F-101');

-- +goose StatementEnd
