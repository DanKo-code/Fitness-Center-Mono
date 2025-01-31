CREATE DATABASE "fitness_center_db";
CREATE USER "fitness_center_admin" WITH PASSWORD 'FitnessCenter2024';
GRANT ALL PRIVILEGES ON DATABASE "fitness_center_db" TO "fitness_center_admin";
\c "fitness_center_db";
GRANT ALL PRIVILEGES ON SCHEMA public TO "fitness_center_admin";