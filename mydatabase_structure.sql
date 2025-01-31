CREATE DATABASE "fitness_center_db_";
CREATE USER "fitness_center_admin" WITH PASSWORD 'FitnessCenter2024';
GRANT ALL PRIVILEGES ON DATABASE "fitness_center_db_" TO "fitness_center_admin";
\c "fitness_center_db_";
GRANT ALL PRIVILEGES ON SCHEMA public TO "fitness_center_admin";

CREATE TABLE public.abonement (
    id uuid NOT NULL,
    title character varying(255),
    validity character varying(255),
    visiting_time character varying(255),
    photo character varying(255),
    price integer,
    created_time timestamp with time zone,
    updated_time timestamp with time zone,
    stripe_price_id character varying
);

ALTER TABLE public.abonement OWNER TO fitness_center_admin;

CREATE TABLE public.abonement_service (
    abonement_id uuid NOT NULL,
    service_id uuid NOT NULL
);

ALTER TABLE public.abonement_service OWNER TO fitness_center_admin;

CREATE TABLE public.coach (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    description character varying(255) NOT NULL,
    photo character varying(255),
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


ALTER TABLE public.coach OWNER TO fitness_center_admin;

CREATE TABLE public.coach_review (
    coach_id uuid NOT NULL,
    review_id uuid NOT NULL
);


ALTER TABLE public.coach_review OWNER TO fitness_center_admin;

CREATE TABLE public.coach_service (
    coach_id uuid NOT NULL,
    service_id uuid NOT NULL
);


ALTER TABLE public.coach_service OWNER TO fitness_center_admin;

CREATE TABLE public."order" (
    id uuid NOT NULL,
    abonement_id uuid,
    user_id uuid,
    status character varying NOT NULL,
    created_time timestamp with time zone,
    updated_time timestamp with time zone,
    expiration_time timestamp with time zone
);


ALTER TABLE public."order" OWNER TO fitness_center_admin;

CREATE TABLE public.refresh_sessions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    refresh_token character varying(400) NOT NULL,
    finger_print character varying(32) NOT NULL,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);


ALTER TABLE public.refresh_sessions OWNER TO fitness_center_admin;

CREATE TABLE public.review (
    id uuid NOT NULL,
    body character varying(255) NOT NULL,
    created_time timestamp with time zone,
    updated_time timestamp with time zone,
    user_id uuid
);

ALTER TABLE public.review OWNER TO fitness_center_admin;

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);

ALTER TABLE public.schema_migrations OWNER TO fitness_center_admin;

CREATE TABLE public.service (
    id uuid NOT NULL,
    title character varying(255),
    photo character varying(255),
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);

ALTER TABLE public.service OWNER TO fitness_center_admin;

CREATE TABLE public."user" (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    role character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    photo character varying(255),
    name character varying(255) NOT NULL,
    created_time timestamp with time zone,
    updated_time timestamp with time zone
);

ALTER TABLE public."user" OWNER TO fitness_center_admin;

ALTER TABLE ONLY public.abonement
    ADD CONSTRAINT abonement_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.abonement_service
    ADD CONSTRAINT abonement_service_pkey PRIMARY KEY (abonement_id, service_id);

ALTER TABLE ONLY public.coach
    ADD CONSTRAINT coach_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.coach_review
    ADD CONSTRAINT coach_review_pkey PRIMARY KEY (coach_id, review_id);

ALTER TABLE ONLY public.coach_service
    ADD CONSTRAINT coach_service_pkey PRIMARY KEY (coach_id, service_id);

ALTER TABLE ONLY public.review
    ADD CONSTRAINT comment_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public."order"
    ADD CONSTRAINT order_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.refresh_sessions
    ADD CONSTRAINT refresh_sessions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);

ALTER TABLE ONLY public.service
    ADD CONSTRAINT service_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.abonement_service
    ADD CONSTRAINT abonement_service_abonement_id_fkey FOREIGN KEY (abonement_id) REFERENCES public.abonement(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.abonement_service
    ADD CONSTRAINT abonement_service_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.service(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.coach_review
    ADD CONSTRAINT coach_review_coach_id_fkey FOREIGN KEY (coach_id) REFERENCES public.coach(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.coach_review
    ADD CONSTRAINT coach_review_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.review(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.coach_service
    ADD CONSTRAINT coach_service_coach_id_fkey FOREIGN KEY (coach_id) REFERENCES public.coach(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.coach_service
    ADD CONSTRAINT coach_service_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.service(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.review
    ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;

ALTER TABLE ONLY public."order"
    ADD CONSTRAINT order_abonement_id_fkey FOREIGN KEY (abonement_id) REFERENCES public.abonement(id) ON DELETE CASCADE;

ALTER TABLE ONLY public."order"
    ADD CONSTRAINT order_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.refresh_sessions
    ADD CONSTRAINT refresh_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;

GRANT ALL ON SCHEMA public TO fitness_center_admin;