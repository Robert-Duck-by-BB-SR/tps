defmodule TPS.Repo.User do
  def user, do: "select * from user where key=?1"
end
