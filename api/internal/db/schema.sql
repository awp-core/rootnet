CREATE TABLE users (
    chain_id      BIGINT NOT NULL DEFAULT 0,
    address       CHAR(42) NOT NULL,
    bound_to      VARCHAR(42) NOT NULL DEFAULT '',
    recipient     VARCHAR(42) NOT NULL DEFAULT '',
    registered_at BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (chain_id, address)
);

CREATE INDEX idx_users_bound_to ON users(chain_id, bound_to) WHERE bound_to != '';

CREATE TABLE subnets (
    subnet_id        NUMERIC(78,0) PRIMARY KEY,
    chain_id         BIGINT NOT NULL DEFAULT 0,
    owner            CHAR(42) NOT NULL,
    name             VARCHAR(64) NOT NULL,
    symbol           VARCHAR(16) NOT NULL,
    subnet_contract  CHAR(42) NOT NULL,
    skills_uri       TEXT,
    metadata_uri     TEXT,
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
CREATE INDEX idx_subnets_chain ON subnets(chain_id);
CREATE INDEX idx_subnets_chain_status ON subnets(chain_id, status);

CREATE TABLE stake_allocations (
    chain_id      BIGINT NOT NULL DEFAULT 0,
    user_address  CHAR(42) NOT NULL,
    agent_address CHAR(42) NOT NULL,
    subnet_id     NUMERIC(78,0) NOT NULL,
    amount        NUMERIC(78,0) NOT NULL DEFAULT 0,
    frozen        BOOLEAN NOT NULL DEFAULT FALSE,
    updated_block BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (chain_id, user_address, agent_address, subnet_id)
);
CREATE INDEX idx_sa_subnet ON stake_allocations(subnet_id);
CREATE INDEX idx_sa_frozen ON stake_allocations(chain_id, frozen) WHERE frozen = TRUE;

CREATE TABLE user_balances (
    chain_id        BIGINT NOT NULL DEFAULT 0,
    user_address    CHAR(42) NOT NULL,
    total_allocated NUMERIC(78,0) NOT NULL DEFAULT 0,
    updated_block   BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (chain_id, user_address)
);

CREATE TABLE stake_positions (
    chain_id      BIGINT NOT NULL DEFAULT 0,
    token_id      BIGINT NOT NULL,
    owner         CHAR(42) NOT NULL,
    amount        NUMERIC(78,0) NOT NULL,
    lock_end_time BIGINT NOT NULL,
    created_at    BIGINT NOT NULL,
    burned        BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (chain_id, token_id)
);
CREATE INDEX idx_sp_owner ON stake_positions(owner);

CREATE TABLE vanity_salts (
    id          SERIAL PRIMARY KEY,
    chain_id    BIGINT NOT NULL DEFAULT 0,
    salt        CHAR(66) NOT NULL,              -- bytes32 hex with 0x prefix
    address     CHAR(42) NOT NULL,              -- predicted Alpha token address
    used        BOOLEAN NOT NULL DEFAULT FALSE,
    subnet_id   NUMERIC(78,0),                   -- set when used by a subnet registration
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (chain_id, salt)
);
CREATE INDEX idx_vs_available ON vanity_salts(used) WHERE used = FALSE;

CREATE TABLE epochs (
    chain_id        BIGINT NOT NULL DEFAULT 0,
    epoch_id        BIGINT NOT NULL,
    start_time      BIGINT NOT NULL,
    daily_emission  NUMERIC(78,0) NOT NULL,
    dao_emission    NUMERIC(78,0),
    PRIMARY KEY (chain_id, epoch_id)
);

CREATE TABLE recipient_awp_distributions (
    id         SERIAL PRIMARY KEY,
    chain_id   BIGINT NOT NULL DEFAULT 0,
    epoch_id   BIGINT NOT NULL,
    recipient  CHAR(42) NOT NULL,
    awp_amount NUMERIC(78,0) NOT NULL
);
CREATE UNIQUE INDEX idx_rad_epoch_recipient ON recipient_awp_distributions(chain_id, epoch_id, recipient);
CREATE INDEX idx_rad_recipient ON recipient_awp_distributions(recipient);

CREATE TABLE proposals (
    chain_id      BIGINT NOT NULL DEFAULT 0,
    proposal_id   VARCHAR(66) NOT NULL,
    proposer      CHAR(42) NOT NULL,
    description   TEXT,
    status        VARCHAR(16) NOT NULL,
    votes_for     NUMERIC(78,0) NOT NULL DEFAULT 0,
    votes_against NUMERIC(78,0) NOT NULL DEFAULT 0,
    PRIMARY KEY (chain_id, proposal_id)
);
CREATE INDEX idx_proposals_proposer ON proposals(proposer);

CREATE TABLE sync_states (
    chain_id      BIGINT NOT NULL DEFAULT 0,
    contract_name VARCHAR(64) NOT NULL,
    last_block    BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (chain_id, contract_name)
);

CREATE TABLE chains (
    chain_id      BIGINT PRIMARY KEY,
    name          VARCHAR(64) NOT NULL,
    rpc_url       TEXT NOT NULL,
    dex           VARCHAR(32) NOT NULL DEFAULT '',
    explorer      TEXT NOT NULL DEFAULT '',
    status        VARCHAR(16) NOT NULL DEFAULT 'active',
    awp_registry  CHAR(42) NOT NULL DEFAULT '',
    awp_token     CHAR(42) NOT NULL DEFAULT '',
    awp_emission  CHAR(42) NOT NULL DEFAULT '',
    staking_vault CHAR(42) NOT NULL DEFAULT '',
    stake_nft     CHAR(42) NOT NULL DEFAULT '',
    subnet_nft    CHAR(42) NOT NULL DEFAULT '',
    dao_address   CHAR(42) NOT NULL DEFAULT '',
    lp_manager    CHAR(42) NOT NULL DEFAULT '',
    pool_manager  CHAR(42) NOT NULL DEFAULT '',
    deploy_block  BIGINT NOT NULL DEFAULT 0,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Block hash chain for reorg detection (optimistic indexing with parent hash verification)
CREATE TABLE indexed_blocks (
    chain_id     BIGINT NOT NULL DEFAULT 0,
    block_number BIGINT NOT NULL,
    block_hash   CHAR(66) NOT NULL,  -- 0x-prefixed hex
    PRIMARY KEY (chain_id, block_number)
);
