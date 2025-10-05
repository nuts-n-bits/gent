# Wire Format Spec



# Proto Format Spec

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
