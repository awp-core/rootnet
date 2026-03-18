schema "public" {}

table "users" {
  schema = schema.public
  column "address" {
    type = char(42)
  }
  column "registered_at" {
    type = bigint
    null = false
  }
  primary_key {
    columns = [column.address]
  }
}

table "agents" {
  schema = schema.public
  column "agent_address" {
    type = char(42)
  }
  column "owner_address" {
    type = char(42)
    null = false
  }
  column "is_manager" {
    type    = boolean
    null    = false
    default = false
  }
  column "removed" {
    type    = boolean
    null    = false
    default = false
  }
  column "removed_at" {
    type = bigint
    null = true
  }
  primary_key {
    columns = [column.agent_address]
  }
  index "idx_agents_owner" {
    columns = [column.owner_address]
  }
}

table "user_reward_recipients" {
  schema = schema.public
  column "user_address" {
    type = char(42)
  }
  column "recipient_address" {
    type = char(42)
    null = false
  }
  primary_key {
    columns = [column.user_address]
  }
}

table "subnets" {
  schema = schema.public
  column "subnet_id" {
    type = integer
  }
  column "owner" {
    type = char(42)
    null = false
  }
  column "name" {
    type = varchar(64)
    null = false
  }
  column "symbol" {
    type = varchar(16)
    null = false
  }
  column "metadata_uri" {
    type = text
    null = true
  }
  column "governance_weight" {
    type    = integer
    null    = false
    default = 0
  }
  column "subnet_contract" {
    type = char(42)
    null = false
  }
  column "coordinator_url" {
    type = text
    null = true
  }
  column "alpha_token" {
    type = char(42)
    null = false
  }
  column "lp_pool" {
    type = char(42)
    null = true
  }
  column "status" {
    type    = varchar(16)
    null    = false
    default = "Pending"
  }
  column "created_at" {
    type = bigint
    null = false
  }
  column "activated_at" {
    type = bigint
    null = true
  }
  column "immunity_ends_at" {
    type = bigint
    null = true
  }
  column "burned" {
    type    = boolean
    null    = false
    default = false
  }
  primary_key {
    columns = [column.subnet_id]
  }
  index "idx_subnets_owner" {
    columns = [column.owner]
  }
  index "idx_subnets_status" {
    columns = [column.status]
  }
}

table "stake_allocations" {
  schema = schema.public
  column "user_address" {
    type = char(42)
    null = false
  }
  column "agent_address" {
    type = char(42)
    null = false
  }
  column "subnet_id" {
    type = integer
    null = false
  }
  column "amount" {
    type    = numeric(78,0)
    null    = false
    default = 0
  }
  column "frozen" {
    type    = boolean
    null    = false
    default = false
  }
  primary_key {
    columns = [column.user_address, column.agent_address, column.subnet_id]
  }
  index "idx_sa_subnet" {
    columns = [column.subnet_id]
  }
}

table "user_balances" {
  schema = schema.public
  column "user_address" {
    type = char(42)
  }
  column "total_balance" {
    type    = numeric(78,0)
    null    = false
    default = 0
  }
  column "total_allocated" {
    type    = numeric(78,0)
    null    = false
    default = 0
  }
  primary_key {
    columns = [column.user_address]
  }
}

table "withdraw_requests" {
  schema = schema.public
  column "user_address" {
    type = char(42)
  }
  column "amount" {
    type = numeric(78,0)
    null = false
  }
  column "available_at" {
    type = bigint
    null = false
  }
  primary_key {
    columns = [column.user_address]
  }
}

table "epochs" {
  schema = schema.public
  column "epoch_id" {
    type = integer
  }
  column "start_time" {
    type = bigint
    null = false
  }
  column "daily_emission" {
    type = numeric(78,0)
    null = false
  }
  column "subnet_emission" {
    type = numeric(78,0)
    null = true
  }
  column "dao_emission" {
    type = numeric(78,0)
    null = true
  }
  primary_key {
    columns = [column.epoch_id]
  }
}

table "subnet_awp_distributions" {
  schema = schema.public
  column "id" {
    type = serial
  }
  column "epoch_id" {
    type = integer
    null = false
  }
  column "subnet_id" {
    type = integer
    null = false
  }
  column "awp_amount" {
    type = numeric(78,0)
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_sad_epoch" {
    columns = [column.epoch_id]
  }
  index "idx_sad_subnet" {
    columns = [column.subnet_id]
  }
}

table "proposals" {
  schema = schema.public
  column "proposal_id" {
    type = varchar(66)
  }
  column "proposer" {
    type = char(42)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "status" {
    type = varchar(16)
    null = false
  }
  column "votes_for" {
    type    = numeric(78,0)
    null    = false
    default = 0
  }
  column "votes_against" {
    type    = numeric(78,0)
    null    = false
    default = 0
  }
  primary_key {
    columns = [column.proposal_id]
  }
  index "idx_proposals_proposer" {
    columns = [column.proposer]
  }
}

table "sync_states" {
  schema = schema.public
  column "contract_name" {
    type = varchar(64)
  }
  column "last_block" {
    type    = bigint
    null    = false
    default = 0
  }
  primary_key {
    columns = [column.contract_name]
  }
}
