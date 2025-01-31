CREATE TABLE abonement (
                           id UUID NOT NULL PRIMARY KEY,
                           title VARCHAR(255),
                           validity VARCHAR(255),
                           visiting_time VARCHAR(255),
                           photo VARCHAR(255),
                           price INT
);

CREATE TABLE service (
                         id UUID NOT NULL PRIMARY KEY,
                         title VARCHAR(255),
                         photo VARCHAR(255)
);

CREATE TABLE "user" (
                        id UUID NOT NULL PRIMARY KEY,
                        first_name VARCHAR(255) NOT NULL,
                        last_name VARCHAR(255) NOT NULL,
                        email VARCHAR(255) NOT NULL,
                        role VARCHAR(255) NOT NULL,
                        password_hash VARCHAR(255) NOT NULL,
                        photo VARCHAR(255),
                        created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE coach (
                       id UUID NOT NULL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       description VARCHAR(255) NOT NULL,
                       photo VARCHAR(255)
);

CREATE TABLE abonement_service (
                                   abonement_id UUID NOT NULL,
                                   service_id UUID NOT NULL,
                                   PRIMARY KEY (abonement_id, service_id),
                                   FOREIGN KEY (abonement_id) REFERENCES abonement(id) ON DELETE CASCADE,
                                   FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE
);

CREATE TABLE coach_service (
                               coach_id UUID NOT NULL,
                               service_id UUID NOT NULL,
                               PRIMARY KEY (coach_id, service_id),
                               FOREIGN KEY (coach_id) REFERENCES coach(id) ON DELETE CASCADE,
                               FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE
);

CREATE TABLE comment (
                         id UUID NOT NULL PRIMARY KEY,
                         comment_body VARCHAR(255) NOT NULL,
                         user_id UUID NOT NULL,
                         coach_id UUID NOT NULL,
                         create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
                         FOREIGN KEY (coach_id) REFERENCES coach(id) ON DELETE CASCADE
);

CREATE TABLE "order" (
                         id UUID NOT NULL PRIMARY KEY,
                         abonement_id UUID,
                         user_id UUID,
                         status INT NOT NULL,
                         FOREIGN KEY (abonement_id) REFERENCES abonement(Id),
                         FOREIGN KEY (user_id) REFERENCES "user"(id)
);

CREATE TABLE refresh_sessions (
                                  id UUID NOT NULL PRIMARY KEY,
                                  user_id UUID NOT NULL,
                                  refresh_token VARCHAR(400) NOT NULL,
                                  finger_print VARCHAR(32) NOT NULL,
                                  FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);