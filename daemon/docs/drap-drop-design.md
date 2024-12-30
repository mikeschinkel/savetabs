I have a web API to write in GoLang and Sqlite and I am trying to understand how it should respond in partial error case.

The API moves one or more "links" (represented by a `link_id`) from one group (represented by a `from_id`) to another group (represented by a `to_id`). The API can move multiple links thus meaning it expects an array of `link_ids`.

For `link_id`s, there are three cases:
1. No `link_id`s find a matching link in the DB.
2. Some `link_id`s find a matching link in the DB, but not all.
3. All `link_id`s find a matching link in the DB.

For `from_id`, there are two cases:
2. `from_id` does not find a matching group in the DB.
1. `from_id` finds a matching group in the DB and all `link_ids` are associated with it.
3. `from_id` finds a matching group in the DB and no `link_ids` are associated with it.
4. `from_id` finds a matching group in the DB and some `link_ids` are associated with it.

For `to_id`, there are two cases:
1. `to_id` does not find a matching group in the DB.
2. `to_id` finds a matching group in the DB and no `link_ids` are associated with it.
3. `to_id` finds a matching group in the DB and all `link_ids` are associated with it.
4. `to_id` finds a matching group in the DB and some `link_ids` are associated with it.          

That results in 3*4*4 or 48 potential scenarios where `n-n-n` means case for `link_id`-`from_id`-`to_id`:
        
3-1-2: Normal case, all `link_id`s moved from `from_id` to `to_id`, returns http status 200.

2-1-2: Semi-normal case, some `link_id`s moved from `from_id` to `to_id`, returns http status 202.

2-1-3: Acceptable case, some `link_id`s disassociated with `from_id` vs. moved to `to_id`, returns http status 202.
2-1-4: Acceptable case, some `link_id`s disassociated with `from_id`, some move to `to_id`, returns http status 202.

2-2-2: Acceptable case, some `link_id`s moved or associated with `to_id`, returns http status 202.
2-2-3: Acceptable case, some does nothing, returns http status 202.
2-2-4: Acceptable case, some `link_id`s associated with `to_id`, returns http status 202.

2-3-2: Acceptable case, some `link_id`s associated with `to_id`, returns http status 202.
2-3-3: Acceptable case, some does nothing, returns http status 202.
2-3-4: Acceptable case, some `link_id`s associated with `to_id`, returns http status 202.

2-4-2: Acceptable case, some `link_id`s associated with `from_id` moved or associated with `to_id`, returns http status 202.
2-4-3: Acceptable case, some `link_id`s disassociated with `from_id`, returns http status 202.
2-4-4: Acceptable case, some `link_id`s disassociated with `from_id`, some associated with `to_id`, some moved to `to_id`, returns http status 202.
 

3-1-3: Acceptable case, all `link_id`s disassociated with `from_id`, returns http status 202.
3-1-4: Acceptable case, some `link_id`s disassociated with `from_id`, some move to `to_id`, returns http status 202.

3-2-2: Acceptable case, all `link_id`s moved or associated with `to_id`, returns http status 202.
3-2-3: Acceptable case, does nothing, returns http status 202.
3-2-4: Acceptable case, some `link_id`s associated with `to_id`, returns http status 202.

3-3-2: Acceptable case, all `link_id`s associated with `to_id`, returns http status 202.
3-3-3: Acceptable case, does nothing, returns http status 202.
3-3-4: Acceptable case, some `link_id`s associated with `to_id`, returns http status 202.

3-4-2: Acceptable case, all `link_id`s associated with `from_id` moved or associated with `to_id`, returns http status 202.
3-4-3: Acceptable case, some `link_id`s disassociated with `from_id`, returns http status 202.
3-4-4: Acceptable case, some `link_id`s disassociated with `from_id`, some associated with `to_id`, some moved to `to_id`, returns http status 202.

1-1-1: Error case, no `link_id` to move, returns http status 400.
1-1-2: Error case, no `link_id` to move, returns http status 400.
1-1-3: Error case, no `link_id` to move, returns http status 400.
1-1-4: Error case, no `link_id` to move, returns http status 400.

1-2-1: Error case, no `link_id` to move, returns http status 400.
1-2-2: Error case, no `link_id` to move, returns http status 400.
1-2-3: Error case, no `link_id` to move, returns http status 400.
1-2-4: Error case, no `link_id` to move, returns http status 400.

1-3-1: Error case, no `link_id` to move, returns http status 400.
1-3-2: Error case, no `link_id` to move, returns http status 400.
1-3-3: Error case, no `link_id` to move, returns http status 400.
1-3-4: Error case, no `link_id` to move, returns http status 400.

1-4-1: Error case, no `link_id` to move, returns http status 400.
1-4-2: Error case, no `link_id` to move, returns http status 400.
1-4-3: Error case, no `link_id` to move, returns http status 400.
1-4-4: Error case, no `link_id` to move, returns http status 400.

2-1-1: Error case, no `to_id` to move `link_ids` to, returns http status 400.
2-2-1: Error case, no `to_id` to move `link_ids` to, returns http status 400.
2-3-1: Error case, no `to_id` to move `link_ids` to, returns http status 400.
2-4-1: Error case, no `to_id` to move `link_ids` to, returns http status 400.

3-1-1: Error case, no `to_id` to move `link_ids` to, returns http status 400.
3-2-1: Error case, no `to_id` to move `link_ids` to, returns http status 400.
3-3-1: Error case, no `to_id` to move `link_ids` to, returns http status 400.
3-4-1: Error case, no `to_id` to move `link_ids` to, returns http status 400.










So my questions are:

1. Does this analysis make sense, or did I miss something?
2. For cases 3-1-1 and 3-2-1 should I generate an error instead of moving the links that matches since other links did not match, or is it better for UX — since this API will be called by drag&drop in a web browser — to move the ones I can?
3. If moving the `link_ids` that exist is acceptable in 3-1-1 and 3-2-1, what should be the HTTP status code returned?  200 seems like it does not provide enough information about what happened.

Note that since this will be driven by drag & drop in the web browser, chances of being anyhing other than 1-1-1 or 1-2-1 are unlikely. Only if a user uses the app in multiple tabs deleting in one tab and moving in another tab are any of the others likely to occur.

See: https://chatgpt.com/share/5a6d7853-a634-43ea-bd8c-6710e0d678f5 for outcome of chat.



