defmodule Tps.Chat.Message do
  alias Tps.Chat.Message
  defstruct [:key, :type, :message]

  def parse_message(<<type::1, key::7, message::binary>>) do
    IO.inspect("parse_message #{type}, #{key}, #{message}")
    %Message{key: key, type: type, message: message}
  end

  def parse_message(_) do
    %Message{key: "SERVER", type: 0, message: "Incorrect message, you wanker"}
  end
end
