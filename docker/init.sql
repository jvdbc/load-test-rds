CREATE TABLE IF NOT EXISTS public.agents (
	id serial NOT NULL,
	name text NOT NULL,
	created timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT agents_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.orders (
	id serial NOT NULL,
	content text NOT NULL,
	created timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	agent_id serial NOT NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id),
    FOREIGN KEY (agent_id) REFERENCES public.agents(id)
);

INSERT INTO public.agents(name) VALUES('Agent 1');