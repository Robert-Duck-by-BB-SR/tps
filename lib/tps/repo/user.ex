defmodule TPS.Repo.User do
  def new, do: "insert into user values (?1, ?2, ?3)"
  def user, do: "select * from user where key=?1"
end
