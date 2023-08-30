table "comments" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "post_id" {
    null = false
    type = integer
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "body" {
    null = false
    type = text
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "post_id_fk" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "user_id_fk" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

table "follows" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "follower_id" {
    null = false
    type = integer
  }
  column "followee_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "followee_id_fk" {
    columns     = [column.followee_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "follower_id_fk" {
    columns     = [column.follower_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "follows_follower_id_followee_id_key" {
    unique  = true
    columns = [column.follower_id, column.followee_id]
  }
}

table "post_likes" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "post_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "post_id_fk" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "user_id_fk" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "post_likes_user_id_post_id_key" {
    unique  = true
    columns = [column.user_id, column.post_id]
  }
}

table "posts" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "image" {
    null = false
    type = character_varying(255)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "user_id" {
    null    = false
    type    = integer
    default = 0
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "user_id_fk" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "avatar" {
    null = false
    type = character_varying(255)
  }
  column "username" {
    null = false
    type = character_varying(14)
  }
  column "email" {
    null = false
    type = character_varying(255)
  }
  column "password" {
    null = false
    type = character_varying(128)
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "email_unique" {
    unique  = true
    columns = [column.email]
  }
  index "username_unique" {
    unique  = true
    columns = [column.username]
  }
}

table "comment_likes" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "comment_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "comment_id_fk" {
    columns     = [column.comment_id]
    ref_columns = [table.comments.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "user_id_fk" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "comment_likes_user_id_comment_id_key" {
    unique  = true
    columns = [column.user_id, column.comment_id]
  }
}

schema "public" {
  comment = "standard public schema"
}
