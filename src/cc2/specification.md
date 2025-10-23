# Wire Format Spec



# Proto Format Spec

## Example

```

// comments start with two forwrad slashes and continues until the end of line

command name-of-command (string[]) { 
        -s --long-name: string[];
}

```

* A proto file consists of a bunch of command blocks
* A command block starts with keyword `command`, followed by an identifier-like literal `name-of-command`
* Following the name of command, is a pair of parenthesis () that specifies the argument type.
  * `command setusername (string) {}` defines a command with name `username` and exactly 1 string argument that follows
    * A buffer of `setusername hinata` will parse successfully.
  * `command setemail (string[]) {}` defines a command with name `setemail` and any number of string arguments.
    * `setemail`, `setemail 1@mail.com` `setemail 1@mail.com 2@mail.com` will all parse.
  * `command setemail () {}` defines a command that does not accept any arguments
* Following the (), is the body of the command that defines the options of that command
  * An option can have many lines, each of them of the form `-s --long: type;`
  * `-s` is the short name of the command, `--long` is the long name of the command. 
  * Both are accepted when appearing in the incoming buffer
  * `setname () { -s --long: string[]; }` defines a repeated string option.
    * `setname -s foo`, `setname -s foo -s bar`, `setname -s foo --long bar` will all parse.

## Identifier-like literal

### Command name

Command names must only use the following charset: 
* Decimal numbers: 0 ~ 9
* Upper and lower case letters: A ~ z
* Dash: `-`
* Underscore: `_`

Restrictions:
* It cannot be empty
* It must start with a letter. 

Normalizing command name to program name when explicit identifier not specified: 
* Convert each dash and underscore to a single space character 0x20
* Insert a space before each uppercase letter
* Convert all continuous spaces to a single space, and trim whitespaces left and right
* Convert the space deliminated fragments to a PascalCase identifier or snake_case identifier

Examples

```
command hello-world { ... } 
        ~~~~~~~~~~~
        command name  
```

### Short option name

Short option names must only use the following charset:
* Decimal numbers: 0 ~ 9
* Upper and lower case letters: A ~ z
* Dash: `-`
* Underscore: `_`

Restictions:
* It cannot be empty
* It must start with a dash. 
* The option name cannot equal to `--`.

### Long option name

Long option name is a short option name with the added restrictions:
* The first non-dash character must be a letter.
* Must include at least one non-dash character.

Normalizing long option name to program name: 
* Convert each dash and underscore to a single space character 0x20
* Insert a space before each uppercase letter
* Convert all continuous spaces to a single space, and trim whitespaces left and right
* Convert the space deliminated fragments to a camelCase or a snake_case identifier depending on the target language.

In the case a long option name is not provided, the short option acquires the long-option-name restrictions, and the 
short option name is normalized into program name instead.

For example, 
* --session-id becomes session_id or sessionId
* --i___e becomes iE or i_e
* --requireAuth becomes requireAuth or require_auth
* --XMLHttpRequest becomes xMLHttpRequest or x_m_l_http_request (pathological, do not imitate)

For consistency, always use either --kebob-case or camelCase for option names.
