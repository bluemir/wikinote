package auth

/*
group:
  unauthoized: guest
  newcomer:
  - viewer
roles:
  admin:
    rules:
    - verbs: # empty means allow all
  editer:
    rules:
	- verbs:
	  - read
	- verbs:
	  - write
	  resources:
	  - editable: true
  viewer:
    rules:
    - verbs:
      - read
      resources:
      - key: value
        key2: ^regexp$
# group automatically bind role having same name
binding:
  "group/guest":
  - viewer
  "user/root":
  - admin
*/
