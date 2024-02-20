create table if not exists cliente (
    id integer not null,
    limite integer not null,
    saldo integer not null,
    primary key (id)
);

create table if not exists transacao (
    id integer not null,
    descricao varchar(10) not null,
    realizada_em timestamp(6) not null,
    tipo char(1) not null,
    valor integer not null,
    cliente_id integer not null,
    primary key (id)
);
create index if not exists CLIENTE_REALIZADA_EM_INDEX
   on transacao (cliente_id, realizada_em desc);

-- create sequence if not exists transacao_seq start with 1 increment by 50;

alter table if exists transacao
   add constraint FK6cqdtt28hwwinbxxayub0wftw
   foreign key (cliente_id)
   references cliente;

insert into cliente values
    ('1', '100000', '0'),
    ('2', '80000', '0'),
    ('3', '1000000', '0'),
    ('4', '10000000', '0'),
    ('5', '500000', '0');