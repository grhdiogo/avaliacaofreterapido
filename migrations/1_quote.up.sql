CREATE TABLE IF NOT EXISTS tb_quote (
  id varchar(36) not null,
  cpf_cnpj varchar(14),
  address_cep varchar(10) not null,
  raw_request jsonb not null,
  raw_response jsonb not null,
  created_at timestamp not null default now(),
  updated_at timestamp,
  deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS tb_quote_volume (
  quote_id varchar(36) not null,
  category int not null,
  amount int not null,
  unitary_weight decimal(8,2) not null,
  price decimal(10,2) not null,
  sku varchar(255) not null,
  height decimal(8,2),
  width decimal(8,2),
  length decimal(8,2),
  created_at timestamp not null default now(),
  updated_at timestamp,
  deleted_at timestamp
);
