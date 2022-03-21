# Reminder
*"Too much to remember? I'll remind stuff for you and notify you whenever you want"*

<br>
<br>
<br>

# 🚀 How it works
Uses Redis TTL to create notifications at a specific time, using a shadowed key<small>(*)</small>, to create a event-driven reminder.
<br>
- <small>(*why the shadowed key strategy: https://stackoverflow.com/questions/18328058/redis-notifications-get-key-and-value-on-expiration
</small>)

![shadow key approach](/assets/shadow_key.png "shadow_key_approach")

<br>
<br>

# 🔧 How to use
## Setup and Listen for notifications
```
	// Setup Service
	r := NewReminder("redis://redis:6379")

	// Listen for notifications
	r.Listen(func(n notification, err error) {
		if err != nil {
			fmt.Println("Received error", err)
		} else {
			fmt.Println("Received notification!", n)
		}
	})
```


## Create Reminders
### *"Remind Me In"*
- **Example**: This will notify you after `2 seconds` with a `"pay bills"`
    ```
    r.Remind("pay bills").In(time.Seconds * 2)
    ```
	or
    ```
    r.RemindIn("pay bills", time.Seconds * 2)
    ```

### *"Remind Me At"*
- **Example**: This will notify you with a `"practice go"` at ` 2022-04-5 at 21:32:01:000 UTC`
    ```
    r.Remind("practice go").At(time.Date(2022, time.April,5, 21, 34, 01, 0, time.UTC))
    ```
	or
    ```
    r.RemindAt("practice go", time.Date(2022, time.April,5, 21, 34, 01, 0, time.UTC))
    ```

<br>
<br>
<br>


# 🚧 What needs to be done
- Tests for more load
- Fail recovery would be nice
- Performance Check and connections
