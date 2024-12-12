defmodule TPS.Repo.Conversation do
  def new_convo, do: "insert into conversation values (?1, ?2)"

  def conversation_by_username, do: "select * from conversation where users like '%?1%'"
end
