defmodule TPS.Chat.Message do
  alias TPS.Chat.Message
  defstruct [:version, :convo, :key, :type, :message]

  def parse_incoming(<<method::8, rest::binary>>) do
    case method do
      0 ->
        {:message, rest}

      1 ->
        {:req, rest}
    end
  end

  def parse_message(
        <<version::8, key_len::8, key::binary-size(key_len), convo::binary-size(36), type::8,
          message::binary>>
      ) do
    {:message, %Message{version: version, convo: convo, key: key, type: type, message: message}}
  end

  def parse_message(_) do
    %Message{key: "SERVER", type: 0, message: "Incorrect message, you wanker"}
  end

  def parse_request(req) do
    req
    |> String.split(";", trim: true)
    |> Enum.map(fn x -> x |> String.split(":", trim: true) end)
  end
end
