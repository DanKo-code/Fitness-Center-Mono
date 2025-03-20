CREATE TYPE coach_shift AS ENUM ('1', '2');

alter table coach
    add column shift coach_shift