defmodule Mix.Tasks.Tps.Users do
  use Mix.Task
  require Bitwise
  @shortdoc "Create a user"
  @moduledoc """
  Provides tasks for managing TPS users.

  ## Commands
    - `mix tps.users create <username>`: Creates a user with the specified username. Returns a key that user should 
    place in their config file.
  """
  alias TPS.Repo
  alias TPS.Supervisor

  def run(["create", username]) do
    Mix.shell().info("Creating a user")
    System.get_env("serv_key")
    Supervisor.start_link([])

    new_id = UUID.uuid4()
    [_, server_secret] = File.read!(".env") |> String.split("=", trim: true)
    encoded_username = TPS.NinjaHash.encode(username, server_secret)

    new_key =
      :crypto.hash(:sha256, encoded_username |> Enum.join(new_id))

    Repo.query_raw(:get, Repo.User.new(), [new_id, username, new_key])
    |> IO.puts()

    IO.inspect(Base.encode16(new_key))
  end
end
