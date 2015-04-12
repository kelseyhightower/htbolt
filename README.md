# htbolt - Manipulate boltdb password databases

htbolt is used to manipulate a boltdb database used to store usernames and passwords for basic authentication of HTTP users. See the [passwd package docs](passwd) for usage with a HTTP server.

All passwords are stored using the bcyrpt hashing algorithm.

## Usage

### Add or update a user entry.

```
htbolt -f .boltpasswd -u kelsey -c "basic auth user account"
```

Print the results to stdout and don't update the database:

```
htbolt -n -f .boltpasswd -u kelsey -c "basic auth user account"
```

```
{
  "Comment": "basic auth user account",
  "PasswordHash": "JDJhJDEwJDBwZVhuSmZwMVRNL2EvaEhmWTdrZmUwUXNkenhlOWhiWHJiSmd6djJOSkkzTWdEQ09vNEpl",
  "Username": "kelsey"
}
```

### List all users

```
htbolt -l -f .boltpasswd
```
```
kelsey # basic auth user account 
```

### Verify an username and password.

If valid exit code will be set to 0.

```
htbolt -v -f .boltpasswd -u kelsey
```

### Delete a user

```
htbolt -x -f .boltpasswd -u kelsey
```
