defmodule TPS.Repo.Message do
  def new_message, do: "insert into message values (?1, ?2, ?3, ?4, ?5)"

  def messages_in_convo, do: "select * from message where conversation=?1"
end
