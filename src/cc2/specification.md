# Format Spec



# Parser Spec

## Identifier-like literal

### Explicit command name identifier

Explicit command identifiers must use the following charset:
* Decimal numbers: 0 ~ 9
* Upper and lower case letters: A ~ z

Restrictions:
* It must start with an Uppercase letter.

Normalizing command name to program name: 
* The explicit identifier is supposed to be a PascalCase program identifier without having to be normalized.

### Command name

Command names must only use the following charset: 
* Decimal numbers: 0 ~ 9
* Upper and lower case letters: A ~ z
* Dash: `-`
* Underscore: `_`

Restrictions:
* It cannot be empty

In the case where an explicit command name identifier is not specified, an additional restriction applies:
* It must start with a letter. 

Normalizing command name to program name when explicit identifier not specified: 
* Convert each dash and underscore to a single space character 0x20
* Insert a space before each uppercase letter
* Convert all continuous spaces to a single space, and trim whitespaces left and right
* Convert the space deliminated fragments to a PascalCase identifier

Examples

```
command hello-world as HelloWorld { ... }  // with explicit identifier
        ~~~~~~~~~~~    ^^^^^^^^^^
        command name   explicit command name identifier

command hello-world { ... }  // no explicit identifier, but HelloWorld would be inferred.

command foo as Bar { ... }  // explicit identifier can differ from what would be automatically inferred
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
