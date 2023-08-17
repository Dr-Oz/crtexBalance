CREATE SEQUENCE id_sequence
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL,
    balance bigint NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.services
(
    id bigint NOT NULL,
    title character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT services_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.logs
(
    id bigint NOT NULL DEFAULT nextval('id_sequence'::regclass),
    user_id bigint NOT NULL,
    date date NOT NULL,
    description character varying(100) COLLATE pg_catalog."default" NOT NULL,
    amount bigint NOT NULL,
    CONSTRAINT report_pkey PRIMARY KEY (id),
    CONSTRAINT "user" FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

INSERT INTO public.services(
	id, title)
	VALUES (1, 'Услуга 1');

INSERT INTO public.services(
	id, title)
	VALUES (2, 'Услуга 2');

INSERT INTO public.services(
	id, title)
	VALUES (3, 'Услуга 3');

INSERT INTO public.services(
	id, title)
	VALUES (4, 'Услуга 4');

INSERT INTO public.services(
	id, title)
	VALUES (5, 'Услуга 5');
