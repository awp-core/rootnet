CREATE TABLE users (
    address       CHAR(42) PRIMARY KEY,
    registered_at BIGINT NOT NULL
);

CREATE TABLE agents (
    agent_address CHAR(42) PRIMARY KEY,
    owner_address CHAR(42) NOT NULL,
    is_manager    BOOLEAN NOT NULL DEFAULT FALSE,
    removed       BOOLEAN NOT NULL DEFAULT FALSE,
    removed_at    BIGINT
);
CREATE INDEX idx_agents_owner ON agents(owner_address);

CREATE TABLE user_reward_recipients (
    user_address      CHAR(42) PRIMARY KEY,
    recipient_address CHAR(42) NOT NULL
);

CREATE TABLE subnets (
    subnet_id        BIGINT PRIMARY KEY,
    owner            CHAR(42) NOT NULL,
    name             VARCHAR(64) NOT NULL,
    symbol           VARCHAR(16) NOT NULL,
    subnet_contract  CHAR(42) NOT NULL,
    skills_uri       TEXT,
    min_stake        NUMERIC(78,0) NOT NULL DEFAULT 0,
    alpha_token      CHAR(42) NOT NULL,
    lp_pool          CHAR(66),
    status           VARCHAR(16) NOT NULL DEFAULT 'Pending',
    created_at       BIGINT NOT NULL,
    activated_at     BIGINT,
    immunity_ends_at BIGINT,
    burned           BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_subnets_owner ON subnets(owner);
CREATE INDEX idx_subnets_status ON subnets(status);

CREATE TABLE stake_allocations (
    user_address  CHAR(42) NOT NULL,
    agent_address CHAR(42) NOT NULL,
    subnet_id     BIGINT NOT NULL,
    amount        NUMERIC(78,0) NOT NULL DEFAULT 0,
    frozen        BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (user_address, agent_address, subnet_id)
);
CREATE INDEX idx_sa_subnet ON stake_allocations(subnet_id);

CREATE TABLE user_balances (
    user_address    CHAR(42) PRIMARY KEY,
    total_allocated NUMERIC(78,0) NOT NULL DEFAULT 0
);

CREATE TABLE stake_positions (
    token_id      BIGINT PRIMARY KEY,
    owner         CHAR(42) NOT NULL,
    amount        NUMERIC(78,0) NOT NULL,
    lock_end_time BIGINT NOT NULL,
    created_at    BIGINT NOT NULL,
    burned        BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_sp_owner ON stake_positions(owner);

CREATE TABLE vanity_salts (
    id          SERIAL PRIMARY KEY,
    salt        CHAR(66) NOT NULL UNIQUE,     -- bytes32 hex with 0x prefix
    address     CHAR(42) NOT NULL,            -- predicted Alpha token address
    used        BOOLEAN NOT NULL DEFAULT FALSE,
    subnet_id   BIGINT,                       -- set when used by a subnet registration
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_vs_available ON vanity_salts(used) WHERE used = FALSE;

CREATE TABLE epochs (
    epoch_id        BIGINT PRIMARY KEY,
    start_time      BIGINT NOT NULL,
    daily_emission  NUMERIC(78,0) NOT NULL,
    dao_emission    NUMERIC(78,0)
);

CREATE TABLE recipient_awp_distributions (
    id         SERIAL PRIMARY KEY,
    epoch_id   BIGINT NOT NULL,
    recipient  CHAR(42) NOT NULL,
    awp_amount NUMERIC(78,0) NOT NULL
);
CREATE UNIQUE INDEX idx_rad_epoch_recipient ON recipient_awp_distributions(epoch_id, recipient);
CREATE INDEX idx_rad_recipient ON recipient_awp_distributions(recipient);

CREATE TABLE proposals (
    proposal_id   VARCHAR(66) PRIMARY KEY,
    proposer      CHAR(42) NOT NULL,
    description   TEXT,
    status        VARCHAR(16) NOT NULL,
    votes_for     NUMERIC(78,0) NOT NULL DEFAULT 0,
    votes_against NUMERIC(78,0) NOT NULL DEFAULT 0
);
CREATE INDEX idx_proposals_proposer ON proposals(proposer);

CREATE TABLE sync_states (
    contract_name VARCHAR(64) PRIMARY KEY,
    last_block    BIGINT NOT NULL DEFAULT 0
);
