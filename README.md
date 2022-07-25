# SlitheringJake

The basic premise is that the bot will actively build a markov chain based on the chat message. Occasionally it will generate a message from the markov chain and send it to chat.

## Commands

Commands are usable by subscribers and above (VIP, Mods, Broadcaster)

 - **generate** - Generates a sentence (minimum of five words) and outputs it to chat
 - **question** - Uses the configuration option `StreamerName` and generates something that looks like a question directed at that user.
 - **contains [words...]** - Generates a sentence containing the expected words or variants of those words.
 - **quote** - This uses a secondary markov chain based off and is provided as an example of extends the bots commands.
   - As this command is implemented separately from the ChatBot, it uses an extra (optional) argument `-quotes <filepath>` which should provide the path to the quotes file used to create an alternative markov chain.

## Usage (Arguments)

`./bot -config <filepath> -quotes <filepath>`

 - `-config <filepath>` - Path to the configuration file, see [example.config](example.config) for an example. . The options of this JSON file are documented in [source](pkg/chatbot/types.go).
 - `-quotes <filepath>` - This is a logfile for quotes to enable the `quote` command. If you want to use this then the quote file should be a plaintext file, with quotes that'll be used to generate a new chain. Each quote or other message **must** end with a `.`, `!`, or `?`. If this is not present the bot will generate run on sentences as it does not recognize newlines as the end of a statement.




