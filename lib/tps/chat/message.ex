defmodule Tps.Chat.Message do
  alias Tps.Chat.Message
  defstruct [:version, :convo, :key, :type, :message]

  def parse_message(
        <<version::8, key_len::8, key::binary-size(key_len), convo::binary-size(36), type::8,
          message::binary>>
      ) do
    %Message{version: version, convo: convo, key: key, type: type, message: message}
  end

  def parse_message(_) do
    %Message{key: "SERVER", type: 0, message: "Incorrect message, you wanker"}
  end
end
